package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httptrace"
	"strings"
)

type tracerStruct struct {
	req         *http.Request
	redirects   int
	currentHost string
}

func (ts *tracerStruct) DNSStart(dnsinfo httptrace.DNSStartInfo) {
	ts.currentHost = dnsinfo.Host
}

func (ts *tracerStruct) ConnectStart(network, addr string) {
	hs, _, _ := net.SplitHostPort(addr)
	Status.Printf("   Trying %s...\n", hs)
}

func (ts *tracerStruct) WroteHeaders() {
	fmt.Fprintln(Outgoing, ts.req.Method, ts.req.URL.Path, ts.req.Proto)
	for k, v := range ts.req.Header {
		fmt.Fprintln(Outgoing, k, v)
	}
	fmt.Fprintln(Outgoing)
}

func (ts *tracerStruct) ConnectDone(network, addr string, err error) {
	Status.Println(" TCP_NODELAY set")
	hs, port, _ := net.SplitHostPort(addr)
	Status.Printf(" Connected to %s (%s) port %s (#%d)%s\n", ts.currentHost, hs, port, ts.redirects, ts.req.RequestURI)
	ts.redirects += 1
}

func (ts *tracerStruct) GotConn(cinfo httptrace.GotConnInfo) {
	if cinfo.Reused {
		hs, port, _ := net.SplitHostPort(cinfo.Conn.RemoteAddr().String())
		Status.Printf(" Re-using existing connection! (#%d) with host %s\n", ts.redirects-1, ts.currentHost)
		Status.Printf(" Connected to %s (%s) port %s (#%d)%s\n", ts.currentHost, hs, port, ts.redirects, ts.req.RequestURI)
	}
}

func (ts *tracerStruct) TLSHandshakeDone(cstate tls.ConnectionState, err error) {
	if !cstate.HandshakeComplete {
		Status.Printf(" TLS Handshake not completed")
		return
	}
	var tlsversion, ciphersuite string
	switch cstate.Version {
	case tls.VersionSSL30:
		tlsversion = "SSLv3"
	case tls.VersionTLS10:
		tlsversion = "TLSv1.0"
	case tls.VersionTLS11:
		tlsversion = "TLSv1.1"
	case tls.VersionTLS12:
		tlsversion = "TLSv1.2"
	}

	ciphersuite = getCipherSuiteString(cstate.CipherSuite)
	Status.Printf(" APLN, server accepted to use %s", cstate.NegotiatedProtocol)
	Status.Printf(" %s, TLS Handshake finished", tlsversion)
	Status.Printf(" SSL connection using %s / %s", tlsversion, ciphersuite)
	if len(cstate.PeerCertificates) > 0 {
		cert := cstate.PeerCertificates[0]
		Status.Println(" Server certificate:")
		Status.Printf("  subject: CN=%s\n", cert.Subject.CommonName)
		Status.Printf("  start date: %s\n", cert.NotBefore.Format("Mon, 02 Jan 2006 15:04:05 MST"))
		Status.Printf("  expire date: %s\n", cert.NotAfter.Format("Mon, 02 Jan 2006 15:04:05 MST"))
		Status.Printf("  issuer: C=%s; O=%s; CN=%s\n",
			strings.Join(cert.Issuer.Country, " "),
			strings.Join(cert.Issuer.Organization, " "),
			cert.Issuer.CommonName)
		Status.Printf("  SSL certificate verify ok.\n")
	}
}

func NewClientTraceForRequest(req *http.Request) *httptrace.ClientTrace {
	if req == nil {
		Status.Fatal("cannot do a verbose trace for a empty request")
	}

	ts := &tracerStruct{req: req}

	return &httptrace.ClientTrace{
		ConnectStart:     ts.ConnectStart,
		ConnectDone:      ts.ConnectDone,
		WroteHeaders:     ts.WroteHeaders,
		GotConn:          ts.GotConn,
		DNSStart:         ts.DNSStart,
		TLSHandshakeDone: ts.TLSHandshakeDone,
	}
}

func getCipherSuiteString(num uint16) string {
	switch num {
	case tls.TLS_RSA_WITH_RC4_128_SHA:
		return "RSA-RC4-128-SHA"
	case tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA:
		return "RSA-3DES-EDE-CBC-SHA"
	case tls.TLS_RSA_WITH_AES_128_CBC_SHA:
		return "RSA-AES-128-CBC-SHA"
	case tls.TLS_RSA_WITH_AES_256_CBC_SHA:
		return "RSA-AES-256-CBC-SHA"
	case tls.TLS_RSA_WITH_AES_128_CBC_SHA256:
		return "RSA-AES-128-CBC-SHA256"
	case tls.TLS_RSA_WITH_AES_128_GCM_SHA256:
		return "RSA-AES-128-GCM-SHA256"
	case tls.TLS_RSA_WITH_AES_256_GCM_SHA384:
		return "RSA-AES-256-GCM-SHA384"
	case tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA:
		return "ECDHE-ECDSA-RC4-128-SHA"
	case tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA:
		return "ECDHE-ECDSA-AES-128-CBC-SHA"
	case tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA:
		return "ECDHE-ECDSA-AES-256-CBC-SHA"
	case tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA:
		return "ECDHE-RSA-RC4-128-SHA"
	case tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA:
		return "ECDHE-RSA-3DES-EDE-CBC-SHA"
	case tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA:
		return "ECDHE-RSA-AES-128-CBC-SHA"
	case tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA:
		return "ECDHE-RSA-AES-256-CBC-SHA"
	case tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256:
		return "ECDHE-ECDSA-AES-128-CBCSHA256"
	case tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256:
		return "ECDHE-RSA-AES-128-CBC-SHA256"
	case tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:
		return "ECDHE-RSA-AES-128-GCM-SHA256"
	case tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256:
		return "ECDHE-ECDSA-AES-128-GCM-SHA256"
	case tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384:
		return "ECDHE-RSA-AES-256-GCM-SHA384"
	case tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384:
		return "ECDHE-ECDSA-AES-256-GCM-SHA384"
	case tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305:
		return "ECDHE-RSA-CHACHA20-POLY1305"
	case tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305:
		return "ECDHE-ECDSA-CHACHA20-POLY1305"
	}
	return "FALLBACK-SCSV"
}
