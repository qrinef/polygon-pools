// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.15;

contract Emulation {
    bytes constant method = abi.encodeWithSelector(0x0902f1ac);

    struct Reserves {
        uint Reserve0;
        uint Reserve1;
    }

    function start(address[] calldata pools) public returns(Reserves[] memory reserves) {
        uint length = pools.length;
        reserves = new Reserves[](length);

        for (uint i = 0; i < length;) {
            (bool success, bytes memory data) = pools[i].call(method);

            if (success && data.length == 96) {
                (uint112 reserve0, uint112 reserve1, ) = abi.decode(data, (uint112, uint112, uint32));
                reserves[i] = Reserves(reserve0, reserve1);
            } else {
                reserves[i] = Reserves(0, 0);
            }

            unchecked { ++i; }
        }
    }
}