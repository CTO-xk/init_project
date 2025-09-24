const { config: dotenvConfig } = require("dotenv");
dotenvConfig();
const hre = require("hardhat");
const fs = require("fs");
const path = require("path");
const { ethers, network, run } = hre;

// 配置文件路径
const CONFIG_PATH = path.join(__dirname, "config.json");

// 初始化日志工具
const logger = {
  info: (msg) => console.log(`[INFO] ${new Date().toISOString()} ${msg}`),
  warn: (msg) => console.log(`[WARN] ${new Date().toISOString()} ${msg}`),
  error: (msg) => console.error(`[ERROR] ${new Date().toISOString()} ${msg}`),
  success: (msg) => console.log(`[SUCCESS] ${new Date().toISOString()} ${msg}`)
};

// 加载配置文件
function loadConfig() {
  try {
    if (fs.existsSync(CONFIG_PATH)) {
      const data = fs.readFileSync(CONFIG_PATH, "utf8");
      return JSON.parse(data);
    }
    return {
      tokenConfig: {
        name: "Enterprise Token",
        symbol: "ENT",
        initialOwner: ""
      },
      networkConfig: {
        sepolia: {
          confirmations: 6,
          verify: true
        },
        "base-sepolia": {
          confirmations: 6,
          verify: true
        },
        localhost: {
          confirmations: 1,
          verify: false
        },
        hardhat: {
          confirmations: 1,
          verify: false
        }
      }
    };
  } catch (error) {
    logger.warn(`加载配置文件失败，使用默认配置: ${error.message}`);
    return { tokenConfig: {}, networkConfig: {} };
  }
}

// 主部署函数
async function main() {
  // 加载配置
  const { tokenConfig, networkConfig } = loadConfig();
  
  // 确定当前网络
  const netName = network.name || (process.env.NETWORK_NAME || "sepolia");
  logger.info(`===== 开始在 ${netName} 网络部署合约 =====`);

  // 部署记录存储路径
  const DEPLOYMENTS_DIR = path.join(process.cwd(), "deployments", netName);
  const DEPLOYMENT_FILE = path.join(DEPLOYMENTS_DIR, "EnterpriseToken.json");
  
  // 检查是否已部署
  const forceRedeploy = process.env.FORCE_DEPLOY === "1" || process.argv.includes("--force");
  const reuseDeployed = process.env.REUSE === "1" || process.argv.includes("--reuse");
  let existingData = null;
  let isReuseMode = false;
  if (fs.existsSync(DEPLOYMENT_FILE)) {
    existingData = JSON.parse(fs.readFileSync(DEPLOYMENT_FILE, "utf8"));
    if (!forceRedeploy) {
      if (reuseDeployed) {
        isReuseMode = true;
        logger.info(`复用已部署合约: ${existingData.address}`);
      } else {
        logger.warn(`当前网络已存在部署记录: ${existingData.address}`);
        logger.warn("如需复用已部署合约并继续后置操作，请追加 --reuse（或设置 REUSE=1）。如需重新部署，请使用 FORCE_DEPLOY=1 或 --force。");
        process.exit(0);
      }
    } else {
      // 备份旧记录文件
      fs.mkdirSync(DEPLOYMENTS_DIR, { recursive: true });
      const backupPath = path.join(
        DEPLOYMENTS_DIR,
        `EnterpriseToken.json.bak-${new Date().toISOString().replace(/[:]/g, "-")}`
      );
      fs.renameSync(DEPLOYMENT_FILE, backupPath);
      logger.warn(`检测到已部署记录，已备份到: ${backupPath}`);
    }
  }

  // 检查网络连接
  try {
    const networkInfo = await ethers.provider.getNetwork();
    logger.info(`网络信息: ${networkInfo.name} (Chain ID: ${networkInfo.chainId})`);
    
    // 检查网络连接
    const blockNumber = await ethers.provider.getBlockNumber();
    logger.info(`当前区块号: ${blockNumber}`);
  } catch (error) {
    logger.error(`网络连接失败: ${error.message}`);
    process.exit(1);
  }

  // 获取部署账号
  const [deployer] = await ethers.getSigners();
  logger.info(`部署账号: ${deployer.address}`);
  logger.info(`部署账号余额: ${ethers.utils.formatEther(await deployer.getBalance())} ETH`);

  // 准备部署参数
  const name = process.env.TOKEN_NAME || tokenConfig.name || "Enterprise Token";
  const symbol = process.env.TOKEN_SYMBOL || tokenConfig.symbol || "ENT";

  // 处理管理员地址
  const ownerEnv = (process.env.INITIAL_OWNER || "").trim();
  const cfgOwner = (tokenConfig.initialOwner || "").trim();
  
  // 地址验证工具函数
  const isValidOwner = (addr) => 
    typeof addr === "string" && 
    ethers.utils.isAddress(addr) && 
    addr !== "0x0000000000000000000000000000000000000000";
  
  const isPlaceholder = (addr) => 
    !addr || /^0x?你的/i.test(addr) || 
    addr === "0x0000000000000000000000000000000000000000";

  // 确定最终管理员地址
  let initialOwner;
  if (!isPlaceholder(ownerEnv) && isValidOwner(ownerEnv)) {
    initialOwner = ownerEnv;
  } else if (isValidOwner(cfgOwner)) {
    initialOwner = cfgOwner;
  } else {
    initialOwner = deployer.address;
    logger.warn("未指定有效管理员地址，使用部署者地址作为管理员");
  }
  logger.info(`合约管理员: ${initialOwner}`);

  // 获取网络配置
  const netConfig = networkConfig[netName] || networkConfig.localhost || {
    confirmations: 6,
    verify: true
  };
  logger.info(`区块确认数: ${netConfig.confirmations}`);

  // 部署或复用合约
  try {
    let token;
    let receipt;
    let contractAddress;

    if (isReuseMode && existingData && existingData.address) {
      logger.info(`跳过部署，直接复用合约地址: ${existingData.address}`);
      token = await ethers.getContractAt("EnterpriseToken", existingData.address);
      contractAddress = existingData.address;
      // 获取当前区块号作为占位信息
      const currentBlock = await ethers.provider.getBlockNumber();
      receipt = { blockNumber: currentBlock };
    } else {
      logger.info(`开始部署合约: name=${name}, symbol=${symbol}`);
      const EnterpriseToken = await ethers.getContractFactory("EnterpriseToken");
      token = await EnterpriseToken.deploy(name, symbol, initialOwner);
      
      // 等待区块确认
      logger.info(`等待 ${netConfig.confirmations} 个区块确认...`);
      
      // 添加重试机制和超时处理
      const maxRetries = 3;
      let retryCount = 0;
      
      while (retryCount < maxRetries) {
        try {
          receipt = await token.deployTransaction.wait(netConfig.confirmations);
          break;
        } catch (error) {
          retryCount++;
          if (error.message.includes('timeout') || error.message.includes('network')) {
            logger.warn(`等待确认超时，重试 ${retryCount}/${maxRetries}: ${error.message}`);
            if (retryCount < maxRetries) {
              await new Promise(resolve => setTimeout(resolve, 5000));
              continue;
            } else {
              logger.warn(`使用更少的确认数重试...`);
              try {
                receipt = await token.deployTransaction.wait(1);
                logger.warn(`使用1个确认数成功，但建议手动验证交易`);
                break;
              } catch (fallbackError) {
                throw error;
              }
            }
          }
          throw error;
        }
      }
      
      contractAddress = token.address;
    }

    logger.success(`合约部署成功！地址: ${contractAddress}`);
    if (token.deployTransaction) {
      logger.info(`部署交易哈希: ${token.deployTransaction.hash}`);
      logger.info(`部署区块号: ${receipt.blockNumber}`);
    } else {
      logger.info(`复用模式：无部署交易。当前区块号: ${receipt.blockNumber}`);
    }

    // 保存或更新部署记录（复用模式也更新时间与最新区块）
    fs.mkdirSync(DEPLOYMENTS_DIR, { recursive: true });
    const deploymentData = {
      network: netName,
      contractName: "EnterpriseToken",
      address: contractAddress,
      deployer: deployer.address,
      owner: initialOwner,
      name,
      symbol,
      deploymentTime: new Date().toISOString(),
      transactionHash: token.deployTransaction ? token.deployTransaction.hash : (existingData?.transactionHash || ""),
      blockNumber: receipt.blockNumber,
      confirmations: netConfig.confirmations
    };
    fs.writeFileSync(DEPLOYMENT_FILE, JSON.stringify(deploymentData, null, 2));
    logger.success(`${isReuseMode ? "部署记录已更新" : "部署记录已保存到"}: ${DEPLOYMENT_FILE}`);

    // 部署后操作：mint / transfer / burn 以产生活动事件
    try {
      const tokenDecimals = 18;
      const amountMint = ethers.utils.parseUnits("100", tokenDecimals); // 铸造 100
      const amountTransfer = ethers.utils.parseUnits("25", tokenDecimals); // 转移 25
      const amountBurn = ethers.utils.parseUnits("10", tokenDecimals); // 销毁 10

      const allSigners = await ethers.getSigners();
      const ownerSigner = deployer;
      const receiverSigner = allSigners[1] || deployer; // 若无第二个账户则回退

      const ownerAddr = await ownerSigner.getAddress();
      const receiverAddr = await receiverSigner.getAddress();

      const tokenWithOwner = token.connect(ownerSigner);

      // 1) mint 给 owner
      logger.info(`开始铸造: 向 ${ownerAddr} 铸造 100 代币`);
      const txMint = await tokenWithOwner.mint(ownerAddr, amountMint);
      const rcMint = await txMint.wait(1);
      logger.success(`铸造完成: tx=${txMint.hash}, block=${rcMint.blockNumber}`);

      // 2) owner -> receiver 转账 25
      logger.info(`开始转账: 从 ${ownerAddr} 转至 ${receiverAddr} 25 代币`);
      const txTransfer = await tokenWithOwner.transfer(receiverAddr, amountTransfer);
      const rcTransfer = await txTransfer.wait(1);
      logger.success(`转账完成: tx=${txTransfer.hash}, block=${rcTransfer.blockNumber}`);

      // 3) 以 owner 身份销毁 receiver 的 10（合约支持 onlyOwner 任意地址燃烧）
      logger.info(`开始销毁: 从 ${receiverAddr} 销毁 10 代币`);
      const txBurn = await tokenWithOwner.burn(receiverAddr, amountBurn);
      const rcBurn = await txBurn.wait(1);
      logger.success(`销毁完成: tx=${txBurn.hash}, block=${rcBurn.blockNumber}`);

      // 余额快照
      const [balOwner, balReceiver] = await Promise.all([
        token.balanceOf(ownerAddr),
        token.balanceOf(receiverAddr),
      ]);
      logger.info(`余额快照: owner=${ethers.utils.formatUnits(balOwner, tokenDecimals)}, receiver=${ethers.utils.formatUnits(balReceiver, tokenDecimals)}`);
    } catch (postErr) {
      logger.warn(`部署后产生活动事件失败（不影响合约可用）: ${postErr.message}`);
    }

    // 自动验证合约（复用模式下可选开启：REVERIFY=1 或 --reverify）
    const wantReverify = process.env.REVERIFY === "1" || process.argv.includes("--reverify");
    if (netConfig.verify && netName !== "localhost" && netName !== "hardhat" && (!isReuseMode || wantReverify)) {
      try {
        logger.info("开始验证合约...");
        await run("verify:verify", {
          address: contractAddress,
          constructorArguments: [name, symbol, initialOwner],
        });
        logger.success("合约验证成功");
      } catch (error) {
        logger.warn(`合约验证失败: ${error.message}`);
        logger.warn("可手动执行验证命令: npx hardhat verify --network ${netName} ${contractAddress} ${name} ${symbol} ${initialOwner}");
      }
    }

    logger.success("===== 部署流程完成 =====");
    return contractAddress;

  } catch (error) {
    logger.error(`部署失败: ${error.message}`);
    process.exit(1);
  }
}

// 执行主函数
main()
  .then(() => process.exit(0))
  .catch((error) => {
    logger.error(`执行出错: ${error.message}`);
    process.exit(1);
  });
