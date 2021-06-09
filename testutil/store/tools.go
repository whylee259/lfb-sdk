package store

import (
	"testing"

	types2 "github.com/line/lfb-sdk/codec/types"
	"github.com/line/lfb-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
)

func ValFmt(i int) *types.BaseAccount {
	return &types.BaseAccount{
		AccountNumber: uint64(i),
	}
}

func CreateTestInterfaceRegistry() types2.InterfaceRegistry {
	interfaceRegistry := types2.NewInterfaceRegistry()
	types.RegisterInterfaces(interfaceRegistry)
	return interfaceRegistry
}

func VerifyValue(t *testing.T, target *types.BaseAccount, source interface{}) {
	sourceValue := source.(*types.BaseAccount)
	require.Equal(t, *target, *sourceValue)
}
