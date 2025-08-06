// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ReverseString{

    function reverse(string memory _str) public  pure  returns (string memory){
        //���ַ���ת���ַ�����
         bytes memory strBytes = bytes(_str);
        //�ַ����鳤��
        bytes memory reversedBytes=new bytes(strBytes.length);
        // ��ת�ߴ��������м佻���ַ�
        for (uint256 i = 0; i < strBytes.length; i++) {
            reversedBytes[i] = strBytes[strBytes.length - 1 - i];
        }
         return string(reversedBytes);
    }
}