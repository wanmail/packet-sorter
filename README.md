# Packet Sorter

When you join a new company, you need to restrict firewall or ACL permissions. How to identify zone-to-zone access traffic and assess the impact of policy changes ? 

This is a Go-based tool for analyzing and sorting network traffic logs, to help you determine how to optimize your policy.

## Features

- Log Type
    - Fortigate firewall logs
    - AWS VPC flow log (TODO)
- Sort network traffic based on source/destination IP zones

## Prerequisites

- Go 1.23.2 or higher
- Make (for building)

## Installation

Clone the repository:
```
git clone https://github.com/wanmail/packet-sorter
```

Build the binary for your platform:
```
make build-x64
```

## Usage

The tool supports several command-line flags:
```
packet-sorter-x64 -h

Usage of packet-sorter-x64:
  -fortigate_policy string
        Fortigate policy ID
  -log_file string
        Log file name
  -log_path string
        Path to the log file
  -log_type string
        Log file type. 
        Supported types: fortigate-log (default "fortigate-log")
  -out_file string
        Output file name (default "flow.csv")
  -zone_file string
        Zone file name
```

### Zone File Format
The zone file should contain CIDR ranges, one per line:
```
10.0.0.0/8
192.168.0.0/16
172.16.0.0/12
```

If no zone file is provided, the tool defaults to using the 10.0.0.0/8 range, split it to 256 * 256 C section.


### Output Format

The tool generates a CSV file with the following columns:
- Sip: Source IP zone
- Dip: Destination IP zone
- Port: Destination port
- Result: Action result
- Count: Number of matching packets
