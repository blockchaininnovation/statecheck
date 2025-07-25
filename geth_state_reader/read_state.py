import rlp
import plyvel
from eth_utils import keccak, to_hex
from rlp.sedes import big_endian_int, Binary
from eth_rlp import HashableRLP

from eth.db.backends.memory import MemoryDB  # ← 修正
from trie.hexary import HexaryTrie


# Gethの LevelDB を開く
leveldb = plyvel.DB("../devchain/geth/chaindata", compression=None)

# MemoryDB にコピー
memdb = MemoryDB()
for k, v in leveldb:
    memdb[k] = v

# ステートルート（例：起動中のgethから `eth.getBlock("latest").stateRoot` で取得）
state_root_hex = "d7c1ef5207349c1dbef799e0bf1d201b7e5cd99bbcfdffe031874252cfccd540"
state_root = bytes.fromhex(state_root_hex)

# トライ構築
trie = HexaryTrie(memdb, root_hash=state_root)

# 対象アドレス
address = "0x5FbDB2315678afecb367f032d93F642f64180aa3"
address_bytes = bytes.fromhex(address[2:])
hashed_address = keccak(address_bytes)

# MPTから該当ノード取得
account_rlp = trie.get(hashed_address)
if account_rlp is None:
    print("⚠️ アカウントが存在しません")
    exit(1)

print(f"📦 Raw account RLP: {account_rlp.hex()}")

# アカウント構造体
class Account(HashableRLP):
    fields = [
        ("nonce", big_endian_int),
        ("balance", big_endian_int),
        ("storage_root", Binary.fixed_length(32)),
        ("code_hash", Binary.fixed_length(32)),
    ]

account = rlp.decode(account_rlp, Account)

print(f"✅ Address: {address}")
print(f"🔢 Nonce: {account.nonce}")
print(f"💰 Balance: {account.balance}")
print(f"📦 Storage Root: {to_hex(account.storage_root)}")
print(f"📄 Code Hash: {to_hex(account.code_hash)}")
