pragma solidity ^0.4.21;

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
    function getStatus() external view returns (uint8);

    /**
     * Set the current contract platform name.
     * @param platformName platform name.
     */
    function setPlatformName(bytes32 platformName) external;

    /**
     * Get the current contract platform name.
     * @return contract platform name.
     */
    function getPlatformName() external view returns (bytes32);

    /**
     * Set the current contract administrator.
     * @param account account of contract administrator.
     */
    function setAdmin(address account) external;

    /**
     * Get the current contract administrator.
     * @return contract administrator.
     */
    function getAdmin() external view returns (address);

    /**
     * Set the Token contract address.
     * @param account contract address.
     */
    function setToken(address account) external;

    /**
     * Get the Token contract address.
     * @return contract address.
     */
    function getToken() external view returns (address);

    /**
     * Set the XCPlugin contract address.
     * @param account contract address.
     */
    function setXCPlugin(address account) external;

    /**
     * Get the XCPlugin contract address.
     * @return contract address.
     */
    function getXCPlugin() external view returns (address);

    /**
     * Set the comparison symbol in the contract.
     * @param symbol comparison symbol ({"-=" : ">" , "+=" : ">=" }).
     */
    function setCompare(bytes2 symbol) external;

    /**
     * Get the comparison symbol in the contract.
     * @return comparison symbol.
     */
    function getCompare() external view returns (bytes2);

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
     * Transfer the money(qtum/eth) from the contract account.
     * @param account the specified account.
     * @param value transfer amount.
     */
    function transfer(address account, uint value) external payable;
}