# statecheck

anvil stateダンプ
```shell
anvil --dump-state state.json
```

build
```
forge build
```

このとき以下のエラーが出る場合がある：
```
Missing dependencies found. Installing now...

Updating dependencies in /home/someone/git/statecheck/simpletoken/lib
The application panicked (crashed).
Message:  failed to extract foundry config:
foundry config error: Unknown evm version: prague for setting `evm_version`

Location: crates/config/src/lib.rs:532

This is a bug. Consider reporting it at https://github.com/foundry-rs/foundry
```

このときはlib/openzeppelin-contracts/foundry.tomlのevm_versionを修正
```
evm_version = 'cancun'
```


SimpleTokenデプロイ
```shell
forge script script/DeploymentSimpleToken.s.sol:DeploymentSimpleToken --rpc-url http://127.0.0.1:8545 --broadcast
```

SimpleNFTデプロイ
```shell
forge script script/DeploymentSimpleNFT.s.sol --rpc-url http://127.0.0.1:8545 --broadcast
```


