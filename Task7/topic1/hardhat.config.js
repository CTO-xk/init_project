require("@nomiclabs/hardhat-waffle");
require("@nomiclabs/hardhat-etherscan");
require("hardhat-gas-reporter");
require("solidity-coverage");

// 加载环境变量（如果使用）
require("dotenv").config();

// 私钥，仅在测试时使用，
const PRIVATE_KEY = process.env.PRIVATE_KEY || "0000000000000000000000000000000000000000000000000000000000000000";

module.exports = {
    solidity: "0.8.17",
    settings: {
        optimizer: {
            enabled: true,
            runs: 200
        }
    },
    //本地开发网络
    networks: {
      localhost: {
        url: "http://127.0.0.1:8545"
      }
    },
    //测试网
    sepolia: {
        url: '',
        accounts: [PRIVATE_KEY]
    },
    // Etherscan配置，用于合约验证
  etherscan: {
    apiKey: process.env.ETHERSCAN_API_KEY
  },
  // Gas报告配置
  gasReporter: {
    enabled: process.env.REPORT_GAS !== undefined,
    currency: "USD",
    outputFile: "gas-report.txt",
    noColors: true,
    coinmarketcap: process.env.COINMARKETCAP_API_KEY,
  },
    // 路径配置
  paths: {
    sources: "./contracts",
    tests: "./test",
    cache: "./cache",
    artifacts: "./artifacts"
  },
  // mocha测试框架配置
  mocha: {
    timeout: 40000
  }
};
