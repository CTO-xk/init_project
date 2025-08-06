// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract RomanToInteger{

    function romanToInt(string memory str) public pure  returns (uint256){
        uint256 result = 0;
        bytes memory sBytes =bytes(str);
        for (uint256 i = 0; i < sBytes.length; i++) {
            uint256 current = getValue(sBytes[i]);
            // ����������һ���ַ����ҵ�ǰֵС����һ��ֵ�����ȥ��ǰֵ
            if (i < sBytes.length - 1 && current < getValue(sBytes[i + 1])) {
                result -= current;
            } else {
                // ������ϵ�ǰֵ
                result += current;
            }
        }
        return result;
    }
    // �������������ص������������ַ���Ӧ��ֵ
    function getValue(bytes1 char) private pure returns (uint256){
        if (char == 'I') return 1;
        if (char == 'V') return 5;
        if (char == 'X') return 10;
        if (char == 'L') return 50;
        if (char == 'C') return 100;
        if (char == 'D') return 500;
        if (char == 'M') return 1000;
        return 0; // ��Ч�ַ�����0
    }
}