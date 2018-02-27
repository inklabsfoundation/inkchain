pragma solidity ^0.4.19;

library Data {

    enum ErrCode {
        Success,
        NotAdmin,
        PlatformTypeInvalid,
        PlatformNameNotNull,
        CatNotOwnerPlatformName,
        NotCredible,
        InsufficientBalance,
        TransferFailed,
        PublicKeyNotExist,
        VoterNotChange,
        WeightNotSatisfied
    }

    struct Admin {
        bytes32 platformName;
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
        address[] publicKeys;
        mapping(bytes32 => Proposal) proposals;
    }
}

interface XCPluginInterface {

    function voter(bytes32 fromPlatform, address fromAccount, address toAccount, uint amount, bytes32 txId, bytes32 r, bytes32 s, bytes v) external returns (Data.ErrCode ErrCode, bool verify);

    function verify(bytes32 fromPlatform, address fromAccount, address toAccount, uint amount, bytes32 txId) external constant returns (Data.ErrCode);

    function deleteProposal(bytes32 platformName, bytes32 txId) external constant returns (Data.ErrCode);

    function setAdmin(bytes32 platformName, address account) external;

    function getAdmin() external constant returns (bytes32, address);

    function addPlatform(uint8 typ, bytes32 name) external returns (Data.ErrCode);

    function deletePlatform(bytes32 name) external constant returns (Data.ErrCode);

    function getPlatform(bytes32 platformName) external returns (Data.ErrCode ErrCode, uint8 typ, bytes32 name, uint weight, address[] publicKeys);

    function existPlatform(bytes32 name) external constant returns (bool);

    function addPublicKey(bytes32 platformName, address publicKey) external constant returns (Data.ErrCode);

    function deletePublicKey(bytes32 platformName, address publicKey) external returns (Data.ErrCode);

    function countOfPublicKey(bytes32 platformName) external constant returns (Data.ErrCode, uint);

    function setWeight(bytes32 name, uint weight) external returns (Data.ErrCode);

    function getWeight(bytes32 name) external constant returns (Data.ErrCode, uint);
}

contract XCPlugin is XCPluginInterface {

    Data.Admin private admin;

    mapping(bytes32 => Data.Platform) private platforms;

    function XCPlugin(bytes32 name) public {
        admin = Data.Admin(name, msg.sender);
    }

    function voter(bytes32 fromPlatform, address fromAccount, address toAccount, uint amount, bytes32 txId, bytes32 r, bytes32 s, byte v) external returns (Data.ErrCode ErrCode, bool verify) {

        if (notExist(fromPlatform)) {
            ErrCode = Data.ErrCode.NotCredible;
            return;
        }

        bytes32 msgHash = signMsg(fromPlatform, fromAccount, admin.platformName, toAccount, amount, txId);

        uint8 v_uint8 = uint8(v);

        if (v_uint8 == 0 || v_uint8 == 1) {
            v_uint8 += 27;
        }

        address publicKey = ecrecover(msgHash, v_uint8, r, s);

        if (!existPublicKey(fromPlatform, publicKey)) {
            ErrCode = Data.ErrCode.PublicKeyNotExist;
            return;
        }

        if (platforms[fromPlatform].proposals[txId].amount == 0) {
            platforms[fromPlatform].proposals[txId].fromAccount = fromAccount;
            platforms[fromPlatform].proposals[txId].toAccount = toAccount;
            platforms[fromPlatform].proposals[txId].amount = amount;
        }

        bool change = changeVoters(fromPlatform, publicKey, txId);

        if (change) {
            ErrCode = Data.ErrCode.Success;
        } else {
            ErrCode = Data.ErrCode.VoterNotChange;
        }

        if (platforms[fromPlatform].proposals[txId].voters.length >= platforms[fromPlatform].weight) {
            verify = true;
        }
        return;
    }

    function verify(bytes32 fromPlatform, address fromAccount, address toAccount, uint amount, bytes32 txId) external constant returns (Data.ErrCode) {

        if (notExist(fromPlatform)) {
            return Data.ErrCode.NotCredible;
        }

        Data.Proposal storage proposal = platforms[fromPlatform].proposals[txId];

        if (proposal.fromAccount != fromAccount || proposal.toAccount != toAccount || proposal.amount != amount) {
            return Data.ErrCode.NotCredible;
        }

        if (proposal.voters.length < platforms[fromPlatform].weight) {
            return Data.ErrCode.WeightNotSatisfied;
        }

        return Data.ErrCode.Success;
    }

    function deleteProposal(bytes32 platformName, bytes32 txId) external constant returns (Data.ErrCode) {

        if (admin.account != msg.sender) {
            return Data.ErrCode.NotAdmin;
        }

        if (platformName == "") {
            return Data.ErrCode.PlatformNameNotNull;
        }

        if (notExist(platformName)) {
            return Data.ErrCode.NotCredible;
        }

        delete platforms[platformName].proposals[txId];

        return Data.ErrCode.Success;
    }

    function setAdmin(bytes32 platformName, address account) external {
        if (admin.account == msg.sender) {
            admin.platformName = platformName;
            admin.account = account;
        }
    }

    function getAdmin() external constant returns (bytes32, address) {
        return (admin.platformName, admin.account);
    }

    function addPlatform(uint8 platformType, bytes32 platformName) external returns (Data.ErrCode) {

        if (admin.account != msg.sender) {
            return Data.ErrCode.NotAdmin;
        }

        if (platformType != 1 && platformType != 2) {
            return Data.ErrCode.PlatformTypeInvalid;
        }

        if (platformName == "") {
            return Data.ErrCode.PlatformNameNotNull;
        }

        if (admin.platformName == platformName) {
            return Data.ErrCode.CatNotOwnerPlatformName;
        }

        if (notExist(platformName)) {
            platforms[platformName].typ = platformType;
            platforms[platformName].name = platformName;
            platforms[platformName].weight = 1;
        }

        return Data.ErrCode.Success;
    }

    function deletePlatform(bytes32 name) external constant returns (Data.ErrCode){

        if (admin.account != msg.sender) {
            return Data.ErrCode.NotAdmin;
        }

        if (name == "") {
            return Data.ErrCode.PlatformNameNotNull;
        }

        if (notExist(name)) {
            return Data.ErrCode.NotCredible;
        }

        delete platforms[name];

        return Data.ErrCode.Success;
    }

    function getPlatform(bytes32 platformName) external returns (Data.ErrCode ErrCode, uint8 typ, bytes32 name, uint weight, address[] publicKeys) {

        if (admin.account != msg.sender) {
            ErrCode = Data.ErrCode.NotAdmin;
            return;
        }

        if (platformName == "") {
            ErrCode = Data.ErrCode.PlatformNameNotNull;
            return;
        }

        if (notExist(platformName)) {
            ErrCode = Data.ErrCode.NotCredible;
            return;
        }

        ErrCode = Data.ErrCode.Success;
        typ = platforms[platformName].typ;
        name = platforms[platformName].name;
        weight = platforms[platformName].weight;
        publicKeys = platforms[platformName].publicKeys;

        return;
    }

    function existPlatform(bytes32 name) external constant returns (bool){
        return (platforms[name].typ != 0);
    }

    function setWeight(bytes32 name, uint weight) external returns (Data.ErrCode) {

        if (admin.account != msg.sender) {
            return Data.ErrCode.NotAdmin;
        }

        if (notExist(name)) {
            return Data.ErrCode.NotCredible;
        }

        platforms[name].weight = weight;

        return Data.ErrCode.Success;
    }

    function getWeight(bytes32 name) external constant returns (Data.ErrCode, uint) {

        if (admin.account != msg.sender) {
            return (Data.ErrCode.NotAdmin, 0);
        }

        if (notExist(name)) {
            return (Data.ErrCode.NotCredible, 0);
        }

        return (Data.ErrCode.Success, platforms[name].weight);
    }

    function addPublicKey(bytes32 platformName, address publicKey) external constant returns (Data.ErrCode) {

        if (admin.account != msg.sender) {
            return Data.ErrCode.NotAdmin;
        }

        if (notExist(platformName)) {
            return Data.ErrCode.NotCredible;
        }

        address[] storage listOfPublicKey = platforms[platformName].publicKeys;

        for (uint i; i < listOfPublicKey.length; i++) {
            if (publicKey == listOfPublicKey[i]) {
                return Data.ErrCode.Success;
            }
        }

        listOfPublicKey.push(publicKey);

        return Data.ErrCode.Success;
    }

    function deletePublicKey(bytes32 platformName, address publickey) external returns (Data.ErrCode) {

        if (admin.account != msg.sender) {
            return Data.ErrCode.NotAdmin;
        }

        if (notExist(platformName)) {
            return Data.ErrCode.NotCredible;
        }

        address[] storage listOfPublicKey = platforms[platformName].publicKeys;

        bool exist;

        for (uint i = 0; i <= listOfPublicKey.length; i++) {
            if (exist) {
                if (i == listOfPublicKey.length) {
                    delete listOfPublicKey[i - 1];
                    listOfPublicKey.length--;
                } else {
                    listOfPublicKey[i - 1] = listOfPublicKey[i];
                }
            } else if (listOfPublicKey[i] == publickey) {
                exist = true;
            }
        }

        return Data.ErrCode.Success;
    }

    function countOfPublicKey(bytes32 platformName) external constant returns (Data.ErrCode, uint){

        if (admin.account != msg.sender) {
            return (Data.ErrCode.NotAdmin, 0);
        }

        if (notExist(platformName)) {
            return (Data.ErrCode.NotCredible, 0);
        }

        return (Data.ErrCode.Success, platforms[platformName].publicKeys.length);
    }

    /**
     *   ######################
     *  #  private function  #
     * ######################
     */
    function existPublicKey(bytes32 platformName, address publicKey) internal constant returns (bool exist) {

        address[] memory listOfPublicKey = platforms[platformName].publicKeys;

        for (uint i = 0; i < listOfPublicKey.length; i++) {
            if (exist) {
                break;
            }
            if (listOfPublicKey[i] == publicKey) {
                exist = true;
            }
        }
    }

    function signMsg(bytes32 fromPlatform, address fromAccount, bytes32 toPlatform, address toAccount, uint amount, bytes32 txId) internal returns (bytes32) {
        return keccak256(bytes32ToStr(fromPlatform), ":0x", uintToStr(uint160(fromAccount), 16), ":", bytes32ToStr(toPlatform), ":0x", uintToStr(uint160(toAccount), 16), ":", uintToStr(amount, 10), ":", bytes32ToStr(txId));
    }

    function notExist(bytes32 name) internal constant returns (bool){
        return (platforms[name].typ == 0);
    }

    function changeVoters(bytes32 platformName, address publicKey, bytes32 txId) internal constant returns (bool change) {

        address[] storage voters = platforms[platformName].proposals[txId].voters;

        change = true;

        for (uint i = 0; i < voters.length; i++) {
            if (voters[i] == publicKey) {
                change = false;
            }
        }

        if (change) {
            voters.push(publicKey);
        }
    }

    function bytes32ToBytes(bytes32 b) internal pure returns (bytes) {

        var length = b.length;

        for (uint i = 0; i < b.length; i++) {
            if (b[b.length - 1 - i] == "") {
                length -= 1;
            } else {
                break;
            }
        }

        bytes memory bs = new bytes(length);

        for (uint j = 0; j < length; j++) {
            bs[j] = b[j];
        }

        return bs;
    }

    function bytes32ToStr(bytes32 b) internal constant returns (string) {
        bytes memory bs = bytes32ToBytes(b);
        return string(bs);
    }

    function uintToStr(uint value, uint base) internal pure returns (string) {

        uint _value = value;

        uint length = 0;

        bytes16 tenStr = "0123456789abcdef";

        while (true) {
            if (_value > 0) {
                length ++;
                _value = _value / base;
            } else {
                break;
            }
        }

        bytes memory bs = new bytes(length);

        for (uint i = 0; i < length; i++) {
            bs[length - 1 - i] = tenStr[value % base];
            value = value / base;
        }

        return string(bs);
    }

}
