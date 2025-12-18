
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

contract MergeSortedArray{

    function merge(uint256[] memory arr1 , uint256[] memory arr2) public pure returns (uint [] memory){
        uint len1 = arr1.length;
        uint len2 = arr2.length;
        require(len1 > 0, "arr1 is empty");
        require(len2 > 0, "arr2 is empty");
        uint256[] memory returnVal = new uint256[](len1 + len2);
        uint256 i = 0 ; 
        uint256 j = 0;
        uint256 k = 0;
        for(; i < len1 && j < len2; ){
            if (arr1[i] > arr2[j]){
                returnVal[k] = arr1[j];
                ++j;
            }else if(arr1[i] < arr2[j]){
                 returnVal[k] = arr1[i];
                 ++i;
            }else{
                returnVal[k] = arr1[i];
                returnVal[++k] = arr1[i];
                ++i;
                ++j;
            }
            ++k;
        }

        while (i < len1) {
            returnVal[k] = arr1[i];
            i++;
            k++;
        }

        while (j < len2) {
            returnVal[k] = arr2[j];
            j++;
            k++;
        }
        return returnVal;
    }
}
