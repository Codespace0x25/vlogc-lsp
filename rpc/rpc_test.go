package rpc_test

import (
	"testing"
	"vlogcc-lsp/rpc"
)

type TestingEncode struct {
	Testing bool
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
	actual := rpc.EncodeMessage(TestingEncode{Testing: true})
	if expected != actual {
		t.Fatalf("Expected %s, Actual %s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	incomingMessage := "Content-Length: 15\r\n\r\n{\"Method\":\"hi\"}"
	method, content, err := rpc.DecodeMessage([]byte(incomingMessage))
	contentlen := len(content)
	if err != nil {
		t.Fatal(err)
	}
	if contentlen != 15 {
		t.Fatalf("Expected 15, Actual %d", contentlen)
	}
	if method != "hi" {
		t.Fatalf("Expected hi, Actual %s", method)
	}
}
