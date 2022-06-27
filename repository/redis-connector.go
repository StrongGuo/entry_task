package repository

import (
	"context"
	"encoding/json"
	"item-admin/model"
	"strconv"
	//"item-admin/model"

	"github.com/fatih/structs"
	"github.com/go-redis/redis"
)

type RedisClient struct {
	client *redis.Client
}

const (
	redisAddr = "localhost:6379"
)

var (
	ctx = context.Background()
)

//NewRedis  redis的初始化函数.
func NewRedis() *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Network:  "tcp",
			Addr:     redisAddr, //:6379
			Password: "",
			DB:       0,
		}),
	}
}

// SetItemCash 缓存商品信息.
func (rs *RedisClient) SetItemCash(itemID int64, i model.ActivityItems) error {
	Iteminfo := structs.Map(&i)
	err := rs.client.HMSet(ctx, strconv.FormatInt(itemID, 10), Iteminfo).Err()
	err2 := rs.client.HMSet(ctx, strconv.FormatInt(itemID, 10), "Valid", "1").Err()
	err3 := rs.client.Expire(ctx, strconv.FormatInt(itemID, 10), 300*1e9).Err() //5分钟过期
	if err != nil || err2 != nil || err3 != nil {
		return err
	}
	return nil
}

// GetItemCash 获取Item信息.
func (rs *RedisClient) GetItemCash(ItemID int64) (i *model.ActivityItems, hasData bool, err error) {
	result, err := rs.client.HGetAll(ctx, strconv.FormatInt(ItemID, 10)).Result()

	if err != nil {
		return i, false, err
	}
	if result["Valid"] != "" {
		hasData = true
	}
	var Iteminfo model.ActivityItems
	//byteresult := []byte(result)
	marshal, err := json.Marshal(result)
	err = json.Unmarshal(marshal, &Iteminfo)

	return &Iteminfo, hasData, nil
}

// InvalidItemCash 手动失效商品信息
func (rs *RedisClient) InvalidItemCash(ItemID int64) error {
	err := rs.client.HSet(ctx, strconv.FormatInt(ItemID, 10), "Valid", "").Err()
	if err != nil {
		return err
	}
	return nil
}
