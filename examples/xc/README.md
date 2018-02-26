# XC

## Data Structure 

```
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

    struct Platform {
        uint8 typ;
        bytes32 name;
        uint weight;
        address[] publickeys;
    }
}

contract XCPlugin {
    Data.Admin private admin;
    mapping(bytes32 => Data.Platform) private platforms;
}

contract XC {
    Data.Admin private admin;
    mapping(bytes32 => uint) public balanceOf;
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
> 5）XC.balanceOf : xc amount；XC.balanceOf[admin.name] : XC contract total lock amount；XC.balanceOf[platformName] : xc platform total amount；
> 

remark：balanceOf Used for checking account transfer between various platforms.


## Event

```
contract XC {
    event lockEvent(bytes32 toPlatform, address toAccount, string amount);
    event unlockEvent(bytes32 txid,bytes32 fromPlatform,address fromAccount ,string amount);
}
```

> XC transaction trigger.

## Interface 

```
interface XCInterface {

    function setAdmin(address account) external;
    function getAdmin() external constant returns (address);
    
    function setINK(address account) external;
    function getINK() external constant returns (address);
    
    function setXCPlugin(address account) external;
    function getXCPlugin() external constant returns (address);
    
    function lockAdmin(bytes32 toPlatform, address toAccount, uint amount) external payable returns (Data.Errcode);
    function unlockAdmin(bytes32 fromPlatform,address fromAccount, address toAccount, uint amount, bytes32 txid) external payable returns (Data.Errcode);
    
    function withdrawal(address account,uint amount) external payable returns (Data.Errcode);
    
    function lock(bytes32 toPlatform, address toAccount, uint amount) external  payable returns (Data.Errcode);
    function unlock(bytes32 fromPlatform,address fromAccount, address toAccount, uint amount, bytes32 txid) external payable returns (Data.Errcode);  
}
```
> 1）setAdmin & getAdmin : Used to maintain contract administrators, transferable.
>
> 2）setINK、getINK、setXCPlugin、getXCPlugin : Used to maintain INK and XCPlugin contracts. can upgrade.
>
> 3）lockAdmin、unlockAdmin : Admin cross - chain ledger usage.
>
> 4）withdrawal : Used to extract the amount of money that someone has mistakenly transferred to a contract.
>
> 5）lock 和 unlock : Used for cross - chain transfer.
> 

remark：1）2）3）4）The administrator maintains the interface. 5）Open interface.

```
interface XCPluginInterface { 

    function voter(bytes32 name, bytes32 txid, bytes sig) external returns (Data.Errcode errcode,bool verify);
    function verify(bytes32 name, bytes32 txid) external constant returns (Data.Errcode);
    function deleteProposal(bytes32 platformName, bytes32 txid) external constant returns (Data.Errcode);
    
    function setAdmin(bytes32 name,address account) external;
    function getAdmin() external constant returns (bytes32,address);
    
    function addPlatform(uint8 typ, bytes32 name) external returns (Data.Errcode);
    function deletePlatfrom(bytes32 name) external constant returns (Data.Errcode);
    function getPlatfrom(bytes32 name) external returns (Data.Errcode errcode, uint8 typ, bytes32 name, uint balance, uint weight, address[] publickeys);
    function existPlatfrom(bytes32 name) external constant returns (bool);
    
    function setWeight(bytes32 name, uint weight) external returns (Data.Errcode);
    function getWeight(bytes32 name) external constant returns (Data.Errcode, uint);
    
    function addPublickey(bytes32 platfromName, address publickey) external constant returns (Data.Errcode);
    function deletePublickey(bytes32 platfromName, address publickey) external returns (Data.Errcode);
    function countOfPublickey(bytes32 platfromName) external constant returns (Data.Errcode, uint);
}
```
> 1）setAdmin & getAdmin : Used to maintain contract administrators, transferable.
>
> 2）addPlatform、deletePlatfrom、getPlatfrom、existPlatfrom : Used to maintain information of trusted chain platform.
>
> 3）addPublickey、deletePublickey、countOfPublickey : Used to maintain the credible public key of the trusted chain platform.
>
> 4）setWeight、getWeight : Used to maintain the weight of trusted chain platform.
>
> 5）voter : Attestation to vote; verify : Verify the legality of cross-chain transactions. deleteProposal : Remove the proposal of vote.
> 

remark：1）2）3）4）The administrator maintains the interface. 5）Open interface.

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