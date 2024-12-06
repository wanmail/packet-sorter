package sorter

import (
	"net"
	"testing"
)

func TestTrie_InsertAndFindNetwork(t *testing.T) {
	trie := NewTrie()

	_, cidr1, _ := net.ParseCIDR("192.168.1.0/24")
	_, cidr2, _ := net.ParseCIDR("192.168.1.0/25")
	_, cidr3, _ := net.ParseCIDR("192.168.1.128/25")
	_, cidr4, _ := net.ParseCIDR("10.0.0.0/8")
	_, cidr5, _ := net.ParseCIDR("10.1.0.0/16")
	_, cidr6, _ := net.ParseCIDR("0.0.0.0/0")

	trie.Insert(cidr1)
	trie.Insert(cidr2)
	trie.Insert(cidr3)
	trie.Insert(cidr4)
	trie.Insert(cidr5)
	trie.Insert(cidr6)

	tests := []struct {
		ip       string
		expected string
	}{
		{"192.168.1.1", "192.168.1.0/25"},
		{"192.168.1.129", "192.168.1.128/25"},
		{"192.168.2.1", "0.0.0.0/0"},
		{"10.0.0.1", "10.0.0.0/8"},
		{"10.1.0.1", "10.1.0.0/16"},
		{"1.1.1.1", "0.0.0.0/0"},
	}

	for _, test := range tests {
		ip := net.ParseIP(test.ip)
		result := trie.FindNetwork(ip)
		if result != test.expected {
			t.Errorf("FindNetwork(%s) = %v; want %v", test.ip, result, test.expected)
		}
	}
}
