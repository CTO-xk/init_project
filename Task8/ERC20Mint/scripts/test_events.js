const { ethers } = require("hardhat");
const { getProvider, getSigner, getContract } = require("../utils/tokenApi");

// 固定参数：最小金额、单次循环、固定接收逻辑
const LOOPS_FIXED = 1;
const AMOUNT_MINT_FIXED = "1";      // 最小基础单位
const AMOUNT_TRANSFER_FIXED = "1";  // 最小基础单位
const AMOUNT_BURN_FIXED = "1";      // 最小基础单位

async function main() {
  const networkName = process.env.HARDHAT_NETWORK || process.env.NETWORK_NAME || "sepolia";
  console.log(`[INFO] network=${networkName}, loops=${LOOPS_FIXED}`);

  // provider + signer（默认使用 SEPOLIA_PK_ONE / BASE_SEPOLIA_PK_ONE 等环境变量）
  const signer = getSigner(
    process.env.SEPOLIA_PK_ONE || process.env.BASE_SEPOLIA_PK_ONE || process.env.PRIVATE_KEY,
    undefined,
    networkName
  );
  const provider = getProvider(undefined, networkName);
  const fromAddress = await signer.getAddress();

  const token = getContract(networkName, signer);
  const tokenRead = token.connect(provider);

  const decimals = await tokenRead.decimals();
  const vMint = ethers.utils.parseUnits(AMOUNT_MINT_FIXED, decimals);
  const vTransfer = ethers.utils.parseUnits(AMOUNT_TRANSFER_FIXED, decimals);
  const vBurn = ethers.utils.parseUnits(AMOUNT_BURN_FIXED, decimals);

  // 目标接收地址：优先参数，其次第二个本地 signer，否则回退为自己
  // 固定接收逻辑：优先环境变量 RECIPIENT，否则回退为自身地址（确保可执行且可预测）
  let toAddress = process.env.RECIPIENT;
  if (!toAddress || !ethers.utils.isAddress(toAddress)) {
    toAddress = fromAddress;
  }

  console.log(`[INFO] from=${fromAddress}, to=${toAddress}`);

  // 初始余额
  const [balFrom0, balTo0] = await Promise.all([
    tokenRead.balanceOf(fromAddress),
    tokenRead.balanceOf(toAddress),
  ]);
  console.log(`[INFO] balances(before): from=${ethers.utils.formatUnits(balFrom0, decimals)}, to=${ethers.utils.formatUnits(balTo0, decimals)}`);

  for (let i = 0; i < LOOPS_FIXED; i++) {
    console.log(`[INFO] loop ${i + 1}/${LOOPS_FIXED}`);

    // 1) mint 给 from
    const txMint = await token.mint(fromAddress, vMint);
    const rcMint = await txMint.wait(1);
    console.log(`[SUCCESS] Mint tx=${txMint.hash} block=${rcMint.blockNumber}`);

    // 2) from -> to 转账
    const txTransfer = await token.transfer(toAddress, vTransfer);
    const rcTransfer = await txTransfer.wait(1);
    console.log(`[SUCCESS] Transfer tx=${txTransfer.hash} block=${rcTransfer.blockNumber}`);

    // 3) 以 owner 身份从 to 销毁
    const txBurn = await token.burn(toAddress, vBurn);
    const rcBurn = await txBurn.wait(1);
    console.log(`[SUCCESS] Burn tx=${txBurn.hash} block=${rcBurn.blockNumber}`);
  }

  const [balFrom1, balTo1] = await Promise.all([
    tokenRead.balanceOf(fromAddress),
    tokenRead.balanceOf(toAddress),
  ]);
  console.log(`[INFO] balances(after): from=${ethers.utils.formatUnits(balFrom1, decimals)}, to=${ethers.utils.formatUnits(balTo1, decimals)}`);
}

main()
  .then(() => process.exit(0))
  .catch((err) => {
    console.error(`[ERROR] ${err.message}`);
    process.exit(1);
  });


