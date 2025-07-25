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
	ERC20_CONTRACT_ADDR  = "0x5FbDB2315678afecb367f032d93F642f64180aa3" // ERC20コントラクトアドレス
	ERC721_CONTRACT_ADDR = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512" // ERC721コントラクトアドレス

	DEPLOYER_ADDR = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266" // デプロイヤー
	DB_PATH       = "../devchain/geth/chaindata"                 // LevelDBパス
)

// address系mapping: mapping(address=>uint256) slotN（ERC20, ERC721 _balancesなど）
func mappingSlotAddress(addr common.Address, slotNum uint64) common.Hash {
	keyBytes := make([]byte, 64)
	copy(keyBytes[12:], addr.Bytes()) // 右詰め
	slot := new(big.Int).SetUint64(slotNum).Bytes()
	copy(keyBytes[32+32-len(slot):], slot) // slotNumを右詰め
	h := sha3.NewLegacyKeccak256()
	h.Write(keyBytes)
	var slotHash common.Hash
	h.Sum(slotHash[:0])
	return slotHash
}

// uint系mapping: mapping(uint256=>address) slotN（ERC721 _ownersなど）
func mappingSlotUint(tokenId *big.Int, slotNum uint64) common.Hash {
	keyBytes := make([]byte, 64)
	tokenIdBytes := tokenId.Bytes()
	copy(keyBytes[32-len(tokenIdBytes):32], tokenIdBytes) // tokenIdを右詰め
	slot := new(big.Int).SetUint64(slotNum).Bytes()
	copy(keyBytes[64-len(slot):], slot) // slotNumを右詰め
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

	contract, shouldReturn := printStateBasic(stateDb, ERC721_CONTRACT_ADDR)
	if shouldReturn {
		return
	}

	deployer := common.HexToAddress(DEPLOYER_ADDR)
	fmt.Printf("👤 Deployer Address: %s\n", deployer.Hex())

	// ========== ERC721: _balances mapping [mapping(address => uint256) at slot 2] ==========
	erc721BalancesSlot := uint64(2)
	slotHash := mappingSlotAddress(deployer, erc721BalancesSlot)
	balanceRaw := stateDb.GetState(contract, slotHash)
	balanceInt := new(big.Int).SetBytes(balanceRaw.Bytes())
	fmt.Printf("    🔑 ERC721 _balances[deployer] (slot %d): %s (hex: %s)\n", erc721BalancesSlot, balanceInt.String(), balanceRaw.Hex())

	// ========== ERC721: _owners mapping [mapping(uint256 => address) at slot 3] ==========
	tokenId := big.NewInt(1)
	erc721OwnersSlot := uint64(3)
	ownerSlotHash := mappingSlotUint(tokenId, erc721OwnersSlot)
	ownerRaw := stateDb.GetState(contract, ownerSlotHash)
	ownerAddr := common.BytesToAddress(ownerRaw.Bytes())
	fmt.Printf("    🔑 ERC721 _owners[%d] (slot %d): %s (hex: %s)\n", tokenId.Int64(), erc721OwnersSlot, ownerAddr.Hex(), ownerRaw.Hex())

	contract, shouldReturn = printStateBasic(stateDb, ERC20_CONTRACT_ADDR)
	if shouldReturn {
		return
	}
	// ========== ERC20: balances mapping [mapping(address => uint256) at slot 0] ==========
	erc20Contract := common.HexToAddress(ERC20_CONTRACT_ADDR)
	erc20BalancesSlot := uint64(0) // ERC20のbalancesはslot 0
	erc20SlotHash := mappingSlotAddress(deployer, erc20BalancesSlot)
	erc20BalanceRaw := stateDb.GetState(erc20Contract, erc20SlotHash)
	erc20BalanceInt := new(big.Int).SetBytes(erc20BalanceRaw.Bytes())
	fmt.Printf("    🔑 ERC20 balances[deployer] (slot %d): %s (hex: %s)\n", erc20BalancesSlot, erc20BalanceInt.String(), erc20BalanceRaw.Hex())
}

func printStateBasic(stateDb *state.StateDB, contractAddress string) (common.Address, bool) {
	contract := common.HexToAddress(contractAddress)
	if stateDb.Empty(contract) {
		fmt.Println("⚠️ Account not found")
		return common.Address{}, true
	}

	fmt.Printf("✅ Contract Address: %s\n", contract.Hex())
	fmt.Printf("🔢 Nonce: %d\n", stateDb.GetNonce(contract))
	fmt.Printf("💰 Balance: %s\n", stateDb.GetBalance(contract))
	fmt.Printf("📄 Code Hash: %x\n", stateDb.GetCodeHash(contract))

	fmt.Println("🧾 Storage Slots (first 10):")
	for i := 0; i < 10; i++ { // slot番号はお好みで
		slot := common.BigToHash(big.NewInt(int64(i)))
		value := stateDb.GetState(contract, slot)
		if value != (common.Hash{}) {
			fmt.Printf("    🔑 slot %d → 📦 value: %s\n", i, value.Hex())
		}
	}
	return contract, false
}
