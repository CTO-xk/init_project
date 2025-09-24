const { ethers } = require("hardhat");

async function main() {
    console.log("🚀 开始部署智能合约...");

    // 获取部署账户
    const [deployer] = await ethers.getSigners();
    console.log("部署账户:", deployer.address);
    console.log("账户余额:", ethers.utils.formatEther(await deployer.getBalance()));

    // 部署原始算术合约
    console.log("\n📝 部署原始算术合约...");
    const Arithmetic = await ethers.getContractFactory("Arithmetic");
    const arithmetic = await Arithmetic.deploy();
    await arithmetic.deployed();
    console.log("原始算术合约已部署到:", arithmetic.address);

    // 部署优化后的算术合约
    console.log("\n⚡ 部署优化后的算术合约...");
    const ArithmeticOptimized = await ethers.getContractFactory("ArithmeticOptimized");
    const arithmeticOptimized = await ArithmeticOptimized.deploy();
    await arithmeticOptimized.deployed();
    console.log("优化算术合约已部署到:", arithmeticOptimized.address);

    // 验证合约功能
    console.log("\n🔍 验证合约功能...");
    
    // 测试原始合约
    console.log("测试原始合约...");
    const addResult = await arithmetic.add(5, 3);
    console.log("5 + 3 =", addResult.toString());
    
    const subResult = await arithmetic.subtract(10, 4);
    console.log("10 - 4 =", subResult.toString());
    
    const mulResult = await arithmetic.multiply(6, 7);
    console.log("6 * 7 =", mulResult.toString());
    
    const divResult = await arithmetic.divide(20, 5);
    console.log("20 / 5 =", divResult.toString());
    
    const modResult = await arithmetic.modulo(17, 5);
    console.log("17 % 5 =", modResult.toString());

    // 测试优化合约
    console.log("\n测试优化合约...");
    const addResultOpt = await arithmeticOptimized.add(5, 3);
    console.log("5 + 3 =", addResultOpt.toString());
    
    const subResultOpt = await arithmeticOptimized.subtract(10, 4);
    console.log("10 - 4 =", subResultOpt.toString());
    
    const mulResultOpt = await arithmeticOptimized.multiply(6, 7);
    console.log("6 * 7 =", mulResultOpt.toString());
    
    const divResultOpt = await arithmeticOptimized.divide(20, 5);
    console.log("20 / 5 =", divResultOpt.toString());
    
    const modResultOpt = await arithmeticOptimized.modulo(17, 5);
    console.log("17 % 5 =", modResultOpt.toString());

    // 测试批量操作
    console.log("\n测试批量操作...");
    const a = [1, 2, 3];
    const b = [4, 5, 6];
    const operationTypes = [
        await arithmeticOptimized.ADD(),
        await arithmeticOptimized.SUBTRACT(),
        await arithmeticOptimized.MULTIPLY()
    ];
    
    const batchResults = await arithmeticOptimized.batchExecute(a, b, operationTypes);
    console.log("批量操作结果:", batchResults.map(r => r.toString()));

    // 部署总结
    console.log("\n" + "=".repeat(60));
    console.log("✅ 合约部署完成！");
    console.log("=".repeat(60));
    console.log("原始算术合约地址:", arithmetic.address);
    console.log("优化算术合约地址:", arithmeticOptimized.address);
    console.log("部署账户:", deployer.address);
    console.log("=".repeat(60));

    // 返回合约地址供后续使用
    return {
        arithmetic: arithmetic.address,
        arithmeticOptimized: arithmeticOptimized.address,
        deployer: deployer.address
    };
}

// 如果直接运行此脚本
if (require.main === module) {
    main()
        .then(() => process.exit(0))
        .catch((error) => {
            console.error("❌ 部署失败:", error);
            process.exit(1);
        });
}

module.exports = { main };
