async function main(){
    //获取部署账户
    const [deployer] = await ethers.getSigners();
    console.log("Deploying contracts with the account:",deployer.address);
    console.log("Account balance:",(await deployer.getBalance()).toString());

    //定义代币参数
    const name = "MemeToken";
    const symbol = "MEME";
    const totalSupply = ethers.utils.parseUnits("1000000000");
     //钱包和流动性池地址
    const devWallet = deployer.address; // 临时使用部署者地址作为开发者钱包
    const liquidityPool = deployer.address; // 临时使用部署者地址作为流动性池

    //部署合约
    const MemeToken = await ethers.getContractFactory("MemeToken");
    const memeToken = await MemeToken.deploy(name,symbol,totalSupply,devWallet,liquidityPool);
    await memeToken.deployed();
    console.log("MemeToken deployed to:",memeToken.address);
    //激活交易功能
    await memeToken.activateTrading();
    console.log("Trading activated");
    //部署后的基本信息
    console.log("Token name:", await memeToken.name());
    console.log("Token symbol:", await memeToken.symbol());
    console.log("Total supply:", ethers.utils.formatEther(await memeToken.totalSupply()));
    console.log("Deployer balance:", ethers.utils.formatEther(await memeToken.balanceOf(deployer.address)));
    //设置税率 (买入税7%, 卖出税10%, 流动性税3%, 奖励税2%, 开发税2%)
    await memeToken.updateTaxRates(7,10,3,2,2);
    console.log("Tax rates updated");
    //设置交易限制 (最大交易额20M, 最大钱包50M, 每日最大交易次数10)
    await memeToken.updateTransactionLimits(
        ethers.utils.parseEther("20000000"), // 20M tokens
        ethers.utils.parseEther("50000000"), // 50M tokens
        10
    );
    console.log("Transaction limits updated");
}
main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
