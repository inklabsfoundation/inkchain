package ink

//GetAccount  provides a way  to get orderers' account.
type Accounts struct {
	ordererType string
	accounts    []string
}

var instance *Accounts

//TBD not account for orderer now

func GetSoloAccounts() []string {
	if instance == nil {
		instance = &Accounts{}
	}
	if instance.ordererType != "solo" {
		instance.accounts = []string{"0x..."}
		instance.ordererType = "solo"
	}
	return instance.accounts
}

func GetKafkaAccounts() []string {
	if instance == nil {
		instance = &Accounts{}
	}
	if instance.ordererType != "kafka" {
		instance.accounts = []string{"0x...", "0x..."}
		instance.ordererType = "kafka"
	}
	return instance.accounts
}
