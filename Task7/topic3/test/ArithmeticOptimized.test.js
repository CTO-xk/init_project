const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("ArithmeticOptimized Contract", function () {
    let arithmetic;
    let owner;
    let addr1;
    let addr2;

    beforeEach(async function () {
        // 获取合约工厂
        const ArithmeticOptimized = await ethers.getContractFactory("ArithmeticOptimized");
        
        // 部署合约
        [owner, addr1, addr2] = await ethers.getSigners();
        arithmetic = await ArithmeticOptimized.deploy();
        await arithmetic.deployed();
    });

    describe("基础功能测试", function () {
        it("应该正确执行加法运算", async function () {
            const result = await arithmetic.add(5, 3);
            expect(result).to.equal(8);
            
            const lastResult = await arithmetic.lastResult();
            expect(lastResult).to.equal(8);
        });

        it("应该正确执行减法运算", async function () {
            const result = await arithmetic.subtract(10, 4);
            expect(result).to.equal(6);
            
            const lastResult = await arithmetic.lastResult();
            expect(lastResult).to.equal(6);
        });

        it("应该正确执行乘法运算", async function () {
            const result = await arithmetic.multiply(6, 7);
            expect(result).to.equal(42);
            
            const lastResult = await arithmetic.lastResult();
            expect(lastResult).to.equal(42);
        });

        it("应该正确执行除法运算", async function () {
            const result = await arithmetic.divide(20, 5);
            expect(result).to.equal(4);
            
            const lastResult = await arithmetic.lastResult();
            expect(lastResult).to.equal(4);
        });

        it("应该正确执行模运算", async function () {
            const result = await arithmetic.modulo(17, 5);
            expect(result).to.equal(2);
            
            const lastResult = await arithmetic.lastResult();
            expect(lastResult).to.equal(2);
        });
    });

    describe("边界条件测试", function () {
        it("应该正确处理加零的情况", async function () {
            const result = await arithmetic.add(5, 0);
            expect(result).to.equal(5);
        });

        it("应该正确处理减零的情况", async function () {
            const result = await arithmetic.subtract(10, 0);
            expect(result).to.equal(10);
        });

        it("应该正确处理乘以零的情况", async function () {
            const result = await arithmetic.multiply(5, 0);
            expect(result).to.equal(0);
        });

        it("应该拒绝除零操作", async function () {
            await expect(arithmetic.divide(10, 0))
                .to.be.revertedWith("Div by zero");
        });

        it("应该拒绝模零操作", async function () {
            await expect(arithmetic.modulo(10, 0))
                .to.be.revertedWith("Mod by zero");
        });
    });

    describe("溢出测试", function () {
        it("应该拒绝加法溢出", async function () {
            const maxValue = ethers.constants.MaxUint128;
            await expect(arithmetic.add(maxValue, 1))
                .to.be.revertedWith("Overflow");
        });

        it("应该拒绝减法下溢", async function () {
            await expect(arithmetic.subtract(5, 10))
                .to.be.revertedWith("Underflow");
        });

        it("应该拒绝乘法溢出", async function () {
            const maxValue = ethers.constants.MaxUint128;
            await expect(arithmetic.multiply(maxValue, 2))
                .to.be.revertedWith("Overflow");
        });
    });

    describe("操作历史测试", function () {
        it("应该正确记录操作历史", async function () {
            // 执行多个操作
            await arithmetic.add(1, 2);
            await arithmetic.subtract(5, 3);
            await arithmetic.multiply(4, 6);

            const count = await arithmetic.getOperationsCount();
            expect(count).to.equal(3);

            // 检查第一个操作
            const [a, b, opType, result, timestamp] = await arithmetic.getOperation(0);
            expect(a).to.equal(1);
            expect(b).to.equal(2);
            expect(opType).to.equal(await arithmetic.ADD());
            expect(result).to.equal(3);
            expect(timestamp).to.be.gt(0);
        });

        it("应该正确清除操作历史", async function () {
            await arithmetic.add(1, 2);
            await arithmetic.subtract(5, 3);

            let count = await arithmetic.getOperationsCount();
            expect(count).to.equal(2);

            await arithmetic.clearHistory();

            count = await arithmetic.getOperationsCount();
            expect(count).to.equal(0);

            const lastResult = await arithmetic.lastResult();
            expect(lastResult).to.equal(0);
        });
    });

    describe("事件测试", function () {
        it("应该正确发出计算事件", async function () {
            await expect(arithmetic.add(5, 3))
                .to.emit(arithmetic, "CalculationPerformed")
                .withArgs(5, 3, await arithmetic.ADD(), 8, await time());
        });
    });

    describe("错误处理测试", function () {
        it("应该拒绝访问越界的操作历史", async function () {
            await expect(arithmetic.getOperation(0))
                .to.be.revertedWith("Out of bounds");
        });
    });

    describe("批量操作测试", function () {
        it("应该正确执行批量操作", async function () {
            const a = [5, 10, 6];
            const b = [3, 4, 7];
            const operationTypes = [
                await arithmetic.ADD(),
                await arithmetic.SUBTRACT(),
                await arithmetic.MULTIPLY()
            ];

            const results = await arithmetic.batchExecute(a, b, operationTypes);

            expect(results.length).to.equal(3);
            expect(results[0]).to.equal(8);  // 5 + 3
            expect(results[1]).to.equal(6);  // 10 - 4
            expect(results[2]).to.equal(42); // 6 * 7
        });

        it("应该拒绝长度不匹配的批量操作", async function () {
            const a = [1, 2];
            const b = [3, 4, 5];
            const operationTypes = [1, 2, 3];

            await expect(arithmetic.batchExecute(a, b, operationTypes))
                .to.be.revertedWith("Length mismatch");
        });
    });

    describe("Gas 消耗测试", function () {
        it("应该测量优化后加法操作的 Gas 消耗", async function () {
            const tx = await arithmetic.add(100, 200);
            const receipt = await tx.wait();
            
            console.log("优化后加法操作 Gas 消耗:", receipt.gasUsed.toString());
            
            // 验证 Gas 消耗在合理范围内（应该比原始版本低）
            expect(receipt.gasUsed).to.be.lt(80000);
        });

        it("应该比较优化后不同操作的 Gas 消耗", async function () {
            // 测试加法
            const addTx = await arithmetic.add(1, 2);
            const addReceipt = await addTx.wait();
            
            // 测试减法
            const subTx = await arithmetic.subtract(5, 3);
            const subReceipt = await subTx.wait();
            
            // 测试乘法
            const mulTx = await arithmetic.multiply(4, 6);
            const mulReceipt = await mulTx.wait();
            
            // 测试除法
            const divTx = await arithmetic.divide(20, 5);
            const divReceipt = await divTx.wait();

            console.log("优化后 Gas 消耗对比:");
            console.log("加法:", addReceipt.gasUsed.toString());
            console.log("减法:", subReceipt.gasUsed.toString());
            console.log("乘法:", mulReceipt.gasUsed.toString());
            console.log("除法:", divReceipt.gasUsed.toString());

            // 验证所有操作的 Gas 消耗都在合理范围内
            expect(addReceipt.gasUsed).to.be.lt(80000);
            expect(subReceipt.gasUsed).to.be.lt(80000);
            expect(mulReceipt.gasUsed).to.be.lt(80000);
            expect(divReceipt.gasUsed).to.be.lt(80000);
        });
    });

    describe("Gas 优化效果测试", function () {
        it("应该测试批量操作的 Gas 优化效果", async function () {
            // 测试单次操作
            const singleTx = await arithmetic.add(1, 2);
            const singleReceipt = await singleTx.wait();
            const singleGas = singleReceipt.gasUsed;

            // 测试批量操作
            const a = [1, 1, 1];
            const b = [2, 2, 2];
            const operationTypes = [
                await arithmetic.ADD(),
                await arithmetic.ADD(),
                await arithmetic.ADD()
            ];

            const batchTx = await arithmetic.batchExecute(a, b, operationTypes);
            const batchReceipt = await batchTx.wait();
            const batchGas = batchReceipt.gasUsed;

            console.log("单次操作 Gas:", singleGas.toString());
            console.log("批量操作 Gas:", batchGas.toString());
            console.log("每次操作平均 Gas:", (batchGas / 3).toString());
            console.log("Gas 节省:", (singleGas - (batchGas / 3)).toString());

            // 批量操作应该比单独操作更高效
            expect(batchGas / 3).to.be.lt(singleGas);
        });
    });

    // 辅助函数：获取当前时间戳
    async function time() {
        const blockNum = await ethers.provider.getBlockNumber();
        const block = await ethers.provider.getBlock(blockNum);
        return block.timestamp;
    }
});
