package service

import (
	"context"
	"fmt"

	"goods-service/internal/good/domain"
)

type GoodCache interface {
	SetGoodsList(ctx context.Context, goodsList domain.GoodsList) (err error)
	GetGoodsList(ctx context.Context) (goodsList domain.GoodsList, err error)
	DeleteGoodsList(ctx context.Context) (err error)
}

type GoodStorage interface {
	CreateGood(ctx context.Context, createGood domain.CreateGood) (good domain.Good, err error)
	UpdateGood(ctx context.Context, updateGood domain.UpdateGood) (good domain.Good, err error)
	DeleteGood(ctx context.Context, deleteGood domain.DeleteGood) (err error)
	ListGoods(ctx context.Context, listGoods domain.ListGoods) (goodsList domain.GoodsList, err error)
	ReprioritizeGood(ctx context.Context, reprioritizeGood domain.ReprioritizeGood) (
		goodsPriorities []domain.GoodPriority, err error)
}

type GoodsService struct {
	cache   GoodCache
	storage GoodStorage
}

func (s *GoodsService) Create(ctx context.Context, createGood domain.CreateGood) (good domain.Good, err error) {
	err = validateCreateGood(createGood)
	if err != nil {
		err = fmt.Errorf("validate create good: %w", err)
		return
	}
	good, err = s.storage.CreateGood(ctx, createGood)
	if err != nil {
		err = fmt.Errorf("create good: %w", err)
		return
	}
	err = s.cache.DeleteGoodsList(ctx)
	if err != nil {
		err = fmt.Errorf("delete goods list: %w", err)
		return
	}
	return
}

func (s *GoodsService) Update(ctx context.Context, updateGood domain.UpdateGood) (good domain.Good, err error) {
	err = validateUpdateGood(updateGood)
	if err != nil {
		err = fmt.Errorf("validate update good: %w", err)
		return
	}
	good, err = s.storage.UpdateGood(ctx, updateGood)
	if err != nil {
		err = fmt.Errorf("update good: %w", err)
		return
	}
	err = s.cache.DeleteGoodsList(ctx)
	if err != nil {
		err = fmt.Errorf("delete goods list: %w", err)
		return
	}
	return
}

func (s *GoodsService) Delete(ctx context.Context, deleteGood domain.DeleteGood) (err error) {
	err = validateDeleteGood(deleteGood)
	if err != nil {
		err = fmt.Errorf("validate delete good: %w", err)
		return
	}
	err = s.storage.DeleteGood(ctx, deleteGood)
	if err != nil {
		err = fmt.Errorf("delete good: %w", err)
		return
	}
	err = s.cache.DeleteGoodsList(ctx)
	if err != nil {
		err = fmt.Errorf("delete goods list: %w", err)
		return
	}
	return
}

func (s *GoodsService) List(ctx context.Context, listGoods domain.ListGoods) (goodsList domain.GoodsList, err error) {
	err = validateListGoods(listGoods)
	if err != nil {
		err = fmt.Errorf("validate list goods: %w", err)
		return
	}
	goodsList, err = s.storage.ListGoods(ctx, listGoods)
	if err != nil {
		err = fmt.Errorf("list goods: %w", err)
		return
	}
	err = s.cache.SetGoodsList(ctx, goodsList)
	if err != nil {
		err = fmt.Errorf("delete good: %w", err)
		return
	}
	return
}

func (s *GoodsService) Reprioritize(ctx context.Context, reprioritizeGood domain.ReprioritizeGood) (
	goodsPriorities []domain.GoodPriority, err error) {
	err = validateReprioritizeGood(reprioritizeGood)
	if err != nil {
		err = fmt.Errorf("validate reprioritize good: %w", err)
		return
	}
	goodsPriorities, err = s.storage.ReprioritizeGood(ctx, reprioritizeGood)
	if err != nil {
		err = fmt.Errorf("reprioritize good: %w", err)
		return
	}
	err = s.cache.DeleteGoodsList(ctx)
	if err != nil {
		err = fmt.Errorf("delete goods list: %w", err)
		return
	}
	return
}

func NewGoodService(cache GoodCache, storage GoodStorage) (service *GoodsService) {
	service = &GoodsService{
		cache:   cache,
		storage: storage,
	}
	return
}
