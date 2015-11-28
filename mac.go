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

func getMAC() (MAC [6]byte) {
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
			copy(MAC[:], itf.HardwareAddr)
			return
		case 8: // EUI-64
			if bytes.Equal(itf.HardwareAddr, zeroMAC[:]) {
				continue
			}
			copy(MAC[:3], itf.HardwareAddr)
			copy(MAC[3:], itf.HardwareAddr[5:])
			return
		}
	}

	return genMAC()
}

func genMAC() (MAC [6]byte) {
	rand.Read(MAC[:])
	MAC[0] |= 0x01 // multicast
	return
}
