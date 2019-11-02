package cli

import (
	"github.com/everett-protocol/everett-hackathon/x/lsp/internal/types"
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
		Short:                      "Liquid staking position",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(client.PostCommands(GetCmdTransfer(cdc))...)

	return txCmd
}

func GetCmdTransfer(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "delegate [transfer_chan_id] [ia_chan_id] [amount] [validator]",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			transferChanId := args[0]
			iaChanId := args[1]

			amount, err := sdk.ParseCoin(args[2])
			if err != nil {
				return err
			}

			// TODO: Can parse bech32 according to counterparty chain.
			validator, err := sdk.ValAddressFromBech32(args[3])
			if err != nil {
				return err
			}

			msg := types.NewMsgOpenLiquidStakingPosition(transferChanId, iaChanId, amount, validator, cliCtx.GetFromAddress())
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
