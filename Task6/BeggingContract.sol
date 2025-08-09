// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;
/**
 * @title BeggingContract
 * @dev 一个简单的讨饭合约，支持用户捐赠以太币、记录捐赠信息，以及所有者提取资金
 */
contract BeggingContract{
    // 记录每个地址的捐赠总金额
    mapping (address => uint256) private _donations;
    // 合约所有者地址
    address private  immutable _owner;

    // 仅所有者可调用的修饰符
    modifier onlyOwner(){
        require(_owner == msg.sender,"caller is not the owner");
        _;
    }
    /**
     * @dev 构造函数，设置合约部署者为所有者
     */
    constructor(){
        _owner = msg.sender;
    }
    /**
     * @dev 捐赠函数，允许用户向合约发送以太币
     * 自动记录捐赠者地址和金额
     */
     function donate() public  payable{
        // 确保捐赠金额大于0
        require(msg.value > 0,"donation amount must be greater than 0");
        // 累加捐赠金额（支持同一地址多次捐赠）
        _donations[msg.sender] += msg.value;
     }
    /**
     * @dev 提款函数，允许所有者提取合约中的所有资金
     * 仅合约所有者可调用
     */
     function withdraw() public  onlyOwner{
        // 确保合约中有可提取的资金
        require(address(this). balance > 0,"no funds to withdraw");
         // 将合约中的所有余额转移给所有者
         uint256 balance =address(this).balance;
         payable(_owner).transfer(balance);
     }
    /**
     * @dev 查询指定地址的捐赠总金额
     * @param donor 捐赠者地址
     * @return 该地址的累计捐赠金额（以wei为单位）
     */
    function getDonation(address donor) public  view returns (uint256){
        return  _donations[donor];
    }
    /**
     * @dev 获取合约当前的总余额
     * @return 合约中以太币的总余额（以wei为单位）
     */
     function getContractBalance() public  view  returns (uint256){
        return address(this).balance;
     }
    /**
     * @dev 获取合约所有者地址
     * @return 所有者地址
     */
     function owner() public  view  returns  (address){
        return _owner;
     }
}