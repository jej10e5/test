package BLC

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

//const targetBites = 20 //채굴난이도
var targetBites, _ = rand.Int(rand.Reader, big.NewInt(12))

type ProofOfWork struct {
	block  *Block
	target *big.Int //요구사항
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.setTxsHash(),
		IntToHex(pow.block.Timestamp),
		IntToHex(targetBites.Int64()),
		IntToHex(int64(nonce)),
	}, []byte{})
	return data
}
func (pow *ProofOfWork) Run() (int, []byte, int64) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce, hash[:], targetBites.Int64()
}

func newProofOfWork(block *Block) *ProofOfWork {
	targetBites, _ = rand.Int(rand.Reader, big.NewInt(12))
	target := big.NewInt(1)
	target.Lsh(target, uint(256-int(targetBites.Int64())))
	pow := &ProofOfWork{block, target}
	return pow
}
