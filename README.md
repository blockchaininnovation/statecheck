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

gethを初期化
```
geth --datadir ./devchain init genesis.json
```

.envのカギをもとにgethにアカウント作成
```
cd simpletoken
./import_key_from_env.sh 
```

geth起動
```

```

state読み込み・出力。gethが起動していたら止めてから実行：
```
cd geth_state_reader

pipenv install
pipenv run python read_state.py
```