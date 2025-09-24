const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("Gas Analysis Tests", function () {
    let arithmetic;
    let arithmeticOptimized;
    let owner;

    beforeEach(async function () {
        [owner] = await ethers.getSigners();
        
        // 部署原始合约
        const Arithmetic = await ethers.getContractFactory("Arithmetic");
        arithmetic = await Arithmetic.deploy();
        await arithmetic.deployed();
        
        // 部署优化合约
        const ArithmeticOptimized = await ethers.getContractFactory("ArithmeticOptimized");
        arithmeticOptimized = await ArithmeticOptimized.deploy();
        await arithmeticOptimized.deployed();
    });

    describe("Gas 消耗对比分析", function () {
        it("应该对比加法操作的 Gas 消耗", async function () {
            console.log("\n=== 加法操作 Gas 消耗对比 ===");
            
            // 原始合约加法
            const originalTx = await arithmetic.add(100, 200);
            const originalReceipt = await originalTx.wait();
            const originalGas = originalReceipt.gasUsed;
            
            // 优化合约加法
            const optimizedTx = await arithmeticOptimized.add(100, 200);
            const optimizedReceipt = await optimizedTx.wait();
            const optimizedGas = optimizedReceipt.gasUsed;
            
            console.log("原始合约加法 Gas:", originalGas.toString());
            console.log("优化合约加法 Gas:", optimizedGas.toString());
            console.log("Gas 节省:", (originalGas - optimizedGas).toString());
            console.log("节省比例:", ((originalGas - optimizedGas) / originalGas * 100).toFixed(2) + "%");
            
            // 验证优化效果
            expect(optimizedGas).to.be.lt(originalGas);
            expect(originalGas - optimizedGas).to.be.gt(0);
        });

        it("应该对比减法操作的 Gas 消耗", async function () {
            console.log("\n=== 减法操作 Gas 消耗对比 ===");
            
            // 原始合约减法
            const originalTx = await arithmetic.subtract(200, 100);
            const originalReceipt = await originalTx.wait();
            const originalGas = originalReceipt.gasUsed;
            
            // 优化合约减法
            const optimizedTx = await arithmeticOptimized.subtract(200, 100);
            const optimizedReceipt = await optimizedTx.wait();
            const optimizedGas = optimizedReceipt.gasUsed;
            
            console.log("原始合约减法 Gas:", originalGas.toString());
            console.log("优化合约减法 Gas:", optimizedGas.toString());
            console.log("Gas 节省:", (originalGas - optimizedGas).toString());
            console.log("节省比例:", ((originalGas - optimizedGas) / originalGas * 100).toFixed(2) + "%");
            
            // 验证优化效果
            expect(optimizedGas).to.be.lt(originalGas);
        });

        it("应该对比乘法操作的 Gas 消耗", async function () {
            console.log("\n=== 乘法操作 Gas 消耗对比 ===");
            
            // 原始合约乘法
            const originalTx = await arithmetic.multiply(50, 60);
            const originalReceipt = await originalTx.wait();
            const originalGas = originalReceipt.gasUsed;
            
            // 优化合约乘法
            const optimizedTx = await arithmeticOptimized.multiply(50, 60);
            const optimizedReceipt = await optimizedTx.wait();
            const optimizedGas = optimizedReceipt.gasUsed;
            
            console.log("原始合约乘法 Gas:", originalGas.toString());
            console.log("优化合约乘法 Gas:", optimizedGas.toString());
            console.log("Gas 节省:", (originalGas - optimizedGas).toString());
            console.log("节省比例:", ((originalGas - optimizedGas) / originalGas * 100).toFixed(2) + "%");
            
            // 验证优化效果
            expect(optimizedGas).to.be.lt(originalGas);
        });

        it("应该对比除法操作的 Gas 消耗", async function () {
            console.log("\n=== 除法操作 Gas 消耗对比 ===");
            
            // 原始合约除法
            const originalTx = await arithmetic.divide(200, 50);
            const originalReceipt = await originalTx.wait();
            const originalGas = originalReceipt.gasUsed;
            
            // 优化合约除法
            const optimizedTx = await arithmeticOptimized.divide(200, 50);
            const optimizedReceipt = await optimizedTx.wait();
            const optimizedGas = optimizedReceipt.gasUsed;
            
            console.log("原始合约除法 Gas:", originalGas.toString());
            console.log("优化合约除法 Gas:", optimizedGas.toString());
            console.log("Gas 节省:", (originalGas - optimizedGas).toString());
            console.log("节省比例:", ((originalGas - optimizedGas) / originalGas * 100).toFixed(2) + "%");
            
            // 验证优化效果
            expect(optimizedGas).to.be.lt(originalGas);
        });

        it("应该对比模运算的 Gas 消耗", async function () {
            console.log("\n=== 模运算 Gas 消耗对比 ===");
            
            // 原始合约模运算
            const originalTx = await arithmetic.modulo(157, 50);
            const originalReceipt = await originalTx.wait();
            const originalGas = originalReceipt.gasUsed;
            
            // 优化合约模运算
            const optimizedTx = await arithmeticOptimized.modulo(157, 50);
            const optimizedReceipt = await optimizedTx.wait();
            const optimizedGas = optimizedReceipt.gasUsed;
            
            console.log("原始合约模运算 Gas:", originalGas.toString());
            console.log("优化合约模运算 Gas:", optimizedGas.toString());
            console.log("Gas 节省:", (originalGas - optimizedGas).toString());
            console.log("节省比例:", ((originalGas - optimizedGas) / originalGas * 100).toFixed(2) + "%");
            
            // 验证优化效果
            expect(optimizedGas).to.be.lt(originalGas);
        });
    });

    describe("批量操作 Gas 优化分析", function () {
        it("应该分析批量操作的 Gas 优化效果", async function () {
            console.log("\n=== 批量操作 Gas 优化分析 ===");
            
            // 测试 3 次单独操作
            let totalSingleGas = 0;
            
            for (let i = 0; i < 3; i++) {
                const tx = await arithmeticOptimized.add(1, 2);
                const receipt = await tx.wait();
                totalSingleGas += receipt.gasUsed.toNumber();
            }
            
            // 测试批量操作
            const a = [1, 1, 1];
            const b = [2, 2, 2];
            const operationTypes = [
                await arithmeticOptimized.ADD(),
                await arithmeticOptimized.ADD(),
                await arithmeticOptimized.ADD()
            ];
            
            const batchTx = await arithmeticOptimized.batchExecute(a, b, operationTypes);
            const batchReceipt = await batchTx.wait();
            const batchGas = batchReceipt.gasUsed;
            
            console.log("3 次单独操作总 Gas:", totalSingleGas);
            console.log("批量操作总 Gas:", batchGas.toString());
            console.log("每次单独操作平均 Gas:", (totalSingleGas / 3).toFixed(0));
            console.log("每次批量操作平均 Gas:", (batchGas / 3).toFixed(0));
            console.log("每次操作 Gas 节省:", ((totalSingleGas / 3) - (batchGas / 3)).toFixed(0));
            console.log("批量操作节省比例:", (((totalSingleGas / 3) - (batchGas / 3)) / (totalSingleGas / 3) * 100).toFixed(2) + "%");
            
            // 验证批量操作的优势
            expect(batchGas / 3).to.be.lt(totalSingleGas / 3);
        });

        it("应该分析不同数量批量操作的 Gas 效率", async function () {
            console.log("\n=== 不同数量批量操作 Gas 效率分析 ===");
            
            const testSizes = [3, 5, 10];
            
            for (const size of testSizes) {
                console.log(`\n--- ${size} 次操作分析 ---`);
                
                // 单独操作
                let totalSingleGas = 0;
                for (let i = 0; i < size; i++) {
                    const tx = await arithmeticOptimized.add(1, 2);
                    const receipt = await tx.wait();
                    totalSingleGas += receipt.gasUsed.toNumber();
                }
                
                // 批量操作
                const a = new Array(size).fill(1);
                const b = new Array(size).fill(2);
                const operationTypes = new Array(size).fill(await arithmeticOptimized.ADD());
                
                const batchTx = await arithmeticOptimized.batchExecute(a, b, operationTypes);
                const batchReceipt = await batchTx.wait();
                const batchGas = batchReceipt.gasUsed;
                
                const singleAvg = totalSingleGas / size;
                const batchAvg = batchGas / size;
                const savings = singleAvg - batchAvg;
                const savingsPercent = (savings / singleAvg * 100);
                
                console.log(`单独操作平均 Gas: ${singleAvg.toFixed(0)}`);
                console.log(`批量操作平均 Gas: ${batchAvg.toFixed(0)}`);
                console.log(`每次操作节省: ${savings.toFixed(0)}`);
                console.log(`节省比例: ${savingsPercent.toFixed(2)}%`);
                
                // 验证批量操作的优势
                expect(batchAvg).to.be.lt(singleAvg);
            }
        });
    });

    describe("存储成本分析", function () {
        it("应该分析存储操作的 Gas 成本", async function () {
            console.log("\n=== 存储操作 Gas 成本分析 ===");
            
            // 测试多次操作后的存储成本
            const operations = [10, 50, 100];
            
            for (const count of operations) {
                console.log(`\n--- ${count} 次操作存储成本分析 ---`);
                
                // 重新部署合约以确保干净状态
                const Arithmetic = await ethers.getContractFactory("Arithmetic");
                const freshArithmetic = await Arithmetic.deploy();
                await freshArithmetic.deployed();
                
                const ArithmeticOptimized = await ethers.getContractFactory("ArithmeticOptimized");
                const freshArithmeticOptimized = await ArithmeticOptimized.deploy();
                await freshArithmeticOptimized.deployed();
                
                // 执行多次操作
                for (let i = 0; i < count; i++) {
                    await freshArithmetic.add(i, i + 1);
                    await freshArithmeticOptimized.add(i, i + 1);
                }
                
                // 获取操作历史数量
                const originalCount = await freshArithmetic.getOperationsCount();
                const optimizedCount = await freshArithmeticOptimized.getOperationsCount();
                
                console.log(`原始合约操作历史数量: ${originalCount}`);
                console.log(`优化合约操作历史数量: ${optimizedCount}`);
                
                // 验证操作历史记录正确
                expect(originalCount).to.equal(count);
                expect(optimizedCount).to.equal(count);
            }
        });
    });

    describe("综合 Gas 分析报告", function () {
        it("应该生成综合 Gas 分析报告", async function () {
            console.log("\n" + "=".repeat(60));
            console.log("🚀 智能合约 Gas 优化分析报告");
            console.log("=".repeat(60));
            
            // 测试所有基本操作
            const operations = [
                { name: "加法", func: "add", args: [100, 200] },
                { name: "减法", func: "subtract", args: [200, 100] },
                { name: "乘法", func: "multiply", args: [50, 60] },
                { name: "除法", func: "divide", args: [200, 50] },
                { name: "模运算", func: "modulo", args: [157, 50] }
            ];
            
            let totalOriginalGas = 0;
            let totalOptimizedGas = 0;
            
            for (const op of operations) {
                // 原始合约
                const originalTx = await arithmetic[op.func](...op.args);
                const originalReceipt = await originalTx.wait();
                const originalGas = originalReceipt.gasUsed;
                
                // 优化合约
                const optimizedTx = await arithmeticOptimized[op.func](...op.args);
                const optimizedReceipt = await optimizedTx.wait();
                const optimizedGas = optimizedReceipt.gasUsed;
                
                totalOriginalGas += originalGas.toNumber();
                totalOptimizedGas += optimizedGas.toNumber();
                
                const savings = originalGas - optimizedGas;
                const savingsPercent = (savings / originalGas * 100);
                
                console.log(`${op.name}: 原始 ${originalGas}, 优化 ${optimizedGas}, 节省 ${savings} (${savingsPercent.toFixed(1)}%)`);
            }
            
            const totalSavings = totalOriginalGas - totalOptimizedGas;
            const totalSavingsPercent = (totalSavings / totalOriginalGas * 100);
            
            console.log("-".repeat(60));
            console.log(`总计: 原始 ${totalOriginalGas}, 优化 ${totalOptimizedGas}, 节省 ${totalSavings} (${totalSavingsPercent.toFixed(1)}%)`);
            console.log("=".repeat(60));
            
            // 验证总体优化效果
            expect(totalOptimizedGas).to.be.lt(totalOriginalGas);
            expect(totalSavings).to.be.gt(0);
        });
    });
});
