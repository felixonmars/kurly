package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptrace"
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

func NewClientTraceForRequest(req *http.Request) *httptrace.ClientTrace {
	if req == nil {
		Status.Fatal("cannot do a verbose trace for a empty request")
	}

	ts := &tracerStruct{req: req}

	return &httptrace.ClientTrace{
		ConnectStart: ts.ConnectStart,
		ConnectDone:  ts.ConnectDone,
		WroteHeaders: ts.WroteHeaders,
		GotConn:      ts.GotConn,
		DNSStart:     ts.DNSStart,
	}
}
