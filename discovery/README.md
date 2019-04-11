## Service discovery
In order to dynamically obtain informations of the peers and orderers on the network, such as peers status, configurations, ledger, policy, etc. we add the function of service discovery. 

The service discovery improves this process by having the peers compute the needed information dynamically and present it to the SDK in a consumable manner.

## Capabilities of the discovery service

The discovery service can respond to the following queries:

**Configuration query**: Returns the MSPConfig of all organizations in the channel along with the orderer endpoints of the channel.

**Peer membership query**: Returns the peers that have joined the channel.

**Endorsement query**: Returns an endorsement descriptor for given chaincode(s) in a channel.

**Local peer membership query**: Returns the local membership information of the peer that responds to the query. By default the client needs to be an administrator for the peer to respond to this query.

