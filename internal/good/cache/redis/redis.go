package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"goods-service/internal/good/domain"
)

const (
	goodsListKey = "goodsList"
	ttl          = time.Minute
)

type Cache struct {
	client *redis.Client
}

func (c *Cache) SetGoodsList(ctx context.Context, goodsList domain.GoodsList) (err error) {
	var jsonData []byte
	list := toRedis(goodsList)
	jsonData, err = json.Marshal(&list)
	if err != nil {
		err = fmt.Errorf("json marshal: %w", err)
		return
	}
	err = c.client.Set(ctx, goodsListKey, jsonData, ttl).Err()
	return
}

func (c *Cache) GetGoodsList(ctx context.Context) (goodsList domain.GoodsList, err error) {
	var jsonData []byte
	jsonData, err = c.client.Get(ctx, goodsListKey).Bytes()
	if err != nil {
		err = fmt.Errorf("get: %w", err)
		return
	}
	var list listOfGoods
	err = json.Unmarshal(jsonData, &list)
	if err != nil {
		err = fmt.Errorf("json unmarshal: %w", err)
		return
	}
	goodsList = fromRedis(list)
	return
}

func (c *Cache) DeleteGoodsList(ctx context.Context) (err error) {
	err = c.client.Del(ctx, goodsListKey).Err()
	if err != nil {
		err = fmt.Errorf("del: %w", err)
		return
	}
	return
}

func NewCache(client *redis.Client) (cache *Cache) {
	cache = &Cache{
		client: client,
	}
	return
}

func toRedis(goodsList domain.GoodsList) (list listOfGoods) {
	goods := make([]good, 0, len(goodsList.Goods))
	for _, item := range goodsList.Goods {
		goods = append(goods, good{
			ID:          item.ID,
			ProjectID:   item.ProjectID,
			Name:        item.Name,
			Description: item.Description,
			Priority:    item.Priority,
			CreatedAt:   item.CreatedAt,
		})
	}
	list = listOfGoods{
		Meta: meta{
			Total:   goodsList.Meta.Total,
			Removed: goodsList.Meta.Removed,
			Limit:   goodsList.Meta.Limit,
			Offset:  goodsList.Meta.Offset,
		},
		Goods: goods,
	}
	return
}

func fromRedis(list listOfGoods) (goodsList domain.GoodsList) {
	goods := make([]domain.Good, 0, len(goodsList.Goods))
	for _, item := range list.Goods {
		goods = append(goods, domain.Good{
			ID:          item.ID,
			ProjectID:   item.ProjectID,
			Name:        item.Name,
			Description: item.Description,
			Priority:    item.Priority,
			CreatedAt:   item.CreatedAt,
		})
	}
	goodsList = domain.GoodsList{
		Meta: domain.Meta{
			Total:   list.Meta.Total,
			Removed: list.Meta.Removed,
			Limit:   list.Meta.Limit,
			Offset:  list.Meta.Offset,
		},
		Goods: goods,
	}
	return
}
