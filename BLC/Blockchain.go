package BLC

import (
	"TX"
	w "Wallet"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/google/go-cmp/cmp"
)

type Blockchain struct {
	Blocks []*Block
}

type BcJson struct {
	Blocks []*BJson
}

//새 블록을 만들어서 기존의 블록체인(블록들)에 추가하는 함수
func (blockchain *Blockchain) AddBlock(txs *TX.TxsData) {
	//------------------------------
	// 채우기
	//------------------------------
	//이전 블록을 찾아서 해시값을 알아내야함
	//기존의 블록체인의 블록들 중 가장 끝에거이므로 len함수를 사용해서 가장 끝의 블록을 가져온다.
	prev := blockchain.Blocks[len(blockchain.Blocks)-1] //--2
	//새로운 블록 만드는 코드
	//data와 이전블록의 해시값 필요 -> 이전 hash값 구해야함 --1
	nB := NewBlock(txs, prev.Hash[:], prev.Height)
	//기존의 블록체인에다가 새 블록을 추가하기
	//구조체
	//blockchain이라는 구조체 내의 Blocks에 값을 넣는거
	blockchain.Blocks = append(blockchain.Blocks, nB) //--3

}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func (bc *Blockchain) FindBlock(id []byte) *Block {
	for _, v := range bc.Blocks {
		if v.EqualHash(id) {
			return v
		}
	}
	return nil
}

func (bc *Blockchain) BcJson() {
	bcj, _ := json.MarshalIndent(bc, "", " ")
	err := ioutil.WriteFile("./test.json", bcj, os.FileMode(0644))
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (bc *Blockchain) GetJson() BcJson {
	b, err := ioutil.ReadFile("./test.json") //파일읽기
	if err != nil {
		fmt.Println(err)
	}

	var bjs BcJson
	json.Unmarshal(b, &bjs)
	return bjs
}

func (bc *Blockchain) Comp(bcj BcJson) bool {
	bj := bcj.Blocks
	for i, b := range bc.Blocks {
		/*if !b.EqualHash(bj[i].Hash) {
			return false
		}*/
		if cmp.Equal(b, bj[i]) { //두 구조체 자료 모두 비교
			return false
		}
	}
	return true
}
func RandData(size int, datasize int64) []int {
	num := make([]int, size)
	for i := 0; i < size; i++ {
		check := 0
		nr, _ := rand.Int(rand.Reader, big.NewInt(datasize))
		num[i] = int(nr.Int64())
		if i > 0 {
			for {
				check = 0
				for j := 0; j < i; j++ {

					if num[i] == num[j] {
						nrr, _ := rand.Int(rand.Reader, big.NewInt(datasize))
						num[i] = int(nrr.Int64())
						check = 1
					}
				}
				if check != 1 {
					break
				}
			}
		}
	}
	return num

}

func (b *Blockchain) GetTxCount(w *w.Wallet) int {
	cnt := 1
	for _, v := range b.Blocks {
		txs := v.Txs
		for _, t := range txs.Txs {
			if bytes.Equal(t.From, []byte(w.Address)) {
				cnt++
			}
		}
	}
	return cnt
}

func (b *Blockchain) FindTx(id []byte) *TX.Tx {
	//inHash := b.Blocks[len(b.Blocks)-1].Hash
	for _, v := range b.Blocks {
		//if bytes.Equal(v.Hash, inHash) {
		for _, t := range v.Txs.Txs {
			if bytes.Equal(t.Hash, id) {
				return t
			}
		}
		//inHash = blc.PrevBlockHash
		//if bytes.Equal(inHash, b.Blocks[0].Hash) {
		//	break
		//}
	}
	return nil
}
