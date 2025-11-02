package network

import (
	"fmt"
	"net"
	"sync"
)

// Peer represents a network peer
type Peer struct {
	Address    string
	Connection net.Conn
	Version    int
	Height     int
}

// PeerList manages known peers
type PeerList struct {
	peers map[string]*Peer
	mu    sync.RWMutex
}

// NewPeerList creates a new peer list
func NewPeerList() *PeerList {
	return &PeerList{
		peers: make(map[string]*Peer),
	}
}

// Add adds a peer to the list
func (pl *PeerList) Add(address string, conn net.Conn) *Peer {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	peer := &Peer{
		Address:    address,
		Connection: conn,
	}
	pl.peers[address] = peer

	return peer
}

// Remove removes a peer from the list
func (pl *PeerList) Remove(address string) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	if peer, exists := pl.peers[address]; exists {
		if peer.Connection != nil {
			peer.Connection.Close()
		}
		delete(pl.peers, address)
	}
}

// Get retrieves a peer by address
func (pl *PeerList) Get(address string) (*Peer, bool) {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	peer, exists := pl.peers[address]
	return peer, exists
}

// GetAll returns all peers
func (pl *PeerList) GetAll() []*Peer {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	peers := make([]*Peer, 0, len(pl.peers))
	for _, peer := range pl.peers {
		peers = append(peers, peer)
	}

	return peers
}

// Count returns the number of peers
func (pl *PeerList) Count() int {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	return len(pl.peers)
}

// GetAddresses returns all peer addresses
func (pl *PeerList) GetAddresses() []string {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	addresses := make([]string, 0, len(pl.peers))
	for addr := range pl.peers {
		addresses = append(addresses, addr)
	}

	return addresses
}

// SendData sends data to a peer
func (p *Peer) SendData(data []byte) error {
	if p.Connection == nil {
		return fmt.Errorf("peer %s has no active connection", p.Address)
	}

	_, err := p.Connection.Write(data)
	return err
}

// UpdateInfo updates peer information
func (p *Peer) UpdateInfo(version, height int) {
	p.Version = version
	p.Height = height
}

