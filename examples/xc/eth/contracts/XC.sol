pragma solidity ^0.4.19;

import "./XCInterface.sol";

import "./INK.sol";

import "./XCPlugin.sol";

contract XC is XCInterface {

    /**
     * Contract Administrator
     * @field status Contract external service status.
     * @field platformName Current contract platform name.
     * @field account Current contract administrator.
     */
    struct Admin {

        uint8 status;

        bytes32 platformName;

        address account;
    }

    Admin private admin;

    uint public lockBalance;

    INK private inkToken;

    XCPlugin private xcPlugin;

    event lockEvent(bytes32 toPlatform, address toAccount, string value);

    event unlockEvent(string txid, bytes32 fromPlatform, address fromAccount, string value);

    function XC(bytes32 name) public payable {

        admin = Admin(0, name, msg.sender);
    }

    function setStatus(uint8 status) external {

        require(admin.account == msg.sender);

        require(status == 0 || status == 1 || status == 2 || status == 3);

        if (admin.status != status) {

            admin.status = status;
        }
    }

    function getStatus() external constant returns (uint8) {

        return admin.status;
    }

    function kill() external {

        require(admin.account == msg.sender);

        selfdestruct(admin.account);
    }

    function setPlatformName(bytes32 platformName) external {

        require(admin.account == msg.sender);

        if (admin.platformName != platformName) {

            admin.platformName = platformName;
        }
    }

    function getPlatformName() external constant returns (bytes32) {

        return admin.platformName;
    }

    function setAdmin(address account) external {

        require(admin.account == msg.sender);

        admin.account = account;
    }

    function getAdmin() external constant returns (address) {

        return admin.account;
    }

    function setINK(address account) external {

        require(admin.account == msg.sender);

        if (inkToken != account) {

            inkToken = INK(account);
        }
    }

    function getINK() external constant returns (address) {

        return inkToken;
    }

    function setXCPlugin(address account) external {

        require(admin.account == msg.sender);

        if (xcPlugin != account) {

            xcPlugin = XCPlugin(account);
        }
    }

    function getXCPlugin() external constant returns (address) {

        return xcPlugin;
    }

    function lock(bytes32 toPlatform, address toAccount, uint value) external payable {

        require(admin.status == 2 || admin.status == 3);

        require(xcPlugin.getStatus());

        require(xcPlugin.existPlatform(toPlatform));

        require(value > 0);

        //get user approve the contract quota
        uint allowance = inkToken.allowance(msg.sender, this);

        require(allowance >= value);

        //do transferFrom
        bool success = inkToken.transferFrom(msg.sender, this, value);

        require(success);

        //record the amount of local platform turn out
        lockBalance += value;

        //trigger lockEvent
        lockEvent(toPlatform, toAccount, uintAppendToString(value));
    }


    //turn in
    function unlock(string txid, bytes32 fromPlatform, address fromAccount, address toAccount, uint value) external payable {

        require(admin.status == 1 || admin.status == 3);

        require(xcPlugin.getStatus());

        require(xcPlugin.existPlatform(fromPlatform));

        require(value > 0);

        //verify args by function xcPlugin.verify
        var (complete,verify) = xcPlugin.verifyProposal(fromPlatform, fromAccount, toAccount, value, txid);

        require(verify && !complete);

        //get contracts balance
        uint balance = inkToken.balanceOf(this);

        //validate the balance of contract were less than amount
        require(balance >= value);

        bool success = inkToken.transfer(toAccount, value);

        require(success);

        success = xcPlugin.commitProposal(fromPlatform, txid);

        require(success);

        lockBalance -= value;

        unlockEvent(txid, fromPlatform, fromAccount, uintAppendToString(value));
    }

    function withdraw(address account, uint value) external payable {

        require(admin.account == msg.sender);

        require(value > 0);

        uint balance = inkToken.balanceOf(this);

        require(balance - lockBalance >= value);

        bool success = inkToken.transfer(account, value);

        require(success);
    }

    function lockAdmin(bytes32 toPlatform, address toAccount, uint value) external payable {

        require(admin.account == msg.sender);

        require(admin.platformName == toPlatform || xcPlugin.existPlatform(toPlatform));

        require(value > 0);

        uint allowance = inkToken.allowance(msg.sender, this);

        require(allowance >= value);

        bool success = inkToken.transferFrom(msg.sender, this, value);

        require(success);

        lockBalance += value;

        if (admin.platformName != toPlatform && xcPlugin.existPlatform(toPlatform)) {

            lockEvent(toPlatform, toAccount, uintAppendToString(value));
        }
    }

    function unlockAdmin(string txid, bytes32 fromPlatform, address fromAccount, address toAccount, uint value) external payable {

        require(admin.account == msg.sender);

        require(xcPlugin.existPlatform(fromPlatform) || fromPlatform == admin.platformName);

        require(value > 0);

        uint balance = inkToken.balanceOf(this);

        require(balance >= 0);

        bool success = inkToken.transfer(toAccount, value);

        require(success);

        lockBalance -= value;

        if (fromPlatform != admin.platformName) {

            unlockEvent(txid, fromPlatform, fromAccount, uintAppendToString(value));
        }
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