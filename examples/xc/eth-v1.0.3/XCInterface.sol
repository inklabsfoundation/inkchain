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
    function getStatus() external view returns (uint8);

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
     * Transfer out of cross chain.
     * @param toAccount account of to platform.
     * @param value transfer amount.
     */
    function lock(address toAccount, uint value) external;

    /**
     * Transfer in of cross chain.
     * @param txid transaction id.
     * @param fromAccount ame of to platform.
     * @param toAccount account of to platform.
     * @param value transfer amount.
     */
    function unlock(string txid, address fromAccount, address toAccount, uint value) external;

    /**
     * Transfer the misoperation to the amount of the contract account to the specified account.
     * @param account the specified account.
     * @param value transfer amount.
     */
    function withdraw(address account, uint value) external;
}