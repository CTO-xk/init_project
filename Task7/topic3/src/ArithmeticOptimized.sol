// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title ArithmeticOptimized
 * @dev 优化后的基本算术运算智能合约
 * @notice 应用多种 Gas 优化策略，减少 Gas 消耗
 */
contract ArithmeticOptimized {
    // 使用 uint128 减少存储成本
    uint128 public lastResult;
    
    // 优化：使用紧凑的结构体减少存储槽位
    struct Operation {
        uint128 a;
        uint128 b;
        uint8 operationType; // 使用 uint8 替代 string，节省大量 Gas
        uint128 result;
        uint32 timestamp; // 使用 uint32 足够表示时间戳
    }
    
    // 优化：使用 mapping 替代数组，减少遍历成本
    mapping(uint256 => Operation) public operations;
    uint256 public operationsCount;
    
    // 优化：使用事件索引减少 Gas 消耗
    event CalculationPerformed(
        uint128 indexed a,
        uint128 indexed b,
        uint8 indexed operationType,
        uint128 result,
        uint32 timestamp
    );
    
    // 操作类型常量，避免重复的字符串存储
    uint8 constant ADD = 1;
    uint8 constant SUBTRACT = 2;
    uint8 constant MULTIPLY = 3;
    uint8 constant DIVIDE = 4;
    uint8 constant MODULO = 5;
    
    /**
     * @dev 执行加法运算
     * @param a 第一个操作数
     * @param b 第二个操作数
     * @return 加法结果
     */
    function add(uint128 a, uint128 b) public returns (uint128) {
        // 优化：内联溢出检查，减少函数调用
        uint128 result = a + b;
        require(result >= a, "Overflow");
        
        _updateState(a, b, ADD, result);
        return result;
    }
    
    /**
     * @dev 执行减法运算
     * @param a 被减数
     * @param b 减数
     * @return 减法结果
     */
    function subtract(uint128 a, uint128 b) public returns (uint128) {
        require(b <= a, "Underflow");
        
        uint128 result = a - b;
        _updateState(a, b, SUBTRACT, result);
        return result;
    }
    
    /**
     * @dev 执行乘法运算
     * @param a 第一个操作数
     * @param b 第二个操作数
     * @return 乘法结果
     */
    function multiply(uint128 a, uint128 b) public returns (uint128) {
        // 优化：早期返回，避免不必要的计算
        if (a == 0 || b == 0) {
            _updateState(a, b, MULTIPLY, 0);
            return 0;
        }
        
        uint128 result = a * b;
        require(result / a == b, "Overflow");
        
        _updateState(a, b, MULTIPLY, result);
        return result;
    }
    
    /**
     * @dev 执行除法运算
     * @param a 被除数
     * @param b 除数
     * @return 除法结果
     */
    function divide(uint128 a, uint128 b) public returns (uint128) {
        require(b > 0, "Div by zero");
        
        uint128 result = a / b;
        _updateState(a, b, DIVIDE, result);
        return result;
    }
    
    /**
     * @dev 执行模运算
     * @param a 被除数
     * @param b 除数
     * @return 模运算结果
     */
    function modulo(uint128 a, uint128 b) public returns (uint128) {
        require(b > 0, "Mod by zero");
        
        uint128 result = a % b;
        _updateState(a, b, MODULO, result);
        return result;
    }
    
    /**
     * @dev 获取操作历史数量
     * @return 操作历史总数
     */
    function getOperationsCount() public view returns (uint256) {
        return operationsCount;
    }
    
    /**
     * @dev 获取指定索引的操作历史
     * @param index 索引
     * @return a 第一个操作数
     * @return b 第二个操作数
     * @return operationType 操作类型
     * @return result 计算结果
     * @return timestamp 时间戳
     */
    function getOperation(uint256 index) public view returns (
        uint128 a,
        uint128 b,
        uint8 operationType,
        uint128 result,
        uint32 timestamp
    ) {
        require(index < operationsCount, "Out of bounds");
        
        Operation memory op = operations[index];
        return (op.a, op.b, op.operationType, op.result, op.timestamp);
    }
    
    /**
     * @dev 清除操作历史
     */
    function clearHistory() public {
        // 优化：只重置计数器，保留映射数据（节省 Gas）
        operationsCount = 0;
        lastResult = 0;
    }
    
    /**
     * @dev 内部函数：更新合约状态（优化版本）
     * @param a 第一个操作数
     * @param b 第二个操作数
     * @param operationType 操作类型
     * @param result 计算结果
     */
    function _updateState(
        uint128 a,
        uint128 b,
        uint8 operationType,
        uint128 result
    ) internal {
        lastResult = result;
        
        // 优化：直接存储到映射，避免数组操作
        operations[operationsCount] = Operation({
            a: a,
            b: b,
            operationType: operationType,
            result: result,
            timestamp: uint32(block.timestamp)
        });
        
        operationsCount++;
        
        emit CalculationPerformed(a, b, operationType, result, uint32(block.timestamp));
    }
    
    /**
     * @dev 批量执行多个操作（Gas 优化：减少多次交易）
     * @param a 第一个操作数数组
     * @param b 第二个操作数数组
     * @param operationTypes 操作类型数组
     * @return results 结果数组
     */
    function batchExecute(
        uint128[] calldata a,
        uint128[] calldata b,
        uint8[] calldata operationTypes
    ) public returns (uint128[] memory results) {
        require(
            a.length == b.length && b.length == operationTypes.length,
            "Length mismatch"
        );
        
        results = new uint128[](a.length);
        
        for (uint256 i = 0; i < a.length; i++) {
            if (operationTypes[i] == ADD) {
                results[i] = add(a[i], b[i]);
            } else if (operationTypes[i] == SUBTRACT) {
                results[i] = subtract(a[i], b[i]);
            } else if (operationTypes[i] == MULTIPLY) {
                results[i] = multiply(a[i], b[i]);
            } else if (operationTypes[i] == DIVIDE) {
                results[i] = divide(a[i], b[i]);
            } else if (operationTypes[i] == MODULO) {
                results[i] = modulo(a[i], b[i]);
            }
        }
        
        return results;
    }
}
