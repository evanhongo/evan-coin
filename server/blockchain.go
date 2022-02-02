package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	mathRand "math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/evanhongo/blockchain-demo/pkg/util/array"
	logger "github.com/sirupsen/logrus"
)

type Transaction struct {
	Sender   string
	Receiver string
	Amount   int
	Fee      int
	Message  string
}

type Block struct {
	Hash         string
	Transactions []Transaction
	PreviousHash string
	Timestamp    time.Time
	Pow          uint64 //Proof of work
	Difficulty   int
	Miner        string
	MiningReward int
}

type BlockChain struct {
	Chain                      []Block
	BlockSize                  int //Maximum number of transactions a block can contain
	CurrentDifficulty          int
	MiningReward               int
	ExpectedMiningTime         int
	BlockNumToAdjustDifficulty int //Adjust the difficulty every how many blocks
	PendingTransactions        []Transaction
}

func (bc *BlockChain) generateGenesisBlock() {
	logger.Infoln("Generate genesis block🎉")
	newBlock := Block{PreviousHash: "Hello World!", Difficulty: bc.CurrentDifficulty, Miner: "God", MiningReward: bc.MiningReward}
	newBlock.Hash = bc.calculateHash(&newBlock)
	bc.Chain = append(bc.Chain, newBlock)
}

func (bc *BlockChain) movePendingTransactionsToBlock(b *Block) {
	//Get the transaction with highest fee by block size
	var selected []Transaction
	sort.Slice(bc.PendingTransactions[:], func(i, j int) bool { return bc.PendingTransactions[j].Fee < bc.PendingTransactions[i].Fee })
	if len(bc.PendingTransactions) > bc.BlockSize {
		selected = bc.PendingTransactions[:bc.BlockSize]
		bc.PendingTransactions = bc.PendingTransactions[bc.BlockSize:]
	} else {
		selected = bc.PendingTransactions
		bc.PendingTransactions = []Transaction{}
	}
	b.Transactions = selected
}

func (bc *BlockChain) calculateHash(b *Block) string {
	transactions, _ := json.Marshal(b.Transactions)
	blockData := b.PreviousHash + string(transactions) + b.Timestamp.String() + strconv.FormatUint(b.Pow, 10)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

func (bc *BlockChain) mine(miner string) {
	logger.Infoln("Mining⛏️")
	start := time.Now()
	lastBlock := bc.Chain[len(bc.Chain)-1]
	mathRand.Seed(time.Now().UnixNano())
	newBlock := Block{
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now(),
		Pow:          mathRand.Uint64(), //Prevent the solver from always being the node with highest computing power
		Difficulty:   bc.CurrentDifficulty,
		Miner:        miner,
		MiningReward: bc.MiningReward}
	bc.movePendingTransactionsToBlock(&newBlock)
	for !strings.HasPrefix(newBlock.Hash, strings.Repeat("0", bc.CurrentDifficulty)) {
		newBlock.Pow++
		newBlock.Hash = bc.calculateHash(&newBlock)
	}
	duration := time.Since(start)
	bc.Chain = append(bc.Chain, newBlock)
	logger.Infof("Solved🎆, pow: %d, time cost⌚: %.4fms", newBlock.Pow, float64(duration)/1000000)
}

func (bc *BlockChain) isChainValid() bool {
	for i := range bc.Chain[1:] {
		previousBlock := bc.Chain[i]
		currentBlock := bc.Chain[i+1]
		if currentBlock.Hash != bc.calculateHash(&currentBlock) || currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}
	}
	return true
}

func (bc *BlockChain) getBalance(account string) int {
	balance := 0
	for _, b := range bc.Chain[1:] {
		if b.Miner == account {
			balance += b.MiningReward
			feeArr := array.Map(b.Transactions, func(t Transaction) int { return t.Fee })
			balance += array.Sum(feeArr)
		}
		for _, t := range b.Transactions {
			if t.Sender == account {
				balance -= t.Amount + t.Fee
			} else if t.Receiver == account {
				balance += t.Amount
			}
		}
	}
	for _, t := range bc.PendingTransactions {
		if t.Sender == account {
			balance -= t.Amount + t.Fee
		} else if t.Receiver == account {
			balance += t.Amount
		}
	}
	return balance
}

func (bc *BlockChain) addPendingTransaction(t *Transaction) {
	bc.PendingTransactions = append(bc.PendingTransactions, *t)
}

func (bc *BlockChain) verifyTransaction(t *Transaction, signature []byte) (isValid bool) {
	isValid = true
	var bytes []byte
	var err error
	if bc.getBalance(t.Sender) < t.Amount+t.Fee {
		logger.Errorln("Balance is not enough")
		return false
	}
	publicKey := t.Sender
	if bytes, err = base64.StdEncoding.DecodeString(publicKey); err != nil {
		logger.Errorln(err)
		return false
	}
	pbk, err := x509.ParsePKIXPublicKey(bytes)
	if err != nil {
		logger.Errorln(err)
		return false
	}
	tStr, _ := json.Marshal(t)
	hashed := sha256.Sum256(tStr)
	if err := rsa.VerifyPKCS1v15(pbk.(*rsa.PublicKey), crypto.SHA256, hashed[:], signature); err != nil {
		logger.Errorln(err)
		return false
	}

	return
}
