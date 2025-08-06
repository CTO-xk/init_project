// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ReverseString{

    function reverse(string memory _str) public  pure  returns (string memory){
        //将字符串转成字符数组
         bytes memory strBytes = bytes(_str);
        //字符数组长度
        bytes memory reversedBytes=new bytes(strBytes.length);
        // 反转逻从两端向中间交换字符
        for (uint256 i = 0; i < strBytes.length; i++) {
            reversedBytes[i] = strBytes[strBytes.length - 1 - i];
        }
         return string(reversedBytes);
    }
}