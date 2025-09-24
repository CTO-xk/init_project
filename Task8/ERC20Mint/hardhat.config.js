const { config: dotenvConfig } = require("dotenv");
dotenvConfig();
require("@nomiclabs/hardhat-ethers");
require("@nomiclabs/hardhat-etherscan");

const SEPOLIA_RPC_URL = process.env.SEPOLIA_RPC_URL || (process.env.INFURA_API_KEY ? `https://sepolia.infura.io/v3/${process.env.INFURA_API_KEY}` : "");
const SEPOLIA_PRIVATE_KEY = process.env.SEPOLIA_PRIVATE_KEY || process.env.PRIVATE_KEY || "";

// Alchemy 相关环境变量与多账户支持
const SEPOLIA_ALCHEMY_AK = process.env.SEPOLIA_ALCHEMY_AK || "";
const SEPOLIA_PK_ONE = process.env.SEPOLIA_PK_ONE || "";
const SEPOLIA_PK_TWO = process.env.SEPOLIA_PK_TWO || "";

// 简单的私钥格式校验（0x + 64 hex）
const isHexPriv = (v) => /^0x[0-9a-fA-F]{64}$/.test(v || "");

// Base Sepolia 多链支持
const BASE_SEPOLIA_ALCHEMY_AK = process.env.BASE_SEPOLIA_ALCHEMY_AK || "";
const BASE_SEPOLIA_PK_ONE = process.env.BASE_SEPOLIA_PK_ONE || "";
const BASE_SEPOLIA_PK_TWO = process.env.BASE_SEPOLIA_PK_TWO || "";

/** @type {import('hardhat/config').HardhatUserConfig} */
const config = {
  paths: {
    sources: "./contracts",
    tests: "./test",
    cache: "./cache",
    artifacts: "./artifacts",
  },
  solidity: {
    version: "0.8.28",
    settings: {
      optimizer: { enabled: true, runs: 200 },
    },
  },
  networks: (() => {
    const networks = {
      hardhat: { chainId: 31337 },
      localhost: { url: "http://127.0.0.1:8545", chainId: 31337 },
    };

    // 优先使用 Alchemy，其次使用显式的 SEPOLIA_RPC_URL / INFURA
    let sepoliaUrl = "";
    let sepoliaAccounts = [];

    if (SEPOLIA_ALCHEMY_AK) {
      sepoliaUrl = `https://eth-sepolia.g.alchemy.com/v2/${SEPOLIA_ALCHEMY_AK}`;
      sepoliaAccounts = [SEPOLIA_PK_ONE, SEPOLIA_PK_TWO].filter(isHexPriv);
    } else if (SEPOLIA_RPC_URL) {
      sepoliaUrl = SEPOLIA_RPC_URL;
      sepoliaAccounts = SEPOLIA_PRIVATE_KEY ? [SEPOLIA_PRIVATE_KEY] : [];
    }

    if (sepoliaUrl) {
      networks.sepolia = {
        url: sepoliaUrl,
        accounts: sepoliaAccounts,
        chainId: 11155111,
      };
    }

    // Base Sepolia
    if (BASE_SEPOLIA_ALCHEMY_AK) {
      networks["base-sepolia"] = {
        url: `https://base-sepolia.g.alchemy.com/v2/${BASE_SEPOLIA_ALCHEMY_AK}`,
        accounts: [BASE_SEPOLIA_PK_ONE, BASE_SEPOLIA_PK_TWO].filter(isHexPriv),
        chainId: 84532,
      };
    }

    return networks;
  })(),
  etherscan: {
    apiKey: {
      sepolia: process.env.ETHERSCAN_API_KEY || "",
      "base-sepolia": process.env.BASESCAN_API_KEY || "",
    },
    customChains: [
      {
        network: "sepolia",
        chainId: 11155111,
        urls: {
          apiURL: "https://api-sepolia.etherscan.io/api",
          browserURL: "https://sepolia.etherscan.io",
        },
      },
      {
        network: "base-sepolia",
        chainId: 84532,
        urls: {
          apiURL: "https://api-sepolia.basescan.org/api",
          browserURL: "https://sepolia.basescan.org",
        },
      },
    ],
  },
};

module.exports = config;
