// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

contract Voting{

    struct Voter {
        bool isVoted; //是否已经投过票
    }

    struct Candidate {
        uint256 voteCount; //候选人票数
        //是否被投过票
        bool exists;
    }

    //记录票数
    mapping(address => Candidate ) public votesReceived;
    //记录所有候选人地址
    address[] public candidateList;

    //记录已经投过票的人
    mapping (address => Voter) public voters;


    function voteTo(address to) public {
        require(to != address(0), "Cannot vote to zero address");
        require(to != msg.sender, "Self-delegation is disallowed."); 
        require(!voters[msg.sender].isVoted, "Already voted.");
        voters[msg.sender].isVoted = true; 
        votesReceived[to].voteCount++;
        if(!votesReceived[to].exists){
            candidateList.push(to);
            votesReceived[to].exists = true; 
        }
    }

    function getVotes(address to) public view returns(uint256){
        return votesReceived[to].voteCount;
    }

    function resetVotes() public {
         for (uint256 i = 0; i < candidateList.length; i++) {
            votesReceived[candidateList[i]].voteCount = 0;
            voters[msg.sender].isVoted = false; 
        }
    }
}