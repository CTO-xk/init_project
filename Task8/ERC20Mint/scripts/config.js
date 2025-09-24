/**
 * 部署配置：集中管理所有环境的部署参数（优先读取 .env）
 */
const tokenConfig = {
  name: process.env.TOKEN_NAME || "EnterpriseToken",
  symbol: process.env.TOKEN_SYMBOL || "ET",
  initialOwner: process.env.INITIAL_OWNER || "0x0000000000000000000000000000000000000000",
};

const networkConfig = {
  localhost: {
    name: "localhost",
    chainId: 31337,
    url: process.env.LOCALHOST_RPC_URL || "http://127.0.0.1:8545",
  },
  mainnet: {
    name: "mainnet",
    chainId: 1,
    url:
      process.env.MAINNET_RPC_URL ||
      (process.env.MAINNET_ALCHEMY_AK
        ? `https://eth-mainnet.g.alchemy.com/v2/${process.env.MAINNET_ALCHEMY_AK}`
        : ""),
  },
  sepolia: {
    name: "sepolia",
    chainId: 11155111,
    url:
      process.env.SEPOLIA_RPC_URL ||
      (process.env.SEPOLIA_ALCHEMY_AK
        ? `https://eth-sepolia.g.alchemy.com/v2/${process.env.SEPOLIA_ALCHEMY_AK}`
        : ""),
  },
};

module.exports = { tokenConfig, networkConfig };