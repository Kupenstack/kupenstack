package vn

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

type Uint128 struct {
	Hi, Lo uint64
}

func BigEndianUint128(b []byte) Uint128 {
	hi := binary.BigEndian.Uint64(b[:8])
	lo := binary.BigEndian.Uint64(b[8:])
	return Uint128{hi, lo}
}

func translateTargetIPInCIDRv4Range(targetIP net.IP, firstIP net.IP, maskIP net.IPMask) net.IP {
	first := binary.BigEndian.Uint32(firstIP)
	mask := binary.BigEndian.Uint32(maskIP)
	target := binary.BigEndian.Uint32(targetIP)

	notMask := (mask ^ 0xffffffff)

	// Takes first part of IP from network and second part of the ip from targetIP
	target = (first & mask) | (target & notMask)

	// returns 4-byte net.IP with required targetIP
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, target)
	return ip
}

func translateTargetIPInCIDRv6Range(targetIP net.IP, firstIP net.IP, maskIP net.IPMask) net.IP {
	first := BigEndianUint128(firstIP)
	mask := BigEndianUint128(maskIP)
	target := BigEndianUint128(targetIP)

	var notMask Uint128
	notMask.Hi = (mask.Hi ^ 0xffffffffffffffff)
	notMask.Lo = (mask.Lo ^ 0xffffffffffffffff)

	target.Hi = (first.Hi & mask.Hi) | (target.Hi & notMask.Hi)
	target.Lo = (first.Lo & mask.Lo) | (target.Lo & notMask.Lo)

	// returns 16-byte net.IP with required targetIP
	ip := make(net.IP, 16)
	binary.BigEndian.PutUint64(ip[:8], target.Hi)
	binary.BigEndian.PutUint64(ip[8:], target.Lo)
	return ip
}

func translateTargetIPInCIDR(IPpattern, cidr string) (string, error) {

	// assume the network is v4 by default.
	CidrIsV4 := true

	_, ipnetwork, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}

	if ipnetwork.IP.To4() == nil {
		CidrIsV4 = false
	}

	formatedIP := strings.ReplaceAll(IPpattern, "*", "0")
	targetIP := net.ParseIP(formatedIP)

	// if not valid IP pattern
	if targetIP.To4() == nil && CidrIsV4 == true {
		CidrIsV4 = false
		return "", fmt.Errorf("Given IP pattern %+v does match given v4 cidr %+v format", IPpattern, cidr)
	}

	// Convert IP to 4-byte or 16-byte based on v4 or v6
	if CidrIsV4 == true {
		targetIP = targetIP.To4()
	} else {
		targetIP = targetIP.To16()

		// if not valid v6 IP
		if targetIP == nil {
			return "", fmt.Errorf("Invalid v6 IP pattern %+v", IPpattern)
		}
	}

	// Translate IP pattern into IP in CIDR Range
	var ip net.IP
	if CidrIsV4 == true {
		ip = translateTargetIPInCIDRv4Range(targetIP, ipnetwork.IP, ipnetwork.Mask)
	} else {
		ip = translateTargetIPInCIDRv6Range(targetIP, ipnetwork.IP, ipnetwork.Mask)
	}

	return ip.String(), nil
}

func getStartIPOf(cidr string) (string, error) {
	_, n, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}
	for i := range n.IP {
		n.IP[i] &= n.Mask[i]
	}
	return n.IP.String(), nil
}

func getEndIPOf(cidr string) (string, error) {
	_, n, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}
	for i := range n.IP {
		n.IP[i] |= ^n.Mask[i]
	}
	return n.IP.String(), nil
}
