// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

contract MyERC20 {
    string private _name;
    string private _symbol;
    uint private immutable _decimals = 18;
    uint256 private _totalSupply;
    address private _ower;

    modifier onlyOwer() {
        require(msg.sender == _ower, "Not owner");
        _;
    }

    //记录账户余额
    mapping(address => uint256) private _balances;

    //记录账户授权额度
    mapping(address => mapping(address => uint256)) _allowances;

    constructor(
        string memory name_,
        string memory symbol_,
        uint256 totalSupply_
    ) {
        _name = name_;
        _symbol = symbol_;
        _totalSupply = totalSupply_ * 10 ** _decimals;
        _ower = msg.sender;
        _balances[_ower] = _totalSupply;
        emit Transfer(address(0), _ower, _totalSupply);
    }

    function name() public view returns (string memory) {
        return _name;
    }

    function symbol() public view returns (string memory) {
        return _symbol;
    }

    function decimals() public pure returns (uint256) {
        return _decimals;
    }

    function totalSupply() public view returns (uint256) {
        return _totalSupply;
    }

    event Approval(address indexed from, address indexed to, uint256 value);
    event Transfer(address indexed from, address indexed to, uint256 value);

    function mint(address _to, uint256 _value) public onlyOwer {
        require(_to != address(0), "Invalid address");
        _totalSupply += _value;
        _balances[_to] += _value;

        emit Transfer(address(0), _to, _value);
    }

    function burn(uint256 _value) public onlyOwer {
        require(msg.sender != address(0), "Invalid address");
        _totalSupply -= _value;
        _balances[msg.sender] -= _value;

        emit Transfer(msg.sender, address(0), _value);
    }

    function balanceOf(address from) public view returns (uint256 balance) {
        require(from != address(0), "Invalid address");
        return _balances[from];
    }

    // 用户自己调用
    function transfer(
        address _to,
        uint256 _value
    ) public returns (bool success) {
        require(_to != address(0), "Invalid address");
        require(_balances[msg.sender] >= _value, "Insufficient balance");
        _balances[msg.sender] -= _value;
        _balances[_to] += _value;
        emit Transfer(msg.sender, _to, _value);
        return true;
    }

    function transferFrom(
        address _from,
        address _to,
        uint256 _value
    ) public returns (bool success) {
        require(_from != address(0), "Invalid address");
        require(_to != address(0), "Invalid address");
        require(_balances[_from] >= _value, "Insufficient balance");
        require(
            _allowances[_from][msg.sender] >= _value,
            "Insufficient allowance"
        );
        _balances[_from] -= _value;
        _balances[_to] += _value;
        _allowances[_from][msg.sender] -= _value;
        emit Transfer(_from, _to, _value);
        return true;
    }

    function approve(
        address _spender,
        uint256 _value
    ) public returns (bool success) {
        require(_spender != address(0), "Invalid address");
        // require(balance[_spender] >= _value, "Insufficient balance");
        _allowances[msg.sender][_spender] = _value;
        emit Approval(msg.sender, _spender, _value);
        return true;
    }

    function allowance(
        address _owner,
        address _spender
    ) public view returns (uint256 remaining) {
        return _allowances[_owner][_spender];
    }
}
