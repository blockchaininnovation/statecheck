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


gethインストール(1.13.15)、PoAで使いたいので古いバージョンを使用。
```
wget https://github.com/ethereum/go-ethereum/archive/refs/tags/v1.13.15.tar.gz
tar -xvzf v1.13.15.tar.gz
cd go-ethereum-1.13.15/
make geth
sudo cp build/bin/geth /usr/local/bin/geth

geth version
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
(上記でインポートしているアカウントは0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266としている。そこで設定したパスワードをpassword.txtに記載)
```
geth --datadir ./devchain   --networkid 1337   --http --http.api eth,web3,net,personal   --unlock 0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266   --password ./password.txt   --mine   --miner.etherbase 0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266   --nodiscover   --verbosity 3   --allow-insecure-unlock
 2063  history
```

ここでコントラクト（SimpleToken, SimpleNFT）デプロイ。


state読み込み・出力。gethが起動していたら止めてから実行：
```
cd geth_state_reader_go

go run main.go
```