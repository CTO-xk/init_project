// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "hardhat/console.sol";
import "../src/Arithmetic.sol";
import "../src/ArithmeticOptimized.sol";

contract GasAnalysisTest {
    Arithmetic public arithmetic;
    ArithmeticOptimized public arithmeticOptimized;
    
    function setUp() public {
        arithmetic = new Arithmetic();
        arithmeticOptimized = new ArithmeticOptimized();
    }
    
    function testGasComparison() public {
        console.log("=== Gas Consumption Comparison ===");
        
        // Test addition
        uint256 gasBefore = gasleft();
        arithmetic.add(100, 200);
        uint256 gasUsedOriginal = gasBefore - gasleft();
        
        gasBefore = gasleft();
        arithmeticOptimized.add(100, 200);
        uint256 gasUsedOptimized = gasBefore - gasleft();
        
        console.log("Original contract addition gas:", gasUsedOriginal);
        console.log("Optimized contract addition gas:", gasUsedOptimized);
        console.log("Gas savings:", gasUsedOriginal - gasUsedOptimized);
        console.log("Savings percentage:", ((gasUsedOriginal - gasUsedOptimized) * 100) / gasUsedOriginal, "%");
        
        // Verify optimization effect
        assert(gasUsedOptimized < gasUsedOriginal);
    }
    
    function testBatchOperationGasEfficiency() public {
        console.log("=== Batch Operation Gas Efficiency ===");
        
        // Test 3 individual operations
        uint256 totalSingleGas = 0;
        for (uint256 i = 0; i < 3; i++) {
            uint256 gasBefore = gasleft();
            arithmeticOptimized.add(1, 2);
            totalSingleGas += gasBefore - gasleft();
        }
        
        // Test batch operation
        uint256 gasBefore = gasleft();
        uint128[] memory a = new uint128[](3);
        uint128[] memory b = new uint128[](3);
        uint8[] memory operationTypes = new uint8[](3);
        
        for (uint256 i = 0; i < 3; i++) {
            a[i] = 1;
            b[i] = 2;
            operationTypes[i] = 1; // ADD
        }
        
        arithmeticOptimized.batchExecute(a, b, operationTypes);
        uint256 batchGas = gasBefore - gasleft();
        
        console.log("3 individual operations total gas:", totalSingleGas);
        console.log("Batch operation total gas:", batchGas);
        console.log("Average gas per individual operation:", totalSingleGas / 3);
        console.log("Average gas per batch operation:", batchGas / 3);
        console.log("Gas savings per operation:", (totalSingleGas / 3) - (batchGas / 3));
        console.log("Batch efficiency percentage:", (((totalSingleGas / 3) - (batchGas / 3)) * 100) / (totalSingleGas / 3), "%");
        
        // Verify batch operation efficiency
        assert(batchGas / 3 < totalSingleGas / 3);
    }
    
    function testAllOperationsGasComparison() public {
        console.log("=== All Operations Gas Comparison ===");
        
        uint256 totalOriginalGas = 0;
        uint256 totalOptimizedGas = 0;
        
        // Test all operations
        uint256 gasBefore = gasleft();
        arithmetic.add(100, 200);
        totalOriginalGas += gasBefore - gasleft();
        
        gasBefore = gasleft();
        arithmeticOptimized.add(100, 200);
        totalOptimizedGas += gasBefore - gasleft();
        
        gasBefore = gasleft();
        arithmetic.subtract(200, 100);
        totalOriginalGas += gasBefore - gasleft();
        
        gasBefore = gasleft();
        arithmeticOptimized.subtract(200, 100);
        totalOptimizedGas += gasBefore - gasleft();
        
        gasBefore = gasleft();
        arithmetic.multiply(50, 60);
        totalOriginalGas += gasBefore - gasleft();
        
        gasBefore = gasleft();
        arithmeticOptimized.multiply(50, 60);
        totalOptimizedGas += gasBefore - gasleft();
        
        gasBefore = gasleft();
        arithmetic.divide(200, 50);
        totalOriginalGas += gasBefore - gasleft();
        
        gasBefore = gasleft();
        arithmeticOptimized.divide(200, 50);
        totalOptimizedGas += gasBefore - gasleft();
        
        gasBefore = gasleft();
        arithmetic.modulo(157, 50);
        totalOriginalGas += gasBefore - gasleft();
        
        gasBefore = gasleft();
        arithmeticOptimized.modulo(157, 50);
        totalOptimizedGas += gasBefore - gasleft();
        
        console.log("Total original contract gas:", totalOriginalGas);
        console.log("Total optimized contract gas:", totalOptimizedGas);
        console.log("Total gas savings:", totalOriginalGas - totalOptimizedGas);
        console.log("Total savings percentage:", ((totalOriginalGas - totalOptimizedGas) * 100) / totalOriginalGas, "%");
        
        // Verify overall optimization effect
        assert(totalOptimizedGas < totalOriginalGas);
    }
}
