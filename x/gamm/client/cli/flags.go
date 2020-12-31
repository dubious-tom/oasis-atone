package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	// Will be parsed to []sdk.Coin
	FlagPoolRecordTokens = "record-tokens"
	// Will be parsed to []sdk.Int
	FlagPoolRecordTokenWeights = "record-tokens-weight"
	// Will be parsed to sdk.Dec
	FlagSwapFee = "swap-fee"
	// Will be parsed to sdk.Dec
	FlagExitFee = "exit-fee"

	FlagPoolId = "pool-id"
	// Will be parsed to sdk.Int
	FlagShareAmountOut = "share-amount-out"
	// Will be parsed to []sdk.Coin
	FlagMaxAountsIn = "max-amounts-in"

	// Will be parsed to sdk.Int
	FlagShareAmountIn = "share-amount-in"
	// Will be parsed to []sdk.Coin
	FlagMinAmountsOut = "min-amounts-out"
)

func FlagSetCreatePool() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.StringArray(FlagPoolRecordTokens, []string{""}, "The tokens to be provided to the pool initially")
	fs.StringArray(FlagPoolRecordTokenWeights, []string{""}, "The weights of the tokens in the pool")
	fs.String(FlagSwapFee, "", "Swap fee of the pool")
	fs.String(FlagExitFee, "", "Exit fee of the pool")
	return fs
}

func FlagSetJoinPool() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.Uint64(FlagPoolId, 0, "The id of pool")
	fs.String(FlagShareAmountOut, "", "TODO: add description")
	fs.StringArray(FlagMaxAountsIn, []string{""}, "TODO: add description")

	return fs
}

func FlagSetExitPool() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.Uint64(FlagPoolId, 0, "The id of pool")
	fs.String(FlagShareAmountIn, "", "TODO: add description")
	fs.StringArray(FlagMinAmountsOut, []string{""}, "TODO: add description")

	return fs
}
