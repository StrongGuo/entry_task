package main

import (
	"context"
	"item-admin/proto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func createItem(c *gin.Context) {
	var item proto.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, cc := connectToServer()

	defer cc.Close()

	res, err := client.CreateItem(context.Background(),
		&proto.CreateItemRequest{Item: &item})

	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"result": res,
		"code":   http.StatusCreated,
	})
}

func updateItem(c *gin.Context) {
	var updateitem proto.UpdateItem
	if err := c.ShouldBindJSON(&updateitem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, cc := connectToServer()

	defer cc.Close()

	res, err := client.UpdateItem(context.Background(),
		&proto.UpdateItemRequest{UpdateItem: &updateitem})

	if err != nil {
		log.Fatal(err)
	}

	code := http.StatusOK
	c.JSON(code, gin.H{
		"result": res,
		"code":   code,
	})
}

func removeItem(c *gin.Context) {
	var removeitem proto.UniqueItemID
	if err := c.ShouldBindJSON(&removeitem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, cc := connectToServer()

	defer cc.Close()

	res, err := client.RemoveItem(context.Background(),
		&proto.RemoveItemRequest{UniqueItemID: &removeitem})

	if err != nil {
		log.Fatal(err)
	}

	code := http.StatusOK
	c.JSON(code, gin.H{
		"result": res,
		"code":   code,
	})
}

//func getItem(c *gin.Context) {
//	var searchitem proto.UniqueItemID
//	if err := c.ShouldBindJSON(&searchitem); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	client, cc := connectToServer()
//
//	defer cc.Close()
//
//	var activityitems []model.ActivityItems
//
//	stream, err := client.GetItem(context.Background(),
//		&proto.GetItemRequest{UniqueItemID: &searchitem})
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for {
//		var activityitem model.ActivityItems
//		res, err := stream.Recv()
//
//		if err == io.EOF {
//			break
//		}
//
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		activityitem.ItemID = res.GetActivityItem().GetItemID()
//		activityitem.ItemName = res.GetActivityItem().GetItemName()
//		activityitem.ItemDesc = res.GetActivityItem().GetItemDesc()
//		activityitem.ItemPrice = res.GetActivityItem().GetItemPrice()
//		activityitem.ItemSpecialPrice = res.GetActivityItem().GetItemSpecialPrice()
//		activityitem.PromotionID = res.GetActivityItem().GetPromotionID()
//
//		activityitems = append(activityitems, activityitem)
//	}
//	code := http.StatusOK
//	c.JSON(code, gin.H{
//		"result": activityitems,
//		"code":   code,
//	})
//}
//
//func getPromotion(c *gin.Context) {
//	var searchpromotion proto.UniqueItemID
//	if err := c.ShouldBindJSON(&searchpromotion); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	client, cc := connectToServer()
//
//	defer cc.Close()
//
//	var promotions []model.Promotions
//
//	stream, err := client.GetPromotion(context.Background(),
//		&proto.GetPromotionRequest{UniqueItemID: &searchpromotion})
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for {
//		var promotion model.Promotions
//		res, err := stream.Recv()
//
//		if err == io.EOF {
//			break
//		}
//
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		promotion.PromotionID = res.GetPromotion().GetPromotionID()
//		promotion.PromotionDesc = res.GetPromotion().GetPromotionDesc()
//		promotion.ItemSpecialPrice = res.GetPromotion().GetItemSpecialPrice()
//		promotion.CreatorID = res.GetPromotion().GetCreatorID()
//		promotion.LastModifierID = res.GetPromotion().GetLastModifierID()
//		promotion.StartTime = res.GetPromotion().GetStartTime()
//		promotion.EndTime = res.GetPromotion().GetEndTime()
//		promotion.ItemID = res.GetPromotion().GetItemID()
//
//		promotions = append(promotions, promotion)
//	}
//	code := http.StatusOK
//	c.JSON(code, gin.H{
//		"result": promotions,
//		"code":   code,
//	})
//}

func createPromotion(c *gin.Context) {
	var promotion proto.Promotion
	if err := c.ShouldBindJSON(&promotion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, cc := connectToServer()

	defer cc.Close()

	res, err := client.CreatePromotion(context.Background(),
		&proto.CreatePromotionRequest{Promotion: &promotion})

	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"result": res,
		"code":   http.StatusCreated,
	})
}

func removePromotion(c *gin.Context) {
	var removepromotion proto.Promotion
	if err := c.ShouldBindJSON(&removepromotion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, cc := connectToServer()

	defer cc.Close()

	res, err := client.RemovePromotion(context.Background(),
		&proto.RemovePromotionRequest{Promotion: &removepromotion})

	if err != nil {
		log.Fatal(err)
	}

	code := http.StatusOK
	c.JSON(code, gin.H{
		"result": res,
		"code":   code,
	})
}

func connectToServer() (proto.ItemServiceClient, *grpc.ClientConn) {
	cc, err := grpc.Dial("127.0.0.1:8001", grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	client := proto.NewItemServiceClient(cc)
	return client, cc
}

func main() {

	r := gin.Default()

	r.POST("item/create", createItem)
	r.PUT("item/update", updateItem)
	r.POST("item/remove", removeItem)
	//r.GET("get/item", getItem)
	//r.GET("get/item/promotion", getPromotion)
	r.POST("promotion/create", createPromotion)
	r.POST("promotion/remove", removePromotion)
	r.Run()
}
