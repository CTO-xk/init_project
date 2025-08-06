// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract IntegerToRoman {
    function intToRoman(uint256 num) public pure returns (string memory) {
        // 显式指定每个元素为uint256类型，匹配固定数组类型
        uint256[13] memory fixedValues = [
            uint256(1000), uint256(900), uint256(500), uint256(400),
            uint256(100), uint256(90), uint256(50), uint256(40),
            uint256(10), uint256(9), uint256(5), uint256(4),
            uint256(1)
        ];
        
        // 将固定大小数组转换为动态数组
        uint256[] memory values = new uint256[](13);
        for (uint256 i = 0; i < 13; i++) {
            values[i] = fixedValues[i];
        }
        
        string[13] memory fixedSymbols = [
            "M", "CM", "D", "CD",
            "C", "XC", "L", "XL",
            "X", "IX", "V", "IV",
            "I"
        ];
        
        // 符号数组转换为动态数组
        string[] memory symbols = new string[](13);
        for (uint256 i = 0; i < 13; i++) {
            symbols[i] = fixedSymbols[i];
        }
        
        bytes memory result = new bytes(0);
        
        for (uint256 i = 0; i < values.length; i++) {
            while (values[i] <= num) {
                result = abi.encodePacked(result, symbols[i]);
                num -= values[i];
            }
            if (num == 0) break;
        }
        
        return string(result);
    }
}
