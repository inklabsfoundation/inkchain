const express = require('express');
const app = express();
const bodyParser = require('body-parser');
const ora = require("ora")
// const parseArgs = require("minimist")

const {
    Qtum
} = require("qtumjs")

const repoData = require("../solar.development")
const qtum = new Qtum("http://qtum:test@localhost:3889", repoData)
const contracts = {
    "INK": "contracts/INK.sol",
    "XC": "contracts/XC.sol",
    "XCPlugin": "contracts/XCPlugin.sol",
}

// application/x-www-form-urlencoded
var urlencodedParser = bodyParser.urlencoded({extended: false})

app.use(express.static('public'));

async function call(contract,method,params,fromAddr) {
    var name = contracts[contract];
    var abi = repoData.contracts[name].abi
    for (var i=0; i < abi.length; i++) {
        var a = abi[i]
        if (a.name == method && a.type == "function") {
            if (a.constant){
                var result = await qtum.contract(name).call(method,params)
                console.log("log:", result.outputs[0])
                return result.outputs
            }else{
                const tx = await qtum.contract(name).send(method,params, {
                    senderAddress: fromAddr,
                })
                console.log("transfer tx:", tx.txid)
                console.log(tx)
                // or: await tx.confirm(1)
                const confirmation = tx.confirm(1)
                ora.promise(confirmation, "confirm transfer")
                await confirmation
            }
        }
    }
}

app.post('/xc/out/', urlencodedParser,async function (req, res) {

    var XC = repoData.contracts[contracts["XC"]]
    var INK = repoData.contracts[contracts["INK"]]

    var result = await call("INK","approve",[XC.address,req.body.value],INK.sender)
    console.log("approve:",result)
    var result2 = await call("INK","allowance",[INK.senderHex,XC.address])
    console.log("allowance:",result2[0].toString())

    var result3 = await call("XC","lock",[req.body.toPlatform,req.body.toAccount,req.body.value],INK.sender)
    console.log("lock:",result3)

    var result4 = await call("XC","lockBalance")
    console.log("lock:",result4)

    var response = {
        "lockBalance":result4.toString()
    };
    console.log(response);
    res.end(JSON.stringify(response));
})

app.post('/xc/in/', urlencodedParser,async function (req, res) {
    var INK = repoData.contracts[contracts["INK"]]
    var sign = req.body.sign;
    var r = sign.substr(0,64)
    var s = sign.substr(64,64)
    var v = sign.substr(128,2)
    var _v = 27;
    if ( v == '00' || v == '1b') {
        _v = 27
    } else if ( v == '01' || v == '1c') {
        _v = 28
    }
    var result = await call("XCPlugin","voter",[req.body.fromPlatform,req.body.fromAccount,req.body.toAccount,req.body.value,req.body.txid,r,s,_v],INK.sender)
    console.log("voter:",result)

    var result2 = await call("XC","unlock",[req.body.txid,req.body.fromPlatform,req.body.fromAccount,req.body.toAccount,req.body.value],INK.sender)
    console.log("lock:",result2)

    var result3 = await call("XC","lockBalance")
    console.log("lock:",result3)

    var response = {
        "lockBalance":result3.toString()
    };
    console.log(response);
    res.end(JSON.stringify(response));
})

var server = app.listen(8080, function () {
    var host = server.address().address
    var port = server.address().port
    console.log("http://%s:%s", host, port)
})
