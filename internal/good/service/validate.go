package service

import (
	"fmt"

	"goods-service/internal/good/domain"
)

func validateCreateGood(createGood domain.CreateGood) (err error) {
	if createGood.ProjectID < 0 {
		err = fmt.Errorf("%w: negative project id", domain.ErrBadRequest)
		return
	}
	if createGood.Name == "" {
		err = fmt.Errorf("%w: empty name", domain.ErrBadRequest)
		return
	}
	return
}

func validateUpdateGood(updateGood domain.UpdateGood) (err error) {
	if updateGood.ID < 0 {
		err = fmt.Errorf("%w: negative id", domain.ErrBadRequest)
		return
	}
	if updateGood.ProjectID < 0 {
		err = fmt.Errorf("%w: negative project id", domain.ErrBadRequest)
		return
	}
	if updateGood.Name == "" {
		err = fmt.Errorf("%w: empty name", domain.ErrBadRequest)
		return
	}
	if updateGood.Description == "" {
		err = fmt.Errorf("%w: empty description", domain.ErrBadRequest)
		return
	}
	return
}

func validateDeleteGood(deleteGood domain.DeleteGood) (err error) {
	if deleteGood.ID < 0 {
		err = fmt.Errorf("%w: negative id", domain.ErrBadRequest)
		return
	}
	if deleteGood.ProjectID < 0 {
		err = fmt.Errorf("%w: negative project id", domain.ErrBadRequest)
		return
	}
	return
}

func validateListGoods(listGoods domain.ListGoods) (err error) {
	if listGoods.Limit < 0 {
		err = fmt.Errorf("%w: negative limit", domain.ErrBadRequest)
		return
	}
	if listGoods.Offset < 0 {
		err = fmt.Errorf("%w: negative offset", domain.ErrBadRequest)
		return
	}
	return
}

func validateReprioritizeGood(reprioritizeGood domain.ReprioritizeGood) (err error) {
	if reprioritizeGood.ID < 0 {
		err = fmt.Errorf("%w: negative id", domain.ErrBadRequest)
		return
	}
	if reprioritizeGood.ProjectID < 0 {
		err = fmt.Errorf("%w: negative project id", domain.ErrBadRequest)
		return
	}
	if reprioritizeGood.NewPriority < 1 {
		err = fmt.Errorf("%w: invalid new priority: must be greater than 1", domain.ErrBadRequest)
		return
	}
	return
}
