package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"golang.org/x/crypto/sha3"
)

// ====== ここでCA/デプロイヤー指定 ======
const (
	CONTRACT_ADDR = "0x5FbDB2315678afecb367f032d93F642f64180aa3" // コントラクトアドレス
	DEPLOYER_ADDR = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266" // デプロイヤー
	DB_PATH       = "../devchain/geth/chaindata"                 // LevelDBパス
)

func mappingSlot(addr common.Address, slotNum uint64) common.Hash {
	// mappingのキー = keccak256(pad(addr) ++ pad(slot))
	keyBytes := make([]byte, 64)
	copy(keyBytes[12:], addr.Bytes()) // 左12バイト0で右詰め
	slot := new(big.Int).SetUint64(slotNum).Bytes()
	copy(keyBytes[32+32-len(slot):], slot) // slotNumを右詰め
	h := sha3.NewLegacyKeccak256()
	h.Write(keyBytes)
	var slotHash common.Hash
	h.Sum(slotHash[:0])
	return slotHash
}

func main() {
	ldb, err := leveldb.New(DB_PATH, 0, 0, "", false)
	if err != nil {
		log.Fatalf("❌ Failed to open LevelDB: %v", err)
	}
	defer ldb.Close()

	db := rawdb.NewDatabase(ldb)

	// 最新ブロックのステートルート取得
	headHash := rawdb.ReadHeadBlockHash(db)
	headNumber := rawdb.ReadHeaderNumber(db, headHash)
	if headNumber == nil {
		log.Fatalf("❌ Could not get block number")
	}
	block := rawdb.ReadBlock(db, headHash, *headNumber)
	if block == nil {
		log.Fatalf("❌ Could not get block")
	}

	stateRoot := block.Root()
	stateDb, err := state.New(stateRoot, state.NewDatabase(db), nil)
	if err != nil {
		log.Fatalf("❌ Failed to open statedb: %v", err)
	}

	contract := common.HexToAddress(CONTRACT_ADDR)
	if stateDb.Empty(contract) {
		fmt.Println("⚠️ Account not found")
		return
	}

	fmt.Printf("✅ Contract Address: %s\n", contract.Hex())
	fmt.Printf("🔢 Nonce: %d\n", stateDb.GetNonce(contract))
	fmt.Printf("💰 Balance: %s\n", stateDb.GetBalance(contract))
	fmt.Printf("📄 Code Hash: %x\n", stateDb.GetCodeHash(contract))

	fmt.Println("🧾 Storage Slots:")
	for i := 0; i < 10; i++ { // slot番号はお好みで
		slot := common.BigToHash(big.NewInt(int64(i)))
		value := stateDb.GetState(contract, slot)
		if value != (common.Hash{}) {
			fmt.Printf("    🔑 slot %d → 📦 value: %s\n", i, value.Hex())
		}
	}
	deployer := common.HexToAddress(DEPLOYER_ADDR)
	fmt.Printf("👤 Deployer Address: %s\n", deployer.Hex())

	// ========== balances mapping のstorage slot取得 ==========
	balancesSlot := uint64(0) // balancesは一番最初の状態変数
	slotHash := mappingSlot(deployer, balancesSlot)
	balanceRaw := stateDb.GetState(contract, slotHash)
	balanceInt := new(big.Int).SetBytes(balanceRaw.Bytes())
	fmt.Printf("    🔑 slot %d → Deployer's ERC20 balance: %s (hex: %s)\n", balancesSlot, balanceInt.String(), balanceRaw.Hex())

}
