pragma solidity ^0.4.19;

library Data {

    enum Errcode {
        Success,
        NotOwner,
        PlatformTypeInvalid,
        PlatformNameNotNull,
        CatNotOwenerPlatformName,
        NotCredible,
        InsufficientBalance,
        TransferFailed,
        PublickeyNotExist,
        VoterNotChange,
        WeightNotSatisfied
    }

    struct Admin {
        bytes32 name;
        address account;
    }

    struct Proposal {
        address fromAccount;
        address toAccount;
        uint amount;
        address[] voters;
    }

    struct Platform {
        uint8 typ;
        bytes32 name;
        uint weight;
        address[] publickeys;
        mapping(bytes32 => Proposal) proposals;
    }
}

interface XCPluginInterface {

    function voter(bytes32 fromPlatform,address fromAccount, address toAccount,uint amount, bytes32 txid, bytes sig)  external returns (Data.Errcode errcode,bool verify);
    function verify(bytes32 fromPlatform,address fromAccount, address toAccount, uint amount, bytes32 txid) external constant returns (Data.Errcode);
    function deleteProposal(bytes32 platformName, bytes32 txid) external constant returns (Data.Errcode);

    function setAdmin(bytes32 name,address account) external;
    function getAdmin() external constant returns (bytes32,address);

    function addPlatform(uint8 typ, bytes32 name) external returns (Data.Errcode);
    function deletePlatfrom(bytes32 name) external constant returns (Data.Errcode);
    function getPlatfrom(bytes32 platfromName) external returns (Data.Errcode errcode, uint8 typ, bytes32 name, uint weight, address[] publickeys);
    function existPlatfrom(bytes32 name) external constant returns (bool);

    function addPublickey(bytes32 platfromName, address publickey) external constant returns (Data.Errcode);
    function deletePublickey(bytes32 platfromName, address publickey) external returns (Data.Errcode);
    function countOfPublickey(bytes32 platfromName) external constant returns (Data.Errcode, uint);

    function setWeight(bytes32 name, uint weight) external returns (Data.Errcode);
    function getWeight(bytes32 name) external constant returns (Data.Errcode, uint);
}

contract XCPlugin is XCPluginInterface {

    Data.Admin private admin;
    mapping(bytes32 => Data.Platform) private platforms;

    function XCPlugin(bytes32 name) public {
        admin = Data.Admin(name, msg.sender);
    }

    function voter(bytes32 fromPlatform,address fromAccount, address toAccount,uint amount, bytes32 txid, bytes sig) external returns (Data.Errcode errcode,bool verify) {
        if (notExist(fromPlatform)) {
            errcode = Data.Errcode.NotCredible;
            return;
        }
        bytes32 msgHash = signMsg(fromPlatform,fromAccount,admin.name,toAccount,amount,txid);
        address publickey = ecrecoverDecode(msgHash, sig);
        if (!existPublickey(fromPlatform, publickey)) {
            errcode = Data.Errcode.PublickeyNotExist;
            return;
        }
        if (platforms[fromPlatform].proposals[txid].amount == 0 ) {
            platforms[fromPlatform].proposals[txid].fromAccount = fromAccount;
            platforms[fromPlatform].proposals[txid].toAccount = toAccount;
            platforms[fromPlatform].proposals[txid].amount = amount;
        }
        bool change = changeVoters(fromPlatform, publickey, txid);
        if (change) {
            errcode = Data.Errcode.Success;
        }else{
            errcode = Data.Errcode.VoterNotChange;
        }
        if (platforms[fromPlatform].proposals[txid].voters.length >= platforms[fromPlatform].weight) {
            verify = true;
        }
        return ;
    }

    function verify(bytes32 fromPlatform,address fromAccount, address toAccount, uint amount, bytes32 txid) external constant returns (Data.Errcode) {
        if (notExist(fromPlatform)) {
            return Data.Errcode.NotCredible;
        }
        var proposal = platforms[fromPlatform].proposals[txid];
        if (proposal.fromAccount != fromAccount ||  proposal.toAccount != toAccount || proposal.amount != amount) {
            return Data.Errcode.NotCredible;
        }
        if (platforms[fromPlatform].proposals[txid].voters.length < platforms[fromPlatform].weight) {
            return Data.Errcode.WeightNotSatisfied;
        }
        return Data.Errcode.Success;
    }

    function deleteProposal(bytes32 platformName, bytes32 txid) external constant returns (Data.Errcode) {
        if (admin.account != msg.sender) {
            return Data.Errcode.NotOwner;
        }
        if (platformName == "") {
            return Data.Errcode.PlatformNameNotNull;
        }
        if (notExist(platformName)) {
            return Data.Errcode.NotCredible;
        }
        delete platforms[platformName].proposals[txid];
        return Data.Errcode.Success;
    }

    function setAdmin(bytes32 name,address account) external {
        if (admin.account == msg.sender) {
            admin.name = name;
            admin.account = account;
        }
    }

    function getAdmin() external constant returns (bytes32, address) {
        return (admin.name,admin.account);
    }

    function addPlatform(uint8 platfromType, bytes32 platfromName) external returns (Data.Errcode) {
        if (admin.account != msg.sender) {
            return Data.Errcode.NotOwner;
        }
        if (platfromType != 1 && platfromType != 2) {
            return Data.Errcode.PlatformTypeInvalid;
        }
        if (platfromName == "") {
            return Data.Errcode.PlatformNameNotNull;
        }
        if (admin.name == platfromName) {
            return Data.Errcode.CatNotOwenerPlatformName;
        }
        if (notExist(platfromName)) {
            platforms[platfromName].typ = platfromType;
            platforms[platfromName].name = platfromName;
            platforms[platfromName].weight = 1;
        }
        return Data.Errcode.Success;
    }

    function deletePlatfrom(bytes32 name) external constant returns (Data.Errcode){
        if (admin.account != msg.sender) {
            return Data.Errcode.NotOwner;
        }
        if (name == "") {
            return Data.Errcode.PlatformNameNotNull;
        }
        if (notExist(name)) {
            return Data.Errcode.NotCredible;
        }
        delete platforms[name];
        return Data.Errcode.Success;
    }

    function getPlatfrom(bytes32 platfromName) external returns (Data.Errcode errcode, uint8 typ, bytes32 name, uint weight, address[] publickeys) {
        if (admin.account != msg.sender) {
            errcode = Data.Errcode.NotOwner;
            return;
        }
        if (platfromName == "") {
            errcode = Data.Errcode.PlatformNameNotNull;
            return;
        }
        if (notExist(platfromName)) {
            errcode = Data.Errcode.NotCredible;
            return;
        }
        errcode = Data.Errcode.Success;
        typ = platforms[platfromName].typ;
        name = platforms[platfromName].name;
        weight = platforms[platfromName].weight;
        publickeys = platforms[platfromName].publickeys;
        return;
    }

    function existPlatfrom(bytes32 name) external constant returns (bool){
        return (platforms[name].typ != 0);
    }

    function setWeight(bytes32 name, uint weight) external returns (Data.Errcode) {
        if (admin.account != msg.sender) {
            return Data.Errcode.NotOwner;
        }
        if (notExist(name)) {
            return Data.Errcode.NotCredible;
        }
        platforms[name].weight = weight;
        return Data.Errcode.Success;
    }

    function getWeight(bytes32 name) external constant returns (Data.Errcode, uint) {
        if (admin.account != msg.sender) {
            return (Data.Errcode.NotOwner, 0);
        }
        if (notExist(name)) {
            return (Data.Errcode.NotCredible, 0);
        }
        return (Data.Errcode.Success, platforms[name].weight);
    }

    function addPublickey(bytes32 platfromName, address publickey) external constant returns (Data.Errcode) {
        if (admin.account != msg.sender) {
            return Data.Errcode.NotOwner;
        }
        if (notExist(platfromName)) {
            return Data.Errcode.NotCredible;
        }
        address[] storage listOfPublickey = platforms[platfromName].publickeys;
        for (uint i; i < listOfPublickey.length; i++) {
            if (publickey == listOfPublickey[i]) {
                return Data.Errcode.Success;
            }
        }
        listOfPublickey.push(publickey);
        return Data.Errcode.Success;
    }

    function deletePublickey(bytes32 platfromName, address publickey) external returns (Data.Errcode) {
        if (admin.account != msg.sender) {
            return Data.Errcode.NotOwner;
        }
        if (notExist(platfromName)) {
            return Data.Errcode.NotCredible;
        }
        address[] storage listOfPublickey = platforms[platfromName].publickeys;
        bool exist;
        for (uint i = 0; i <= listOfPublickey.length; i++) {
            if (exist) {
                if (i == listOfPublickey.length) {
                    delete listOfPublickey[i - 1];
                    listOfPublickey.length--;
                } else {
                    listOfPublickey[i - 1] = listOfPublickey[i];
                }
            } else if (listOfPublickey[i] == publickey) {
                exist = true;
            }
        }
        return Data.Errcode.Success;
    }

    function countOfPublickey(bytes32 platfromName) external constant returns (Data.Errcode, uint){
        if (admin.account != msg.sender) {
            return (Data.Errcode.NotOwner, 0);
        }
        if (notExist(platfromName)) {
            return (Data.Errcode.NotCredible, 0);
        }
        return (Data.Errcode.Success, platforms[platfromName].publickeys.length);
    }

    /**
     *   ######################
     *  #  private function  #
     * ######################
     */
    function existPublickey(bytes32 platformName, address publickey) internal constant returns (bool exist) {
        address[] memory listOfPublickey = platforms[platformName].publickeys;
        for (uint i = 0; i < listOfPublickey.length; i++) {
            if (exist) {
                break;
            }
            if (listOfPublickey[i] == publickey) {
                exist = true;
            }
        }
    }

    function signMsg(bytes32 fromPlatform,address fromAccount, bytes32 toPlatform,address toAccount,uint amount, bytes32 txid) internal returns (bytes32) {
        return keccak256(bytes32ToStr(fromPlatform),":0x",uintToStr(uint160(fromAccount),16),":",bytes32ToStr(toPlatform),":0x",uintToStr(uint160(toAccount),16),":",uintToStr(amount,10),":",bytes32ToStr(txid));
    }

    function notExist(bytes32 name) internal constant returns (bool){
        return (platforms[name].typ == 0);
    }

    function changeVoters(bytes32 platformName, address publickey, bytes32 txid) internal constant returns (bool change) {
        address[] storage voters = platforms[platformName].proposals[txid].voters;
        change = true;
        for (uint i = 0; i < voters.length; i++) {
            if (voters[i] == publickey) {
                change = false;
            }
        }
        if (change) {
            voters.push(publickey);
        }
    }

    function signatureSplit(bytes memory signature) internal returns (bytes32 r, bytes32 s, uint8 v) {
        assembly {
            r := mload(add(signature, 32))
            s := mload(add(signature, 64))
            v := and(mload(add(signature, 65)), 0xff)
        }
        if (v == 0 || v == 1) {
            v += 27;
        }
        require(v == 27 || v == 28);
    }

    function ecrecoverDecode(bytes32 msg, bytes sign) internal returns (address publickey){
        var (r, s, v) = signatureSplit(sign);
        publickey = ecrecover(msg, v, r, s);
    }


    function bytes32ToBytes(bytes32 b) internal pure returns  (bytes) {
        var length = b.length;
        for(uint i=0;i<b.length;i++){
            if (b[b.length-1-i]=="") {
                length -= 1;
            }else{
                break;
            }
        }
        bytes memory names = new bytes(length);
        for(uint j=0;j<length;j++){
            names[j] = b[j];
        }
        return names;
    }

    function bytes32ToStr(bytes32 b) internal constant returns  (string) {
        bytes memory names = bytes32ToBytes(b);
        return string(names);
    }

    function uintToStr(uint value,uint base) internal pure returns (string) {
        uint _value = value;
        uint length = 0;
        bytes16 tenStr = "0123456789abcdef";
        while (true) {
            if (_value>0) {
                length ++;
                _value = _value/base;
            }else{
                break;
            }
        }
        bytes memory names = new bytes(length);
        for(uint i=0;i<length;i++) {
            names[length-1-i] = tenStr[value%base];
            value = value/base;
        }
        return string(names);
    }

}
