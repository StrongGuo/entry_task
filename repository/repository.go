package repository

import (
	"context"
	"item-admin/proto"

	"item-admin/model"
)

type ItemRepository interface {
	CreateItem(ctx context.Context, i model.Items) (int64, error)
	UpdateItem(ctx context.Context, i model.Items) (int64, error)
	RemoveItem(ctx context.Context, i model.Items) (int64, error)
	GetItem(ctx context.Context, a model.ActivityItems) (ra model.ActivityItems, err error)
	GetPromotions(ctx context.Context, i model.Items) (p []*proto.Promotion, err error)
	CreatePromotion(ctx context.Context, p model.Promotions) (int64, error)
	RemovePromotion(ctx context.Context, p model.Promotions) (int64, error)
}

type ItemRedisRepository interface {
	SetItemCash(itemID int64, i model.ActivityItems) error
	GetItemCash(ItemID int64) (i *model.ActivityItems, hasData bool, err error)
	InvalidItemCash(ItemID int64) error
}
