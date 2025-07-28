package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func EncodeMessage(msg interface{}) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

type BaseMessage struct {
	Method string `json:"method"`
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, errors.New("Separator Not Found")
	}

	contentlenghtbyte := header[len("Content-Length: "):]
	contentlen, err := strconv.Atoi(string(contentlenghtbyte))
	if err != nil {
		return "", nil, nil
	}
	var baseMessage BaseMessage
	if err := json.Unmarshal(content[:contentlen], &baseMessage); err != nil {
		return "", nil, nil
	}
	_ = content
	return baseMessage.Method, content[:contentlen], nil
}

func Split(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return 0, nil, nil
	}

	contentLenghtByte := header[len("Content-Length: "):]
	contentlen, err := strconv.Atoi(string(contentLenghtByte))
	if err != nil {
		return 0, nil, err
	}
	if len(content) < contentlen {
		return 0, nil, nil
	}
	totalLen := len(header) + 4 + contentlen
	return totalLen, data[:totalLen], nil
}
