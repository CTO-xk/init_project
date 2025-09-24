// hardhat.config.js
require("@nomicfoundation/hardhat-toolbox");
require("@openzeppelin/hardhat-upgrades");
require("dotenv").config(); // 加载.env文件中的环境变量

// 从.env中读取关键信息
const INFURA_API_KEY = process.env.INFURA_API_KEY;
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const ETHERSCAN_API_KEY = process.env.ETHERSCAN_API_KEY;

module.exports = {
  solidity: {
    version: "0.8.20", // 与你的合约编译版本一致（如之前的StakeSystem合约用0.8.20）
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
  networks: {
    // 本地测试网（Hardhat内置）
    hardhat: {},
    // 配置Sepolia测试网（核心：使用你的Infura API地址）
    sepolia: {
      url: `https://sepolia.infura.io/v3/${INFURA_API_KEY}`, // 拼接成完整的Sepolia节点URL
      accounts: [PRIVATE_KEY] || [], // 部署合约的钱包私钥（从.env读取）
      gas: 2100000, // 可选：设置Gas限制
      gasPrice: 8000000000, // 可选：设置Gas价格（单位：wei）
    },
  },
  // 配置Etherscan（可选，用于合约验证）
  etherscan: {
    apiKey: ETHERSCAN_API_KEY,
  },
  // 配置合约编译后的输出目录（可选）
  paths: {
    sources: "./contracts",
    tests: "./test",
    cache: "./cache",
    artifacts: "./artifacts",
  },
};