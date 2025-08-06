// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract BinarySearch{
    function binarySearch(uint[] memory arr, uint target) public pure returns (int) {
        uint left = 0;
        uint right = arr.length - 1;
                while (left <= right) {
            uint mid = left + (right - left) / 2;
                        if (arr[mid] == target) {
                return int(mid);
            }
            else if (arr[mid] < target) {
                left = mid + 1;
            }
            else {
                right = mid - 1;
            }
        }
        return -1;
    }

    function exampleSearch() public  pure  returns (int){
        uint[] memory arr1=new uint[](5);
        arr1[0] = 2;
        arr1[1] = 3;
        arr1[2] = 5;
        arr1[3] = 8;
        arr1[4] = 10;
        uint target = 8;
        return binarySearch(arr1,target);
    } 

    
}