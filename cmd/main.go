package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/wanmail/packet-sorter/pkg/sorter"
	"github.com/wanmail/packet-sorter/pkg/source"
)

var (
	LogPath string
	LogFile string
	LogType string

	ZoneFile string

	OutFile string

	FortigatePolicy string
)

func init() {
	flag.StringVar(&LogPath, "log_path", "", "Path to the log file")
	flag.StringVar(&LogFile, "log_file", "", "Log file name")
	flag.StringVar(&LogType, "log_type", FortigateLog, "Log file type. \nSupported types: fortigate-log")
	flag.StringVar(&ZoneFile, "zone_file", "", "Zone file name")
	flag.StringVar(&OutFile, "out_file", "flow.csv", "Output file name")
	flag.StringVar(&FortigatePolicy, "fortigate_policy", "", "Fortigate policy ID")
}

const (
	FortigateLog = "fortigate-log"
)

func main() {
	flag.Parse()

	if LogPath == "" && LogFile == "" {
		panic("Log path or file name are required")
	}

	if LogType == "" {
		panic("Log type is required")
	}

	trie := sorter.NewTrie()
	if ZoneFile == "" {
		fmt.Println("Zone file is not provided. Using default 10.0.0.0/8")
		for i := 0; i <= 255; i++ {
			for j := 0; j <= 255; j++ {
				ipRange := &net.IPNet{
					IP:   net.IPv4(10, byte(i), byte(j), 0),
					Mask: net.CIDRMask(24, 32),
				}
				trie.Insert(ipRange)
			}
		}
		trie.Insert(&net.IPNet{
			IP:   net.IPv4(0, 0, 0, 0),
			Mask: net.CIDRMask(0, 32),
		})
	} else {
		// Read zone file
		fd, err := os.Open(ZoneFile)
		if err != nil {
			panic(err)
		}
		defer fd.Close()

		scanner := bufio.NewScanner(fd)

		for scanner.Scan() {
			line := scanner.Text()
			_, cidr, err := net.ParseCIDR(line)
			if err != nil {
				fmt.Printf("Error parsing CIDR: %s %s\n", line, err)
			}
			trie.Insert(cidr)
		}
	}

	s := sorter.NewSorter(trie)

	wg := sync.WaitGroup{}
	ch := make(chan sorter.Packet)

	wg.Add(1)
	go func() {
		err := s.SortPackets(ch)
		if err != nil {
			fmt.Println(err)
		}
		wg.Done()
	}()

	switch LogType {
	case FortigateLog:
		wg.Add(1)
		go func() {
			if LogFile != "" {
				err := source.ParseFortigateLogFile(LogPath+LogFile, FortigatePolicy, ch)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				err := source.ParseFortigateLogPath(LogPath, FortigatePolicy, ch)
				if err != nil {
					fmt.Println(err)
				}
			}

			wg.Done()
			close(ch)
		}()

	default:
		panic("Unsupported log type")
	}

	wg.Wait()

	out, err := os.Create(OutFile)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	writer := csv.NewWriter(out)
	defer writer.Flush()

	writer.Write([]string{"Sip", "Dip", "Port", "Result", "Count"})

	for result, flows := range s.Traffics {
		for path, conn := range flows {
			var sip, dip string
			l := strings.Split(path, " -> ")
			if len(l) != 2 {
				sip = path
				dip = path
			} else {
				sip = l[0]
				dip = l[1]
			}

			for port, count := range conn {
				writer.Write([]string{sip, dip, fmt.Sprintf("%d", port), result, fmt.Sprintf("%d", count)})
			}
		}
	}
}
