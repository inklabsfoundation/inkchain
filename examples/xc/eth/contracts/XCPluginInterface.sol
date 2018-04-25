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
    function existPlatform(bytes32 name) external view returns (bool);

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
    function existPublicKey(bytes32 platformName, address publicKey) external view returns (bool);

    /**
     * Get the count of public key for the trusted platform.
     * @param platformName a platform name.
     * @return count of public key.
     */
    function countOfPublicKey(bytes32 platformName) external view returns (uint);

    /**
     * Get the list of public key for the trusted platform.
     * @param platformName a platform name.
     * @return list of public key.
     */
    function publicKeys(bytes32 platformName) external view returns (address[]);

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
    function getWeight(bytes32 platformName) external view returns (uint);

    /**
     * Initiate and vote on the transaction proposal.
     * @param fromPlatform name of form platform.
     * @param fromAccount name of to platform.
     * @param toAccount account of to platform.
     * @param value transfer amount.
     * @param txid transaction id.
     * @param r transaction signature (signature[0:64]).
     * @param s transaction signature (signature[64:128]).
     * @param v transaction signature (uint8(signature[128:130])).
     */
    function voteProposal(bytes32 fromPlatform, address fromAccount, address toAccount, uint value, string txid, bytes32 r, bytes32 s, uint8 v) external;

    /**
     * Verify that the transaction proposal is valid.
     * @param fromPlatform name of form platform.
     * @param fromAccount name of to platform.
     * @param toAccount account of to platform.
     * @param value transfer amount.
     * @param txid transaction id.
     */
    function verifyProposal(bytes32 fromPlatform, address fromAccount, address toAccount, uint value, string txid) external view returns (bool, bool);

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
     * @return status completion status of proposal.
     * @return fromAccount account of to platform.
     * @return toAccount account of to platform.
     * @return value transfer amount.
     * @return voters notarial voters.
     * @return weight The weight value of the completed time.
     */
    function getProposal(bytes32 platformName, string txid) external view returns (bool status, address fromAccount, address toAccount, uint value, address[] voters, uint weight);

    /**
     * Delete the transaction proposal information.
     * @param platformName a platform name.
     * @param txid transaction id.
     */
    function deleteProposal(bytes32 platformName, string txid) external;

    /**
     * Transfer the money(qtum/eth) from the contract account.
     * @param account the specified account.
     * @param value transfer amount.
     */
    function transfer(address account, uint value) external payable;
}