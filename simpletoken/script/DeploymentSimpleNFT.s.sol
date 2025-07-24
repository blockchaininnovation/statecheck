// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Script.sol";
import "../src/SimpleNFT.sol";

contract DeploymentSimpleNFT is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        SimpleNFT s = new SimpleNFT();

        vm.stopBroadcast();

        bytes memory encodedData = abi.encodePacked(
            "[SimpleNFT] deployed address: ",
            vm.toString(address(s))
        );
        console.log(string(encodedData));
    }
}
