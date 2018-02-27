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

contract INK {

    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;

    function transferFrom(address _from, address _to, uint256 value) public returns (bool success);
    function transfer(address _to, uint256 value) public returns (bool success);
}

contract XCPlugin {

    function existPlatfrom(bytes32 name) external constant returns (bool);
    function verify(bytes32 fromPlatform,address fromAccount, address toAccount, uint amount, bytes32 txid) external constant returns (Data.Errcode);
    function deleteProposal(bytes32 platformName, bytes32 txid) external constant returns (Data.Errcode);
}


interface XCInterface {

    function setAdmin(bytes32 name,address account) external;
    function getAdmin() external constant returns (bytes32,address);

    function setINK(address account) external;
    function getINK() external constant returns (address);

    function setXCPlugin(address account) external;
    function getXCPlugin() external constant returns (address);

    function lock(bytes32 toPlatform, address toAccount, uint amount) external payable returns (Data.Errcode);
    function unlock(bytes32 fromPlatform,address fromAccount, address toAccount, uint amount, bytes32 txid) external payable returns (Data.Errcode);

    function withdrawal(address account,uint amount) external payable returns (Data.Errcode);

    function lockAdmin(bytes32 toPlatform, address toAccount, uint amount) external payable returns (Data.Errcode);
    function unlockAdmin(bytes32 fromPlatform,address fromAccount, address toAccount, uint amount, bytes32 txid) external payable returns (Data.Errcode);
}

contract XC is XCInterface {

    Data.Admin private admin;
    mapping(bytes32 => uint) public balanceOf;

    INK private inkToken;
    XCPlugin private xcPlugin;

    event lockEvent(bytes32 toPlatform, address toAccount, string amount);
    event unlockEvent(bytes32 txid,bytes32 fromPlatform,address fromAccount ,string amount);

    function XC(bytes32 name) public payable {
        admin = Data.Admin(name, msg.sender);
    }

    function setAdmin(bytes32 name, address account) external {
        if (admin.account == msg.sender) {
            admin.name = name;
            admin.account = account;
        }
    }

    function getAdmin() external constant returns (bytes32,address) {
        return (admin.name,admin.account);
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

    function lock(bytes32 toPlatform, address toAccount, uint amount) external payable returns (Data.Errcode) {

        if (!xcPlugin.existPlatfrom(toPlatform)) {
            return Data.Errcode.NotCredible;
        }

        uint allowance = inkToken.allowance(msg.sender, this);
        if (allowance < amount) {
            return Data.Errcode.InsufficientBalance;
        }

        bool success = inkToken.transferFrom(msg.sender, this, amount);
        if (!success) {
            return Data.Errcode.TransferFailed;
        }

        balanceOf[admin.name] += amount;

        balanceOf[toPlatform] += amount;

        string memory value = uintAppendToString(amount);

        lockEvent(toPlatform, toAccount, value);
        return Data.Errcode.Success;
    }

    function unlock(bytes32 fromPlatform,address fromAccount, address toAccount, uint amount, bytes32 txid) external payable returns (Data.Errcode) {

        if (!xcPlugin.existPlatfrom(fromPlatform)) {
            return Data.Errcode.NotCredible;
        }

        Data.Errcode errcode = xcPlugin.verify(fromPlatform,fromAccount,toAccount,amount, txid);
        if (errcode == Data.Errcode.Success) {
            return errcode;
        }

        uint balanceOfContract = inkToken.balanceOf(this);
        if (balanceOfContract < amount) {
            return Data.Errcode.InsufficientBalance;
        }

        bool success = inkToken.transfer(toAccount, amount);
        if (!success) {
            return Data.Errcode.TransferFailed;
        }

        errcode = xcPlugin.deleteProposal(fromPlatform, txid);
        if (errcode == Data.Errcode.Success) {
            return errcode;
        }

        balanceOf[admin.name] -= amount;

        balanceOf[fromPlatform] -= amount;

        string memory value = uintAppendToString(amount);

        unlockEvent(txid, fromPlatform, fromAccount, value);
        return Data.Errcode.Success;
    }

    function withdrawal(address account,uint amount) external payable returns (Data.Errcode) {
        if (admin.account != msg.sender) {
            return Data.Errcode.NotOwner;
        }

        uint balanceOfContract = inkToken.balanceOf(this);
        uint balance = balanceOf[admin.name];
        if (balanceOfContract - balance < amount) {
            return Data.Errcode.InsufficientBalance;
        }

        bool success = inkToken.transfer(account, amount);
        if (!success) {
            return Data.Errcode.TransferFailed;
        }
        return Data.Errcode.Success;
    }

    function lockAdmin(bytes32 toPlatform, address toAccount, uint amount) external payable returns (Data.Errcode) {

        if (admin.account != msg.sender) {
            return Data.Errcode.NotOwner;
        }

        if (!xcPlugin.existPlatfrom(toPlatform) && admin.name != toPlatform) {
            return Data.Errcode.NotCredible;
        }

        uint allowance = inkToken.allowance(msg.sender, this);
        if (allowance < amount) {
            return Data.Errcode.InsufficientBalance;
        }

        bool success = inkToken.transferFrom(msg.sender, this, amount);
        if (!success) {
            return Data.Errcode.TransferFailed;
        }

        balanceOf[admin.name] += amount;

        if (admin.name != toPlatform && xcPlugin.existPlatfrom(toPlatform)) {
            balanceOf[toPlatform] += amount;
            string memory value = uintAppendToString(amount);
            lockEvent(toPlatform, toAccount, value);
        }

        return Data.Errcode.Success;
    }

    function unlockAdmin(bytes32 fromPlatform,address fromAccount, address toAccount, uint amount, bytes32 txid) external payable returns (Data.Errcode) {

        if (admin.account != msg.sender) {
            return Data.Errcode.NotOwner;
        }

        if (!xcPlugin.existPlatfrom(fromPlatform) && fromPlatform !=admin.name) {
            return Data.Errcode.NotCredible;
        }

        uint balanceOfContract = inkToken.balanceOf(this);
        if (balanceOfContract < amount) {
            return Data.Errcode.InsufficientBalance;
        }

        bool success = inkToken.transfer(toAccount, amount);
        if (!success) {
            return Data.Errcode.TransferFailed;
        }

        balanceOf[admin.name] -= amount;
        if (fromPlatform !=admin.name) {
            balanceOf[fromPlatform] -= amount;
            string memory value = uintAppendToString(amount);
            unlockEvent(txid, fromPlatform, fromAccount, value);
        }
        return Data.Errcode.Success;
    }


    /**
     *   ######################
     *  #  private function  #
     * ######################
     */

    function uintAppendToString(uint v) pure internal returns (string){
        uint maxlength = 100;
        bytes memory reversed = new bytes(maxlength);
        bytes32 sixTeenStr = "0123456789abcdef";

        uint i = 0;
        while (v != 0) {
            uint remainder = v % 16;
            v = v / 16;
            reversed[i++] = byte(sixTeenStr[remainder]);

        }
        string memory bytesList = "0000000000000000000000000000000000000000000000000000000000000000";
        bytes memory strb = bytes(bytesList);

        for (uint j = 0; j < i; j++) {
            strb[strb.length - j - 1] = reversed[i - j - 1];
        }
        return string(strb);
    }
}