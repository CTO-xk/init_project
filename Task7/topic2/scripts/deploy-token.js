const hre=require("hardhat");
const { ethers } = hre;

async function main(){
    console.log("Deploying MetaNodeToken to Sepolia...");

    // 部署奖励代币
    const MetaNodeToken=await ethers.getContractFactory("MetaNodeToken");
    const metaNodeToken = await MetaNodeToken.deploy(
        "MetaNode Token",  // 名称
        "MNT",             // 符号
        hre.ethers.parseEther("100000000"),  // 初始供应量：1亿
        (await hre.ethers.getSigners())[0].address  // 初始所有者
      );
      await metaNodeToken.waitForDeployment();
      const tokenAddress = await metaNodeToken.getAddress();
      console.log(`MetaNodeToken deployed to: ${tokenAddress}`);
      if(process.env.ETHERSCAN_API_KEY){
        console.log("Waiting for block confirmations...");
        await metaNodeToken.deploymentTransaction().wait(6);
        await hre.run("verify:verify", {
            address: tokenAddress,
            constructorArguments: [
              "MetaNode Token",
              "MNT",
              hre.ethers.parseEther("100000000"),
              (await hre.ethers.getSigners())[0].address
            ],
          });
          console.log("Contract verified on Etherscan");
        }
}
main()
.then(()=>process.exit(0))
.catch((error)=>{
    console.error(error);
    process.exit(1);
});
