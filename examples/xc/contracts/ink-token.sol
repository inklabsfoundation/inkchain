pragma solidity ^0.4.18;

contract Token {
    //Gets the number of tokens that the account _owner owns.
    function balanceOf(address _owner) constant returns (uint256 balance);

    //The token from the message sender's account to the to account number is value.
    function transfer(address _to, uint256 _value) returns (bool success);

    //The message sending account setting account spender can transfer the amount of value from the sending account.
    function transferFrom(address _from, address _to, uint256 _value) returns (bool success);

    //The message sending account setting account spender can transfer the amount of value from the sending account.
    function approve(address _spender, uint256 _value) returns (bool success);

    //The access account spender can transfer the number of tokens from the account owner.
    function allowance(address _owner, address _spender) constant returns (uint256 remaining);

    //Events that must be triggered when a transfer occurs.
    event Transfer(address indexed _from, address indexed _to, uint256 _value);

    //The event that must be triggered when the function approve(address spender, uint256 value) executes successfully.
    event Approval(address indexed _owner, address indexed _spender, uint256 _value);

    //The total token amount, by default, generates a getter function interface for the public variable, named totalSupply().
    uint256 public totalSupply;
}

contract StandardToken is Token {

    mapping(address => uint256) balances;

    mapping(address => mapping(address => uint256)) allowed;

    function balanceOf(address _owner) constant returns (uint256 balance) {
        return balances[_owner];
    }

    function transfer(address _to, uint256 _value) returns (bool success) {

        require(balances[msg.sender] >= _value && balances[_to] + _value > balances[_to]);

        balances[msg.sender] -= _value;

        balances[_to] += _value;

        Transfer(msg.sender, _to, _value);

        return true;
    }

    function transferFrom(address _from, address _to, uint256 _value) returns (bool success) {

        require(balances[_from] >= _value && allowed[_from][msg.sender] >= _value && balances[_to] + _value > balances[_to]);

        balances[_to] += _value;

        balances[_from] -= _value;

        allowed[_from][msg.sender] -= _value;

        Transfer(_from, _to, _value);

        return true;
    }

    function approve(address _spender, uint256 _value) returns (bool success) {

        allowed[msg.sender][_spender] = _value;

        Approval(msg.sender, _spender, _value);

        return true;
    }

    function increaseApproval(address _spender, uint _addedValue) public returns (bool) {

        require(_addedValue > 0);

        allowed[msg.sender][_spender] = allowed[msg.sender][_spender] + _addedValue;

        Approval(msg.sender, _spender, allowed[msg.sender][_spender]);

        return true;
    }

    function decreaseApproval(address _spender, uint _subtractedValue) public returns (bool) {

        require(_subtractedValue > 0);

        uint oldValue = allowed[msg.sender][_spender];

        if (_subtractedValue >= oldValue) {
            allowed[msg.sender][_spender] = 0;
        } else {
            allowed[msg.sender][_spender] = oldValue - _subtractedValue;
        }

        Approval(msg.sender, _spender, allowed[msg.sender][_spender]);

        return true;
    }

    function allowance(address _owner, address _spender) constant returns (uint256 remaining) {
        return allowed[_owner][_spender];
    }
}

contract INK is StandardToken {

    // token name
    string public name;

    uint8  public decimals;

    // token symbol
    string public symbol;

    function(){
        revert();
    }

    function INK() {
        name = 'INK Coin';
        symbol = 'INK';
        decimals = 9;
        totalSupply = 10 * (10 ** 8) * (10 ** uint256(decimals));
        balances[msg.sender] = totalSupply;
    }

    /* Approves and then calls the receiving contract */
    function approveAndCall(address _spender, uint256 _value, bytes _extraData) returns (bool success) {
        allowed[msg.sender][_spender] = _value;
        Approval(msg.sender, _spender, _value);
        require(_spender.call(bytes4(bytes32(sha3("receiveApproval(address,uint256,address,bytes)"))), msg.sender, _value, this, _extraData));
        return true;
    }
}