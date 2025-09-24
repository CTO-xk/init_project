/** @type import('hardhat/config').HardhatUserConfig */
const config = {
  solidity: "0.8.19",
  networks: {
    hardhat: {
      chainId: 1337,
      type: "edr-simulated"
    }
  },
  paths: {
    sources: "./src",
    tests: "./test",
    cache: "./cache",
    artifacts: "./artifacts"
  }
};

export default config;
