// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MetaNodeToken is ERC20 ,Ownable {
    constructor(
        string memory name,
        string memory symbol,
        uint256 initialSupply,
        address initialOwner
    )  ERC20(name, symbol) Ownable() {
        _mint(initialOwner, initialSupply);
    }
      // 允许合约所有者铸造新代币
      function mint(address to, uint256 amount) external onlyOwner {
        _mint(to, amount);
      }
      // 允许合约所有者销毁代币
    function burn(address from, uint256 amount) external onlyOwner {
        _burn(from, amount);
    }
}   