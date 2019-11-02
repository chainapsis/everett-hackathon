package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeAccountAlreadyExist sdk.CodeType = 101
)

func ErrAccountAlreadyExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeAccountAlreadyExist, "account already exists")
}
