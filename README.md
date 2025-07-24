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


gethインストール
```
sudo add-apt-repository -y ppa:ethereum/ethereum
sudo apt update
sudo apt install ethereum
```

gethをdevモードで起動
```
geth --datadir ./devchain --dev --http --http.api eth,net,web3,debug
```

.envのカギをもとにgethにアカウント作成
```
cd simpletoken
./import_key_from_env.sh 
```

作ったアカウントに開発用自動生成アカウントから送金
まずはgethコンソール起動
```
cd ..
geth attach ipc:./devchain/geth.ipc
```

gethコンソール内で：
```
eth.sendTransaction({from: eth.accounts[0], to: "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266", value: web3.toWei(1000, "ether")})
```

残高確認
```
web3.fromWei(eth.getBalance("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"), "ether")
```