/*
Copyright Ziggurat Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorser

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/golang/protobuf/proto"
	"github.com/inklabsfoundation/inkchain/common/flogging"
	"golang.org/x/net/context"

	"errors"

	"github.com/inklabsfoundation/inkchain/common/policies"
	"github.com/inklabsfoundation/inkchain/common/util"
	"github.com/inklabsfoundation/inkchain/core/chaincode"
	"github.com/inklabsfoundation/inkchain/core/chaincode/shim"
	"github.com/inklabsfoundation/inkchain/core/common/ccprovider"
	"github.com/inklabsfoundation/inkchain/core/common/validation"
	"github.com/inklabsfoundation/inkchain/core/ledger"
	"github.com/inklabsfoundation/inkchain/core/peer"
	"github.com/inklabsfoundation/inkchain/core/policy"
	syscc "github.com/inklabsfoundation/inkchain/core/scc"
	"github.com/inklabsfoundation/inkchain/core/wallet"
	"github.com/inklabsfoundation/inkchain/core/wallet/ink"
	"github.com/inklabsfoundation/inkchain/core/wallet/ink/impl"
	"github.com/inklabsfoundation/inkchain/msp"
	"github.com/inklabsfoundation/inkchain/msp/mgmt"
	"github.com/inklabsfoundation/inkchain/protos/common"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	putils "github.com/inklabsfoundation/inkchain/protos/utils"
)

// >>>>> begin errors section >>>>>
//chaincodeError is a inkchain error signifying error from chaincode
type chaincodeError struct {
	status int32
	msg    string
}

func (ce chaincodeError) Error() string {
	return fmt.Sprintf("chaincode error (status: %d, message: %s)", ce.status, ce.msg)
}

// <<<<< end errors section <<<<<<

var endorserLogger = flogging.MustGetLogger("endorser")

// The Jira issue that documents Endorser flow along with its relationship to

// Endorser provides the Endorser service ProcessProposal
type Endorser struct {
	policyChecker policy.PolicyChecker
	inkCalculator ink.InkAlg
}

// NewEndorserServer creates and returns a new Endorser server instance.
func NewEndorserServer() pb.EndorserServer {
	e := new(Endorser)
	e.policyChecker = policy.NewPolicyChecker(
		peer.NewChannelPolicyManagerGetter(),
		mgmt.GetLocalMSP(),
		mgmt.NewLocalMSPPrincipalGetter(),
	)
	e.inkCalculator = impl.NewSimpleInkAlg()
	return e
}

// checkACL checks that the supplied proposal complies
// with the writers policy of the chain
func (e *Endorser) checkACL(signedProp *pb.SignedProposal, chdr *common.ChannelHeader, shdr *common.SignatureHeader, hdrext *pb.ChaincodeHeaderExtension) error {
	return e.policyChecker.CheckPolicy(chdr.ChannelId, policies.ChannelApplicationWriters, signedProp)
}

//TODO - check for escc and vscc
func (*Endorser) checkEsccAndVscc(prop *pb.Proposal) error {
	return nil
}

//INKCHAIN - check invoke signature
func (*Endorser) checkInvokeSigAndSetSender(cis *pb.ChaincodeInvocationSpec, txsim ledger.TxSimulator) error {
	if cis.Sig == nil || cis.SenderSpec == nil {
		return nil
	}
	hash, err := wallet.GetInvokeHash(cis.ChaincodeSpec, cis.IdGenerationAlg, cis.SenderSpec)
	if err != nil {
		return err
	}
	address, err := wallet.GetSenderFromSignature(hash, cis.Sig)
	if err != nil {
		return err
	}
	sender := address.ToString()
	if string(cis.SenderSpec.Sender[:]) != sender {
		return fmt.Errorf("endorser: invalid signature")
	}
	txsim.SetSender(sender)
	return nil
}

func (e *Endorser) checkCounterAndInk(cis *pb.ChaincodeInvocationSpec, txsim ledger.TxSimulator, byteCount int, account *wallet.Account) error {
	if cis.Sig == nil || cis.SenderSpec == nil {
		return nil
	}
	if cis.SenderSpec.Counter < account.Counter {
		return fmt.Errorf("endorser: invalid counter")
	}
	inkFee, err := e.inkCalculator.CalcInk(byteCount)
	fee := big.NewInt(inkFee)
	if err != nil {
		return fmt.Errorf("endorser: error when calculating ink")
	}
	inkBalance, ok := account.Balance[wallet.MAIN_BALANCE_NAME]
	if !ok || inkBalance.Cmp(fee) < 0 {
		return fmt.Errorf("endorser: insufficient balance for ink consumption")
	}
	inkLimit, ok := new(big.Int).SetString(string(cis.SenderSpec.InkLimit), 10)
	if !ok {
		return fmt.Errorf("endorser: invalid inklimit.")
	}
	if fee.Cmp(inkLimit) > 0 {
		return fmt.Errorf("endorser: fee exceeds inkLimit.")
	}
	if fee.Cmp(wallet.InkMinimumFee) < 0 {
		return fmt.Errorf("endorser: fee too low")
	}
	return nil
}

func (*Endorser) getTxSimulator(ledgername string) (ledger.TxSimulator, error) {
	lgr := peer.GetLedger(ledgername)
	if lgr == nil {
		return nil, fmt.Errorf("channel does not exist: %s", ledgername)
	}
	return lgr.NewTxSimulator()
}

func (*Endorser) getHistoryQueryExecutor(ledgername string) (ledger.HistoryQueryExecutor, error) {
	lgr := peer.GetLedger(ledgername)
	if lgr == nil {
		return nil, fmt.Errorf("channel does not exist: %s", ledgername)
	}
	return lgr.NewHistoryQueryExecutor()
}

//call specified chaincode (system or user)
func (e *Endorser) callChaincode(ctxt context.Context, chainID string, version string, txid string, signedProp *pb.SignedProposal, prop *pb.Proposal, cis *pb.ChaincodeInvocationSpec, cid *pb.ChaincodeID, txsim ledger.TxSimulator) (*pb.Response, *pb.ChaincodeEvent, error) {
	endorserLogger.Debugf("Entry - txid: %s channel id: %s version: %s", txid, chainID, version)
	defer endorserLogger.Debugf("Exit")
	var err error
	var res *pb.Response
	var ccevent *pb.ChaincodeEvent

	if txsim != nil {
		ctxt = context.WithValue(ctxt, chaincode.TXSimulatorKey, txsim)
	}

	//is this a system chaincode
	scc := syscc.IsSysCC(cid.Name)

	cccid := ccprovider.NewCCContext(chainID, cid.Name, version, txid, scc, signedProp, prop)

	res, ccevent, err = chaincode.ExecuteChaincode(ctxt, cccid, cis.ChaincodeSpec.Input.Args)

	if err != nil {
		return nil, nil, err
	}

	//per doc anything < 400 can be sent as TX.
	//inkchain errors will always be >= 400 (ie, unambiguous errors )
	//"lscc" will respond with status 200 or 500 (ie, unambiguous OK or ERROR)
	if res.Status >= shim.ERRORTHRESHOLD {
		return res, nil, nil
	}

	//----- BEGIN -  SECTION THAT MAY NEED TO BE DONE IN LSCC ------
	//if this a call to deploy a chaincode, We need a mechanism
	//to pass TxSimulator into LSCC. Till that is worked out this
	//special code does the actual deploy, upgrade here so as to collect
	//all state under one TxSimulator
	//
	//NOTE that if there's an error all simulation, including the chaincode
	//table changes in lscc will be thrown away
	if cid.Name == "lscc" && len(cis.ChaincodeSpec.Input.Args) >= 3 && (string(cis.ChaincodeSpec.Input.Args[0]) == "deploy" || string(cis.ChaincodeSpec.Input.Args[0]) == "upgrade") {
		var cds *pb.ChaincodeDeploymentSpec
		cds, err = putils.GetChaincodeDeploymentSpec(cis.ChaincodeSpec.Input.Args[2])
		if err != nil {
			return nil, nil, err
		}

		//this should not be a system chaincode
		if syscc.IsSysCC(cds.ChaincodeSpec.ChaincodeId.Name) {
			return nil, nil, fmt.Errorf("attempting to deploy a system chaincode %s/%s", cds.ChaincodeSpec.ChaincodeId.Name, chainID)
		}

		cccid = ccprovider.NewCCContext(chainID, cds.ChaincodeSpec.ChaincodeId.Name, cds.ChaincodeSpec.ChaincodeId.Version, txid, false, signedProp, prop)

		_, _, err = chaincode.Execute(ctxt, cccid, cds)
		if err != nil {
			return nil, nil, fmt.Errorf("%s", err)
		}
	}
	//----- END -------

	return res, ccevent, err
}

//TO BE REMOVED WHEN JAVA CC IS ENABLED
//disableJavaCCInst if trying to install, instantiate or upgrade Java CC
func (e *Endorser) disableJavaCCInst(cid *pb.ChaincodeID, cis *pb.ChaincodeInvocationSpec) error {
	//if not lscc we don't care
	if cid.Name != "lscc" {
		return nil
	}

	//non-nil spec ? leave it to callers to handle error if this is an error
	if cis.ChaincodeSpec == nil || cis.ChaincodeSpec.Input == nil {
		return nil
	}

	//should at least have a command arg, leave it to callers if this is an error
	if len(cis.ChaincodeSpec.Input.Args) < 1 {
		return nil
	}

	var argNo int
	switch string(cis.ChaincodeSpec.Input.Args[0]) {
	case "install":
		argNo = 1
	case "deploy", "upgrade":
		argNo = 2
	default:
		//what else can it be ? leave it caller to handle it if error
		return nil
	}

	if argNo >= len(cis.ChaincodeSpec.Input.Args) {
		return errors.New("Too few arguments passed")
	}

	//the inner dep spec will contain the type
	cds, err := putils.GetChaincodeDeploymentSpec(cis.ChaincodeSpec.Input.Args[argNo])
	if err != nil {
		return err
	}

	//finally, if JAVA error out
	if cds.ChaincodeSpec.Type == pb.ChaincodeSpec_JAVA {
		return fmt.Errorf("Java chaincode is work-in-progress and disabled")
	}

	//not a java install, instantiate or upgrade op
	return nil
}

//simulate the proposal by calling the chaincode
func (e *Endorser) simulateProposal(ctx context.Context, chainID string, txid string, signedProp *pb.SignedProposal, prop *pb.Proposal, cid *pb.ChaincodeID, txsim ledger.TxSimulator) (*ccprovider.ChaincodeData, *pb.Response, []byte, *pb.ChaincodeEvent, error) {
	endorserLogger.Debugf("Entry - txid: %s channel id: %s", txid, chainID)
	defer endorserLogger.Debugf("Exit")
	//we do expect the payload to be a ChaincodeInvocationSpec
	//if we are supporting other payloads in future, this be glaringly point
	//as something that should change
	cis, err := putils.GetChaincodeInvocationSpec(prop)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	//disable Java install,instantiate,upgrade for now
	if err = e.disableJavaCCInst(cid, cis); err != nil {
		return nil, nil, nil, nil, err
	}
	//---1. check Invoke Signature
	if err = e.checkInvokeSigAndSetSender(cis, txsim); err != nil {
		return nil, nil, nil, nil, err
	}
	//---2. check ESCC and VSCC for the chaincode
	if err = e.checkEsccAndVscc(prop); err != nil {
		return nil, nil, nil, nil, err
	}

	var cdLedger *ccprovider.ChaincodeData
	var version string

	if !syscc.IsSysCC(cid.Name) {
		cdLedger, err = e.getCDSFromLSCC(ctx, chainID, txid, signedProp, prop, cid.Name, txsim)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("%s - make sure the chaincode %s has been successfully instantiated and try again", err, cid.Name)
		}
		version = cdLedger.Version

		err = ccprovider.CheckInsantiationPolicy(cid.Name, version, cdLedger)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	} else {
		version = util.GetSysCCVersion()
	}

	//*** added: get account for checking ink & counter. GetState() could not be called after GetTxSimulationResults()
	var account *wallet.Account
	if cis.Sig != nil && cis.SenderSpec != nil {
		account = &wallet.Account{}
		var resBytes []byte
		resBytes, err = txsim.GetState(wallet.WALLET_NAMESPACE, string(cis.SenderSpec.Sender[:]))
		if err != nil {
			return nil, nil, nil, nil, err
		}
		err = json.Unmarshal(resBytes, account)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}

	//---3. execute the proposal and get simulation results
	var simResult []byte
	var res *pb.Response
	var ccevent *pb.ChaincodeEvent
	res, ccevent, err = e.callChaincode(ctx, chainID, version, txid, signedProp, prop, cis, cid, txsim)
	if err != nil {
		endorserLogger.Errorf("failed to invoke chaincode %s on transaction %s, error: %s", cid, txid, err)
		return nil, nil, nil, nil, err
	}

	if txsim != nil {
		if simResult, err = txsim.GetTxSimulationResults(); err != nil {
			return nil, nil, nil, nil, err
		}
	}

	//---4. check counter and ink
	txLength := len(simResult)
	if cis.SenderSpec != nil {
		txLength += len(cis.SenderSpec.String())
	}
	err = e.checkCounterAndInk(cis, txsim, txLength, account)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return cdLedger, res, simResult, ccevent, nil
}

func (e *Endorser) getCDSFromLSCC(ctx context.Context, chainID string, txid string, signedProp *pb.SignedProposal, prop *pb.Proposal, chaincodeID string, txsim ledger.TxSimulator) (*ccprovider.ChaincodeData, error) {
	ctxt := ctx
	if txsim != nil {
		ctxt = context.WithValue(ctx, chaincode.TXSimulatorKey, txsim)
	}

	return chaincode.GetChaincodeDataFromLSCC(ctxt, txid, signedProp, prop, chainID, chaincodeID)
}

//endorse the proposal by calling the ESCC
func (e *Endorser) endorseProposal(ctx context.Context, chainID string, txid string, signedProp *pb.SignedProposal, proposal *pb.Proposal, response *pb.Response, simRes []byte, event *pb.ChaincodeEvent, visibility []byte, ccid *pb.ChaincodeID, txsim ledger.TxSimulator, cd *ccprovider.ChaincodeData) (*pb.ProposalResponse, error) {
	endorserLogger.Debugf("Entry - txid: %s channel id: %s chaincode id: %s", txid, chainID, ccid)
	defer endorserLogger.Debugf("Exit")

	isSysCC := cd == nil
	// 1) extract the name of the escc that is requested to endorse this chaincode
	var escc string
	//ie, not "lscc" or system chaincodes
	if isSysCC {
		// FIXME: getCDSFromLSCC seems to fail for lscc - not sure this is expected?
		// TODO: who should endorse a call to LSCC?
		escc = "escc"
	} else {
		escc = cd.Escc
		if escc == "" { // this should never happen, LSCC always fills this field
			panic("No ESCC specified in ChaincodeData")
		}
	}

	endorserLogger.Debugf("info: escc for chaincode id %s is %s", ccid, escc)

	// marshalling event bytes
	var err error
	var eventBytes []byte
	if event != nil {
		eventBytes, err = putils.GetBytesChaincodeEvent(event)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal event bytes - %s", err)
		}
	}

	resBytes, err := putils.GetBytesResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response bytes - %s", err)
	}

	// set version of executing chaincode
	if isSysCC {
		// if we want to allow mixed inkchain levels we should
		// set syscc version to ""
		ccid.Version = util.GetSysCCVersion()
	} else {
		ccid.Version = cd.Version
	}

	ccidBytes, err := putils.Marshal(ccid)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ChaincodeID - %s", err)
	}

	// 3) call the ESCC we've identified
	// arguments:
	// args[0] - function name (not used now)
	// args[1] - serialized Header object
	// args[2] - serialized ChaincodeProposalPayload object
	// args[3] - ChaincodeID of executing chaincode
	// args[4] - result of executing chaincode
	// args[5] - binary blob of simulation results
	// args[6] - serialized events
	// args[7] - payloadVisibility
	args := [][]byte{[]byte(""), proposal.Header, proposal.Payload, ccidBytes, resBytes, simRes, eventBytes, visibility}
	version := util.GetSysCCVersion()
	ecccis := &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_GOLANG, ChaincodeId: &pb.ChaincodeID{Name: escc}, Input: &pb.ChaincodeInput{Args: args}}}
	res, _, err := e.callChaincode(ctx, chainID, version, txid, signedProp, proposal, ecccis, &pb.ChaincodeID{Name: escc}, txsim)
	if err != nil {
		return nil, err
	}

	if res.Status >= shim.ERRORTHRESHOLD {
		return &pb.ProposalResponse{Response: res}, nil
	}

	prBytes := res.Payload
	// Note that we do not extract any simulation results from
	// the call to ESCC. This is intentional becuse ESCC is meant
	// to endorse (i.e. sign) the simulation results of a chaincode,
	// but it can't obviously sign its own. Furthermore, ESCC runs
	// on private input (its own signing key) and so if it were to
	// produce simulationr results, they are likely to be different
	// from other ESCCs, which would stand in the way of the
	// endorsement process.

	//3 -- respond
	pResp, err := putils.GetProposalResponse(prBytes)
	if err != nil {
		return nil, err
	}

	return pResp, nil
}

// ProcessProposal process the Proposal
func (e *Endorser) ProcessProposal(ctx context.Context, signedProp *pb.SignedProposal) (*pb.ProposalResponse, error) {
	endorserLogger.Debugf("Entry")
	defer endorserLogger.Debugf("Exit")
	// at first, we check whether the message is valid
	prop, hdr, hdrExt, err := validation.ValidateProposalMessage(signedProp)
	if err != nil {
		return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, err
	}

	chdr, err := putils.UnmarshalChannelHeader(hdr.ChannelHeader)
	if err != nil {
		return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, err
	}

	shdr, err := putils.GetSignatureHeader(hdr.SignatureHeader)
	if err != nil {
		return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, err
	}

	// block invocations to security-sensitive system chaincodes
	if syscc.IsSysCCAndNotInvokableExternal(hdrExt.ChaincodeId.Name) {
		endorserLogger.Errorf("Error: an attempt was made by %#v to invoke system chaincode %s",
			shdr.Creator, hdrExt.ChaincodeId.Name)
		err = fmt.Errorf("Chaincode %s cannot be invoked through a proposal", hdrExt.ChaincodeId.Name)
		return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, err
	}

	chainID := chdr.ChannelId

	// Check for uniqueness of prop.TxID with ledger
	// Notice that ValidateProposalMessage has already verified
	// that TxID is computed properly
	txid := chdr.TxId
	if txid == "" {
		err = errors.New("Invalid txID. It must be different from the empty string.")
		return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, err
	}
	endorserLogger.Debugf("processing txid: %s", txid)
	if chainID != "" {
		// here we handle uniqueness check and ACLs for proposals targeting a chain
		lgr := peer.GetLedger(chainID)
		if lgr == nil {
			return nil, fmt.Errorf("failure while looking up the ledger %s", chainID)
		}
		if _, err := lgr.GetTransactionByID(txid); err == nil {
			return nil, fmt.Errorf("Duplicate transaction found [%s]. Creator [%x]. [%s]", txid, shdr.Creator, err)
		}

		// check ACL only for application chaincodes; ACLs
		// for system chaincodes are checked elsewhere
		if !syscc.IsSysCC(hdrExt.ChaincodeId.Name) {
			// check that the proposal complies with the channel's writers
			if err = e.checkACL(signedProp, chdr, shdr, hdrExt); err != nil {
				return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, err
			}
		}
	} else {
		// chainless proposals do not/cannot affect ledger and cannot be submitted as transactions
		// ignore uniqueness checks; also, chainless proposals are not validated using the policies
		// of the chain since by definition there is no chain; they are validated against the local
		// MSP of the peer instead by the call to ValidateProposalMessage above
	}

	// obtaining once the tx simulator for this proposal. This will be nil
	// for chainless proposals
	// Also obtain a history query executor for history queries, since tx simulator does not cover history
	var txsim ledger.TxSimulator
	var historyQueryExecutor ledger.HistoryQueryExecutor
	if chainID != "" {
		if txsim, err = e.getTxSimulator(chainID); err != nil {
			return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, err
		}
		if historyQueryExecutor, err = e.getHistoryQueryExecutor(chainID); err != nil {
			return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, err
		}
		// Add the historyQueryExecutor to context
		// TODO shouldn't we also add txsim to context here as well? Rather than passing txsim parameter
		// around separately, since eventually it gets added to context anyways
		ctx = context.WithValue(ctx, chaincode.HistoryQueryExecutorKey, historyQueryExecutor)

		defer txsim.Done()
	}
	//this could be a request to a chainless SysCC

	// TODO: if the proposal has an extension, it will be of type ChaincodeAction;
	//       if it's present it means that no simulation is to be performed because
	//       we're trying to emulate a submitting peer. On the other hand, we need
	//       to validate the supplied action before endorsing it

	//1 -- simulate
	cd, res, simulationResult, ccevent, err := e.simulateProposal(ctx, chainID, txid, signedProp, prop, hdrExt.ChaincodeId, txsim)
	if err != nil {
		return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, err
	}
	if res != nil {
		if res.Status >= shim.ERROR {
			endorserLogger.Errorf("simulateProposal() resulted in chaincode response status %d for txid: %s", res.Status, txid)
			var cceventBytes []byte
			if ccevent != nil {
				cceventBytes, err = putils.GetBytesChaincodeEvent(ccevent)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal event bytes - %s", err)
				}
			}
			pResp, err := putils.CreateProposalResponseFailure(prop.Header, prop.Payload, res, simulationResult, cceventBytes, hdrExt.ChaincodeId, hdrExt.PayloadVisibility)
			if err != nil {
				return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, err
			}

			return pResp, &chaincodeError{res.Status, res.Message}
		}
	}

	//2 -- endorse and get a marshalled ProposalResponse message
	var pResp *pb.ProposalResponse

	//TODO till we implement global ESCC, CSCC for system chaincodes
	//chainless proposals (such as CSCC) don't have to be endorsed
	if chainID == "" {
		pResp = &pb.ProposalResponse{Response: res}
	} else {
		pResp, err = e.endorseProposal(ctx, chainID, txid, signedProp, prop, res, simulationResult, ccevent, hdrExt.PayloadVisibility, hdrExt.ChaincodeId, txsim, cd)
		if err != nil {
			return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, err
		}
		if pResp != nil {
			if res.Status >= shim.ERRORTHRESHOLD {
				endorserLogger.Debugf("endorseProposal() resulted in chaincode error for txid: %s", txid)
				return pResp, &chaincodeError{res.Status, res.Message}
			}
		}
	}

	// Set the proposal response payload - it
	// contains the "return value" from the
	// chaincode invocation
	pResp.Response.Payload = res.Payload

	return pResp, nil
}

// Only exposed for testing purposes - commit the tx simulation so that
// a deploy transaction is persisted and that chaincode can be invoked.
// This makes the endorser test self-sufficient
func (e *Endorser) commitTxSimulation(proposal *pb.Proposal, chainID string, signer msp.SigningIdentity, pResp *pb.ProposalResponse, blockNumber uint64) error {
	tx, err := putils.CreateSignedTx(proposal, signer, pResp)
	if err != nil {
		return err
	}

	lgr := peer.GetLedger(chainID)
	if lgr == nil {
		return fmt.Errorf("failure while looking up the ledger")
	}

	txBytes, err := proto.Marshal(tx)
	if err != nil {
		return err
	}
	block := common.NewBlock(blockNumber, []byte{}, []byte{}, 0)
	block.Data.Data = [][]byte{txBytes}
	block.Header.DataHash = block.Data.Hash()
	if err = lgr.Commit(block); err != nil {
		return err
	}

	return nil
}
