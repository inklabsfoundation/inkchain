## Smart Contract interfaces for developers

GetArgs() [][]byte
> GetArgs returns the arguments intended for the chaincode Init and Invoke as an array of byte arrays.

GetStringArgs() []string
> GetStringArgs returns the arguments intended for the chaincode Init and Invoke as a string array. Only use GetStringArgs if the client passes arguments intended to be used as strings.

GetFunctionAndParameters() (string, []string)
> GetFunctionAndParameters returns the first argument as the function name and the rest of the arguments as parameters in a string array. Only use GetFunctionAndParameters if the client passes arguments intended to be used as strings.

GetArgsSlice() ([]byte, error)
> GetArgsSlice returns the arguments intended for the chaincode Init and Invoke as a byte array.

GetTxID() string
> GetTxID returns the tx_id of the transaction proposal (see ChannelHeader in protos/common/common.proto)

InvokeChaincode(chaincodeName string, args [][]byte, channel string) pb.Response
> InvokeChaincode locally calls the specified chaincode `Invoke` using the same transaction context; that is, chaincode calling chaincode doesn't create a new transaction message. If the called chaincode is on the same channel, it simply adds the called chaincode read set and write set to the calling transaction. If the called chaincode is on a different channel, only the Response is returned to the calling chaincode; any PutState calls from the called chaincode will not have any effect on the ledger; that is, the called chaincode on a different channel will not have its read set and write set applied to the transaction. Only the calling chaincode's read set and write set will be applied to the transaction. Effectively the called chaincode on a different channel is a `Query`, which does not participate in state validation checks in subsequent commit phase. If `channel` is empty, the caller's channel is assumed.

GetState(key string) ([]byte, error)
> GetState returns the value of the specified `key` from the ledger. Note that GetState doesn't read data from the writeset, which has not been committed to the ledger. In other words, GetState doesn't consider data modified by PutState that has not been committed. If the key does not exist in the state database, (nil, nil) is returned.

PutState(key string, value []byte) error
> PutState puts the specified `key` and `value` into the transaction's writeset as a data-write proposal. PutState doesn't effect the ledger until the transaction is validated and successfully committed. Simple keys must not be an empty string and must not start with null character (0x00), in order to avoid range query collisions with composite keys, which internally get prefixed with 0x00 as composite key namespace.

DelState(key string) error
> DelState records the specified `key` to be deleted in the writeset of the transaction proposal. The `key` and its value will be deleted from the ledger when the transaction is validated and successfully committed.

GetStateByRange(startKey, endKey string) (StateQueryIteratorInterface, error)
>  GetStateByRange returns a range iterator over a set of keys in the ledger. The iterator can be used to iterate over all keys between the startKey (inclusive) and endKey (exclusive). The keys are returned by the iterator in lexical order. Note that startKey and endKey can be empty string, which implies unbounded range query on start or end. Call Close() on the returned StateQueryIteratorInterface object when done. The query is re-executed during validation phase to ensure result set has not changed since transaction endorsement (phantom reads detected).

GetStateByPartialCompositeKey(objectType string, keys []string) (StateQueryIteratorInterface, error)
> GetStateByPartialCompositeKey queries the state in the ledger based on a given partial composite key. This function returns an iterator which can be used to iterate over all composite keys whose prefix matches the given partial composite key. The `objectType` and attributes are expected to have only valid utf8 strings and should not contain U+0000 (nil byte) and U+10FFFF (biggest and unallocated code point). See related functions SplitCompositeKey and CreateCompositeKey. Call Close() on the returned StateQueryIteratorInterface object when done. The query is re-executed during validation phase to ensure result set has not changed since transaction endorsement (phantom reads detected).

CreateCompositeKey(objectType string, attributes []string) (string, error)
> CreateCompositeKey combines the given `attributes` to form a composite key. The objectType and attributes are expected to have only valid utf8 strings and should not contain U+0000 (nil byte) and U+10FFFF (biggest and unallocated code point). The resulting composite key can be used as the key in PutState().

SplitCompositeKey(compositeKey string) (string, []string, error)
> SplitCompositeKey splits the specified key into attributes on which the composite key was formed. Composite keys found during range queries or partial composite key queries can therefore be split into their composite parts.

GetQueryResult(query string) (StateQueryIteratorInterface, error)
> GetQueryResult performs a "rich" query against a state database. It is only supported for state databases that support rich query, e.g.CouchDB. The query string is in the native syntax of the underlying state database. An iterator is returned which can be used to iterate (next) over the query result set. The query is NOT re-executed during validation phase, phantom reads are not detected. That is, other committed transactions may have added, updated, or removed keys that impact the result set, and this would not be detected at validation/commit time.  Applications susceptible to this should therefore not use GetQueryResult as part of transactions that update ledger, and should limit use to read-only chaincode operations.

GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error)
> GetHistoryForKey returns a history of key values across time. For each historic key update, the historic value and associated transaction id and timestamp are returned. The timestamp is the timestamp provided by the client in the proposal header. GetHistoryForKey requires peer configuration core.ledger.history.enableHistoryDatabase to be true. The query is NOT re-executed during validation phase, phantom reads are not detected. That is, other committed transactions may have updated the key concurrently, impacting the result set, and this would not be detected at validation/commit time. Applications susceptible to this should therefore not use GetHistoryForKey as part of transactions that update ledger, and should limit use to read-only chaincode operations.

GetCreator() ([]byte, error)
> GetCreator returns `SignatureHeader.Creator` (e.g. an identity) of the `SignedProposal`. This is the identity of the agent (or user) submitting the transaction.

GetTransient() (map[string][]byte, error)
> GetTransient returns the `ChaincodeProposalPayload.Transient` field. It is a map that contains data (e.g. cryptographic material) that might be used to implement some form of application-level confidentiality. The contents of this field, as prescribed by `ChaincodeProposalPayload`, are supposed to always be omitted from the transaction and excluded from the ledger.

GetBinding() ([]byte, error)
> GetBinding returns the transaction binding.

GetSignedProposal() (*pb.SignedProposal, error)
> GetSignedProposal returns the SignedProposal object, which contains all data elements part of a transaction proposal.

GetTxTimestamp() (*timestamp.Timestamp, error)
> GetTxTimestamp returns the timestamp when the transaction was created. This is taken from the transaction ChannelHeader, therefore it will indicate the client's timestamp, and will have the same value across all endorsers.

SetEvent(name string, payload []byte) error
> SetEvent allows the chaincode to propose an event on the transaction proposal. If the transaction is validated and successfully committed, the event will be delivered to the current event listeners.

Transfer(to string, balanceType string, amount *big.Int) error
> Tranfer implements atomic balance changes. It allows an transaction of a specific tpye of token (e.g., INK) from invoker's account to another one.

MultiTransfer(trans *kvtranset.KVTranSet) error

GetAccount(address string) (*wallet.Account, error)
> GetAccount returns the account information of the given address. Account information includes its address, balances of different kinds of tokens, and a counter.

GetSender() (string, error)
> GetSender returns the sender's address. The address is revealed from his/her signature.