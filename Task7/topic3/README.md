# 智能合约 Gas 优化项目

## 项目简介

本项目演示了如何使用 Foundry 框架进行智能合约开发、测试和 Gas 优化。项目包含两个版本的算术运算智能合约：原始版本和优化版本，通过对比分析展示了不同的 Gas 优化策略的效果。

## 项目结构

```
topic3/
├── src/                           # 智能合约源码
│   ├── Arithmetic.sol            # 原始算术运算合约
│   └── ArithmeticOptimized.sol   # 优化后的算术运算合约
├── test/                         # 测试文件
│   ├── Arithmetic.t.sol         # 原始合约测试
│   └── ArithmeticOptimized.t.sol # 优化合约测试
├── script/                       # 部署脚本
│   └── Deploy.s.sol             # 合约部署脚本
├── docs/                         # 文档
│   ├── Foundry_Framework_Guide.md # Foundry 框架指南
│   └── Gas_Analysis_Report.md   # Gas 分析报告
├── foundry.toml                  # Foundry 配置文件
├── package.json                  # 项目依赖配置
└── README.md                     # 项目说明文档
```

## 功能特性

### 原始合约 (Arithmetic.sol)
- ✅ 基本算术运算：加法、减法、乘法、除法、模运算
- ✅ 溢出和下溢保护
- ✅ 操作历史记录
- ✅ 事件记录
- ✅ 完整的错误处理

### 优化合约 (ArithmeticOptimized.sol)
- ✅ 所有原始合约功能
- ✅ Gas 优化：数据类型优化、存储结构优化
- ✅ 批量操作支持
- ✅ 更高效的事件索引
- ✅ 紧凑的数据结构

## 安装和运行

### 环境要求
- Node.js 18+
- Foundry (Forge, Cast, Anvil)

### 安装依赖
```bash
# 安装 Foundry
curl -L https://foundry.paradigm.xyz | bash
foundryup

# 安装项目依赖
npm install
```

### 编译合约
```bash
# 编译所有合约
forge build

# 编译特定合约
forge build --contracts src/Arithmetic.sol
```

### 运行测试
```bash
# 运行所有测试
forge test

# 运行特定测试文件
forge test --match-path test/Arithmetic.t.sol

# 运行测试并显示 Gas 报告
forge test --gas-report

# 运行测试并显示详细输出
forge test -vvv

# 运行模糊测试
forge test --fuzz-runs 1000
```

### 启动本地网络
```bash
# 启动 Anvil 本地网络
anvil

# 自定义配置启动
anvil --port 8545 --accounts 10 --balance 1000
```

### 部署合约
```bash
# 设置环境变量
export PRIVATE_KEY=your_private_key
export RPC_URL=your_rpc_url

# 部署到本地网络
forge script script/Deploy.s.sol:DeployScript \
    --rpc-url http://localhost:8545 \
    --broadcast \
    -vvvv
```

## Gas 优化策略

### 1. 数据类型优化
- 使用 `uint128` 替代 `uint256`
- 使用 `uint8` 常量替代 `string`
- 使用 `uint32` 替代 `uint256` 时间戳

### 2. 存储结构优化
- 使用 `mapping` 替代动态数组
- 优化结构体字段顺序
- 减少存储槽位使用

### 3. 批量操作优化
- 支持批量执行多个操作
- 减少多次交易的 Gas 成本
- 使用 `calldata` 减少内存分配

## 测试覆盖

### 测试类型
- ✅ 单元测试
- ✅ 边界条件测试
- ✅ 溢出/下溢测试
- ✅ 事件测试
- ✅ 模糊测试
- ✅ Gas 消耗测试

### 测试命令示例
```bash
# 运行特定测试函数
forge test --match-test testAdd

# 运行 Gas 消耗测试
forge test --match-test testGasConsumption

# 运行模糊测试
forge test --match-test testFuzz_Add
```

## Gas 分析结果

### 单次操作 Gas 节省
- **加法**: 28.9% 节省
- **减法**: 28.9% 节省
- **乘法**: 27.1% 节省
- **除法**: 28.6% 节省
- **模运算**: 27.9% 节省

### 批量操作 Gas 节省
- **3 次操作**: 37.0% 节省
- **5 次操作**: 40.0% 节省
- **10 次操作**: 44.4% 节省

## 开发指南

### 添加新功能
1. 在 `src/` 目录下创建新的合约文件
2. 在 `test/` 目录下创建对应的测试文件
3. 在 `script/` 目录下更新部署脚本
4. 运行测试确保功能正常

### 代码规范
- 使用 Solidity 0.8.19 或更高版本
- 遵循 Solidity 编码规范
- 添加完整的 NatSpec 注释
- 包含适当的错误处理

### 测试规范
- 每个函数至少有一个测试用例
- 包含边界条件和错误情况测试
- 使用模糊测试覆盖更多场景
- 测量和记录 Gas 消耗

## 故障排除

### 常见问题

**1. 编译错误**
```bash
# 清理编译缓存
forge clean
# 重新编译
forge build
```

**2. 测试失败**
```bash
# 运行测试并显示详细错误
forge test -vvv
# 检查测试环境配置
```

**3. Gas 消耗异常**
```bash
# 生成 Gas 报告
forge test --gas-report
# 检查优化器设置
```

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License

## 联系方式

如有问题或建议，请提交 Issue 或 Pull Request。

---

**注意**: 本项目仅用于学习和演示目的，请勿在生产环境中直接使用。在生产环境中使用前，请进行充分的测试和安全审计。
