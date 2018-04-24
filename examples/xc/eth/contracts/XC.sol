pragma solidity ^0.4.19;

import "./XCInterface.sol";

import "./Token.sol";

import "./XCPlugin.sol";

import "./SafeMath.sol";

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

    Token private token;

    XCPlugin private xcPlugin;

    bytes2 private compareSymbol;

    event Lock(bytes32 toPlatform, address toAccount, string value);

    event Unlock(string txid, bytes32 fromPlatform, address fromAccount, string value);

    function XC() public payable {
        //TODO
        bytes32 name = "ETH";
        //TODO
        compareSymbol = "+=";

        //TODO totalSupply = 10 * (10 ** 8) * (10 ** 9);
        lockBalance = 10 * (10 ** 8) * (10 ** 9);

        admin = Admin(0, name, msg.sender);
    }

    function setStatus(uint8 status) external {

        require(admin.account == msg.sender);

        require(status == 0 || status == 1 || status == 2 || status == 3);

        if (admin.status != status) {

            admin.status = status;
        }
    }

    function getStatus() external view returns (uint8) {

        return admin.status;
    }

    function setPlatformName(bytes32 platformName) external {

        require(admin.account == msg.sender);

        if (admin.platformName != platformName) {

            admin.platformName = platformName;
        }
    }

    function getPlatformName() external view returns (bytes32) {

        return admin.platformName;
    }

    function setAdmin(address account) external {

        require(account != address(0));

        require(admin.account == msg.sender);

        if (admin.account != account) {

            admin.account = account;
        }
    }

    function getAdmin() external view returns (address) {

        return admin.account;
    }

    function setToken(address account) external {

        require(admin.account == msg.sender);

        if (token != account) {

            token = Token(account);
        }
    }

    function getToken() external view returns (address) {

        return token;
    }

    function setXCPlugin(address account) external {

        require(admin.account == msg.sender);

        if (xcPlugin != account) {

            xcPlugin = XCPlugin(account);
        }
    }

    function getXCPlugin() external view returns (address) {

        return xcPlugin;
    }

    function setCompare(bytes2 symbol) external {

        require(admin.account == msg.sender);

        require(symbol == "+=" || symbol == "-=");

        if (compareSymbol != symbol) {

            compareSymbol = symbol;
        }
    }

    function getCompare() external view returns (bytes2){

        require(admin.account == msg.sender);

        return compareSymbol;
    }

    function lock(bytes32 toPlatform, address toAccount, uint value) external payable {

        require(admin.status == 2 || admin.status == 3);

        require(xcPlugin.getStatus());

        require(xcPlugin.existPlatform(toPlatform));

        require(toAccount != address(0));

        // require(token.totalSupply >= value && value > 0);
        require(value > 0);

        //get user approve the contract quota
        uint allowance = token.allowance(msg.sender, this);

        require(toCompare(allowance, value));

        //do transferFrom
        bool success = token.transferFrom(msg.sender, this, value);

        require(success);

        //record the amount of local platform turn out
        lockBalance = SafeMath.add(lockBalance, value);
        // require(token.totalSupply >= lockBalance);

        //trigger Lock
        emit Lock(toPlatform, toAccount, uintAppendToString(value));
    }

    //turn in
    function unlock(string txid, bytes32 fromPlatform, address fromAccount, address toAccount, uint value) external payable {

        require(admin.status == 1 || admin.status == 3);

        require(xcPlugin.getStatus());

        require(xcPlugin.existPlatform(fromPlatform));

        require(toAccount != address(0));

        // require(token.totalSupply >= value && value > 0);
        require(value > 0);

        //verify args by function xcPlugin.verify
        bool complete;

        bool verify;

        (complete, verify) = xcPlugin.verifyProposal(fromPlatform, fromAccount, toAccount, value, txid);

        require(verify && !complete);

        //get contracts balance
        uint balance = token.balanceOf(this);

        //validate the balance of contract were less than amount
        require(toCompare(balance, value));

        require(token.transfer(toAccount, value));

        require(xcPlugin.commitProposal(fromPlatform, txid));

        lockBalance = SafeMath.sub(lockBalance, value);

        emit Unlock(txid, fromPlatform, fromAccount, uintAppendToString(value));
    }

    function withdraw(address account, uint value) external payable {

        require(admin.account == msg.sender);

        require(account != address(0));

        // require(token.totalSupply >= value && value > 0);
        require(value > 0);

        uint balance = token.balanceOf(this);

        require(toCompare(SafeMath.sub(balance, lockBalance), value));

        bool success = token.transfer(account, value);

        require(success);
    }

    function transfer(address account, uint value) external payable {

        require(admin.account == msg.sender);

        require(account != address(0));

        require(value > 0 && value >= this.balance);

        this.transfer(account, value);
    }

    /**
     *   ######################
     *  #  private function  #
     * ######################
     */

    function toCompare(uint f, uint s) internal view returns (bool) {

        if (compareSymbol == "-=") {

            return f > s;
        } else if (compareSymbol == "+=") {

            return f >= s;
        } else {

            return false;
        }
    }

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