
## overview

The configtxlator tool was created to support reconfiguration independent of SDKs.

The standard usage is expected to be configtxlator:

1. Proto translation
2. Config update computation


The binary will start an http server listening on the designated port(7059) and is now ready to process request.

First prepare some tools:

```bash

apt-get update && apt-get install vim && apt-get install jq
```

## To start the configtxlator server

docker exec -it cli bash

configtxlator start &

## Dynamically modify the maximum number txs and block  of each block / block time


lounch a new terminal.

1. here we will introduce how to re-configuration config.block, first fetch the block and translate it to json.

```bash
ORDERER_CA=/opt/gopath/src/github.com/inklabsfoundation/inkchain/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
CHANNEL_NAME=mychannel

peer channel fetch config config_block.pb -o orderer.example.com:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_CA

curl -X POST --data-binary @config_block.pb http://127.0.0.1:7059/protolator/decode/common.Block > config_block.json
```

2. Extract the config section from the block:

```bash
jq .data.data[0].payload.data.config config_block.json > config.json
```

3. edit the config.json, set the batch size to 20, or reset timeout=5s and saving it as update_config.json

```bash
jq ".channel_group.groups.Orderer.values.BatchSize.value.max_message_count = 20" config.json  > updated_config.json

jq ".channel_group.groups.Orderer.values.BatchTimeout.value.timeout=\"5s\"" config.json > updated_config.json
```

4. Re-encode both the original config, and the updated config into proto:

```bash
curl -X POST --data-binary @config.json http://127.0.0.1:7059/protolator/encode/common.Config > config.pb
curl -X POST --data-binary @updated_config.json http://127.0.0.1:7059/protolator/encode/common.Config > updated_config.pb
```

5. send them to the configtxlator service to compute the config update which transitions between the two.

```bash
curl -X POST -F original=@config.pb -F updated=@updated_config.pb http://127.0.0.1:7059/configtxlator/compute/update-from-configs -F channel=mychannel > config_update.pb
```

6. we decode the ConfigUpdate so that we may work with it as text:
```bash
curl -X POST --data-binary @config_update.pb http://127.0.0.1:7059/protolator/decode/common.ConfigUpdate > config_update.json
```

7. Then, we wrap it in an envelope message:

```bash
echo '{"payload":{"header":{"channel_header":{"channel_id":"mychannel", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' > config_update_as_envelope.json
```

8. Next, convert it back into the proto form of a full fledged config transaction:

```bash
curl -X POST --data-binary @config_update_as_envelope.json http://127.0.0.1:7059/protolator/encode/common.Envelope > config_update_as_envelope.pb
````

9. Finally, submit the config update transaction to ordering to perform a config update.

```bash
CORE_PEER_LOCALMSPID=OrdererMSP
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/inklabsfoundation/inkchain/peer/crypto/ordererOrganizations/example.com/users/Admin@example.com/msp

peer channel update -o orderer.example.com:7050 -c mychannel -f config_update_as_envelope.pb --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA
```

10. verify

export MAXBATCHSIZEPATH=".data.data[0].payload.data.config.channel_group.groups.Orderer.values.BatchSize.value.max_message_count"
export MAXTIMEOUT=".data.data[0].payload.data.config.channel_group.groups.Orderer.values.BatchTimeout.value.timeout"

peer channel fetch config config_new_block.pb -o orderer.example.com:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_CA

curl -X POST --data-binary @config_new_block.pb http://127.0.0.1:7059/protolator/decode/common.Block > config_new_block.json

jq $MAXTIMEOUT config_new_block.json

jq $MAXBATCHSIZEPATH config_new_block.json
