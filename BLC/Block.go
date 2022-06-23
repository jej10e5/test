package BLC

import (
	"TX"
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	Hash          []byte      `json:"Hash"`
	PrevBlockHash []byte      `json:"PrevBlockHash"`
	Timestamp     int64       `json:"Timestamp"`
	Pow           []byte      `json:"Pow"`
	Nonce         int         `json:"Nonce"`
	Bit           int64       `json:"Bit"`
	Txs           *TX.TxsData `json:"Txs"`
	Height        int         `json:"Height"`
}

type BJson struct {
	Hash          []byte      `json:"Hash"`
	PrevBlockHash []byte      `json:"PrevBlockHash"`
	Timestamp     int64       `json:"Timestamp"`
	Pow           []byte      `json:"Pow"`
	Nonce         int         `json:"Nonce"`
	Bit           int64       `json:"Bit"`
	Txs           *TX.TxsData `json:"Txs"`
	Height        int         `json:"Height"`
}

func (block *Block) setHash() {

	timestamp := strconv.FormatInt(block.Timestamp, 10)
	timeBytes := []byte(timestamp)

	//---------------------------------
	// blockBytes 값 채우기
	//---------------------------------

	blockBytes := bytes.Join([][]byte{
		block.PrevBlockHash,
		block.setTxsHash(),
		timeBytes,
		block.Pow,
		IntToHex(int64(block.Nonce)),
		IntToHex(int64(block.Bit)),
	}, []byte{})
	hash := sha256.Sum256(blockBytes)
	block.Hash = hash[:]

}

//블록 생성시 data값과 이전블록의hash값이 필요.
//Blockchain.go에서 AddBlock메서드에서 사용됨
func NewBlock(txs *TX.TxsData, prevBlockHash []byte, pHeight int) *Block {

	block := &Block{}
	//----------------------------
	//  block element 값 채우기
	//----------------------------
	block.PrevBlockHash = prevBlockHash
	block.Timestamp = time.Now().UTC().Unix() //utc기준 시간
	block.Txs = txs                           //나중에처리
	pow := newProofOfWork(block)
	nonce, hash, bits := pow.Run()
	block.Pow = hash[:]
	block.Nonce = nonce
	block.Bit = bits
	block.Height = pHeight + 1
	block.setHash()

	return block
}

func NewGenesisBlock() *Block { //hash값 32byte
	gt := TX.NewTxs()
	gt.AddTx(TX.NewGenesisTx())
	return NewBlock(gt, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 0)
} //previous hash값--genesis니까 없음

func (block *Block) Bprint() {
	fmt.Printf("Height:%d\n", block.Height)
	for i, t := range block.Txs.Txs {
		fmt.Printf("Data[%d]:{TxHash:%x}\n", i, t.Hash)
		fmt.Printf("\t{TxFrom:%x}\n", t.From)
		fmt.Printf("\t{TxTo:%x}\n", t.To)
		fmt.Printf("\t{TxItem:%s}\n", t.Item)
		fmt.Printf("\t{TxPrice:%d}\n", t.Price)
		fmt.Printf("\t{TxNonce:%d}\n", t.Nonce)
		fmt.Printf("\t{TxSig:%x}\n", t.Sig)
	}
	fmt.Printf("Pre:%x\n", block.PrevBlockHash)
	fmt.Printf("Hash:%x\n", block.Hash)
	fmt.Printf("bits:%d\n", block.Bit)
	fmt.Println("---------------------------------------------------------------")
}

func (b *Block) EqualHash(e []byte) bool {
	return bytes.Equal(b.Hash, e)
}

/*
func (b *Block) EqualData(e []byte) bool {
	return bytes.Equal(b.Txs, e)
}
*/
func (b *Block) IsGenBlock() bool {
	return bytes.Equal(b.PrevBlockHash, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
func (b *Block) setTxsHash() []byte {
	hash := []byte{}
	for _, t := range b.Txs.Txs {
		hash = append(hash, t.Hash...)
	}
	return hash[:]
}
