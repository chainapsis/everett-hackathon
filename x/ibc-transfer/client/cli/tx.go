package cli

import (
	"github.com/everett-protocol/everett-hackathon/x/ibc-transfer/internal/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Send assets between block chain",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(client.PostCommands(GetCmdTransfer(cdc))...)

	return txCmd
}

func GetCmdTransfer(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "transfer [from_key_or_address] [to_address] [amount] [to_chain_id]",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			// TODO: Can parse bech32 according to counterparty chain.
			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgTransfer(args[3], coins, cliCtx.GetFromAddress(), to)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
