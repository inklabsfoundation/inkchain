pragma solidity ^0.4.19;

import "./TokenInterface.sol";
import "./SafeMath.sol";

contract Token is TokenInterface {

    string public name;

    string public symbol;

    uint8 public decimals;

    uint public totalSupply;

    bool private status;

    address private admin;

    mapping(address => uint) private balances;

    mapping(address => mapping(address => uint)) private allowed;

    event Transfer(address indexed from, address indexed to, uint value);

    event Approval(address indexed owner, address indexed spender, uint value);

    function Token() public {
        //TODO
        name = "INK Coin";
        //TODO
        symbol = "INK";
        //TODO
        decimals = 9;
        //TODO
        uint initSupply = 10 * (10 ** 8);

        totalSupply = SafeMath.mul(initSupply, (10 ** uint(decimals)));

        balances[msg.sender] = totalSupply;

        admin = msg.sender;

        status = true;
    }

    function setStatus(bool _status) external {

        require(msg.sender == admin);

        if (status != _status) {

            status = _status;
        }
    }

    function getStatus() external view returns (bool) {

        return status;
    }

    function setAdmin(address account) external {

        require(msg.sender == admin);

        if (admin != account) {

            admin = account;
        }
    }

    function getAdmin() external view returns (address) {

        require(msg.sender == admin);

        return admin;
    }

    function balanceOf(address owner) external view returns (uint) {

        return balances[owner];
    }

    function transfer(address to, uint value) external returns (bool) {

        require(status);

        require(to != address(0));

        require(totalSupply >= value && value > 0);

        require(balances[msg.sender] >= value);

        balances[msg.sender] = SafeMath.sub(balances[msg.sender], value);

        balances[to] = SafeMath.add(balances[to], value);

        emit Transfer(msg.sender, to, value);

        return true;
    }

    function transferFrom(address from, address to, uint value) external returns (bool) {

        require(status);

        require(to != address(0));

        require(totalSupply >= value && value > 0);

        require(balances[from] >= value);

        require(allowed[from][msg.sender] >= value);

        balances[from] = SafeMath.sub(balances[from], value);

        balances[to] = SafeMath.add(balances[to], value);

        allowed[from][msg.sender] = SafeMath.sub(allowed[from][msg.sender], value);

        emit Transfer(from, to, value);

        return true;
    }

    function allowance(address owner, address spender) external view returns (uint) {

        return allowed[owner][spender];
    }

    function approve(address spender, uint value) external returns (bool) {

        return _approve(spender, value);
    }

    function _approve(address spender, uint value) internal returns (bool) {

        allowed[msg.sender][spender] = value;

        emit Approval(msg.sender, spender, value);

        return true;
    }

    function increaseApproval(address spender, uint value) external returns (bool) {

        return _increaseApproval(spender, value);
    }

    function _increaseApproval(address spender, uint value) internal returns (bool) {

        allowed[msg.sender][spender] = SafeMath.add(allowed[msg.sender][spender], value);

        emit Approval(msg.sender, spender, allowed[msg.sender][spender]);

        return true;
    }

    function decreaseApproval(address spender, uint value) external returns (bool) {

        return _decreaseApproval(spender, value);
    }

    function _decreaseApproval(address spender, uint value) internal returns (bool) {

        if (value > allowed[msg.sender][spender]) {

            allowed[msg.sender][spender] = 0;
        } else {

            allowed[msg.sender][spender] = SafeMath.sub(allowed[msg.sender][spender], value);
        }

        emit Approval(msg.sender, spender, allowed[msg.sender][spender]);

        return true;
    }

    function approveAndCall(byte _symbol, address spender, uint value, bytes extraData) external returns (bool success) {

        require(spender != address(0));

        require(_symbol == '-' || _symbol == '' || _symbol == '+');

        bool _status;

        if (_symbol == '-') {

            _status = _decreaseApproval(spender, value);
        }

        if (_symbol == '') {

            _status = _approve(spender, value);
        }

        if (_symbol == '+') {

            _status = _increaseApproval(spender, value);
        }

        if (_status) {

            if (!spender.call(bytes4(bytes32(sha3("receiveApproval(address,uint,address,bytes)"))), msg.sender, value, this, extraData)) {

                revert();
            }

            return true;
        }
    }
}