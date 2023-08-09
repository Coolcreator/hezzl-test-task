package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"goods-service/internal/good/domain"
)

const (
	goodIDParam    = "id"
	projectIDParam = "projectId"
	limitParam     = "limit"
	offsetParam    = "offset"
)

type GoodService interface {
	Create(ctx context.Context, createGood domain.CreateGood) (good domain.Good, err error)
	Update(ctx context.Context, updateGood domain.UpdateGood) (good domain.Good, err error)
	Delete(ctx context.Context, deleteGood domain.DeleteGood) (err error)
	List(ctx context.Context, listGoods domain.ListGoods) (goodsList domain.GoodsList, err error)
	Reprioritize(ctx context.Context, reprioritizeGood domain.ReprioritizeGood) (
		goodPriorities []domain.GoodPriority, err error)
}

type Controller struct {
	service GoodService
}

func (c *Controller) create(w http.ResponseWriter, r *http.Request) (err error) {
	projectIDStr := chi.URLParam(r, projectIDParam)
	if projectIDStr == "" {
		err = fmt.Errorf("%w: missing url param: projectId", domain.ErrBadRequest)
		return
	}
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			err = fmt.Errorf("%w: projectId has invalid syntax", domain.ErrBadRequest)
			return
		}
		err = fmt.Errorf("parse int: %w", err)
		return
	}
	req := new(createGoodRequest)
	err = render.DecodeJSON(r.Body, req)
	if err != nil {
		err = fmt.Errorf("decode json: %w", err)
		return
	}
	defer r.Body.Close()
	good, err := c.service.Create(r.Context(), domain.CreateGood{
		ProjectID: projectID,
		Name:      req.Name,
	})
	if err != nil {
		err = fmt.Errorf("service create: %w", err)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, goodResult{
		ID:          good.ID,
		ProjectID:   good.ProjectID,
		Name:        good.Name,
		Description: good.Description,
		Priority:    good.Priority,
		CreatedAt:   good.CreatedAt,
	})
	return
}

func (c *Controller) update(w http.ResponseWriter, r *http.Request) (err error) {
	var (
		goodID    int64
		projectID int64
	)
	goodID, projectID, err = getURLParams(r)
	if err != nil {
		err = fmt.Errorf("get url params: %w", err)
		return
	}
	req := new(updateGoodRequest)
	err = render.DecodeJSON(r.Body, req)
	if err != nil {
		err = fmt.Errorf("decode json: %w", err)
		return
	}
	defer r.Body.Close()
	good, err := c.service.Update(r.Context(), domain.UpdateGood{
		ID:          goodID,
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		err = fmt.Errorf("service update: %w", err)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, goodResult{
		ID:          good.ID,
		ProjectID:   good.ProjectID,
		Name:        good.Name,
		Description: good.Description,
		Priority:    good.Priority,
		CreatedAt:   good.CreatedAt,
	})
	return
}

func (c *Controller) remove(w http.ResponseWriter, r *http.Request) (err error) {
	var (
		goodID    int64
		projectID int64
	)
	goodID, projectID, err = getURLParams(r)
	if err != nil {
		err = fmt.Errorf("get url params: %w", err)
		return
	}
	err = c.service.Delete(r.Context(), domain.DeleteGood{
		ID:        goodID,
		ProjectID: projectID,
	})
	if err != nil {
		err = fmt.Errorf("delete good: %w", err)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, removeGoodResponse{
		ID:        goodID,
		ProjectID: projectID,
		Removed:   true,
	})
	return
}

func (c *Controller) list(w http.ResponseWriter, r *http.Request) (err error) {
	var (
		limit  int64
		offset int64
	)
	limitStr := chi.URLParam(r, limitParam)
	if limitStr == "" {
		err = fmt.Errorf("%w: missing url param: limit", domain.ErrBadRequest)
		return
	}
	limit, err = strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			err = fmt.Errorf("%w: limit has invalid syntax", domain.ErrBadRequest)
			return
		}
		err = fmt.Errorf("parse int: %w", err)
		return
	}
	offsetStr := chi.URLParam(r, offsetParam)
	if offsetStr == "" {
		err = fmt.Errorf("%w: missing url param: offset", domain.ErrBadRequest)
		return
	}
	offset, err = strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			err = fmt.Errorf("%w: offset has invalid syntax", err)
			return
		}
		err = fmt.Errorf("parse int: %w", err)
		return
	}
	goodsList, err := c.service.List(r.Context(), domain.ListGoods{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		err = fmt.Errorf("service list: %w", err)
		return
	}
	goodResults := make([]goodResult, 0, len(goodsList.Goods))
	for _, good := range goodsList.Goods {
		goodResults = append(goodResults, goodResult{
			ID:          good.ID,
			ProjectID:   good.ProjectID,
			Name:        good.Name,
			Description: good.Description,
			Priority:    good.Priority,
			CreatedAt:   good.CreatedAt,
		})
	}
	meta := meta{
		Total:   goodsList.Meta.Total,
		Removed: goodsList.Meta.Removed,
		Limit:   goodsList.Meta.Limit,
		Offset:  goodsList.Meta.Offset,
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, goodsListResult{
		Meta:  meta,
		Goods: goodResults,
	})
	return
}

func (c *Controller) Register(r chi.Router) {
	eh := errorHandler{}
	r.Post("/good/create", eh.wrap(c.create))
	r.Patch("/good/update", eh.wrap(c.update))
	r.Delete("/good/remove", eh.wrap(c.remove))
	r.Get("/good/list", eh.wrap(c.list))
	r.Patch("/good/reprioritize", eh.wrap(c.reprioritize))
}

func (c *Controller) reprioritize(w http.ResponseWriter, r *http.Request) (err error) {
	var (
		goodID    int64
		projectID int64
	)
	goodID, projectID, err = getURLParams(r)
	if err != nil {
		err = fmt.Errorf("get url params: %w", err)
		return
	}
	req := new(reprioritizeRequest)
	err = render.DecodeJSON(r.Body, req)
	if err != nil {
		err = fmt.Errorf("decode json: %w", err)
		return
	}
	defer r.Body.Close()
	goodPriorities, err := c.service.Reprioritize(r.Context(), domain.ReprioritizeGood{
		ID:          goodID,
		ProjectID:   projectID,
		NewPriority: req.NewPriority,
	})
	if err != nil {
		err = fmt.Errorf("reprioritize: %w", err)
		return
	}

	result := reprioritizeResult{}
	result.Priotities = make([]goodPriority, len(goodPriorities))
	for _, priority := range goodPriorities {
		result.Priotities = append(result.Priotities, goodPriority{
			ID:       priority.ID,
			Priority: priority.Priority,
		})
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, result)
	return
}

func NewController(service GoodService) (controller *Controller) {
	controller = &Controller{
		service: service,
	}
	return
}

func getURLParams(r *http.Request) (goodID int64, projectID int64, err error) {
	goodIDStr := chi.URLParam(r, goodIDParam)
	if goodIDStr == "" {
		err = fmt.Errorf("%w: missing url param: id", domain.ErrBadRequest)
		return
	}
	goodID, err = strconv.ParseInt(goodIDStr, 10, 64)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			err = fmt.Errorf("%w: id has invalid syntax", domain.ErrBadRequest)
			return
		}
		err = fmt.Errorf("parse int: %w", err)
		return
	}
	projectIDStr := chi.URLParam(r, projectIDParam)
	if projectIDStr == "" {
		err = fmt.Errorf("%w: missing url param: projectId", domain.ErrBadRequest)
		return
	}
	projectID, err = strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			err = fmt.Errorf("%w: projectId has invalid syntax", domain.ErrBadRequest)
			return
		}
		err = fmt.Errorf("parse int: %w", err)
		return
	}
	return
}
