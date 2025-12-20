// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

contract BeggingContract {
    mapping(address => uint256) private _balances;
    address private _owner;
    uint256 private _startTime;
    uint256 private _endTime;

    constructor(uint256 startTime, uint256 endTime) {
        _owner = msg.sender;
        _startTime = startTime;
        _endTime = endTime;
    }

    modifier onlyOwner() {
        require(msg.sender == _owner, "Only owner can call");
        _;
    }

    modifier mustInDonationDuration() {
        require(
            block.timestamp >= _startTime,
            "Donation period has not started yet."
        );
        require(block.timestamp <= _endTime, "Donation period has ended.");
        _;
    }

    event Donation(address indexed donator, uint256 amount);

    function donate() public payable mustInDonationDuration {
        uint256 amount = msg.value;
        require(amount > 0, "You must send some Ether to donate");
        _balances[msg.sender] += amount;
        emit Donation(msg.sender, amount);
    }

    function withdraw() external onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "Insufficient balance");
        payable(_owner).transfer(balance);

        // 使用 .send() 发送 ETH
        //bool success = payable(msg.sender).send(balance);
        //require(success, "Failed to send Ether");

        //使用call
        //(bool success, ) = msg.sender.call{value: balance}("");
        //require(success, "Failed to send Ether");
    }

    function getDonation(address donator) public view returns (uint256) {
        return _balances[donator];
    }
}
