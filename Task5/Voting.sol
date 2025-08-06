// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract Voting{
    //储存候选人的票数
    mapping (string =>uint) private  votes;
    //记录所有出现候选过的候选人
    string [] candidates;

    //判断候选人是否被记录
    mapping (string =>bool) private isCandidateRegistered;

    //vote函数，允许用户投票给某个候选人
    function vote (string memory candidate) public {
        if(!isCandidateRegistered[candidate]){
           candidates.push(candidate); 
           isCandidateRegistered[candidate] = true;
        }
        votes[candidate]++;
    }
    //getVotes函数，返回某个候选人的得票数
    function getVotes (string memory candidate)public  view returns (uint){
        return votes[candidate];
    }
    //resetVotes函数，重置所有候选人的得票数
    function resetVotes()public {
        for (uint i=0;i<candidates.length;i++){
            votes[candidates[i]] = 0;
        }
    }
}