package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

// CommandLength is the fixed length for command names
const CommandLength = 12

// Message types
const (
	CmdVersion     = "version"
	CmdGetBlocks   = "getblocks"
	CmdInv         = "inv"
	CmdGetData     = "getdata"
	CmdBlock       = "block"
	CmdTx          = "tx"
	CmdAddr        = "addr"
	CmdPing        = "ping"
	CmdPong        = "pong"
)

// Inventory types
const (
	InvTypeBlock = "block"
	InvTypeTx    = "tx"
)

// Version message for handshake
type Version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

// GetBlocks requests blocks from a peer
type GetBlocks struct {
	AddrFrom string
}

// Inv inventory message
type Inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

// GetData requests specific data
type GetData struct {
	AddrFrom string
	Type     string
	ID       []byte
}

// Block message
type BlockMsg struct {
	AddrFrom string
	Block    []byte
}

// Tx transaction message
type TxMsg struct {
	AddrFrom    string
	Transaction []byte
}

// Addr peer address message
type Addr struct {
	AddrList []string
}

// Ping message
type Ping struct{}

// Pong response
type Pong struct{}

// CmdToBytes converts command to fixed-length byte array
func CmdToBytes(cmd string) []byte {
	var bytes [CommandLength]byte

	for i, c := range cmd {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

// BytesToCmd extracts command from byte array
func BytesToCmd(bytes []byte) string {
	var cmd []byte

	for _, b := range bytes {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}

	return fmt.Sprintf("%s", cmd)
}

// GobEncode encodes data using gob
func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// ExtractCmd extracts command from request
func ExtractCmd(request []byte) []byte {
	return request[:CommandLength]
}

