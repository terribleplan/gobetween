package proxyprotocol

import (
	"fmt"
	"net"

	proxyproto "github.com/pires/go-proxyproto"
)

func isIp4(addr net.Addr) (bool, error) {
	ipString, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		return false, err
	}

	ip := net.ParseIP(ipString)
	if ip == nil {
		return false, fmt.Errorf("Could not parse IP")
	}

	return ip.To4() != nil, nil
}

/// SendProxyProtocolV1 sends a proxy protocol v1 header to initialize the connection
/// https://www.haproxy.org/download/1.8/doc/proxy-protocol.txt
func SendProxyProtocolV1(client net.Conn, backend net.Conn) error {
	remoteAddr := client.RemoteAddr()
	h := proxyproto.Header{
		Version:         1,
		SourceAddr:      client.RemoteAddr(),
		DestinationAddr: client.LocalAddr(),
	}

	if ip4, err := isIp4(remoteAddr); err != nil {
		return err
	} else if ip4 {
		h.TransportProtocol = proxyproto.TCPv4
	} else {
		h.TransportProtocol = proxyproto.TCPv6
	}

	_, err := h.WriteTo(backend)
	if err != nil {
		return nil
	}
	return nil
}
