// SPDX-License-Identifier: MIT 
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable2Step.sol";
import { IEnterpriseToken } from "../interfaces/IEnterpriseToken.sol";
/**
 * @title EnterpriseToken
 * @dev ERC20代币，支持铸造和销毁功能
 * 采用Ownable2Step实现安全的权限管理
 */
 contract EnterpriseToken is ERC20, Ownable2Step, IEnterpriseToken {
    //事件
    event Mint(address indexed to, uint256 amount,uint256 timestamp);
    event Burn(address indexed from, uint256 amount,uint256 timestamp);
   constructor(
   string memory name,
   string memory symbol,
   address initialOwner
   ) ERC20(name, symbol) Ownable(initialOwner) {}
      
     /**
     * @dev 铸造代币（仅管理员）
     * @param to 接收地址
     * @param amount 铸造数量（18位小数）
     */
     function mint(address to, uint256 amount) external onlyOwner {
        //参数校验
        require(to!=address(0), "Invalid address");
        require(amount>0, "Invalid amount");
        _mint(to, amount);
        emit Mint(to, amount, block.timestamp);
     }

     /**
     * @dev 销毁代币（仅管理员）
     * @param from 销毁地址
     * @param amount 销毁数量
     */
     function burn(address from, uint256 amount) external onlyOwner {
        //参数校验
        require(from!=address(0), "Invalid address");
        require(amount>0, "Invalid amount");
        require(balanceOf(from)>=amount, "Insufficient balance");
        _burn(from, amount);
        emit Burn(from, amount, block.timestamp);
     }
 }