package lsp

import "errors"

type Symbol struct {
	Name string
	Type string
	DataType string
	Parameters string
}

func NewSymbol(name, symbolType, dataType, parameters, methods string) (Symbol, error) {
	panic("TODO: Not implementated")
}

func (s Symbol) Validate() error {
	panic("TODO: Not implementated")
}

