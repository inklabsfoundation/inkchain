# XC

## Data Structs 

```

/**
 * Data Structs
 */
library Data {

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
     */
    struct Proposal {
        bool status;
        address fromAccount;
        address toAccount;
        uint value;
        address[] voters;
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
}

contract XCPlugin {
    Data.Admin private admin;
    mapping(bytes32 => Data.Platform) private platforms;
}

contract XC {
    Data.Admin private admin;
    uint public lockBalance;
}
```

> 1）Data.Admin : contract admin info
> 
> 2）Data.Platform : xc platform info
>
> 3）XCPlugin.platforms : trusted platforms
> 
> 4）XC.lockBalance : XC contract total lock amount
> 

remark：balanceOf Used for checking account transfer between various platforms.


## Event

```
contract XC {
    event lockEvent(bytes32 toPlatform, address toAccount, string value);
    event unlockEvent(bytes32 txId, bytes32 fromPlatform, address fromAccount, string value);
}
```

> XC transaction trigger.

## Interface 

```
/**
 * XC Contract Interface.
 */
interface XCInterface {

    /**
     * Open the contract service status.
     */
    function start() external;

    /**
     * Close the contract service status.
     */
    function stop() external;

    /**
     * Get contract service status.
     * @return contract service status.
     */
    function getStatus() external constant returns (bool);

    /**
     * Destruction of the contract.
     */
    function kill() external;

    /**
     * Set the current contract platform name.
     * @param platformName platform name.
     */
    function setPlatformName(bytes32 platformName) external;

    /**
     * Get the current contract platform name.
     * @return contract platform name.
     */
    function getPlatformName() external constant returns (bytes32);

    /**
     * Set the current contract administrator.
     * @param account account of contract administrator.
     */
    function setAdmin(address account) external;

    /**
     * Get the current contract administrator.
     * @return contract administrator.
     */
    function getAdmin() external constant returns (address);

    /**
     * Set the INK contract address.
     * @param account contract address.
     */
    function setINK(address account) external;

    /**
     * Get the INK contract address.
     * @return contract address.
     */
    function getINK() external constant returns (address);

    /**
     * Set the XCPlugin contract address.
     * @param account contract address.
     */
    function setXCPlugin(address account) external;

    /**
     * Get the XCPlugin contract address.
     * @return contract address.
     */
    function getXCPlugin() external constant returns (address);

    /**
     * Transfer out of cross chain.
     * @param toPlatform name of to platform.
     * @param toAccount account of to platform.
     * @param value transfer amount.
     */
    function lock(bytes32 toPlatform, address toAccount, uint value) external payable;

    /**
     * Transfer in of cross chain.
     * @param txid transaction id.
     * @param fromPlatform name of form platform.
     * @param fromAccount ame of to platform.
     * @param toAccount account of to platform.
     * @param value transfer amount.
     */
    function unlock(string txid, bytes32 fromPlatform, address fromAccount, address toAccount, uint value) external payable;

    /**
     * Transfer the misoperation to the amount of the contract account to the specified account.
     * @param account the specified account.
     * @param value transfer amount.
     */
    function withdrawal(address account, uint value) external payable;

    /**
     * Administrator’s transfer out of cross chain.
     * @param toPlatform name of to platform.
     * @param toAccount account of to platform.
     * @param value transfer amount.
     */
    function lockAdmin(bytes32 toPlatform, address toAccount, uint value) external payable;

    /**
     * Administrator’s transfer in of cross chain.
     * @param txid transaction id.
     * @param fromPlatform name of form platform.
     * @param fromAccount ame of to platform.
     * @param toAccount account of to platform.
     * @param value transfer amount.
     */
    function unlockAdmin(string txid, bytes32 fromPlatform, address fromAccount, address toAccount, uint value) external payable;
}
```
> 1）setPlatformName & getPlatformName：Used to maintain the name of the contract platform.
>
> 2）start & stop & getStatus：Used to maintain contract service status.
>
> 3）setAdmin & getAdmin : Used to maintain contract administrators, transferable.
>
> 4）setINK、getINK、setXCPlugin、getXCPlugin : Used to maintain INK and XCPlugin contracts. can upgrade.
>
> 5）lockAdmin、unlockAdmin : Admin cross - chain ledger usage.
>
> 6）withdrawal : Used to extract the amount of money that someone has mistakenly transferred to a contract.
>
> 7）lock 和 unlock : Used for cross - chain transfer.
> 

remark：1）2）3）4）5）6）The administrator maintains the interface. 7）Open interface.

```
/**
 * XC Plugin Contract Interface.
 */
interface XCPluginInterface {

    /**
     * Open the contract service status.
     */
    function start() external;

    /**
     * Close the contract service status.
     */
    function stop() external;

    /**
     * Get contract service status.
     * @return contract service status.
     */
    function getStatus() external constant returns (bool);

    /**
     * Destruction of the contract.
     */
    function kill() external;

    /**
     * Set the current contract platform name.
     * @param platformName platform name.
     */
    function setPlatformName(bytes32 platformName) external;

    /**
     * Get the current contract platform name.
     * @return contract platform name.
     */
    function getPlatformName() external constant returns (bytes32);

    /**
     * Set the current contract administrator.
     * @param account account of contract administrator.
     */
    function setAdmin(address account) external;

    /**
     * Get the current contract administrator.
     * @return contract administrator.
     */
    function getAdmin() external constant returns (address);

    /**
     * Add a contract trust caller.
     * @param caller account of caller.
     */
    function addCaller(address caller) external;

    /**
     * Delete a contract trust caller.
     * @param caller account of caller.
     */
    function deleteCaller(address caller) external;

    /**
     * Whether the trust caller exists.
     * @param caller account of caller.
     * @return whether exists.
     */
    function existCaller(address caller) public constant returns (bool);

    /**
     * Get all contract trusted callers.
     * @return al lcallers.
     */
    function getCallers() external constant returns (address[]);

    /**
     * Add a trusted platform name.
     * @param name a platform name.
     */
    function addPlatform(bytes32 name) external;

    /**
     * Delete a trusted platform name.
     * @param name a platform name.
     */
    function deletePlatform(bytes32 name) external;

    /**
     * Whether the trusted platform information exists.
     * @param name a platform name.
     * @return whether exists.
     */
    function existPlatform(bytes32 name) public constant returns (bool);

    /**
     * Add the trusted platform public key information.
     * @param platformName a platform name.
     * @param publicKey a public key.
     */
    function addPublicKey(bytes32 platformName, address publicKey) external;

    /**
     * Delete the trusted platform public key information.
     * @param platformName a platform name.
     * @param publicKey a public key.
     */
    function deletePublicKey(bytes32 platformName, address publicKey) external;

    /**
     * Whether the trusted platform public key information exists.
     * @param platformName a platform name.
     * @param publicKey a public key.
     */
    function existPublicKey(bytes32 platformName, address publicKey) public constant returns (bool);

    /**
     * Get the count of public key for the trusted platform.
     * @param platformName a platform name.
     * @return count of public key.
     */
    function countOfPublicKey(bytes32 platformName) external constant returns (uint);

    /**
     * Get the list of public key for the trusted platform.
     * @param platformName a platform name.
     * @return list of public key.
     */
    function publicKeys(bytes32 platformName) external constant returns (address[]);

    /**
     * Set the weight of a trusted platform.
     * @param platformName a platform name.
     * @param weight weight of platform.
     */
    function setWeight(bytes32 platformName, uint weight) external;

    /**
     * Get the weight of a trusted platform.
     * @param platformName a platform name.
     * @return weight of platform.
     */
    function getWeight(bytes32 platformName) external constant returns (uint);

    /**
     * Initiate and vote on the transaction proposal.
     * @param fromPlatform name of form platform.
     * @param fromAccount name of to platform.
     * @param toAccount account of to platform.
     * @param value transfer amount.
     * @param txid transaction id.
     * @param sign transaction signature.
     */
    function voterProposal(bytes32 fromPlatform, address fromAccount, address toAccount, uint value, string txid, bytes sign) external;

    /**
     * Verify that the transaction proposal is valid.
     * @param fromPlatform name of form platform.
     * @param fromAccount name of to platform.
     * @param toAccount account of to platform.
     * @param value transfer amount.
     * @param txid transaction id.
     */
    function verifyProposal(bytes32 fromPlatform, address fromAccount, address toAccount, uint value, string txid) external constant returns (bool);

    /**
     * Commit the transaction proposal.
     * @param platformName a platform name.
     * @param txid transaction id.
     */
    function commitProposal(bytes32 platformName, string txid) external returns (bool);

    /**
     * Get the transaction proposal information.
     * @param platformName a platform name.
     * @param txid transaction id.
     * @return status
     * @return fromAccount
     * @return toAccount
     * @return value
     * @return voters
     */
    function getProposal(bytes32 platformName, string txid) external returns (bool status, address fromAccount, address toAccount, uint value, address[] voters);
}
```

> 1）setPlatformName & getPlatformName：Used to maintain the name of the contract platform.
>
> 2）start & stop & getStatus：Used to maintain contract service status.
>
> 2）setAdmin & getAdmin : Used to maintain contract administrators, transferable.
>
> 3）addPlatform、deletePlatfrom、getPlatfrom、existPlatfrom : Used to maintain information of trusted chain platform.
>
> 4）addPublickey、deletePublickey、countOfPublickey : Used to maintain the credible public key of the trusted chain platform.
>
> 5）setWeight、getWeight : Used to maintain the weight of trusted chain platform.
>
> 6）voter : Attestation to vote; verify : Verify the legality of cross-chain transactions. deleteProposal : Remove the proposal of vote.
> 

remark：1）2）3）4） 5）The administrator maintains the interface. 6）Open interface.

## Contract

```
contract INK {}
contract XCPlugin {}
contract XC {
    INK private inkToken;
    XCPlugin private xcPlugin;
}
```
> 1）INK : INK Token contract. XCPlugin : Cross-chain functional contracts. XC : Across the chain contracts.
>
> 2）INK & XCPlugin are XC function plug-in contracts.
> 

## Using

### Install INK contract
...

### Install XCPlugin contract
```
1）Using addPlatform, deletePlatfrom, getPlatfrom, existPlatfrom method to maintain trusted platform information .
2）Use the addPublickey, deletePublickey, countOfPublickey method to maintain the trusted public key information of the trusted platform.
3）Use setWeight and getWeight method to set the verification weight of each trusted platform.
```

### Install XC contract；

```
On XC contract， Using setINK、setXCPlugin method  setup INK contract & XCPlugin contract.
Use getINK & getXCPlugin method to check.
```