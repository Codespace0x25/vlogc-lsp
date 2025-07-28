package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"vlogcc-lsp/lsp"
	"vlogcc-lsp/rpc"
)

func main() {
	logger := getLogger("/home/xsoder/programming/vlogcc-lsp/log.txt")
	logger.Println("Started Logging")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s", err)
		}
		handleMessage(logger, method, contents)
	}
}

func handleMessage(logger *log.Logger, method string, contents []byte) {
	logger.Printf("Recieved Message with: %s", method)
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Could not parse this text: %s", err)
		}
		logger.Printf("Connected to %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
		msg := lsp.NewInitializeResponse(request.ID)
		reply := rpc.EncodeMessage(msg)
		writer := os.Stdout
		writer.Write([]byte(reply))
		logger.Print("Reply Sent")
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Could not parse this text: %s", err)
		}
		logger.Printf("Opened %s %s", request.Params.TextDocument.URI, request.Params.TextDocument.Text)
	}
}

func getLogger(file string) *log.Logger {
	logfile, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Unable to open the file: %v", err))
	}
	return log.New(logfile, "[vlogcc-lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
