package TX

import (
	w "Wallet"
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type Tx struct {
	Hash      []byte `json:"txId"`
	Timestamp int64  `json:"timestamp"`
	From      []byte `json:"from"`
	To        []byte `json:"to"`
	Item      []byte `json:"item"`
	Price     int    `json:"price"`
	Nonce     int    `json:"nonce"`
	Sig       []byte `json:"sig"`
}

func (tx *Tx) setHash() {

	timestamp := strconv.FormatInt(tx.Timestamp, 10)
	timeBytes := []byte(timestamp)
	txBytes := bytes.Join([][]byte{
		timeBytes,
		tx.From,
		tx.To,
		tx.Item,
		IntToHex(int64(tx.Price)),
		//TODO nonce
	}, []byte{})

	hash := sha256.Sum256(txBytes)
	tx.Hash = hash[:]

}

func NewTx(fw *w.Wallet, tw *w.Wallet, item string, pri int, nonce int) *Tx {

	tx := &Tx{}
	tx.Timestamp = time.Now().UTC().Unix() //utc기준 시간
	tx.Price = pri
	tx.From = []byte(fw.Address)
	tx.To = []byte(tw.Address)
	tx.Item = []byte(item)
	tx.Nonce = nonce
	tx.setHash()
	tx.Sign(fw)

	return tx
}

func NewGenesisTx() *Tx {

	tx := &Tx{}
	tx.Timestamp = time.Now().UTC().Unix() //utc기준 시간
	tx.Price = 0
	tx.From = []byte("Genesis")
	tx.To = []byte("Mine")
	tx.Nonce = 0
	tx.Sig = nil
	tx.setHash()
	return tx
}

func (tx *Tx) Txprint() {
	fmt.Printf("Hash:%x\n", tx.Hash)
	fmt.Printf("From:%d\n", tx.From)
	fmt.Printf("To:%s\n", tx.To)
	fmt.Printf("Price:%d\n", tx.Price)
	fmt.Printf("Sig:%d\n", tx.Sig)
	fmt.Println("---------------------------------------------------------------")
}

func (t *Tx) EqualHash(e []byte) bool {
	return bytes.Equal(t.Hash, e)
}

func (tx *Tx) Sign(w *w.Wallet) {
	privateKey := w.PrivateKey
	sig, err := ecdsa.SignASN1(rand.Reader, &privateKey, tx.Hash[:])
	if err != nil {
		panic(err)
	}
	tx.Sig = sig[:]
}

func (tx *Tx) ValidateTx(w *w.Wallet) bool {
	return ecdsa.VerifyASN1(&w.PrivateKey.PublicKey, tx.Hash[:], tx.Sig[:])
}
