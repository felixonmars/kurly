package main

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestSetHeaders(t *testing.T) {
	t.Log("Testing setHeaders()... (expecting Kurly/1.0)")

	header := []string{"User-Agent: Kurly/1.0"}
	req, err := http.NewRequest("GET", "http://url.com/", nil)
	if err != nil {
		panic(err)
	}

	setHeaders(req, header)

	if len(req.Header) > 0 {
		for _, v := range req.Header {
			userAgentValue := strings.Join(v, "")
			if userAgentValue != "Kurly/1.0" {
				t.Errorf("Expected Kurly/1.0, but got %s", userAgentValue)
			}
		}
	} else {
		t.Error("setHeaders() set no header")
	}
}

func TestMaxTime(t *testing.T) {
	t.Log("Testing maxTime()... (expecting no timeout)")

	maxTime(2)
	time.Sleep(1 * time.Second)
}
