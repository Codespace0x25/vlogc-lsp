package lsp

type InitializeRequest struct {
	Request
	Params InializeRequestParams `json:"params"`
}

type InializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverinfo"`
}

type ServerCapabilities struct {
	TextDocumentSync int `json:"textDocumentSync"`
	HoverProvider bool `json:"hoverProvider"`

}


type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: 1,
				HoverProvider : true,
			},
			ServerInfo: ServerInfo{
				Name:    "vlogcc",
				Version: "0.1.0-beta",
			},
		},
	}
}
