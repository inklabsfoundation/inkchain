pragma solidity ^0.4.19;

import "./XCPluginInterface.sol";

contract XCPlugin is XCPluginInterface {

    /**
     * Contract Administrator
     * @field status Contract external service status.
     * @field platformName Current contract platform name.
     * @field account Current contract administrator.
     */
    struct Admin {

        bool status;

        bytes32 platformName;

        address account;
    }

    /**
     * Transaction Proposal
     * @field status Transaction proposal status(false:pending,true:complete).
     * @field fromAccount Account of form platform.
     * @field toAccount Account of to platform.
     * @field value Transfer amount.
     * @field voters Proposers.
     * @field weight The weight value of the completed time.
     */
    struct Proposal {

        bool status;

        address fromAccount;

        address toAccount;

        uint value;

        address[] voters;

        uint weight;
    }

    /**
     * Trusted Platform
     * @field status Trusted platform state(false:no trusted,true:trusted).
     * @field weight weight of platform.
     * @field publicKeys list of public key.
     * @field proposals list of proposal.
     */
    struct Platform {

        bool status;

        uint weight;

        address[] publicKeys;

        mapping(string => Proposal) proposals;
    }

    Admin private admin;

    address[] private callers;

    mapping(bytes32 => Platform) private platforms;

    function XCPlugin() public {
        //TODO
        bytes32 name = "ETH";

        admin = Admin(false, name, msg.sender);
    }

    function start() external {

        require(admin.account == msg.sender);

        if (!admin.status) {

            admin.status = true;
        }
    }

    function stop() external {

        require(admin.account == msg.sender);

        if (admin.status) {

            admin.status = false;
        }
    }

    function getStatus() external view returns (bool) {

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

    function addCaller(address caller) external {

        require(admin.account == msg.sender);

        if (!_existCaller(caller)) {

            callers.push(caller);
        }
    }

    function deleteCaller(address caller) external {

        require(admin.account == msg.sender);

        if (_existCaller(caller)) {

            bool exist;

            for (uint i = 0; i <= callers.length; i++) {

                if (exist) {

                    if (i == callers.length) {

                        delete callers[i - 1];

                        callers.length--;
                    } else {

                        callers[i - 1] = callers[i];
                    }
                } else if (callers[i] == caller) {

                    exist = true;
                }
            }

        }
    }

    function existCaller(address caller) external view returns (bool) {

        return _existCaller(caller);
    }

    function _existCaller(address caller) internal view returns (bool) {

        for (uint i = 0; i < callers.length; i++) {

            if (callers[i] == caller) {

                return true;
            }
        }

        return false;
    }

    function getCallers() external view returns (address[]) {

        require(admin.account == msg.sender);

        return callers;
    }

    function addPlatform(bytes32 name) external {

        require(admin.account == msg.sender);

        require(name != "");

        require(name != admin.platformName);

        if (!_existPlatform(name)) {

            platforms[name].status = true;

            if (platforms[name].weight == 0) {

                platforms[name].weight = 1;
            }
        }
    }

    function deletePlatform(bytes32 name) external {

        require(admin.account == msg.sender);

        require(name != admin.platformName);

        if (_existPlatform(name)) {

            platforms[name].status = false;
        }
    }

    function existPlatform(bytes32 name) external view returns (bool){

        return _existPlatform(name);
    }

    function _existPlatform(bytes32 name) internal view returns (bool){

        return platforms[name].status;
    }

    function setWeight(bytes32 platformName, uint weight) external {

        require(admin.account == msg.sender);

        require(_existPlatform(platformName));

        require(weight > 0);

        if (platforms[platformName].weight != weight) {

            platforms[platformName].weight = weight;
        }
    }

    function getWeight(bytes32 platformName) external view returns (uint) {

        require(admin.account == msg.sender);

        require(_existPlatform(platformName));

        return platforms[platformName].weight;
    }

    function addPublicKey(bytes32 platformName, address publicKey) external {

        require(admin.account == msg.sender);

        require(_existPlatform(platformName));

        require(publicKey != address(0));

        address[] storage listOfPublicKey = platforms[platformName].publicKeys;

        for (uint i; i < listOfPublicKey.length; i++) {

            if (publicKey == listOfPublicKey[i]) {

                return;
            }
        }

        listOfPublicKey.push(publicKey);
    }

    function deletePublicKey(bytes32 platformName, address publickey) external {

        require(admin.account == msg.sender);

        require(_existPlatform(platformName));

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
    }

    function existPublicKey(bytes32 platformName, address publicKey) external view returns (bool) {

        require(admin.account == msg.sender);

        return _existPublicKey(platformName, publicKey);
    }

    function _existPublicKey(bytes32 platformName, address publicKey) internal view returns (bool) {

        address[] memory listOfPublicKey = platforms[platformName].publicKeys;

        for (uint i = 0; i < listOfPublicKey.length; i++) {

            if (listOfPublicKey[i] == publicKey) {

                return true;
            }
        }

        return false;
    }

    function countOfPublicKey(bytes32 platformName) external view returns (uint){

        require(admin.account == msg.sender);

        require(_existPlatform(platformName));

        return platforms[platformName].publicKeys.length;
    }

    function publicKeys(bytes32 platformName) external view returns (address[]){

        require(admin.account == msg.sender);

        require(_existPlatform(platformName));

        return platforms[platformName].publicKeys;
    }

    function voteProposal(bytes32 fromPlatform, address fromAccount, address toAccount, uint value, string txid, bytes32 r, bytes32 s, uint8 v) external {

        require(admin.status);

        require(_existPlatform(fromPlatform));

        bytes32 msgHash = hashMsg(fromPlatform, fromAccount, admin.platformName, toAccount, value, txid);

        address publicKey = ecrecover(msgHash, v, r, s);

        require(_existPublicKey(fromPlatform, publicKey));

        Proposal storage proposal = platforms[fromPlatform].proposals[txid];

        if (proposal.value == 0) {

            proposal.fromAccount = fromAccount;

            proposal.toAccount = toAccount;

            proposal.value = value;
        } else {

            require(proposal.fromAccount == fromAccount && proposal.toAccount == toAccount && proposal.value == value);
        }

        changeVoters(fromPlatform, publicKey, txid);
    }

    function verifyProposal(bytes32 fromPlatform, address fromAccount, address toAccount, uint value, string txid) external view returns (bool, bool) {

        require(admin.status);

        require(_existPlatform(fromPlatform));

        Proposal storage proposal = platforms[fromPlatform].proposals[txid];

        if (proposal.status) {

            return (true, (proposal.voters.length >= proposal.weight));
        }

        if (proposal.value == 0) {

            return (false, false);
        }

        require(proposal.fromAccount == fromAccount && proposal.toAccount == toAccount && proposal.value == value);

        return (false, (proposal.voters.length >= platforms[fromPlatform].weight));
    }

    function commitProposal(bytes32 platformName, string txid) external returns (bool) {

        require(admin.status);

        require(_existCaller(msg.sender) || msg.sender == admin.account);

        require(_existPlatform(platformName));

        require(!platforms[platformName].proposals[txid].status);

        platforms[platformName].proposals[txid].status = true;

        platforms[platformName].proposals[txid].weight = platforms[platformName].proposals[txid].voters.length;

        return true;
    }

    function getProposal(bytes32 platformName, string txid) external view returns (bool status, address fromAccount, address toAccount, uint value, address[] voters, uint weight){

        require(admin.status);

        require(_existPlatform(platformName));

        fromAccount = platforms[platformName].proposals[txid].fromAccount;

        toAccount = platforms[platformName].proposals[txid].toAccount;

        value = platforms[platformName].proposals[txid].value;

        voters = platforms[platformName].proposals[txid].voters;

        status = platforms[platformName].proposals[txid].status;

        weight = platforms[platformName].proposals[txid].weight;

        return;
    }

    function deleteProposal(bytes32 platformName, string txid) external {

        require(msg.sender == admin.account);

        require(_existPlatform(platformName));

        delete platforms[platformName].proposals[txid];
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

    function hashMsg(bytes32 fromPlatform, address fromAccount, bytes32 toPlatform, address toAccount, uint value, string txid) internal pure returns (bytes32) {

        return sha256(bytes32ToStr(fromPlatform), ":0x", uintToStr(uint160(fromAccount), 16), ":", bytes32ToStr(toPlatform), ":0x", uintToStr(uint160(toAccount), 16), ":", uintToStr(value, 10), ":", txid);
    }

    function changeVoters(bytes32 platformName, address publicKey, string txid) internal {

        address[] storage voters = platforms[platformName].proposals[txid].voters;

        bool change = true;

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

        uint length = b.length;

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

    function bytes32ToStr(bytes32 b) internal pure returns (string) {

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

        if (base == 16) {
            length = 40;
        }

        bytes memory bs = new bytes(length);

        for (uint i = 0; i < length; i++) {

            bs[length - 1 - i] = tenStr[value % base];

            value = value / base;
        }

        return string(bs);
    }

}
