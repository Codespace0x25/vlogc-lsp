package lsp

import (
	"fmt"
)

type Symbol struct {
	Name       string
	Type       string
	DataType   string
	Parameters string
	Mutable    bool
}

func NewSymbol(name, symbolType, dataType, parameters string, mutable bool) Symbol {
	return Symbol{
		Name:       name,
		Type:       symbolType,
		DataType:   dataType,
		Parameters: parameters,
		Mutable:    mutable,
	}
}

func (s Symbol) Validate() error {
	if s.Name == "" || s.Type == "" {
		return fmt.Errorf("symbol missing required fields")
	}
	return nil
}

