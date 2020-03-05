package types

import (
	"encoding/json"
)

const (
	TokenTypeLength    = 8
	SmallestAlphanum   = "0"
	FungibleFlag       = SmallestAlphanum
	ReservedEmpty      = "00000000"
	SmallestFTType     = "00000001"
	ReservedEmptyNFT   = "10000000"
	SmallestNFTType    = "10000001"
	SmallestTokenIndex = "00000001"
)

type TokenType interface {
	GetName() string
	SetName(string) TokenType
	GetContractID() string
	GetTokenType() string
	String() string
}

type BaseTokenType struct {
	ContractID string `json:"contract_id"`
	TokenType  string `json:"token_type"`
	Name       string `json:"name"`
}

func NewBaseTokenType(contractID, tokenType, name string) TokenType {
	return &BaseTokenType{
		ContractID: contractID,
		TokenType:  tokenType,
		Name:       name,
	}
}
func (t BaseTokenType) GetName() string { return t.Name }
func (t BaseTokenType) SetName(name string) TokenType {
	t.Name = name
	return t
}
func (t BaseTokenType) GetContractID() string { return t.ContractID }
func (t BaseTokenType) GetTokenType() string  { return t.TokenType }
func (t BaseTokenType) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type TokenTypes []TokenType

func (ts TokenTypes) String() string {
	b, err := json.Marshal(ts)
	if err != nil {
		panic(err)
	}
	return string(b)
}
