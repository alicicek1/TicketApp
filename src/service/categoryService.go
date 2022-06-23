package service

import (
	"github.com/google/uuid"
	"ticketApp/src/repository"
	"ticketApp/src/type/entity"
	"ticketApp/src/type/util"
	"time"
)

type CategoryServiceType struct {
	CategoryRepository repository.CategoryRepository
}

type CategoryService interface {
	CategoryServiceInsert(user entity.Category) (*entity.CategoryPostResponseModel, *util.Error)
	CategoryServiceGetById(id string) (*entity.Category, *util.Error)
	CategoryServiceDeleteById(id string) (util.DeleteResponseType, *util.Error)
	CategoryServiceGetAll(filter util.Filter) (*entity.CategoryGetResponseModel, *util.Error)
}

func NewCategoryService(r repository.CategoryRepository) CategoryServiceType {
	return CategoryServiceType{CategoryRepository: r}
}

func (c CategoryServiceType) CategoryServiceInsert(category entity.Category) (*entity.CategoryPostResponseModel, *util.Error) {
	if category.Id == "" {
		isSuccess, err := util.CheckCategoryModel(category)
		if !isSuccess {
			return nil, err
		}
	}

	category.Id = uuid.New().String()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	result, err := c.CategoryRepository.CategoryRepoInsert(category)

	return result, err
}
func (c CategoryServiceType) CategoryServiceGetById(id string) (*entity.Category, *util.Error) {
	result, err := c.CategoryRepository.CategoryRepoGetById(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (c CategoryServiceType) CategoryServiceDeleteById(id string) (util.DeleteResponseType, *util.Error) {
	result, err := c.CategoryRepository.CategoryRepoDeleteById(id)
	if err != nil || result.IsSuccess == false {
		return util.DeleteResponseType{IsSuccess: false}, err
	}
	return util.DeleteResponseType{IsSuccess: true}, nil
}
func (c CategoryServiceType) CategoryServiceGetAll(filter util.Filter) (*entity.CategoryGetResponseModel, *util.Error) {
	result, err := c.CategoryRepository.CategoryRepositoryGetAll(filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
