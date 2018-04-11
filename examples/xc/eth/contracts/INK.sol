pragma solidity ^0.4.19;

contract Token {

    string public name;

    string public symbol;

    uint8 public decimals;

    uint public totalSupply;

    mapping(address => uint) private balances;

    mapping(address => mapping(address => uint)) private allowed;

    event Transfer(address indexed from, address indexed to, uint value);

    event Approval(address indexed owner, address indexed spender, uint value);

    function Token() public {

        name = "INK Coin";

        symbol = "INK";

        decimals = 9;

        totalSupply = 10 * (10 ** 8) * (10 ** uint(decimals));

        balances[msg.sender] = totalSupply;
    }

    function balanceOf(address owner) public constant returns (uint) {

        return balances[owner];
    }

    function transfer(address to, uint value) public returns (bool) {

        require(to != address(0));

        require(value <= balances[msg.sender]);

        balances[msg.sender] -= value;

        balances[to] += value;

        emit Transfer(msg.sender, to, value);

        return true;
    }

    function transferFrom(address from, address to, uint value) public returns (bool) {

        require(to != address(0));

        require(value <= balances[from]);

        require(value <= allowed[from][msg.sender]);

        balances[from] -= value;

        balances[to] += value;

        allowed[from][msg.sender] -= value;

        emit Transfer(from, to, value);

        return true;
    }

    function allowance(address owner, address spender) public constant returns (uint) {

        return allowed[owner][spender];
    }

    function approve(address spender, uint value) public returns (bool) {

        allowed[msg.sender][spender] = value;

        emit Approval(msg.sender, spender, value);

        return true;
    }

    function increaseApproval(address spender, uint value) public returns (bool) {

        allowed[msg.sender][spender] += value;

        emit Approval(msg.sender, spender, allowed[msg.sender][spender]);

        return true;
    }

    function decreaseApproval(address spender, uint value) public returns (bool) {

        if (value > allowed[msg.sender][spender]) {

            allowed[msg.sender][spender] = 0;
        } else {

            allowed[msg.sender][spender] -= value;
        }

        emit Approval(msg.sender, spender, allowed[msg.sender][spender]);

        return true;
    }

    function approveAndCall(byte symb, address spender, uint value, bytes extraData) returns (bool success) {

        require(spender != address(0));

        require(symb == '-' || symb == '' || symb == '+');

        bool status;

        if (symb == '-') {

            status = decreaseApproval(spender, value);
        }

        if (symb == '') {

            status = approve(spender, value);
        }

        if (symb == '+') {

            status = increaseApproval(spender, value);
        }

        if (status) {

            if (!spender.call(bytes4(bytes32(sha3("receiveApproval(address,uint,address,bytes)"))), msg.sender, value, this, extraData)) {

                revert();
            }

            return true;
        }
    }
}