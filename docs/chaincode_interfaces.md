## Smart Contract interfaces for developers

GetArgs() [][]byte

GetStringArgs() []string

GetFunctionAndParameters() (string, []string)

GetArgsSlice() ([]byte, error)

GetTxID() string

InvokeChaincode(chaincodeName string, args [][]byte, channel string) pb.Response

GetState(key string) ([]byte, error)

PutState(key string, value []byte) error

DelState(key string) error

GetStateByRange(startKey, endKey string) (StateQueryIteratorInterface, error)

GetStateByPartialCompositeKey(objectType string, keys []string) (StateQueryIteratorInterface, error)

CreateCompositeKey(objectType string, attributes []string) (string, error)

SplitCompositeKey(compositeKey string) (string, []string, error)

GetQueryResult(query string) (StateQueryIteratorInterface, error)

GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error)

GetCreator() ([]byte, error)

GetTransient() (map[string][]byte, error)

GetBinding() ([]byte, error)

GetSignedProposal() (*pb.SignedProposal, error)

GetTxTimestamp() (*timestamp.Timestamp, error)

SetEvent(name string, payload []byte) error

Transfer(to string, balanceType string, amount *big.Int) error

MultiTransfer(trans *kvtranset.KVTranSet) error

GetAccount(address string) (*wallet.Account, error)

GetSender() (string, error)