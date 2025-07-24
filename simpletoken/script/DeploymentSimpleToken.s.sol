// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Script.sol";
import "../src/SimpleToken.sol";

contract DeploymentSimpleToken is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        SimpleToken s = new SimpleToken();

        vm.stopBroadcast();

        bytes memory encodedData = abi.encodePacked(
            "[SimpleToken] deployed address: ",
            vm.toString(address(s))
        );
        console.log(string(encodedData));
    }
}
