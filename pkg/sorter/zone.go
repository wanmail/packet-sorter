package sorter

import (
	"net"
)

type Locator interface {
	Locate(net.IP) string
}

// TrieNode represents a node in the prefix tree
type TrieNode struct {
	children [2]*TrieNode // Two child nodes: 0 and 1
	network  string       // If this node is the end of a network range, store the network information
	ipnet    *net.IPNet
}

// Trie represents the prefix tree
type Trie struct {
	root *TrieNode
}

// NewTrie creates a new prefix tree
func NewTrie() *Trie {
	return &Trie{root: &TrieNode{}}
}

// Insert inserts an IP network range into the prefix tree
func (t *Trie) Insert(cidr *net.IPNet) {
	ip := cidr.IP.To4()
	mask, _ := cidr.Mask.Size()

	node := t.root
	for i := 0; i < mask; i++ {
		bit := (ip[i/8] >> (7 - i%8)) & 1
		if node.children[bit] == nil {
			node.children[bit] = &TrieNode{}
		}
		node = node.children[bit]
	}
	node.network = cidr.String()
	node.ipnet = cidr
}

// FindNetwork finds the most specific network range for a given IP
func (t *Trie) FindNetwork(ip net.IP) string {
	ip = ip.To4()
	node := t.root
	network := ip.String()

	// 0.0.0.0/0
	if t.root.ipnet != nil {
		if t.root.ipnet.Contains(ip) {
			network = t.root.network
		}
	}

	for i := 0; i < 32; i++ {
		bit := (ip[i/8] >> (7 - i%8)) & 1
		if node.children[bit] == nil {
			break
		}
		node = node.children[bit]
		if node.network != "" {
			network = node.network
		}
	}

	return network
}

func (t *Trie) Locate(ip net.IP) string {
	return t.FindNetwork(ip)
}
