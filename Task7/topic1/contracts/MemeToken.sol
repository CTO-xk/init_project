// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

contract MemeToken is ERC20, Ownable, ReentrancyGuard {
    // 税收配置
    uint256 public buyTax = 5;      // 买入税率 (%)
    uint256 public sellTax = 10;     // 卖出税率 (%)
    uint256 public liquidityTax = 3; // 用于流动性的税率 (%)
    uint256 public rewardTax = 2;    // 用于持有者奖励的税率 (%)
    uint256 public devTax = 2;       // 用于开发者的税率 (%)
    
    // 交易限制配置
    uint256 public maxTransactionAmount; // 最大交易金额
    uint256 public maxWalletAmount;      // 最大钱包持有量
    uint256 public dailyTransactionLimit = 10; // 每日最大交易次数
    
    // 流动性池地址
    address public liquidityPool;
    
    // 开发者钱包地址
    address public devWallet;
    
    // 记录地址的交易信息
    mapping(address => uint256) public transactionCount;
    mapping(address => uint256) public lastTransactionDay;
    
    // 税收钱包
    address public rewardWallet;
    address public liquidityWallet;
    
    // 交易开关
    bool public tradingActive = false;
    
    // 已排除税收的地址
    mapping(address => bool) public isExcludedFromTax;
    
    event UpdateTaxRates(
        uint256 newBuyTax,
        uint256 newSellTax,
        uint256 newLiquidityTax,
        uint256 newRewardTax,
        uint256 newDevTax
    );
    
    event UpdateTransactionLimits(
        uint256 newMaxTransactionAmount,
        uint256 newMaxWalletAmount,
        uint256 newDailyTransactionLimit
    );
    
    event AddLiquidity(
        address indexed provider,
        uint256 amountToken,
        uint256 amountETH,
        uint256 liquidity
    );
    
    constructor(
        string memory name,
        string memory symbol,
        uint256 totalSupply,
        address _devWallet,
        address _liquidityPool
    ) ERC20(name, symbol) {
        _mint(msg.sender, totalSupply);
        
        devWallet = _devWallet;
        liquidityPool = _liquidityPool;
        
        // 创建税收钱包 - 使用合约创建者作为临时地址
        rewardWallet = msg.sender;
        liquidityWallet = msg.sender;
        
        // 计算最大交易和钱包限额 (总供应量的2%)
        maxTransactionAmount = (totalSupply * 2) / 100;
        maxWalletAmount = (totalSupply * 5) / 100;
        
        // 排除合约所有者和税收钱包的税收
        isExcludedFromTax[msg.sender] = true;
        isExcludedFromTax[devWallet] = true;
        isExcludedFromTax[rewardWallet] = true;
        isExcludedFromTax[liquidityWallet] = true;
    }
    
    /**
     * @dev 激活交易功能
     */
    function activateTrading() external onlyOwner {
        require(!tradingActive, "Trading already active");
        tradingActive = true;
    }
    
    /**
     * @dev 更新税率
     */
    function updateTaxRates(
        uint256 _buyTax,
        uint256 _sellTax,
        uint256 _liquidityTax,
        uint256 _rewardTax,
        uint256 _devTax
    ) external onlyOwner {
        require(
            _liquidityTax + _rewardTax + _devTax <= _buyTax,
            "Buy tax components exceed total buy tax"
        );
        require(
            _liquidityTax + _rewardTax + _devTax <= _sellTax,
            "Sell tax components exceed total sell tax"
        );
        
        buyTax = _buyTax;
        sellTax = _sellTax;
        liquidityTax = _liquidityTax;
        rewardTax = _rewardTax;
        devTax = _devTax;
        
        emit UpdateTaxRates(_buyTax, _sellTax, _liquidityTax, _rewardTax, _devTax);
    }
    
    /**
     * @dev 更新交易限制
     */
    function updateTransactionLimits(
        uint256 _maxTransactionAmount,
        uint256 _maxWalletAmount,
        uint256 _dailyTransactionLimit
    ) external onlyOwner {
        maxTransactionAmount = _maxTransactionAmount;
        maxWalletAmount = _maxWalletAmount;
        dailyTransactionLimit = _dailyTransactionLimit;
        
        emit UpdateTransactionLimits(
            _maxTransactionAmount,
            _maxWalletAmount,
            _dailyTransactionLimit
        );
    }
    
    /**
     * @dev 更新流动性池地址
     */
    function setLiquidityPool(address _liquidityPool) external onlyOwner {
        liquidityPool = _liquidityPool;
    }
    
    /**
     * @dev 更新开发者钱包地址
     */
    function setDevWallet(address _devWallet) external onlyOwner {
        devWallet = _devWallet;
    }
    
    /**
     * @dev 排除/包含地址的税收
     */
    function setExcludedFromTax(address _address, bool _excluded) external onlyOwner {
        isExcludedFromTax[_address] = _excluded;
    }
    
    /**
     * @dev 检查是否是卖出交易
     */
    function isSellTransaction(address to) internal view returns (bool) {
        return to == liquidityPool;
    }
    
    /**
     * @dev 检查是否是买入交易
     */
    function isBuyTransaction(address from) internal view returns (bool) {
        return from == liquidityPool;
    }
    
    /**
     * @dev 检查交易限制
     */
    function checkTransactionLimits(address sender, address recipient, uint256 amount) internal {
        // 检查交易是否已激活
        require(tradingActive || isExcludedFromTax[sender], "Trading not active");
        
        // 检查交易金额限制
        if (!isExcludedFromTax[sender] && !isExcludedFromTax[recipient]) {
            require(amount <= maxTransactionAmount, "Amount exceeds max transaction");
            
            // 检查钱包持有量限制
            if (isBuyTransaction(sender)) {
                require(
                    balanceOf(recipient) + amount <= maxWalletAmount,
                    "Wallet exceeds max holding"
                );
            }
            
            // 检查每日交易次数限制
            uint256 currentDay = block.timestamp / 1 days;
            if (lastTransactionDay[sender] != currentDay) {
                lastTransactionDay[sender] = currentDay;
                transactionCount[sender] = 0;
            }
            
            require(
                transactionCount[sender] < dailyTransactionLimit,
                "Exceeded daily transaction limit"
            );
            
            transactionCount[sender]++;
        }
    }
    
    /**
     * @dev 计算并分配税费
     */
    function calculateAndDistributeTaxes(
        address sender,
        address recipient,
        uint256 amount
    ) internal returns (uint256) {
        if (isExcludedFromTax[sender] || isExcludedFromTax[recipient]) {
            return amount; // 不征税
        }
        
        // 确定税率
        uint256 taxRate = isSellTransaction(recipient) ? sellTax : 
                         isBuyTransaction(sender) ? buyTax : 0;
                         
        if (taxRate == 0) {
            return amount; // 非买卖交易不征税
        }
        
        // 计算总税费
        uint256 totalTax = (amount * taxRate) / 100;
        uint256 taxFreeAmount = amount - totalTax;
        
        // 计算各项税费
        uint256 liquidityAmount = (amount * liquidityTax) / 100;
        uint256 rewardAmount = (amount * rewardTax) / 100;
        uint256 devAmount = (amount * devTax) / 100;
        
        // 分配税费
        if (liquidityAmount > 0) {
            _transfer(sender, liquidityWallet, liquidityAmount);
        }
        
        if (rewardAmount > 0) {
            _transfer(sender, rewardWallet, rewardAmount);
        }
        
        if (devAmount > 0) {
            _transfer(sender, devWallet, devAmount);
        }
        
        return taxFreeAmount;
    }
    
    /**
     * @dev 向流动性池添加流动性
     */
    function addLiquidity(uint256 tokenAmount) external payable nonReentrant {
        require(liquidityPool != address(0), "Liquidity pool not set");
        require(tokenAmount > 0, "Amount must be greater than 0");
        require(msg.value > 0, "ETH amount must be greater than 0");
        
        // 转移代币到流动性池
        _transfer(msg.sender, liquidityPool, tokenAmount);
        
        // 这里应该有与流动性池合约的交互，实际实现需要集成具体的DEX协议
        
        emit AddLiquidity(msg.sender, tokenAmount, msg.value, 0); // 实际流动性值需要从DEX获取
    }
    
    /**
     * @dev 从流动性池移除流动性
     */
    function removeLiquidity(uint256 liquidityAmount) external nonReentrant onlyOwner {
        require(liquidityPool != address(0), "Liquidity pool not set");
        require(liquidityAmount > 0, "Amount must be greater than 0");
        
        // 实际实现需要集成具体的DEX协议
    }
    
    /**
     * @dev 分发持有者奖励
     */
    function distributeRewards() external nonReentrant {
        uint256 rewardBalance = balanceOf(rewardWallet);
        require(rewardBalance > 0, "No rewards to distribute");
        
    }
    
    /**
     * @dev 重写transfer函数，加入税收和交易限制逻辑
     */
    function transfer(address recipient, uint256 amount) public override returns (bool) {
        // 检查交易限制
        checkTransactionLimits(_msgSender(), recipient, amount);
        
        // 计算税后金额
        uint256 taxFreeAmount = calculateAndDistributeTaxes(_msgSender(), recipient, amount);
        
        // 执行转账
        _transfer(_msgSender(), recipient, taxFreeAmount);
        return true;
    }
}
