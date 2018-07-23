pragma solidity ^0.4.19;

interface TokenInterface {

    function setStatus(bool _status) external;

    function getStatus() external view returns (bool);

    function setAdmin(address account) external;

    function getAdmin() external view returns (address);

    function balanceOf(address owner) external view returns (uint);

    function transfer(address to, uint value) external returns (bool);

    function transferFrom(address from, address to, uint value) external returns (bool);

    function allowance(address owner, address spender) external view returns (uint);

    function approve(address spender, uint value) external returns (bool);

    function increaseApproval(address spender, uint value) external returns (bool);

    function decreaseApproval(address spender, uint value) external returns (bool);

    function approveAndCall(byte symb, address spender, uint value, bytes extraData) external returns (bool success);
}