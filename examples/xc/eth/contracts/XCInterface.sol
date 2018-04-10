pragma solidity ^0.4.19;

/**
 * XC Contract Interface.
 */
interface XCInterface {

    /**
     * Set contract service status.
     * @param status contract service status (0:closed;1:only-closed-lock;2:only-closed-unlock;3:opened;).
     */
    function setStatus(uint8 status) external;

    /**
     * Get contract service status.
     * @return contract service status.
     */
    function getStatus() external constant returns (uint8);

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
    function withdraw(address account, uint value) external payable;

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
     * @param fromAccount account of to platform.
     * @param toAccount account of to platform.
     * @param value transfer amount.
     */
    function unlockAdmin(string txid, bytes32 fromPlatform, address fromAccount, address toAccount, uint value) external payable;
}