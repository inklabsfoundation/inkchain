pragma solidity ^0.4.13;

contract Token {
    uint256 public totalSupply;

    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    function _transfer(address _from, address _to, uint _value) internal {
        require(balanceOf[_from] > _value);
        require(balanceOf[_to] + _value > balanceOf[_to]);
        balanceOf[_from] -= _value;
        balanceOf[_to] += _value;
        Transfer(_from, _to, _value);
    }

    function transfer(address _to, uint256 _value) {
        _transfer(msg.sender, _to, _value);
    }

    function transferFrom(address _from, address _to, uint256 _value) returns (bool success) {
        require(_value < allowance[_from][msg.sender]);
        allowance[_from][msg.sender] -= _value;
        _transfer(_from, _to, _value);
        return true;
    }

    function approve(address _spender, uint256 _value)
    returns (bool success) {
        allowance[msg.sender][_spender] = _value;
        Approval(msg.sender, _spender, _value);
        return true;
    }
}

contract INK is Token {
    uint8 public constant decimals = 9;
    uint256 public constant initialSupply = 10 * (10 ** 8) * (10 ** uint256(decimals));
    string public constant name = 'INK Coin';
    string public constant symbol = 'INK';


    function(){
        revert();
    }

    function INK() {
        balanceOf[msg.sender] = initialSupply;
        totalSupply = initialSupply;
    }

    function approveAndCall(address _spender, uint256 _value, bytes _extraData)
    returns (bool success) {
        if (approve(_spender, _value)) {
            if (!_spender.call(bytes4(bytes32(sha3("receiveApproval(address,uint256,address,bytes)"))), msg.sender, _value, this, _extraData)) {
                revert();
            }
            return true;
        }
    }

}
