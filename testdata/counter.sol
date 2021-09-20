pragma solidity ^0.8.6;

contract Counter {

    uint256 count;

    function tick() external {
        count += 1;
    }
}
