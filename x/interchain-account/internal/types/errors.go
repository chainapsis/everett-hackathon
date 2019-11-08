package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeAccountAlreadyExist  sdk.CodeType = 101
	CodeUnsupportedChianType sdk.CodeType = 102
)

func ErrAccountAlreadyExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeAccountAlreadyExist, "account already exists")
}

func ErrUnsupportedChainType(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeUnsupportedChianType, "unsupported chain type")
}
