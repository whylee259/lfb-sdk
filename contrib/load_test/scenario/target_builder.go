package scenario

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

const (
	TxURL            = "/txs"
	QuerySimulateURL = TxURL + "/simulate"
)

type TargetBuilder struct {
	cdc    *codec.Codec
	LCDURL string
}

func NewTargetBuilder(cdc *codec.Codec, lcdURL string) *TargetBuilder {
	return &TargetBuilder{
		cdc:    cdc,
		LCDURL: lcdURL,
	}
}

func (tb *TargetBuilder) MakeQueryTarget(url string) *vegeta.Target {
	return &vegeta.Target{
		Method: "GET",
		URL:    tb.LCDURL + url,
	}
}

func (tb *TargetBuilder) MakeTxTarget(stdTx auth.StdTx, mode string) (target *vegeta.Target, err error) {
	return tb.makeTxTarget(stdTx, mode, TxURL)
}

func (tb *TargetBuilder) MakeQuerySimulateTarget(stdTx auth.StdTx, mode string) (target *vegeta.Target, err error) {
	return tb.makeTxTarget(stdTx, mode, QuerySimulateURL)
}

func (tb *TargetBuilder) makeTxTarget(stdTx auth.StdTx, mode string, url string) (target *vegeta.Target, err error) {
	bz, err := tb.cdc.MarshalJSON(rest.BroadcastReq{Mode: mode, Tx: stdTx})
	if err != nil {
		return
	}

	return &vegeta.Target{
		Method: "POST",
		URL:    tb.LCDURL + url,
		Body:   bz,
	}, nil
}
