// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract ToWatch {
    uint256 public num;

    event NumberChanged(uint256 indexed newNum);

    function setNum(uint256 _num) external {
        num = _num;
        emit NumberChanged(_num);
    }
}