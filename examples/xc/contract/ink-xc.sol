pragma solidity ^0.4.19;

//data structs
library Data {

    //error code
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

    //admin info
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

    //platform info
    struct Platform {
        uint8 typ;
        bytes32 name;
        uint weight;
        address[] publicKeys;
        mapping(bytes32 => Proposal) proposals;
    }
}

contract INK {

    mapping(address => uint256) public balanceOf;

    mapping(address => mapping(address => uint256)) public allowance;

    function transferFrom(address _from, address _to, uint256 value) public returns (bool success);

    function transfer(address _to, uint256 value) public returns (bool success);
}

contract XCPlugin {

    function existPlatform(bytes32 name) external constant returns (bool);

    function verify(bytes32 fromPlatform, address fromAccount, address toAccount, uint amount, bytes32 txId) external constant returns (Data.ErrCode);

    function deleteProposal(bytes32 platformName, bytes32 txId) external constant returns (Data.ErrCode);
}


interface XCInterface {

    function setAdmin(bytes32 platformName, address account) external;

    function getAdmin() external constant returns (bytes32, address);

    function setINK(address account) external;

    function getINK() external constant returns (address);

    function setXCPlugin(address account) external;

    function getXCPlugin() external constant returns (address);

    function lock(bytes32 toPlatform, address toAccount, uint amount) external payable returns (Data.ErrCode);

    function unlock(bytes32 fromPlatform, address fromAccount, address toAccount, uint amount, bytes32 txId) external payable returns (Data.ErrCode);

    function withdrawal(address account, uint amount) external payable returns (Data.ErrCode);

    function lockAdmin(bytes32 toPlatform, address toAccount, uint amount) external payable returns (Data.ErrCode);

    function unlockAdmin(bytes32 fromPlatform, address fromAccount, address toAccount, uint amount, bytes32 txId) external payable returns (Data.ErrCode);
}

contract XC is XCInterface {

    Data.Admin private admin;

    mapping(bytes32 => uint) public balanceOf;

    INK private inkToken;

    XCPlugin private xcPlugin;

    event lockEvent(bytes32 toPlatform, address toAccount, string amount);

    event unlockEvent(bytes32 txId, bytes32 fromPlatform, address fromAccount, string amount);

    //init contract
    function XC(bytes32 name) public payable {
        admin = Data.Admin(name, msg.sender);
    }

    //reset admin info
    function setAdmin(bytes32 platformName, address account) external {
        if (admin.account == msg.sender) {
            admin.platformName = platformName;
            admin.account = account;
        }
    }

    //get admin info  platformName and account
    function getAdmin() external constant returns (bytes32, address) {
        return (admin.platformName, admin.account);
    }

    //instantiation inkTOKEN
    function setINK(address account) external {
        inkToken = INK(account);
    }

    //get variable inkToken
    function getINK() external constant returns (address) {
        return inkToken;
    }

    //instantiation variable xcPlugin
    function setXCPlugin(address account) external {
        xcPlugin = XCPlugin(account);
    }

    //get private variable xcPlugin
    function getXCPlugin() external constant returns (address) {
        return xcPlugin;
    }

    //turn out
    function lock(bytes32 toPlatform, address toAccount, uint amount) external payable returns (Data.ErrCode) {
        //determine whether toPlatform exist in xcPlugin's existPlatform
        if (!xcPlugin.existPlatform(toPlatform)) {
            return Data.ErrCode.NotCredible;
        }

        //get user approve the contract quota
        uint allowance = inkToken.allowance(msg.sender, this);

        //judge whether the amount authorized by user to the contract is less than amount
        if (allowance < amount) {
            return Data.ErrCode.InsufficientBalance;
        }

        //do transferFrom
        bool success = inkToken.transferFrom(msg.sender, this, amount);
        if (!success) {
            return Data.ErrCode.TransferFailed;
        }

        //record the amount of local platform turn out
        balanceOf[admin.platformName] += amount;
        //record the amount of local platform turn to toPlatform
        balanceOf[toPlatform] += amount;
        //trigger lockEvent
        lockEvent(toPlatform, toAccount, uintAppendToString(amount));

        return Data.ErrCode.Success;
    }

    //turn in
    function unlock(bytes32 fromPlatform, address fromAccount, address toAccount, uint amount, bytes32 txId) external payable returns (Data.ErrCode) {
        //determine whether fromPlatform exist in xcPlugin's existPlatform
        if (!xcPlugin.existPlatform(fromPlatform)) {
            return Data.ErrCode.NotCredible;
        }

        //verify args by function xcPlugin.verify
        Data.ErrCode ErrCode = xcPlugin.verify(fromPlatform, fromAccount, toAccount, amount, txId);

        if (ErrCode == Data.ErrCode.Success) {
            return ErrCode;
        }
        //get contracts balance
        uint balanceOfContract = inkToken.balanceOf(this);
        //validate the balance of contract were less than amount
        if (balanceOfContract < amount) {
            return Data.ErrCode.InsufficientBalance;
        }

        bool success = inkToken.transfer(toAccount, amount);

        if (!success) {
            return Data.ErrCode.TransferFailed;
        }

        ErrCode = xcPlugin.deleteProposal(fromPlatform, txId);

        if (ErrCode == Data.ErrCode.Success) {
            return ErrCode;
        }

        balanceOf[admin.platformName] -= amount;

        balanceOf[fromPlatform] -= amount;
        //trigger unlockEvent
        unlockEvent(txId, fromPlatform, fromAccount, uintAppendToString(amount));

        return Data.ErrCode.Success;
    }

    //refund
    function withdrawal(address account, uint amount) external payable returns (Data.ErrCode) {

        if (admin.account != msg.sender) {
            return Data.ErrCode.NotAdmin;
        }
        //get balance of contract
        uint balanceOfContract = inkToken.balanceOf(this);

        uint balance = balanceOf[admin.platformName];
        //validate availability of non-cross-chain balance were less than amount
        if (balanceOfContract - balance < amount) {
            return Data.ErrCode.InsufficientBalance;
        }

        bool success = inkToken.transfer(account, amount);

        if (!success) {
            return Data.ErrCode.TransferFailed;
        }

        return Data.ErrCode.Success;
    }

    //build union chain consume record before cross-chain publish
    function lockAdmin(bytes32 toPlatform, address toAccount, uint amount) external payable returns (Data.ErrCode) {

        if (admin.account != msg.sender) {
            return Data.ErrCode.NotAdmin;
        }

        if (!xcPlugin.existPlatform(toPlatform) && admin.platformName != toPlatform) {
            return Data.ErrCode.NotCredible;
        }

        uint allowance = inkToken.allowance(msg.sender, this);

        if (allowance < amount) {
            return Data.ErrCode.InsufficientBalance;
        }

        bool success = inkToken.transferFrom(msg.sender, this, amount);

        if (!success) {
            return Data.ErrCode.TransferFailed;
        }

        balanceOf[admin.platformName] += amount;

        if (admin.platformName != toPlatform && xcPlugin.existPlatform(toPlatform)) {
            balanceOf[toPlatform] += amount;
            lockEvent(toPlatform, toAccount, uintAppendToString(amount));
        }

        return Data.ErrCode.Success;
    }

    function unlockAdmin(bytes32 fromPlatform, address fromAccount, address toAccount, uint amount, bytes32 txid) external payable returns (Data.ErrCode) {

        if (admin.account != msg.sender) {
            return Data.ErrCode.NotAdmin;
        }

        if (!xcPlugin.existPlatform(fromPlatform) && fromPlatform != admin.platformName) {
            return Data.ErrCode.NotCredible;
        }

        uint balanceOfContract = inkToken.balanceOf(this);

        if (balanceOfContract < amount) {
            return Data.ErrCode.InsufficientBalance;
        }

        bool success = inkToken.transfer(toAccount, amount);

        if (!success) {
            return Data.ErrCode.TransferFailed;
        }

        balanceOf[admin.platformName] -= amount;

        if (fromPlatform != admin.platformName) {
            balanceOf[fromPlatform] -= amount;
            unlockEvent(txid, fromPlatform, fromAccount, uintAppendToString(amount));
        }

        return Data.ErrCode.Success;
    }

    //private function change uint to string
    function uintAppendToString(uint v) pure internal returns (string){

        uint length = 100;

        bytes memory reversed = new bytes(length);

        bytes16 sixTeenStr = "0123456789abcdef";

        uint i = 0;

        while (v != 0) {
            uint remainder = v % 16;
            v = v / 16;
            reversed[i++] = byte(sixTeenStr[remainder]);
        }

        string memory bytesList = "0000000000000000000000000000000000000000000000000000000000000000";

        bytes memory str = bytes(bytesList);

        for (uint j = 0; j < i; j++) {
            str[str.length - j - 1] = reversed[i - j - 1];
        }

        return string(str);
    }
}