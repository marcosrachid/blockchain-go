package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"

	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version        = byte(0x00) // Address version (similar to Bitcoin)
)

// getWalletFile returns the wallet file path, checking for Docker environment first
func getWalletFile() string {
	// Check if we're in Docker environment by looking for the data directory
	dockerPath := "/app/data/tmp/wallets.dat"
	dockerDir := "/app/data/tmp"

	// Create directory if it doesn't exist (Docker environment)
	if _, err := os.Stat("/app/data"); err == nil {
		os.MkdirAll(dockerDir, 0755)
		log.Printf("🔑 Using Docker wallet path: %s", dockerPath)
		return dockerPath
	}

	// Fallback to local development path
	// Create local tmp directory if needed
	if _, err := os.Stat("./tmp"); os.IsNotExist(err) {
		os.MkdirAll("./tmp", 0755)
	}
	log.Printf("🔑 Using local wallet path: ./tmp/wallets.dat")
	return "./tmp/wallets.dat"
}

// Wallet stores private and public keys (ECDSA cryptography)
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// serializableWallet is a serializable version of Wallet
type serializableWallet struct {
	D         []byte
	X         []byte
	Y         []byte
	PublicKey []byte
}

// Wallets stores a collection of wallets
type Wallets struct {
	Wallets map[string]*Wallet
}

// MarshalBinary implements encoding.BinaryMarshaler
func (w *Wallet) MarshalBinary() ([]byte, error) {
	sw := serializableWallet{
		D:         w.PrivateKey.D.Bytes(),
		X:         w.PrivateKey.X.Bytes(),
		Y:         w.PrivateKey.Y.Bytes(),
		PublicKey: w.PublicKey,
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(sw)
	return buf.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (w *Wallet) UnmarshalBinary(data []byte) error {
	var sw serializableWallet
	buf := bytes.NewReader(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&sw); err != nil {
		return err
	}

	curve := elliptic.P256()
	w.PrivateKey.PublicKey.Curve = curve
	w.PrivateKey.D = new(big.Int).SetBytes(sw.D)
	w.PrivateKey.X = new(big.Int).SetBytes(sw.X)
	w.PrivateKey.Y = new(big.Int).SetBytes(sw.Y)
	w.PublicKey = sw.PublicKey

	return nil
}

// NewWallet creates a new wallet
func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

// Address returns the wallet address (similar to Bitcoin addresses)
func (w Wallet) Address() []byte {
	pubHash := HashPubKey(w.PublicKey)

	versionedHash := append([]byte{version}, pubHash...)
	checksum := Checksum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)

	return address
}

// newKeyPair generates a new key pair using ECDSA
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pub
}

// HashPubKey hashes the public key (SHA256 + RIPEMD160, like in Bitcoin)
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// Checksum generates a checksum for a payload
func Checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:checksumLength]
}

// ValidateAddress validates a Bitcoin-like address
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
	targetChecksum := Checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Equal(actualChecksum, targetChecksum)
}

// NewWallets creates a new collection of wallets
func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFile()

	return &wallets, err
}

// AddWallet adds a wallet to the collection
func (ws *Wallets) AddWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	ws.Wallets[address] = wallet

	return address
}

// GetWallet returns a wallet by address
func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

// GetAllAddresses returns all wallet addresses
func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// LoadFile loads wallets from file
func (ws *Wallets) LoadFile() error {
	walletFilePath := getWalletFile()
	if _, err := os.Stat(walletFilePath); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets

	fileContent, err := ioutil.ReadFile(walletFilePath)
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		return err
	}

	ws.Wallets = wallets.Wallets

	return nil
}

// SaveFile saves wallets to file
func (ws *Wallets) SaveFile() {
	var content bytes.Buffer

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	walletFilePath := getWalletFile()
	err = ioutil.WriteFile(walletFilePath, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
