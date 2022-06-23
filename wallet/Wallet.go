package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	crypto "crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"time"

	base58 "github.com/btcsuite/btcd/btcutil/base58"

	"golang.org/x/crypto/ripemd160"
)

var logFile *log.Logger = GetlogFile()

func GetlogFile() *log.Logger {
	f, _ := os.OpenFile("rpc.log", os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
	defer f.Close()

	return log.New(f, "[INFO]", log.LstdFlags)
}

type Wa Wallet

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
	Address    string
	Alias      string
	Timestamp  int64
}

type Wallets struct {
	Wallet map[string]*Wallet
}

func NewWallets() *Wallets {
	return &Wallets{}
}
func (ws *Wallets) AddWallet(w *Wallet) {
	ws.Wallet = map[string]*Wallet{w.Address: w}
}

func newWallet(alias string) *Wa {
	private, public := newKeyPair()
	wallet := &Wa{PrivateKey: private, PublicKey: public, Alias: alias, Timestamp: time.Now().UTC().Unix()}
	version := byte(0x00)
	pubKeyHash := HashPubKey(public)
	b := make([]byte, 0, 1+len(pubKeyHash)+4)
	b = append(b, version)
	b = append(b, pubKeyHash[:]...)
	cksum := checksum(b)
	b = append(b, cksum[:]...)
	addr := base58.Encode(b)
	wallet.Address = addr
	return wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, _ := ecdsa.GenerateKey(curve, crypto.Reader)
	public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, public
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		logFile.Fatal(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

func checksum(check []byte) []byte {
	firstSHA := sha256.Sum256(check)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:4]
}

func ValidateAddress(address string) bool {
	dum := base58.Decode(address)
	version := dum[0]
	keyHash := dum[len(dum)-4:]
	dum = dum[1 : len(dum)-4]
	checkresult := checksum(append([]byte{version}, dum...))

	return bytes.Equal(keyHash, checkresult)
}

func GetAddress(public []byte) string {
	pubKeyHash := HashPubKey(public)
	ver := byte(0x00)
	verHash := append([]byte{ver}, pubKeyHash...)
	cksum := checksum(verHash)
	verHashsum := append(verHash, cksum...)
	address := base58.Encode(verHashsum)

	if ValidateAddress(address) {
		return address
	} else {
		return "지갑 생성 실패"
	}

}

func encode(privateKey *ecdsa.PrivateKey) string {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)                                // RSA 개인키 형식을 PEM블록형식으로 인코딩
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded}) // []byte type 메모리로 인코드

	return string(pemEncoded)
}

func decode(pemEncoded string) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

	return privateKey
}
