const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("MemeToken", function () {
  let memeToken;
  let owner;
  let addr1;
  let addr2;
  let devWallet;
  let liquidityPool;

  const totalSupply = ethers.utils.parseEther("1000000000");

  beforeEach(async function () {
    [owner, addr1, addr2, devWallet, liquidityPool] = await ethers.getSigners();

    // 部署合约
    const MemeToken = await ethers.getContractFactory("MemeToken");
    memeToken = await MemeToken.deploy(
      "MemeToken",
      "MEME",
      totalSupply,
      devWallet.address,
      liquidityPool.address
    );
    await memeToken.deployed();

    // 激活交易
    await memeToken.activateTrading();

    // 向addr1转账一些代币用于测试
    await memeToken.transfer(addr1.address, ethers.utils.parseEther("10000"));
  });

  describe("基本功能测试", function () {
    it("应该正确设置代币元数据", async function () {
      expect(await memeToken.name()).to.equal("MemeToken");
      expect(await memeToken.symbol()).to.equal("MEME");
      expect(await memeToken.totalSupply()).to.equal(totalSupply);
      expect(await memeToken.balanceOf(owner.address)).to.equal(
        ethers.utils.parseEther("999990000") // 总供应量减去转给addr1的10000
      );
    });

    it("应该正确设置初始地址", async function () {
      expect(await memeToken.devWallet()).to.equal(devWallet.address);
      expect(await memeToken.liquidityPool()).to.equal(liquidityPool.address);
    });
  });

  describe("交易税机制测试", function () {
    it("买入交易应该正确征收税费", async function () {
      // 先给流动性池一些代币用于测试
      await memeToken.transfer(liquidityPool.address, ethers.utils.parseEther("10000"));
      
      // 模拟从流动性池买入（addr1从流动性池接收代币）
      const initialBalance = await memeToken.balanceOf(addr1.address);
      const transferAmount = ethers.utils.parseEther("1000");
      
      // 让流动性池转账给addr1（模拟买入）
      await memeToken.connect(liquidityPool).transfer(addr1.address, transferAmount);
      
      // 买入税率为5%，所以实际收到应为950
      const expectedReceived = transferAmount.mul(95).div(100);
      expect(await memeToken.balanceOf(addr1.address)).to.equal(initialBalance.add(expectedReceived));
      
      // 检查税费分配
      const devTax = transferAmount.mul(2).div(100);
      expect(await memeToken.balanceOf(devWallet.address)).to.equal(devTax);
    });

    it("卖出交易应该正确征收税费", async function () {
      // 模拟卖出到流动性池
      const initialBalance = await memeToken.balanceOf(liquidityPool.address);
      const transferAmount = ethers.utils.parseEther("1000");
      
      // addr1向流动性池转账（模拟卖出）
      await memeToken.connect(addr1).transfer(liquidityPool.address, transferAmount);
      
      // 卖出税率为10%，所以流动性池实际收到应为900
      const expectedReceived = transferAmount.mul(90).div(100);
      expect(await memeToken.balanceOf(liquidityPool.address)).to.equal(initialBalance.add(expectedReceived));
    });
  });

  describe("交易限制测试", function () {
    it("应该限制超过最大额度的交易", async function () {
      // 最大交易额度是总供应量的2% = 20,000,000
      const maxAmount = ethers.utils.parseEther("20000000");
      const tooMuchAmount = maxAmount.add(ethers.utils.parseEther("1"));
      
      // 使用addr1来测试，因为owner被排除在税收之外
      await expect(
        memeToken.connect(addr1).transfer(addr2.address, tooMuchAmount)
      ).to.be.revertedWith("Amount exceeds max transaction");
    });

    it("应该限制超过每日交易次数的交易", async function () {
      const transferAmount = ethers.utils.parseEther("100");
      
      // 默认每日最大交易次数是10次
      for (let i = 0; i < 10; i++) {
        await memeToken.connect(addr1).transfer(addr2.address, transferAmount);
      }
      
      // 第11次交易应该失败
      await expect(
        memeToken.connect(addr1).transfer(addr2.address, transferAmount)
      ).to.be.revertedWith("Exceeded daily transaction limit");
    });
  });

  describe("流动性功能测试", function () {
    it("应该允许添加流动性", async function () {
      const tokenAmount = ethers.utils.parseEther("1000");
      const ethAmount = ethers.utils.parseEther("1");
      
      // 先批准代币转账
      await memeToken.connect(addr1).approve(memeToken.address, tokenAmount);
      
      // 添加流动性
      await expect(
        memeToken.connect(addr1).addLiquidity(tokenAmount, { value: ethAmount })
      ).to.emit(memeToken, "AddLiquidity")
        .withArgs(addr1.address, tokenAmount, ethAmount, 0);
    });
  });

  describe("管理员功能测试", function () {
    it("应该允许更新税率", async function () {
      await expect(
        memeToken.updateTaxRates(6, 10, 2, 2, 2)
      ).to.emit(memeToken, "UpdateTaxRates")
        .withArgs(6, 10, 2, 2, 2);
      
      expect(await memeToken.buyTax()).to.equal(6);
      expect(await memeToken.sellTax()).to.equal(10);
    });

    it("应该禁止非管理员更新税率", async function () {
      await expect(
        memeToken.connect(addr1).updateTaxRates(4, 9, 2, 2, 2)
      ).to.be.revertedWith("Ownable: caller is not the owner");
    });
  });
});
