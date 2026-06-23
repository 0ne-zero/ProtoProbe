package dns_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/miekg/dns"
)

// selfSignedCert generates a self-signed TLS certificate for 127.0.0.1.
func selfSignedCert(t *testing.T) tls.Certificate {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("selfSignedCert: generate key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "127.0.0.1"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("selfSignedCert: create cert: %v", err)
	}

	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatalf("selfSignedCert: marshal key: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		t.Fatalf("selfSignedCert: X509KeyPair: %v", err)
	}
	return cert
}

// dnsHandler is a simple DNS handler that responds to all A queries with 127.0.0.1.
func dnsHandler(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	for _, q := range r.Question {
		if q.Qtype == dns.TypeA {
			m.Answer = append(m.Answer, &dns.A{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    60,
				},
				A: net.ParseIP("127.0.0.1"),
			})
		}
	}
	w.WriteMsg(m) //nolint:errcheck
}

// startDNSServer starts a miekg/dns server on a random local port using the
// given network ("udp" or "tcp"). Returns the host and port it bound to.
func startDNSServer(t *testing.T, network string) (host string, port int) {
	t.Helper()

	started := make(chan struct{})

	var srv *dns.Server

	switch network {
	case "udp":
		pc, err := net.ListenPacket("udp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("startDNSServer(udp): ListenPacket: %v", err)
		}
		addr := pc.LocalAddr().(*net.UDPAddr)
		host = addr.IP.String()
		port = addr.Port

		srv = &dns.Server{
			PacketConn:       pc,
			Net:              "udp",
			Handler:          dns.HandlerFunc(dnsHandler),
			NotifyStartedFunc: func() { close(started) },
		}
	case "tcp":
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("startDNSServer(tcp): Listen: %v", err)
		}
		addr := ln.Addr().(*net.TCPAddr)
		host = addr.IP.String()
		port = addr.Port

		srv = &dns.Server{
			Listener:         ln,
			Net:              "tcp",
			Handler:          dns.HandlerFunc(dnsHandler),
			NotifyStartedFunc: func() { close(started) },
		}
	default:
		t.Fatalf("startDNSServer: unknown network %q", network)
	}

	go func() {
		if err := srv.ActivateAndServe(); err != nil {
			// Server stopped; this is expected on cleanup.
		}
	}()

	<-started

	t.Cleanup(func() { srv.Shutdown() }) //nolint:errcheck

	return host, port
}
