package internal

import (
	"bytes"
	"net"

	"github.com/chanxuehong/rand"
)

var MAC [6]byte

func init() {
	MAC = getMAC()
}

var zeroMAC [8]byte

func getMAC() (mac [6]byte) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return genMAC()
	}

	for _, itf := range interfaces {
		if itf.Flags&net.FlagLoopback == net.FlagLoopback ||
			itf.Flags&net.FlagPointToPoint == net.FlagPointToPoint {
			continue
		}

		switch len(itf.HardwareAddr) {
		case 6: // MAC-48, EUI-48
			if bytes.Equal(itf.HardwareAddr, zeroMAC[:6]) {
				continue
			}
			copy(mac[:], itf.HardwareAddr)
			return
		case 8: // EUI-64
			if bytes.Equal(itf.HardwareAddr, zeroMAC[:]) {
				continue
			}
			copy(mac[:3], itf.HardwareAddr)
			copy(mac[3:], itf.HardwareAddr[5:])
			return
		}
	}

	return genMAC()
}

func genMAC() (mac [6]byte) {
	rand.Read(mac[:])
	mac[0] |= 0x01 // multicast
	return
}
