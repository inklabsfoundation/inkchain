pragma solidity ^0.4.19;

library Data {

    enum ErrCode {
        Success,
        NotAdmin,
        StatusClosed,
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
        bool status;
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

    function verifyProposal(bytes32 fromPlatform, address fromAccount, address toAccount, uint amount, bytes32 txId) external constant returns (Data.ErrCode);

    function deleteProposal(bytes32 platformName, bytes32 txId) external constant returns (Data.ErrCode);
}


interface XCInterface {

    function start() external;

    function stop() external;

    function setAdmin(bytes32 platformName, address account) external;

    function getAdmin() external constant returns (bool, bytes32, address);

    function setINK(address account) external;

    function getINK() external constant returns (address);

    function setXCPlugin(address account) external;

    function getXCPlugin() external constant returns (address);

    function lock(bytes32 toPlatform, address toAccount, uint amount) external payable returns (Data.ErrCode);

    function unlock(bytes32 txId, bytes32 fromPlatform, address fromAccount, address toAccount, uint amount) external payable returns (Data.ErrCode);

    function withdrawal(address account, uint amount) external payable returns (Data.ErrCode);

    function lockAdmin(bytes32 toPlatform, address toAccount, uint amount) external payable returns (Data.ErrCode);

    function unlockAdmin(bytes32 txId, bytes32 fromPlatform, address fromAccount, address toAccount, uint amount) external payable returns (Data.ErrCode);
}

contract XC is XCInterface {

    Data.Admin private admin;

    mapping(bytes32 => uint) public balanceOf;

    INK private inkToken;

    XCPlugin private xcPlugin;

    event lockEvent(bytes32 toPlatform, address toAccount, string amount);

    event unlockEvent(bytes32 txId, bytes32 fromPlatform, address fromAccount, string amount);

    function XC(bytes32 name) public payable {
        admin = Data.Admin(false, name, msg.sender);
    }

    function start() external {
        if (admin.account == msg.sender) {
            admin.status = true;
        }
    }

    function stop() external {
        if (admin.account == msg.sender) {
            admin.status = false;
        }
    }

    function setAdmin(bytes32 platformName, address account) external {
        if (admin.account == msg.sender) {
            admin.platformName = platformName;
            admin.account = account;
        }
    }

    function getAdmin() external constant returns (bool, bytes32, address) {
        return (admin.status, admin.platformName, admin.account);
    }

    function setINK(address account) external {
        inkToken = INK(account);
    }

    function getINK() external constant returns (address) {
        return inkToken;
    }

    function setXCPlugin(address account) external {
        xcPlugin = XCPlugin(account);
    }

    function getXCPlugin() external constant returns (address) {
        return xcPlugin;
    }

    function lock(bytes32 toPlatform, address toAccount, uint amount) external payable returns (Data.ErrCode) {

        if (!admin.status) {
            return Data.ErrCode.StatusClosed;
        }

        if (!xcPlugin.existPlatform(toPlatform)) {
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

        balanceOf[toPlatform] += amount;

        lockEvent(toPlatform, toAccount, uintAppendToString(amount));

        return Data.ErrCode.Success;
    }

    function unlock(bytes32 txId, bytes32 fromPlatform, address fromAccount, address toAccount, uint amount) external payable returns (Data.ErrCode) {

        if (!admin.status) {
            return Data.ErrCode.StatusClosed;
        }

        if (!xcPlugin.existPlatform(fromPlatform)) {
            return Data.ErrCode.NotCredible;
        }

        Data.ErrCode ErrCode = xcPlugin.verifyProposal(fromPlatform, fromAccount, toAccount, amount, txId);

        if (ErrCode == Data.ErrCode.Success) {
            return ErrCode;
        }

        uint balanceOfContract = inkToken.balanceOf(this);

        if (balanceOfContract < amount) {
            return Data.ErrCode.InsufficientBalance;
        }

        bool success = inkToken.transfer(toAccount, amount);

        if (!success) {
            return Data.ErrCode.TransferFailed;
        }

        ErrCode = xcPlugin.deleteProposal(fromPlatform, txId);

        if (ErrCode != Data.ErrCode.Success) {
            return ErrCode;
        }

        balanceOf[admin.platformName] -= amount;

        balanceOf[fromPlatform] -= amount;

        unlockEvent(txId, fromPlatform, fromAccount, uintAppendToString(amount));

        return Data.ErrCode.Success;
    }

    function withdrawal(address account, uint amount) external payable returns (Data.ErrCode) {

        if (admin.account != msg.sender) {
            return Data.ErrCode.NotAdmin;
        }

        uint balanceOfContract = inkToken.balanceOf(this);

        uint balance = balanceOf[admin.platformName];

        if (balanceOfContract - balance < amount) {
            return Data.ErrCode.InsufficientBalance;
        }

        bool success = inkToken.transfer(account, amount);

        if (!success) {
            return Data.ErrCode.TransferFailed;
        }

        return Data.ErrCode.Success;
    }

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

    function unlockAdmin(bytes32 txId, bytes32 fromPlatform, address fromAccount, address toAccount, uint amount) external payable returns (Data.ErrCode) {

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
            unlockEvent(txId, fromPlatform, fromAccount, uintAppendToString(amount));
        }

        return Data.ErrCode.Success;
    }


    /**
     *   ######################
     *  #  private function  #
     * ######################
     */

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