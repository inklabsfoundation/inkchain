pragma solidity ^0.4.19;

contract Token {

    function transfer(address to, uint value) external;

    function transferFrom(address from, address to, uint value) external returns (bool);

    function balanceOf(address owner) external view returns (uint);

    function allowance(address owner, address spender) external view returns (uint);
}