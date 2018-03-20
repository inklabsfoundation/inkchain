# XC

## Data Structure 

```
library Data {

    enum ErrCode {
        Success,
        NotAdmin,
        StatusClosed,
        PlatformNameNotNull,
        CatNotOwnerPlatformName,
        NotCredible,
        InvalidTransferAmount,
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

contract XCPlugin {
    Data.Admin private admin;
    mapping(bytes32 => Data.Platform) private platforms;
}

contract XC {
    Data.Admin private admin;
    uint public lockBalance;
}
```

> 1）Data.Errcode : return code
> 
> 2）Data.Admin : contract admin info
> 
> 3）Data.Platform : xc platform info
>
> 4）XCPlugin.platforms : credible platforms
> 
> 5）XC.lockBalance : XC contract total lock amount
> 

remark：balanceOf Used for checking account transfer between various platforms.


## Event

```
contract XC {
    event lockEvent(bytes32 toPlatform, address toAccount, string amount);
    event unlockEvent(bytes32 txId, bytes32 fromPlatform, address fromAccount, string amount);
}
```

> XC transaction trigger.

## Interface 

```
interface XCInterface {

    function start() external;

    function stop() external;

    function getStatus() external constant returns (bool);

    function setPlatformName(bytes32 platformName) external;

    function getPlatformName() external constant returns (bytes32);

    function setAdmin(address account) external;

    function getAdmin() external constant returns (address);

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
interface XCPluginInterface {

    function start() external;

    function stop() external;

    function getStatus() external constant returns (bool);

    function setPlatformName(bytes32 platformName) external;

    function getPlatformName() external constant returns (bytes32);

    function setAdmin(address account) external;

    function getAdmin() external constant returns (address);

    function addCaller(address caller) external;

    function deleteCaller(address caller) external;

    function addPlatform(bytes32 name) external returns (Data.ErrCode);

    function deletePlatform(bytes32 name) external returns (Data.ErrCode);

    function getPlatform(bytes32 platformName) external constant returns (Data.ErrCode ErrCode, bytes32 name, uint weight, address[] publicKeys);

    function existPlatform(bytes32 name) external constant returns (bool);

    function addPublicKey(bytes32 platformName, address publicKey) external returns (Data.ErrCode);

    function deletePublicKey(bytes32 platformName, address publicKey) external returns (Data.ErrCode);

    function countOfPublicKey(bytes32 platformName) external constant returns (Data.ErrCode, uint);

    function setWeight(bytes32 name, uint weight) external returns (Data.ErrCode);

    function getWeight(bytes32 name) external constant returns (Data.ErrCode, uint);

    function voter(bytes32 fromPlatform, address fromAccount, address toAccount, uint amount, bytes32 txId, bytes32 r, bytes32 s, uint8 v) external returns (Data.ErrCode ErrCode, bool verify);

    function verifyProposal(bytes32 fromPlatform, address fromAccount, address toAccount, uint amount, bytes32 txId) external constant returns (Data.ErrCode);

    function deleteProposal(bytes32 platformName, bytes32 txId) external returns (Data.ErrCode);
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