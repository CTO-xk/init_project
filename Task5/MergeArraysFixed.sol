// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract MergeArraysWithFor{
    function merge(uint[] memory arr1,uint[] memory arr2) public  pure returns (uint[] memory){
        uint  len1 =arr1.length;
        uint  len2 =arr2.length;
        uint resLen = len1+len2;
        uint[] memory result =new uint[](resLen);
        uint i = 0;
        uint j = 0;

        for(uint k =0;k<result.length;k++){
            if(i>=len1){
                result[k]=arr2[j];
                j++;
            }else if(j>=len2){
                result[k] = arr1[i];
                i++;
            }else if(arr1[i] <= arr2[j]){
                result[k]=arr1[i];
                i++;
            }else{
                result[k] = arr2[j];
                j++;
            }
        }

        return  result;
    }

    function exampleMerge() public  pure returns (uint[] memory){
        uint[] memory arr1=new uint[](4);
        arr1[0] =1;
        arr1[1] =4;
        arr1[2] =7;
        arr1[3] =9;
        uint[] memory arr2= new uint[](5);
        arr2[0] = 2;
        arr2[1] = 3;
        arr2[2] = 5;
        arr2[3] = 8;
        arr2[4] = 10;
        return merge(arr1,arr2);
    }
}