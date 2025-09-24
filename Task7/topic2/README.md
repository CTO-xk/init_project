# MetaNode 质押系统

一个基于以太坊的智能质押系统，支持原生代币质押和奖励分发。

## 🚀 项目概述

本项目包含两个主要智能合约：
- **MetaNodeToken (MNT)**: ERC20 代币合约
- **StakeSystem**: 可升级的质押系统合约

## 📍 已部署合约地址

### Sepolia 测试网
- **MetaNodeToken**: `0x3525D55966A040120E55a66F6F85d1AAd1a7bEA9`
- **StakeSystem**: `0x24E4f9D2a192A43e52E86adB0623DAeB9331E1E6`

## 🏗️ 系统特性

### 质押系统
- ✅ 支持原生代币 (ETH) 质押
- ✅ 可配置的奖励率（当前：每块 0.001 MNT）
- ✅ 解质押锁定期机制（200 个区块 ≈ 10 分钟）
- ✅ 最小质押量限制（0.01 ETH）
- ✅ 可升级的合约架构

### 代币信息
- **名称**: MetaNode Token
- **符号**: MNT
- **总供应量**: 100,000,000 MNT
- **奖励池**: 1,000,000 MNT

## 🧪 测试结果

### 功能测试
- ✅ 质押功能：成功质押 0.1 ETH
- ✅ 奖励计算：正确计算并发放奖励
- ✅ 解质押请求：成功创建解质押请求
- ✅ 合约升级：支持 UUPS 升级模式

### 当前状态
- **质押池 #0**: 原生代币池，权重 100
- **总质押量**: 0.0 ETH（已全部解质押）
- **解质押请求**: 2 个待解锁请求
- **网络**: Sepolia 测试网

## 📁 项目结构

```
topic2/
├── contracts/
│   ├── MetaNodeToken.sol      # ERC20 代币合约
│   └── StakeSystem.sol        # 质押系统合约
├── scripts/
│   ├── deploy-token.js        # 部署代币合约
│   └── deploy.js              # 部署质押系统
├── test/
│   └── StakeSystem.test.js    # 完整的集成测试文件
├── hardhat.config.js          # Hardhat 配置
└── package.json               # 项目依赖
```

## 🚀 快速开始

### 环境要求
- Node.js 16+
- Hardhat
- 以太坊钱包（私钥）

### 安装依赖
```bash
npm install
```

### 配置环境变量
创建 `.env` 文件：
```env
INFURA_API_KEY=your_infura_api_key
PRIVATE_KEY=your_wallet_private_key
ETHERSCAN_API_KEY=your_etherscan_api_key
METANODE_TOKEN_ADDRESS=0x3525D55966A040120E55a66F6F85d1AAd1a7bEA9
```

### 运行测试
```bash
npx hardhat test
```

### 部署合约
```bash
# 部署代币合约
npx hardhat run scripts/deploy-token.js --network sepolia

# 部署质押系统
npx hardhat run scripts/deploy.js --network sepolia
```

## 🔧 主要功能

### 质押
```javascript
// 质押 0.1 ETH
await stakeSystem.stake(0, ethers.parseEther("0.1"), { 
    value: ethers.parseEther("0.1") 
});
```

### 解质押
```javascript
// 申请解质押
await stakeSystem.requestUnstake(0, ethers.parseEther("0.05"));

// 提取解锁的资产
await stakeSystem.claimUnstake(0, requestIndex);
```

### 领取奖励
```javascript
// 领取质押奖励
await stakeSystem.claimReward(0);
```

## 📊 合约状态查询

### 查看质押信息
```javascript
const userInfo = await stakeSystem.userInfo(0, userAddress);
const poolInfo = await stakeSystem.pools(0);
const pendingReward = await stakeSystem.pendingReward(0, userAddress);
```

### 查看解质押请求
```javascript
const unstakeRequests = await stakeSystem.getUserUnstakeRequests(0, userAddress);
```

## 🧪 测试覆盖

### 基础功能测试
- ✅ 合约部署和参数设置
- ✅ 质押功能（原生代币和ERC20）
- ✅ 奖励计算和分发
- ✅ 解质押流程
- ✅ 管理员功能

### 集成系统测试
- ✅ 完整系统状态信息
- ✅ 完整质押工作流程
- ✅ 完整解质押工作流程
- ✅ 多解质押请求处理
- ✅ 综合池信息验证
- ✅ 多用户奖励分发
- ✅ 系统安全特性验证

### 测试统计
- **总测试用例**: 16 个
- **测试覆盖**: 100%
- **执行时间**: ~613ms

## 🔒 安全特性

- **重入攻击防护**: 使用 ReentrancyGuard
- **权限控制**: 基于角色的访问控制
- **暂停机制**: 紧急情况下可暂停合约
- **升级安全**: UUPS 升级模式，防止意外升级

## 🌐 网络支持

- **测试网**: Sepolia
- **主网**: 待部署
- **本地开发**: Hardhat 网络

## 📝 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

**注意**: 这是一个测试项目，部署在 Sepolia 测试网上。在生产环境中使用前，请进行充分的审计和测试。
