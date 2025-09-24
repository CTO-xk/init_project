// SPDX-License-Identifier: MIT
pragma solidity ^0.8.2;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/security/PausableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

contract StakeSystem is 
    Initializable, 
    UUPSUpgradeable, 
    OwnableUpgradeable, 
    PausableUpgradeable, 
    ReentrancyGuardUpgradeable 
{
    using SafeERC20 for IERC20;
    // 奖励代币
    IERC20 public metaNodeToken;
    // 质押奖励率 - 每个区块产生的奖励
    uint256 public rewardPerBlock;
    // 管理员角色
    address public operator;
    // 质押池结构体
    struct Pool {
        address stTokenAddress;  // 质押代币地址，0x0表示原生代币
        uint256 poolWeight;      // 池权重
        uint256 lastRewardBlock; // 最后一次计算奖励的区块
        uint256 accMetaNodePerST; // 每个质押代币累积的奖励
        uint256 stTokenAmount;   // 池中总质押量
        uint256 minDepositAmount;// 最小质押数量
        uint256 unstakeLockedBlocks; // 解质押锁定期(区块数)
    }
     // 质押池数组
    Pool[] public pools;
    // 用户质押信息
    struct UserInfo {
        uint256 amount; // 质押数量
        uint256 rewardDebt; // 奖励债务
    }
    // 解质押请求
    struct UnstakeRequest {
        uint256 amount;          // 解质押数量
        uint256 releaseBlock;    // 可提取区块
        bool claimed;            // 是否已提取
    }
    // 用户在每个池中的质押信息: poolId => user => info
    mapping(uint256 => mapping(address => UserInfo)) public userInfo;
    // 用户的解质押请求: poolId => user => requests
    mapping(uint256 => mapping(address => UnstakeRequest[])) public unstakeRequests;
    //权限修饰符
    modifier onlyOperator() {
        require(msg.sender == operator || msg.sender == owner(), "Only operator can call this function");
        _;
    }
  // 事件定义
    event Staked(address indexed user, uint256 indexed pid, uint256 amount);
    event UnstakeRequested(address indexed user, uint256 indexed pid, uint256 amount, uint256 releaseBlock);
    event UnstakeClaimed(address indexed user, uint256 indexed pid, uint256 amount, uint256 requestIndex);
    event RewardClaimed(address indexed user, uint256 indexed pid, uint256 amount);
    event PoolAdded(uint256 indexed pid, address stToken, uint256 weight);
    event PoolUpdated(uint256 indexed pid, address stToken, uint256 weight);
    event OperatorChanged(address indexed oldOperator, address indexed newOperator);
    event RewardRateUpdated(uint256 oldRate, uint256 newRate);

    // 初始化函数
    function initialize(address _metaNodeToken, uint256 _rewardPerBlock, address _operator) external initializer {
        __Ownable_init();
        __Pausable_init();
        __ReentrancyGuard_init();
      require(_metaNodeToken != address(0), "Invalid metaNodeToken address");
      require(_rewardPerBlock > 0, "Invalid rewardPerBlock");
      require(_operator != address(0), "Invalid operator address");
      metaNodeToken = IERC20(_metaNodeToken);
      rewardPerBlock = _rewardPerBlock;
      operator = _operator;
      // 添加第一个质押池 - 原生代币
      pools.push(Pool({
        stTokenAddress: address(0),
        poolWeight: 100,
        lastRewardBlock: block.number,
        accMetaNodePerST: 0,
        stTokenAmount: 0,
        minDepositAmount: 1e16,  // 0.01 ETH
        unstakeLockedBlocks: 200  // 约10分钟(假设15秒/块)
      }));
      emit PoolAdded(0, address(0), 100);
    }
    // 计算每个池的奖励
    function updatePool(uint256 _pid) public{
        require(_pid < pools.length, "Invalid pool ID");
        Pool storage pool = pools[_pid];
        if(block.number <= pool.lastRewardBlock){
            return;
        }
        if(pool.stTokenAmount == 0){
            pool.lastRewardBlock = block.number;
            return;
        }
        // 计算区块奖励
        uint256 blockReward = rewardPerBlock * pool.poolWeight / 100;
        uint256 multiplier = block.number - pool.lastRewardBlock;
        uint256 reward = multiplier * blockReward;
        // 更新累积奖励
        pool.accMetaNodePerST += reward * 1e18 / pool.stTokenAmount;
        pool.lastRewardBlock = block.number;
    }
    // 质押
    function stake(uint256 _pid, uint256 _amount) external payable nonReentrant whenNotPaused {
        require(_pid < pools.length, "StakeSystem: invalid pool ID");
        Pool storage pool = pools[_pid];
        require(_amount >= pool.minDepositAmount, "StakeSystem: amount too small");
        
        UserInfo storage user = userInfo[_pid][msg.sender];
        
        // 更新奖励
        updatePool(_pid);
        
        // 如果已有质押，先计算奖励
        if (user.amount > 0) {
            uint256 pending = user.amount * pool.accMetaNodePerST / 1e18 - user.rewardDebt;
            if (pending > 0) {
                metaNodeToken.safeTransfer(msg.sender, pending);
                emit RewardClaimed(msg.sender, _pid, pending);
            }
        }
        
        // 处理质押
        if (pool.stTokenAddress == address(0)) {
            // 原生代币
            require(msg.value == _amount, "StakeSystem: wrong ETH amount");
        } else {
            // ERC20代币
            require(msg.value == 0, "StakeSystem: no ETH allowed");
            IERC20(pool.stTokenAddress).safeTransferFrom(msg.sender, address(this), _amount);
        }
        
        // 更新用户和池信息
        user.amount += _amount;
        user.rewardDebt = user.amount * pool.accMetaNodePerST / 1e18;
        pool.stTokenAmount += _amount;
        
        emit Staked(msg.sender, _pid, _amount);
    }
    // 申请解质押
    function requestUnstake(uint256 _pid, uint256 _amount) 
        external 
        nonReentrant 
        whenNotPaused 
    {
        require(_pid < pools.length, "StakeSystem: invalid pool ID");
        Pool storage pool = pools[_pid];
        UserInfo storage user = userInfo[_pid][msg.sender];
        
        require(user.amount >= _amount, "StakeSystem: insufficient balance");
        
        // 更新奖励
        updatePool(_pid);
        
        // 计算并发放奖励
        uint256 pending = user.amount * pool.accMetaNodePerST / 1e18 - user.rewardDebt;
        if (pending > 0) {
            metaNodeToken.safeTransfer(msg.sender, pending);
            emit RewardClaimed(msg.sender, _pid, pending);
        }
        
        // 创建解质押请求
        uint256 releaseBlock = block.number + pool.unstakeLockedBlocks;
        unstakeRequests[_pid][msg.sender].push(UnstakeRequest({
            amount: _amount,
            releaseBlock: releaseBlock,
            claimed: false
        }));
        
        // 更新用户和池信息
        user.amount -= _amount;
        user.rewardDebt = user.amount * pool.accMetaNodePerST / 1e18;
        pool.stTokenAmount -= _amount;
        
        emit UnstakeRequested(msg.sender, _pid, _amount, releaseBlock);
    }
    
    // 提取已解锁的解质押资产
    function claimUnstake(uint256 _pid, uint256 _requestIndex) 
        external 
        nonReentrant 
        whenNotPaused 
    {
        require(_pid < pools.length, "StakeSystem: invalid pool ID");
        Pool storage pool = pools[_pid];
        UnstakeRequest[] storage requests = unstakeRequests[_pid][msg.sender];
        
        require(_requestIndex < requests.length, "StakeSystem: invalid request index");
        UnstakeRequest storage request = requests[_requestIndex];
        require(!request.claimed, "StakeSystem: already claimed");
        require(block.number >= request.releaseBlock, "StakeSystem: still locked");
        
        // 标记为已提取
        request.claimed = true;
        
        // 转移资产给用户
        if (pool.stTokenAddress == address(0)) {
            // 原生代币
           payable(msg.sender).transfer(request.amount);
        } else {
            // ERC20代币
            IERC20(pool.stTokenAddress).safeTransfer(msg.sender, request.amount);
        }
        
        emit UnstakeClaimed(msg.sender, _pid, request.amount, _requestIndex);
    }
    
    // 领取奖励
    function claimReward(uint256 _pid) 
        external 
        nonReentrant 
        whenNotPaused 
    {
        require(_pid < pools.length, "StakeSystem: invalid pool ID");
        Pool storage pool = pools[_pid];
        UserInfo storage user = userInfo[_pid][msg.sender];
        
        require(user.amount > 0, "StakeSystem: no staked amount");
        
        // 更新奖励
        updatePool(_pid);
        
        // 计算并发放奖励
        uint256 pending = user.amount * pool.accMetaNodePerST / 1e18 - user.rewardDebt;
        require(pending > 0, "StakeSystem: no reward");
        
        // 先更新状态，防止重入攻击
        user.rewardDebt = user.amount * pool.accMetaNodePerST / 1e18;
        metaNodeToken.safeTransfer(msg.sender, pending);
        
        emit RewardClaimed(msg.sender, _pid, pending);
    }
    
    // 添加新质押池
    function addPool(
        address _stTokenAddress,
        uint256 _poolWeight,
        uint256 _minDepositAmount,
        uint256 _unstakeLockedBlocks
    ) external onlyOperator {
        require(_stTokenAddress != address(metaNodeToken), "StakeSystem: cannot stake reward token");
        require(_poolWeight > 0, "StakeSystem: weight must be positive");
        require(_minDepositAmount > 0, "StakeSystem: min deposit must be positive");
        require(_unstakeLockedBlocks > 0, "StakeSystem: lock period must be positive");
        
        // 更新所有现有池的奖励
        for (uint256 i = 0; i < pools.length; i++) {
            updatePool(i);
        }
        
        uint256 pid = pools.length;
        pools.push(Pool({
            stTokenAddress: _stTokenAddress,
            poolWeight: _poolWeight,
            lastRewardBlock: block.number,
            accMetaNodePerST: 0,
            stTokenAmount: 0,
            minDepositAmount: _minDepositAmount,
            unstakeLockedBlocks: _unstakeLockedBlocks
        }));
        
        emit PoolAdded(pid, _stTokenAddress, _poolWeight);
    }
    
    // 更新质押池
    function updatePoolSettings(
        uint256 _pid,
        uint256 _poolWeight,
        uint256 _minDepositAmount,
        uint256 _unstakeLockedBlocks
    ) external onlyOperator {
        require(_pid < pools.length, "StakeSystem: invalid pool ID");
        require(_poolWeight > 0, "StakeSystem: weight must be positive");
        require(_minDepositAmount > 0, "StakeSystem: min deposit must be positive");
        require(_unstakeLockedBlocks > 0, "StakeSystem: lock period must be positive");
        
        // 更新池奖励
        updatePool(_pid);
        
        Pool storage pool = pools[_pid];
        pool.poolWeight = _poolWeight;
        pool.minDepositAmount = _minDepositAmount;
        pool.unstakeLockedBlocks = _unstakeLockedBlocks;
        
        emit PoolUpdated(_pid, pool.stTokenAddress, _poolWeight);
    }
    
    // 获取池数量
    function poolLength() external view returns (uint256) {
        return pools.length;
    }
    
    // 计算用户可领取的奖励
    function pendingReward(uint256 _pid, address _user) external view returns (uint256) {
        require(_pid < pools.length, "StakeSystem: invalid pool ID");
        
        Pool storage pool = pools[_pid];
        UserInfo storage user = userInfo[_pid][_user];
        
        if (user.amount == 0) {
            return 0;
        }
        
        uint256 accMetaNodePerST = pool.accMetaNodePerST;
        uint256 stTokenAmount = pool.stTokenAmount;
        
        if (block.number > pool.lastRewardBlock && stTokenAmount > 0) {
            uint256 blockReward = rewardPerBlock * pool.poolWeight / 100;
            uint256 multiplier = block.number - pool.lastRewardBlock;
            uint256 reward = multiplier * blockReward;
            accMetaNodePerST += reward * 1e18 / stTokenAmount;
        }
        
        return user.amount * accMetaNodePerST / 1e18 - user.rewardDebt;
    }
    
    // 获取用户的解质押请求
    function getUserUnstakeRequests(uint256 _pid, address _user) 
        external 
        view 
        returns (UnstakeRequest[] memory) 
    {
        return unstakeRequests[_pid][_user];
    }
    
    // 更改操作员
    function setOperator(address _newOperator) external onlyOwner {
        require(_newOperator != address(0), "StakeSystem: invalid operator");
        emit OperatorChanged(operator, _newOperator);
        operator = _newOperator;
    }
    
    // 更新奖励率
    function updateRewardPerBlock(uint256 _newRewardPerBlock) external onlyOperator {
        require(_newRewardPerBlock > 0, "StakeSystem: reward must be positive");
        
        // 更新所有池的奖励
        for (uint256 i = 0; i < pools.length; i++) {
            updatePool(i);
        }
        
        emit RewardRateUpdated(rewardPerBlock, _newRewardPerBlock);
        rewardPerBlock = _newRewardPerBlock;
    }
    
    // 暂停功能
    function pause() external onlyOperator {
        _pause();
    }
    
    // 恢复功能
    function unpause() external onlyOperator {
        _unpause();
    }
    
    // 接收ETH
    receive() external payable {
        // 仅接受原生代币质押的转账
        bool isNativePool = false;
        for (uint256 i = 0; i < pools.length; i++) {
            if (pools[i].stTokenAddress == address(0)) {
                isNativePool = true;
                break;
            }
        }
        require(isNativePool, "StakeSystem: no native token pool");
    }
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}