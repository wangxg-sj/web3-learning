// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

contract BinarySearch{

    function binarySearch(uint256[] memory sotredArray, uint target) public pure returns (int) {
        uint left = 0;
        uint256 right = sotredArray.length - 1;
        while(left <= right){
            uint midd = (left + right)/2;
            uint middVal = sotredArray[midd];
            if (middVal == target){
                return int(midd);
            }
            if (middVal < target){
                left = midd + 1;
            }else{
                right = midd - 1;
            }
        }
        return -1 ;
    }
}
