// SPDX-License-Identifier: MIT 
pragma solidity ^0.8.28;

interface IEnterpriseToken {
    /**
     * @dev 铸造代币
     * @param to 接收地址
     * @param amount 铸造数量
     */
     function mint(address to, uint256 amount   ) external;

     /**
     * @dev 销毁代币
     * @param from 销毁地址
     * @param amount 销毁数量
     */
     function burn(address from, uint256 amount) external;
}   
