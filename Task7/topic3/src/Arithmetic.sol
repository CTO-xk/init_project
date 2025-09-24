// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title Arithmetic
 * @dev 基本算术运算智能合约
 * @notice 提供加法、减法、乘法、除法等基本运算功能
 */
contract Arithmetic {
    // 存储计算结果
    uint256 public lastResult;
    
    // 存储操作历史
    struct Operation {
        uint256 a;
        uint256 b;
        string operation;
        uint256 result;
        uint256 timestamp;
    }
    
    Operation[] public operations;
    
    // 事件声明
    event CalculationPerformed(
        uint256 indexed a,
        uint256 indexed b,
        string operation,
        uint256 result,
        uint256 timestamp
    );
    
    event ErrorOccurred(string message);
    
    /**
     * @dev 执行加法运算
     * @param a 第一个操作数
     * @param b 第二个操作数
     * @return 加法结果
     */
    function add(uint256 a, uint256 b) public returns (uint256) {
        uint256 result = a + b;
        
        // 检查溢出
        require(result >= a, "Arithmetic: addition overflow");
        
        _updateState(a, b, "add", result);
        return result;
    }
    
    /**
     * @dev 执行减法运算
     * @param a 被减数
     * @param b 减数
     * @return 减法结果
     */
    function subtract(uint256 a, uint256 b) public returns (uint256) {
        require(b <= a, "Arithmetic: subtraction underflow");
        
        uint256 result = a - b;
        _updateState(a, b, "subtract", result);
        return result;
    }
    
    /**
     * @dev 执行乘法运算
     * @param a 第一个操作数
     * @param b 第二个操作数
     * @return 乘法结果
     */
    function multiply(uint256 a, uint256 b) public returns (uint256) {
        if (a == 0 || b == 0) {
            _updateState(a, b, "multiply", 0);
            return 0;
        }
        
        uint256 result = a * b;
        require(result / a == b, "Arithmetic: multiplication overflow");
        
        _updateState(a, b, "multiply", result);
        return result;
    }
    
    /**
     * @dev 执行除法运算
     * @param a 被除数
     * @param b 除数
     * @return 除法结果
     */
    function divide(uint256 a, uint256 b) public returns (uint256) {
        require(b > 0, "Arithmetic: division by zero");
        
        uint256 result = a / b;
        _updateState(a, b, "divide", result);
        return result;
    }
    
    /**
     * @dev 执行模运算
     * @param a 被除数
     * @param b 除数
     * @return 模运算结果
     */
    function modulo(uint256 a, uint256 b) public returns (uint256) {
        require(b > 0, "Arithmetic: modulo by zero");
        
        uint256 result = a % b;
        _updateState(a, b, "modulo", result);
        return result;
    }
    
    /**
     * @dev 获取操作历史数量
     * @return 操作历史总数
     */
    function getOperationsCount() public view returns (uint256) {
        return operations.length;
    }
    
    /**
     * @dev 获取指定索引的操作历史
     * @param index 索引
     * @return a 第一个操作数
     * @return b 第二个操作数
     * @return operation 操作类型
     * @return result 计算结果
     * @return timestamp 时间戳
     */
    function getOperation(uint256 index) public view returns (
        uint256 a,
        uint256 b,
        string memory operation,
        uint256 result,
        uint256 timestamp
    ) {
        require(index < operations.length, "Arithmetic: index out of bounds");
        
        Operation memory op = operations[index];
        return (op.a, op.b, op.operation, op.result, op.timestamp);
    }
    
    /**
     * @dev 清除操作历史
     */
    function clearHistory() public {
        delete operations;
        lastResult = 0;
    }
    
    /**
     * @dev 内部函数：更新合约状态
     * @param a 第一个操作数
     * @param b 第二个操作数
     * @param operation 操作类型
     * @param result 计算结果
     */
    function _updateState(
        uint256 a,
        uint256 b,
        string memory operation,
        uint256 result
    ) internal {
        lastResult = result;
        
        operations.push(Operation({
            a: a,
            b: b,
            operation: operation,
            result: result,
            timestamp: block.timestamp
        }));
        
        emit CalculationPerformed(a, b, operation, result, block.timestamp);
    }
}
