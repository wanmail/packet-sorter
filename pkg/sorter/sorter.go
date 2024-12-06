package sorter

import "fmt"

type Sorter struct {
	Locator Locator

	Traffics map[string]map[string]map[int]int
}

func (s *Sorter) SortPackets(c <-chan Packet) error {
	for packet := range c {
		src := s.Locator.Locate(packet.SourceIP)
		dst := s.Locator.Locate(packet.DestIP)

		path := fmt.Sprintf("%s -> %s", src, dst)

		flows, ok := s.Traffics[packet.Result]
		if !ok {
			flows = make(map[string]map[int]int)
			s.Traffics[packet.Result] = flows
		}

		if conn, ok := flows[path]; ok {
			conn[packet.DestPort] += 1
		} else {
			flows[path] = map[int]int{packet.DestPort: 1}
		}
	}
	return nil
}

func NewSorter(locator Locator) *Sorter {
	return &Sorter{
		Locator:  locator,
		Traffics: make(map[string]map[string]map[int]int),
	}
}
