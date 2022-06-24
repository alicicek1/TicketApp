package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ticketApp/src/type/entity"
	"ticketApp/src/type/util"
	"time"
)

type CategoryRepositoryType struct {
	CategoryCollection *mongo.Collection
}

func NewCategoryRepository(categoryCollection *mongo.Collection) CategoryRepositoryType {
	return CategoryRepositoryType{CategoryCollection: categoryCollection}
}

type CategoryRepository interface {
	CategoryRepoInsert(category entity.Category) (*entity.CategoryPostResponseModel, *util.Error)
	CategoryRepoGetById(id string) (*entity.Category, *util.Error)
	CategoryRepoDeleteById(id string) (util.DeleteResponseType, *util.Error)
	CategoryRepositoryGetAll(filter util.Filter) (*entity.CategoryGetResponseModel, *util.Error)
}

func (c CategoryRepositoryType) CategoryRepoInsert(category entity.Category) (*entity.CategoryPostResponseModel, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := c.CategoryCollection.InsertOne(ctx, category)
	if err != nil {
		return nil, util.UpsertFailed.ModifyApplicationName("user repository").ModifyErrorCode(4015)
	}
	return &entity.CategoryPostResponseModel{Id: category.Id}, nil
}

func (c CategoryRepositoryType) CategoryRepoGetById(id string) (*entity.Category, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var category entity.Category
	filter := bson.M{"_id": id}
	if err := c.CategoryCollection.FindOne(ctx, filter).Decode(&category); err != nil {
		return nil, util.NotFound.ModifyApplicationName("category repository").ModifyErrorCode(4028)
	}
	return &category, nil
}
func (c CategoryRepositoryType) CategoryRepoDeleteById(id string) (util.DeleteResponseType, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	result, err := c.CategoryCollection.DeleteOne(ctx, filter)
	if err != nil || result.DeletedCount <= 0 {
		return util.DeleteResponseType{IsSuccess: false}, util.DeleteFailed.ModifyApplicationName("category repository").ModifyErrorCode(4029)
	}
	return util.DeleteResponseType{IsSuccess: true}, nil
}
func (c CategoryRepositoryType) CategoryRepositoryGetAll(filter util.Filter) (*entity.CategoryGetResponseModel, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	totalCount, err := c.CategoryCollection.CountDocuments(ctx, filter.Filters)
	if err != nil {
		return nil, util.CountGet.ModifyApplicationName("category repository").ModifyDescription(err.Error()).ModifyErrorCode(3000)
	}
	opts := options.Find().SetSkip(filter.Page).SetLimit(filter.PageSize)
	if filter.SortingField != "" && filter.SortingDirection != 0 {
		opts.SetSort(bson.D{{filter.SortingField, filter.SortingDirection}})
	}

	cur, err := c.CategoryCollection.Find(ctx, filter.Filters, opts)
	if err != nil {

	}
	var categories []entity.Category
	err = cur.All(ctx, &categories)
	return &entity.CategoryGetResponseModel{
		RowCount:   totalCount,
		Categories: categories,
	}, nil
}
