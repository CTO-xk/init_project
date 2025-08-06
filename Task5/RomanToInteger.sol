// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract RomanToInteger{

    function romanToInt(string memory str) public pure  returns (uint256){
        uint256 result = 0;
        bytes memory sBytes =bytes(str);
        for (uint256 i = 0; i < sBytes.length; i++) {
            uint256 current = getValue(sBytes[i]);
            // 如果不是最后一个字符，且当前值小于下一个值，则减去当前值
            if (i < sBytes.length - 1 && current < getValue(sBytes[i + 1])) {
                result -= current;
            } else {
                // 否则加上当前值
                result += current;
            }
        }
        return result;
    }
    // 辅助函数：返回单个罗马数字字符对应的值
    function getValue(bytes1 char) private pure returns (uint256){
        if (char == 'I') return 1;
        if (char == 'V') return 5;
        if (char == 'X') return 10;
        if (char == 'L') return 50;
        if (char == 'C') return 100;
        if (char == 'D') return 500;
        if (char == 'M') return 1000;
        return 0; // 无效字符返回0
    }
}