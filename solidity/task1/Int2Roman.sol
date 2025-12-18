// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

contract Int2Roman {
    function int2Roman(uint256 numb) public pure returns (string memory){
        require(numb > 0 && numb <= 3999, "Out of range");
        string[4] memory thousands = ["", "M", "MM", "MMM"];
        string[10] memory hundreds = ["", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM"];
        string[10] memory tens = ["", "X", "XX", "XXX", "XL", "L", "LX", "LXX", "LXXX", "XC"];
        string[10] memory ones = ["", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX"];
        string memory returnVal = string(abi.encodePacked(thousands[numb/1000] , hundreds[(numb%1000)/100] , tens[(numb%100)/10] , ones[numb%10]));
        return returnVal;
    }
}