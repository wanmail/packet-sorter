package source

import (
	"bufio"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/wanmail/packet-sorter/pkg/sorter"
)

const (
	FortigateForwardLog = `srcip=(?P<srcip>\d+\.\d+\.\d+\.\d+)\s+(?:srcport=(?P<srcport>\d+))?.+dstip=(?P<dstip>\d+\.\d+\.\d+\.\d+)\s+(?:dstport=(?P<dstport>\d+))?.+\s+action="(?P<action>[^"]+)"\s+policyid=(?P<policyid>\d+)`
)

func ParseFortigateLogPath(path string, policy string, ch chan<- sorter.Packet) error {
	return filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		return ParseFortigateLogFile(file, policy, ch)
	})
}

func ParseFortigateLogFile(file string, policy string, ch chan<- sorter.Packet) error {
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	re := regexp.MustCompile(FortigateForwardLog)

	scanner := bufio.NewScanner(fd)

	for scanner.Scan() {
		line := scanner.Text()

		match := re.FindStringSubmatch(line)
		if match == nil {
			continue
			// return errors.Errorf("No match found: %s", line)
		}

		var sip, dip net.IP
		var sport, dport int

		if policy == match[6] {
			sip = net.ParseIP(match[1])
			if match[2] == "" {
				sport = 0
			} else {
				sport, err = strconv.Atoi(match[2])
				if err != nil {
					return err
				}
			}

			dip = net.ParseIP(match[3])

			if match[4] == "" {
				dport = 0
			} else {
				dport, err = strconv.Atoi(match[4])
				if err != nil {
					return err
				}
			}

			ch <- sorter.Packet{
				SourceIP:   sip,
				SourcePort: sport,
				DestIP:     dip,
				DestPort:   dport,
				Result:     match[5],
			}
		}

	}

	return scanner.Err()
}
