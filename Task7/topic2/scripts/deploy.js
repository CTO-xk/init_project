const hre= require("hardhat");
const { ethers,upgrades } = hre;

async function main() {
    console.log("Deploying StakeSystem to Sepolia...");
    // 获取部署者账户
    const [deployer] = await ethers.getSigners();
    console.log(`Deploying with account: ${deployer.address}`);
    // 确保已提供奖励代币地址
    if (!process.env.METANODE_TOKEN_ADDRESS) {
        throw new Error("Please set METANODE_TOKEN_ADDRESS in .env file");
    }
    const metaNodeTokenAddress = process.env.METANODE_TOKEN_ADDRESS;
     // 部署可升级的质押系统合约
     const StakeSystem=await ethers.getContractFactory("StakeSystem");
     const stakeSystem = await upgrades.deployProxy(StakeSystem, [
        metaNodeTokenAddress,          // 奖励代币地址
        ethers.parseEther("0.001"),    // 每个区块的奖励
        deployer.address               // 操作员地址
      ]);
      // 等待部署完成
      await stakeSystem.waitForDeployment();
      console.log("StakeSystem deployed to:", stakeSystem.target);
      const stakeSystemAddress=await stakeSystem.getAddress();
      console.log(`StakeSystem deployed to: ${stakeSystemAddress}`);

    // 向质押合约转移奖励代币
    console.log("Transferring reward tokens to StakeSystem...");
    const MetaNodeToken=await ethers.getContractFactory("MetaNodeToken");
    const metaNodeToken=await MetaNodeToken.attach(metaNodeTokenAddress);
    // 转移100万奖励代币到质押合约
    const transferAmount=ethers.parseEther("1000000");
    await metaNodeToken.transfer(stakeSystemAddress,transferAmount);
    console.log(`Transferred ${ethers.formatEther(transferAmount)} MNT to StakeSystem`);
      // 验证合约（Etherscan API密钥）
  if (process.env.ETHERSCAN_API_KEY) {
    console.log("Waiting for block confirmations...");
    // 等待6个区块确认
    await stakeSystem.deploymentTransaction().wait(6);
    
    console.log("Verifying contract on Etherscan...");
    const implementationAddress = await upgrades.erc1967.getImplementationAddress(stakeSystemAddress);
    
    await hre.run("verify:verify", {
      address: implementationAddress,
      constructorArguments: [],
    });
    console.log("Contract implementation verified on Etherscan");
  }
}
main()
.then(()=>process.exit(0))
.catch((error)=>{
    console.error(error);
    process.exit(1);
});