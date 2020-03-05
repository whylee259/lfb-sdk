package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
	"github.com/line/link/x/contract"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	accountKeeper  types.AccountKeeper
	iamKeeper      types.IamKeeper
	contractKeeper contract.Keeper
	storeKey       sdk.StoreKey
	cdc            *codec.Codec
}

func NewKeeper(cdc *codec.Codec, accountKeeper types.AccountKeeper, iamKeeper types.IamKeeper, contractKeeper contract.Keeper, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		accountKeeper:  accountKeeper,
		iamKeeper:      iamKeeper.WithPrefix(types.ModuleName),
		contractKeeper: contractKeeper,
		storeKey:       storeKey,
		cdc:            cdc,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) UnmarshalJSON(bz []byte, ptr interface{}) error {
	return k.cdc.UnmarshalJSON(bz, ptr)
}

func (k Keeper) MarshalJSON(o interface{}) ([]byte, error) {
	return k.cdc.MarshalJSON(o)
}

func (k Keeper) MarshalJSONIndent(o interface{}) ([]byte, error) {
	return k.cdc.MarshalJSONIndent(o, "", "  ")
}

func (k Keeper) mustEncodeString(str string) []byte {
	return k.cdc.MustMarshalBinaryBare(str)
}

func (k Keeper) mustDecodeString(bz []byte) (str string) {
	k.cdc.MustUnmarshalBinaryBare(bz, &str)
	return str
}
