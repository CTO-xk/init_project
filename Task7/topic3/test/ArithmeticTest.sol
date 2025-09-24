// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "hardhat/console.sol";
import "../src/Arithmetic.sol";

contract ArithmeticTest {
    Arithmetic public arithmetic;
    
    function setUp() public {
        arithmetic = new Arithmetic();
    }
    
    function testAdd() public {
        uint256 result = arithmetic.add(5, 3);
        console.log("5 + 3 =", result);
        assert(result == 8);
    }
    
    function testSubtract() public {
        uint256 result = arithmetic.subtract(10, 4);
        console.log("10 - 4 =", result);
        assert(result == 6);
    }
    
    function testMultiply() public {
        uint256 result = arithmetic.multiply(6, 7);
        console.log("6 * 7 =", result);
        assert(result == 42);
    }
    
    function testDivide() public {
        uint256 result = arithmetic.divide(20, 5);
        console.log("20 / 5 =", result);
        assert(result == 4);
    }
    
    function testModulo() public {
        uint256 result = arithmetic.modulo(17, 5);
        console.log("17 % 5 =", result);
        assert(result == 2);
    }
}
