package service

import (
	"context"
	"fmt"
	"item-admin/model"
	"item-admin/proto"
	"item-admin/repository"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var repo repository.ItemRepository
var reds repository.ItemRedisRepository

func init() {
	repo = repository.NewMySQL()
	reds = repository.NewRedis()

}

type Server struct {
	proto.UnimplementedItemServiceServer
}

func (s *Server) CreateItem(ctx context.Context, req *proto.CreateItemRequest) (*proto.CreateItemResponse, error) {
	item := req.GetItem()
	timeUnix := time.Now().Unix()

	newItem := model.Items{
		ItemID:         timeUnix,
		ItemName:       item.GetItemName(),
		ItemDesc:       item.GetItemDesc(),
		ItemPrice:      item.GetItemPrice(),
		ItemStock:      item.GetItemStock(),
		Status:         item.GetStatus(),
		CreatorID:      item.GetCreatorID(),
		LastModifierID: item.GetLastModifierID(),
	}
	if newItem.ItemName == "" || newItem.ItemDesc == "" || newItem.CreatorID == 0 || newItem.LastModifierID == 0 || newItem.ItemPrice < 0 || newItem.ItemStock < 0 {
		// 校验invalid params
		return &proto.CreateItemResponse{
			Msg: "invalid params",
		}, nil
	}
	_, err := repo.CreateItem(ctx, newItem)

	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("internal error : %v", err))

	}

	return &proto.CreateItemResponse{
		Msg: "success",
	}, nil
}

func (s *Server) UpdateItem(ctx context.Context, req *proto.UpdateItemRequest) (*proto.UpdateItemResponse, error) {
	u := req.GetUpdateItem()

	UpdateItemInfo := model.Items{
		ItemID:         u.GetItemID(),
		ItemName:       u.GetItemName(),
		ItemDesc:       u.GetItemDesc(),
		ItemPrice:      u.GetItemPrice(),
		ItemStock:      u.GetItemStock(),
		Status:         u.GetStatus(),
		LastModifierID: u.GetLastModifierID(),
	}

	if UpdateItemInfo.ItemName == "" || UpdateItemInfo.ItemDesc == "" || UpdateItemInfo.LastModifierID == 0 || UpdateItemInfo.ItemPrice < 0 || UpdateItemInfo.ItemStock < 0 {
		// 校验非空
		return &proto.UpdateItemResponse{
			Msg: "invalid params",
		}, nil
	}
	//更新数据库
	res, err := repo.UpdateItem(ctx, UpdateItemInfo)
	fmt.Println(res)
	if res != 0 {
		return &proto.UpdateItemResponse{
			Msg: "Nothing updated, please check your ItemID or params",
		}, nil
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("internal error : %v", err))
	}

	//删除缓存
	reds.InvalidItemCash(u.GetItemID())

	return &proto.UpdateItemResponse{Msg: "success"}, nil
}

func (s *Server) RemoveItem(ctx context.Context, req *proto.RemoveItemRequest) (*proto.RemoveItemResponse, error) {
	g := req.GetUniqueItemID()

	RemoveItemInfo := model.Items{
		ItemID: g.GetItemID(),
	}

	_, err := repo.RemoveItem(ctx, RemoveItemInfo)
	reds.InvalidItemCash(g.GetItemID()) //清除缓存

	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("internal error : %v", err))

	}

	return &proto.RemoveItemResponse{}, nil
}

func (s *Server) GetItem(ctx context.Context, req *proto.GetItemRequest) (*proto.GetItemResponse, error) {
	r := req.GetUniqueItemID()

	ItemInfo := model.ActivityItems{
		ItemID: r.GetItemID(),
	}

	//查redis
	cash, flag, _ := reds.GetItemCash(ItemInfo.ItemID)
	itemInfo := proto.ActivityItem{
		ItemID:           r.GetItemID(),
		ItemName:         cash.ItemName,
		ItemDesc:         cash.ItemDesc,
		ItemPrice:        cash.ItemPrice,
		ItemSpecialPrice: cash.ItemSpecialPrice,
		PromotionID:      cash.PromotionID,
	}
	if flag {
		fmt.Println("此处为读缓存...")
		return &proto.GetItemResponse{
			ActivityItem: &itemInfo,
		}, nil

	}

	//查数据库
	la, err := repo.GetItem(context.Background(), ItemInfo)

	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("internal error : %v", err))

	}

	newcash := model.ActivityItems{
		ItemID:           r.GetItemID(),
		ItemName:         la.ItemName,
		ItemDesc:         la.ItemDesc,
		ItemPrice:        la.ItemPrice,
		ItemSpecialPrice: la.ItemSpecialPrice,
		PromotionID:      la.PromotionID,
	}

	reds.SetItemCash(newcash.ItemID, newcash)
	//fmt.Println("redis插入失败..")

	fmt.Println("此处为读数据库...")
	return &proto.GetItemResponse{ActivityItem: &proto.ActivityItem{
		ItemID:           r.GetItemID(),
		ItemName:         la.ItemName,
		ItemDesc:         la.ItemDesc,
		ItemPrice:        la.ItemPrice,
		ItemSpecialPrice: la.ItemSpecialPrice,
		PromotionID:      la.PromotionID,
	}}, nil
}

func (s *Server) GetPromotion(ctx context.Context, req *proto.GetPromotionRequest) (*proto.GetPromotionResponse, error) {
	r := req.GetUniqueItemID()

	ItemInfo := model.Items{
		ItemID: r.GetItemID(),
	}

	lp, err := repo.GetPromotions(context.Background(), ItemInfo)

	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("internal error : %v", err))

	}

	promotionList := &proto.GetPromotionResponse{}

	for _, v := range lp {
		promotionList.Promotion = append(promotionList.Promotion, v)
	}

	//fmt.Println(lp[0].ItemID)
	return promotionList, nil

}

func (s *Server) CreatePromotion(ctx context.Context, req *proto.CreatePromotionRequest) (*proto.CreatePromotionResponse, error) {
	promotion := req.GetPromotion()
	newPromotion := model.Promotions{
		PromotionDesc:    promotion.GetPromotionDesc(),
		ItemID:           promotion.GetItemID(),
		ItemSpecialPrice: promotion.GetItemSpecialPrice(),
		CreatorID:        promotion.GetCreatorID(),
		LastModifierID:   promotion.GetLastModifierID(),
		StartTime:        promotion.GetStartTime(),
		EndTime:          promotion.GetEndTime(),
	}

	if newPromotion.EndTime < newPromotion.StartTime {
		return &proto.CreatePromotionResponse{
			Msg: "Invalid params",
		}, nil
	}

	_, err := repo.CreatePromotion(ctx, newPromotion)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("internal error : %v", err))
	}

	//删除缓存
	reds.InvalidItemCash(promotion.GetItemID())

	return &proto.CreatePromotionResponse{
		Msg: "success",
	}, nil
}

func (s *Server) RemovePromotion(ctx context.Context, req *proto.RemovePromotionRequest) (*proto.RemovePromotionResponse, error) {
	r := req.GetPromotion()

	RemovePromotionInfo := model.Promotions{
		PromotionID: r.GetPromotionID(),
		ItemID:      r.GetItemID(),
	}

	res, err := repo.RemovePromotion(ctx, RemovePromotionInfo)
	fmt.Println(res)

	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("internal error : %v", err))

	}
	//删除缓存
	reds.InvalidItemCash(r.GetItemID())

	return &proto.RemovePromotionResponse{}, nil
}
