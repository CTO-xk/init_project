const fs = require("fs");
const path = require("path");
const { ethers } = require("ethers");

// 最小版 ERC20 ABI（仅包含标准事件与常用函数）
const MINIMAL_ERC20_ABI = [
  { "anonymous": false, "inputs": [ { "indexed": true, "internalType": "address", "name": "from", "type": "address" }, { "indexed": true, "internalType": "address", "name": "to", "type": "address" }, { "indexed": false, "internalType": "uint256", "name": "value", "type": "uint256" } ], "name": "Transfer", "type": "event" },
  { "anonymous": false, "inputs": [ { "indexed": true, "internalType": "address", "name": "owner", "type": "address" }, { "indexed": true, "internalType": "address", "name": "spender", "type": "address" }, { "indexed": false, "internalType": "uint256", "name": "value", "type": "uint256" } ], "name": "Approval", "type": "event" },
  { "inputs": [ { "internalType": "address", "name": "owner", "type": "address" } ], "name": "balanceOf", "outputs": [ { "internalType": "uint256", "name": "", "type": "uint256" } ], "stateMutability": "view", "type": "function" },
  { "inputs": [], "name": "decimals", "outputs": [ { "internalType": "uint8", "name": "", "type": "uint8" } ], "stateMutability": "view", "type": "function" },
  { "inputs": [], "name": "symbol", "outputs": [ { "internalType": "string", "name": "", "type": "string" } ], "stateMutability": "view", "type": "function" },
  { "inputs": [ { "internalType": "address", "name": "to", "type": "address" }, { "internalType": "uint256", "name": "value", "type": "uint256" } ], "name": "transfer", "outputs": [ { "internalType": "bool", "name": "", "type": "bool" } ], "stateMutability": "nonpayable", "type": "function" }
];

// -------- Paths --------
const PROJECT_ROOT = path.resolve(__dirname, "..");
const ARTIFACT_PATH = path.join(
  PROJECT_ROOT,
  "artifacts",
  "contracts",
  "core",
  "EnterpriseToken.sol",
  "EnterpriseToken.json"
);

// -------- Internals --------
function loadDeployment(networkName = process.env.HARDHAT_NETWORK || "sepolia") {
  const deploymentFile = path.join(
    PROJECT_ROOT,
    "deployments",
    networkName,
    "EnterpriseToken.json"
  );
  if (!fs.existsSync(deploymentFile)) {
    throw new Error(`未找到部署记录文件: ${deploymentFile}`);
  }
  return JSON.parse(fs.readFileSync(deploymentFile, "utf8"));
}

function loadAbi() {
  if (!fs.existsSync(ARTIFACT_PATH)) {
    throw new Error(`未找到合约构建产物: ${ARTIFACT_PATH}`);
  }
  const artifact = JSON.parse(fs.readFileSync(ARTIFACT_PATH, "utf8"));
  if (!artifact.abi) {
    throw new Error("构建产物缺少 ABI 字段");
  }
  return artifact.abi;
}

function resolveRpcUrl(networkName = "sepolia") {
  // 支持 sepolia 与 base-sepolia，优先 Alchemy
  if (networkName === "sepolia") {
    const alchemyKey = process.env.SEPOLIA_ALCHEMY_AK;
    if (alchemyKey) return `https://eth-sepolia.g.alchemy.com/v2/${alchemyKey}`;
    if (process.env.SEPOLIA_RPC_URL) return process.env.SEPOLIA_RPC_URL;
    if (process.env.INFURA_API_KEY) return `https://sepolia.infura.io/v3/${process.env.INFURA_API_KEY}`;
  }
  if (networkName === "base-sepolia") {
    const baseAk = process.env.BASE_SEPOLIA_ALCHEMY_AK;
    if (baseAk) return `https://base-sepolia.g.alchemy.com/v2/${baseAk}`;
    if (process.env.BASE_SEPOLIA_RPC_URL) return process.env.BASE_SEPOLIA_RPC_URL;
  }
  return undefined;
}

function getProvider(rpcUrl, networkName = process.env.HARDHAT_NETWORK || "sepolia") {
  const finalUrl = rpcUrl || resolveRpcUrl(networkName);
  if (!finalUrl) {
    throw new Error(`缺少 ${networkName} 的 RPC URL，请设置对应的 Alchemy/RPC 环境变量`);
  }
  return new ethers.providers.JsonRpcProvider(finalUrl);
}

function getSigner(privateKey, rpcUrl, networkName = process.env.HARDHAT_NETWORK || "sepolia") {
  if (!privateKey || !/^0x[0-9a-fA-F]{64}$/.test(privateKey)) {
    throw new Error("无效的私钥，请提供形如 0x + 64 hex 的私钥");
  }
  const provider = getProvider(rpcUrl, networkName);
  return new ethers.Wallet(privateKey, provider);
}

function getContract(addressOrNetwork, providerOrSigner) {
  const abi = loadAbi();

  let address = addressOrNetwork;
  if (!address || address.startsWith?.("sepolia") || address === "sepolia" || address === "base-sepolia") {
    const networkName = address === "base-sepolia" ? "base-sepolia" : "sepolia";
    const deployment = loadDeployment(networkName);
    address = deployment.address;
  }

  if (!ethers.utils.isAddress(address)) {
    throw new Error(`无效的合约地址: ${address}`);
  }

  return new ethers.Contract(address, abi, providerOrSigner);
}

function getContractMinimal(addressOrNetwork, providerOrSigner) {
  let address = addressOrNetwork;
  if (!address || address === "sepolia" || address === "base-sepolia") {
    const networkName = address === "base-sepolia" ? "base-sepolia" : "sepolia";
    const deployment = loadDeployment(networkName);
    address = deployment.address;
  }
  if (!ethers.utils.isAddress(address)) {
    throw new Error(`无效的合约地址: ${address}`);
  }
  return new ethers.Contract(address, MINIMAL_ERC20_ABI, providerOrSigner);
}

// -------- Read methods --------
async function getTokenInfo(rpcUrl, networkName = "sepolia") {
  const provider = getProvider(rpcUrl, networkName);
  const contract = getContract(networkName, provider);
  const [name, symbol, decimals, totalSupply] = await Promise.all([
    contract.name(),
    contract.symbol(),
    contract.decimals(),
    contract.totalSupply(),
  ]);
  return { name, symbol, decimals, totalSupply: totalSupply.toString() };
}

async function getBalanceOf(address, rpcUrl, networkName = "sepolia") {
  const provider = getProvider(rpcUrl, networkName);
  const contract = getContract(networkName, provider);
  const balance = await contract.balanceOf(address);
  return balance.toString();
}

async function getOwner(rpcUrl, networkName = "sepolia") {
  const provider = getProvider(rpcUrl, networkName);
  const contract = getContract(networkName, provider);
  return contract.owner();
}

// -------- Write methods (require signer) --------
async function mint(to, amount, options = {}) {
  const { privateKey = process.env.SEPOLIA_PK_ONE, rpcUrl, decimals = 18, networkName = "sepolia" } = options;
  const signer = getSigner(privateKey, rpcUrl, networkName);
  const contract = getContract(networkName, signer);
  const value = ethers.utils.parseUnits(String(amount), decimals);
  const tx = await contract.mint(to, value);
  const receipt = await tx.wait(1);
  return { hash: tx.hash, blockNumber: receipt.blockNumber };
}

async function burn(from, amount, options = {}) {
  const { privateKey = process.env.SEPOLIA_PK_ONE, rpcUrl, decimals = 18, networkName = "sepolia" } = options;
  const signer = getSigner(privateKey, rpcUrl, networkName);
  const contract = getContract(networkName, signer);
  const value = ethers.utils.parseUnits(String(amount), decimals);
  const tx = await contract.burn(from, value);
  const receipt = await tx.wait(1);
  return { hash: tx.hash, blockNumber: receipt.blockNumber };
}

async function transferOwnership(newOwner, options = {}) {
  const { privateKey = process.env.SEPOLIA_PK_ONE, rpcUrl, networkName = "sepolia" } = options;
  const signer = getSigner(privateKey, rpcUrl, networkName);
  const contract = getContract(networkName, signer);
  const tx = await contract.transferOwnership(newOwner);
  const receipt = await tx.wait(1);
  return { hash: tx.hash, blockNumber: receipt.blockNumber };
}

async function acceptOwnership(options = {}) {
  const { privateKey = process.env.SEPOLIA_PK_ONE, rpcUrl, networkName = "sepolia" } = options;
  const signer = getSigner(privateKey, rpcUrl, networkName);
  const contract = getContract(networkName, signer);
  const tx = await contract.acceptOwnership();
  const receipt = await tx.wait(1);
  return { hash: tx.hash, blockNumber: receipt.blockNumber };
}

// -------- Event queries --------
function buildFilters(contract) {
  return {
    Mint: contract.filters.Mint(null, null, null),
    Burn: contract.filters.Burn(null, null, null),
    Transfer: contract.filters.Transfer(null, null, null),
  };
}

async function getPastEvents(eventName, { fromBlock, toBlock, rpcUrl, networkName = "sepolia" } = {}) {
  const provider = getProvider(rpcUrl, networkName);
  const contract = getContract(networkName, provider);
  const latest = await provider.getBlockNumber();
  const start = typeof fromBlock === "number" ? fromBlock : Math.max(latest - 5000, 0);
  const end = typeof toBlock === "number" ? toBlock : latest;
  const filters = buildFilters(contract);
  const filter = filters[eventName];
  if (!filter) throw new Error(`不支持的事件: ${eventName}`);
  const logs = await contract.queryFilter(filter, start, end);
  return logs.map(l => ({
    event: l.event,
    args: l.args ? Array.from(l.args) : [],
    blockNumber: l.blockNumber,
    transactionHash: l.transactionHash,
  }));
}

async function getPastMintEvents(opts) { return getPastEvents("Mint", opts); }
async function getPastBurnEvents(opts) { return getPastEvents("Burn", opts); }
async function getPastTransferEvents(opts) { return getPastEvents("Transfer", opts); }

// -------- Event subscriptions (polling by block) --------
function onEvents({ rpcUrl, networkName = "sepolia", fromBlock }, handlers = {}) {
  const provider = getProvider(rpcUrl, networkName);
  const contract = getContract(networkName, provider);
  const filters = buildFilters(contract);

  let lastSeen = typeof fromBlock === "number" ? fromBlock : undefined;
  let running = true;

  async function handleBlock(blockNumber) {
    if (!running) return;
    const start = lastSeen ? lastSeen + 1 : Math.max(blockNumber - 10, 0);
    const end = blockNumber;

    const toPairs = [
      ["Mint", handlers.onMint, filters.Mint],
      ["Burn", handlers.onBurn, filters.Burn],
      ["Transfer", handlers.onTransfer, filters.Transfer],
    ];

    for (const [name, cb, flt] of toPairs) {
      if (typeof cb !== "function") continue;
      const logs = await contract.queryFilter(flt, start, end);
      for (const l of logs) {
        cb({
          event: l.event,
          args: l.args ? Array.from(l.args) : [],
          blockNumber: l.blockNumber,
          transactionHash: l.transactionHash,
        });
      }
    }

    if (typeof handlers.onAny === "function") {
      const logsAll = await contract.queryFilter({}, start, end);
      for (const l of logsAll) {
        handlers.onAny({
          event: l.event,
          args: l.args ? Array.from(l.args) : [],
          blockNumber: l.blockNumber,
          transactionHash: l.transactionHash,
        });
      }
    }

    lastSeen = end;
  }

  provider.on("block", handleBlock);

  return () => {
    running = false;
    provider.off("block", handleBlock);
  };
}

module.exports = {
  // low-level helpers
  loadDeployment,
  loadAbi,
  getProvider,
  getSigner,
  getContract,
  // read
  getTokenInfo,
  getBalanceOf,
  getOwner,
  // write
  mint,
  burn,
  transferOwnership,
  acceptOwnership,
  // minimal abi
  MINIMAL_ERC20_ABI,
  getContractMinimal,
};


