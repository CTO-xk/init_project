合约部署指南
1 前置条件
    已安装 Node.js (v14+) 和 npm
    已配置 Hardhat 开发环境
    拥有足够 ETH 的钱包（用于支付 gas 费用）
2 部署步骤
 1) 环境准备
 # 克隆项目并安装依赖
    git clone <xxxxxxxx>
    cd TOPIC1
    npm install --legacy-peer-deps
 2. 部署 Meme 代币合约
    # 本地部署
    npx hardhat run scripts/deploy.js --network localhost

    # 测试网部署
    npx hardhat run scripts/deploy.js --network speolia
 3. 验证合约
    npx hardhat verify --network goerli \
    <合约地址> \
    "MemeToken" \
    "MEME" \
    1000000000000000000000000000 \
    <开发者钱包地址> \
    <流动性池地址>
 4. ![image-20250823154648873](/Users/luna_xk/Library/Application Support/typora-user-images/image-20250823154648873.png)
 5. ![image-20250823154805687](/Users/luna_xk/Library/Application Support/typora-user-images/image-20250823154805687.png)
 6. ### 流动性操作注意事项

    1. 添加流动性时，需要同时提供代币和 ETH（比例由系统自动计算）
    2. 移除流动性会按比例返回代币和 ETH
    3. 流动性操作会产生 gas 费用，请确保钱包中有足够 ETH
