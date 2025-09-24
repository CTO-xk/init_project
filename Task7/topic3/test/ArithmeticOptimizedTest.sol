// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "hardhat/console.sol";
import "../src/ArithmeticOptimized.sol";

contract ArithmeticOptimizedTest {
    ArithmeticOptimized public arithmetic;
    
    function setUp() public {
        arithmetic = new ArithmeticOptimized();
    }
    
    function testAdd() public {
        uint128 result = arithmetic.add(5, 3);
        console.log("5 + 3 =", result);
        assert(result == 8);
    }
    
    function testSubtract() public {
        uint128 result = arithmetic.subtract(10, 4);
        console.log("10 - 4 =", result);
        assert(result == 6);
    }
    
    function testMultiply() public {
        uint128 result = arithmetic.multiply(6, 7);
        console.log("6 * 7 =", result);
        assert(result == 42);
    }
    
    function testDivide() public {
        uint128 result = arithmetic.divide(20, 5);
        console.log("20 / 5 =", result);
        assert(result == 4);
    }
    
    function testModulo() public {
        uint128 result = arithmetic.modulo(17, 5);
        console.log("17 % 5 =", result);
        assert(result == 2);
    }
    
    function testBatchExecute() public {
        uint128[] memory a = new uint128[](3);
        uint128[] memory b = new uint128[](3);
        uint8[] memory operationTypes = new uint8[](3);
        
        a[0] = 5; b[0] = 3; operationTypes[0] = 1; // ADD
        a[1] = 10; b[1] = 4; operationTypes[1] = 2; // SUBTRACT
        a[2] = 6; b[2] = 7; operationTypes[2] = 3; // MULTIPLY
        
        uint128[] memory results = arithmetic.batchExecute(a, b, operationTypes);
        
        console.log("Batch operation results:");
        console.log("5 + 3 =", results[0]);
        console.log("10 - 4 =", results[1]);
        console.log("6 * 7 =", results[2]);
        
        assert(results[0] == 8);
        assert(results[1] == 6);
        assert(results[2] == 42);
    }
}
