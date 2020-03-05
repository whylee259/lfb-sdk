package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

type CollectionKeeper interface {
	CreateCollection(ctx sdk.Context, collection types.Collection, owner sdk.AccAddress) sdk.Error
	ExistCollection(ctx sdk.Context, contractID string) bool
	GetCollection(ctx sdk.Context, contractID string) (collection types.Collection, err sdk.Error)
	SetCollection(ctx sdk.Context, collection types.Collection) sdk.Error
	UpdateCollection(ctx sdk.Context, collection types.Collection) sdk.Error
	GetAllCollections(ctx sdk.Context) types.Collections
}

var _ CollectionKeeper = (*Keeper)(nil)

func (k Keeper) NewContractID(ctx sdk.Context) string {
	return k.contractKeeper.NewContractID(ctx)
}

func (k Keeper) HasContractID(ctx sdk.Context, contractID string) bool {
	return k.contractKeeper.HasContractID(ctx, contractID)
}

func (k Keeper) CreateCollection(ctx sdk.Context, collection types.Collection, owner sdk.AccAddress) sdk.Error {
	err := k.SetCollection(ctx, collection)
	if err != nil {
		return err
	}
	k.SetSupply(ctx, types.DefaultSupply(collection.GetContractID()))

	perms := types.Permissions{
		types.NewIssuePermission(collection.GetContractID()),
		types.NewMintPermission(collection.GetContractID()),
		types.NewBurnPermission(collection.GetContractID()),
		types.NewModifyPermission(collection.GetContractID()),
	}
	for _, perm := range perms {
		k.AddPermission(ctx, owner, perm)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateCollection,
			sdk.NewAttribute(types.AttributeKeyContractID, collection.GetContractID()),
			sdk.NewAttribute(types.AttributeKeyName, collection.GetName()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
		),
		sdk.NewEvent(
			types.EventTypeGrantPermToken,
			sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
			sdk.NewAttribute(types.AttributeKeyResource, collection.GetContractID()),
		),
	})
	for _, perm := range perms {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeGrantPermToken,
				sdk.NewAttribute(types.AttributeKeyAction, perm.GetAction()),
			),
		})
	}

	return nil
}

func (k Keeper) ExistCollection(ctx sdk.Context, contractID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.CollectionKey(contractID))
}

func (k Keeper) GetCollection(ctx sdk.Context, contractID string) (collection types.Collection, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CollectionKey(contractID))
	if bz == nil {
		return collection, types.ErrCollectionNotExist(types.DefaultCodespace, contractID)
	}

	collection = k.mustDecodeCollection(bz)
	return collection, nil
}

func (k Keeper) SetCollection(ctx sdk.Context, collection types.Collection) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	if store.Has(types.CollectionKey(collection.GetContractID())) {
		return types.ErrCollectionExist(types.DefaultCodespace, collection.GetContractID())
	}

	store.Set(types.CollectionKey(collection.GetContractID()), k.cdc.MustMarshalBinaryBare(collection))
	k.setNextTokenTypeFT(ctx, collection.GetContractID(), types.ReservedEmpty)
	k.setNextTokenTypeNFT(ctx, collection.GetContractID(), types.ReservedEmptyNFT)
	return nil
}

func (k Keeper) UpdateCollection(ctx sdk.Context, collection types.Collection) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.CollectionKey(collection.GetContractID())) {
		return types.ErrCollectionNotExist(types.DefaultCodespace, collection.GetContractID())
	}

	store.Set(types.CollectionKey(collection.GetContractID()), k.cdc.MustMarshalBinaryBare(collection))
	return nil
}

func (k Keeper) GetAllCollections(ctx sdk.Context) types.Collections {
	var collections types.Collections
	appendCollection := func(collection types.Collection) (stop bool) {
		collections = append(collections, collection)
		return false
	}
	k.iterateCollections(ctx, "", appendCollection)
	return collections
}

func (k Keeper) iterateCollections(ctx sdk.Context, contractID string, process func(types.Collection) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.CollectionKey(contractID))
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		collection := k.mustDecodeCollection(val)
		if process(collection) {
			return
		}
		iter.Next()
	}
}

func (k Keeper) mustDecodeCollection(collectionByte []byte) types.Collection {
	var collection types.Collection
	err := k.cdc.UnmarshalBinaryBare(collectionByte, &collection)
	if err != nil {
		panic(err)
	}
	return collection
}
