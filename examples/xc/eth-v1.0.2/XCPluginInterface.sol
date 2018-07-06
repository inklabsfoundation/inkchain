pragma solidity ^0.4.19;

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
    function getStatus() external view returns (bool);

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
     * Get the current token symbol.
     * @return token symbol.
     */
    function getTokenSymbol() external view returns (bytes32);

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
    function existCaller(address caller) external view returns (bool);

    /**
     * Get all contract trusted callers.
     * @return al lcallers.
     */
    function getCallers() external view returns (address[]);

    /**
     * Get the trusted platform name.
     * @return name a platform name.
     */
    function getTrustPlatform() external view returns (bytes32 name);

    /**
     * Add the trusted platform public key information.
     * @param publicKey a public key.
     */
    function addPublicKey(address publicKey) external;

    /**
     * Delete the trusted platform public key information.
     * @param publicKey a public key.
     */
    function deletePublicKey(address publicKey) external;

    /**
     * Whether the trusted platform public key information exists.
     * @param publicKey a public key.
     */
    function existPublicKey(address publicKey) external view returns (bool);

    /**
     * Get the count of public key for the trusted platform.
     * @return count of public key.
     */
    function countOfPublicKey() external view returns (uint);

    /**
     * Get the list of public key for the trusted platform.
     * @return list of public key.
     */
    function publicKeys() external view returns (address[]);

    /**
     * Set the weight of a trusted platform.
     * @param weight weight of platform.
     */
    function setWeight(uint weight) external;

    /**
     * Get the weight of a trusted platform.
     * @return weight of platform.
     */
    function getWeight() external view returns (uint);

    /**
     * Initiate and vote on the transaction proposal.
     * @param fromAccount name of to platform.
     * @param toAccount account of to platform.
     * @param value transfer amount.
     * @param txid transaction id.
     * @param sig transaction signature.
     */
    function voteProposal(address fromAccount, address toAccount, uint value, string txid, bytes sig) external;

    /**
     * Verify that the transaction proposal is valid.
     * @param fromAccount name of to platform.
     * @param toAccount account of to platform.
     * @param value transfer amount.
     * @param txid transaction id.
     */
    function verifyProposal(address fromAccount, address toAccount, uint value, string txid) external view returns (bool, bool);

    /**
     * Commit the transaction proposal.
     * @param txid transaction id.
     */
    function commitProposal(string txid) external returns (bool);

    /**
     * Get the transaction proposal information.
     * @param txid transaction id.
     * @return status completion status of proposal.
     * @return fromAccount account of to platform.
     * @return toAccount account of to platform.
     * @return value transfer amount.
     * @return voters notarial voters.
     * @return weight The weight value of the completed time.
     */
    function getProposal(string txid) external view returns (bool status, address fromAccount, address toAccount, uint value, address[] voters, uint weight);

    /**
     * Delete the transaction proposal information.
     * @param txid transaction id.
     */
    function deleteProposal(string txid) external;
}