pragma solidity ^0.4.18;

contract Token {
    //获取账户_owner拥有token的数量
    function balanceOf(address _owner) constant returns (uint256 balance);

    //从消息发送者账户中往_to账户转数量为_value的token
    function transfer(address _to, uint256 _value) returns (bool success);

    //从账户_from中往账户_to转数量为_value的token，与approve方法配合使用
    function transferFrom(address _from, address _to, uint256 _value) returns (bool success);

    //消息发送账户设置账户_spender能从发送账户中转出数量为_value的token
    function approve(address _spender, uint256 _value) returns (bool success);

    //获取账户_spender可以从账户_owner中转出token的数量
    function allowance(address _owner, address _spender) constant returns (uint256 remaining);

    //发生转账时必须要触发的事件
    event Transfer(address indexed _from, address indexed _to, uint256 _value);

    //当函数approve(address _spender, uint256 _value)成功执行时必须触发的事件
    event Approval(address indexed _owner, address indexed _spender, uint256 _value);

    //token总量，默认会为public变量生成一个getter函数接口，名称为totalSupply().
    uint256 public totalSupply;
}

contract StandardToken is Token {
    function balanceOf(address _owner) constant returns (uint256 balance) {
        return balances[_owner];
    }

    function transfer(address _to, uint256 _value) returns (bool success) {
        //默认totalSupply 不会超过最大值 (2^256 - 1).
        //如果随着时间的推移将会有新的token生成，则可以用下面这句避免溢出的异常
        require(balances[msg.sender] >= _value && balances[_to] + _value > balances[_to]);
        balances[msg.sender] -= _value;
        //从消息发送者账户中减去token数量_value
        balances[_to] += _value;
        //往接收账户增加token数量_value
        Transfer(msg.sender, _to, _value);
        //触发转币交易事件
        return true;
    }

    function transferFrom(address _from, address _to, uint256 _value) returns (bool success) {
        require(balances[_from] >= _value && allowed[_from][msg.sender] >=
        _value && balances[_to] + _value > balances[_to]);
        balances[_to] += _value;
        //接收账户增加token数量_value
        balances[_from] -= _value;
        //支出账户_from减去token数量_value
        allowed[_from][msg.sender] -= _value;
        //消息发送者可以从账户_from中转出的数量减少_value
        Transfer(_from, _to, _value);
        //触发转币交易事件
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
        //允许_spender从_owner中转出的token数
    }

    mapping(address => uint256) balances;
    mapping(address => mapping(address => uint256)) allowed;
}

contract INK is StandardToken {
    /* Public variables of the token */
    string public name;             // token name
    uint8  public decimals;
    string public symbol;           // token symbol

    function(){
        revert();
    }

    function INK() {
        name = 'INK Coin';
        decimals = 9;
        symbol = 'INK';
        totalSupply = 10 * (10 ** 8) * (10 ** uint256(decimals));
        balances[msg.sender] = totalSupply;
    }

    /* Approves and then calls the receiving contract */
    function approveAndCall(address _spender, uint256 _value, bytes _extraData) returns (bool success) {
        allowed[msg.sender][_spender] = _value;
        Approval(msg.sender, _spender, _value);
        //发送者通知代币合约：1STORE币授权给了服务合约（通过调用代币合约的 approveAndCall()函数）
        //代币合约通知服务合约：1STORE币已经授权给了服务合约（通过调用服务合约的 receiveApproval()函数）
        //服务合约指示代币合约将代币从发送者的账户转移到服务合约的账户（通过调用服务合约的 transferFrom()函数 并且存储信息)
        require(_spender.call(bytes4(bytes32(sha3("receiveApproval(address,uint256,address,bytes)"))), msg.sender, _value, this, _extraData));
        return true;
    }
}