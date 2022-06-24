package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strings"
	"ticketApp/src/type/entity"
	"ticketApp/src/type/util"
	"time"
)

type UserRepositoryType struct {
	UserCollection *mongo.Collection
}

func NewUserRepository(userCollection *mongo.Collection) UserRepositoryType {
	return UserRepositoryType{UserCollection: userCollection}
}

type UserRepository interface {
	UserRepoInsert(user entity.User) (*entity.UserPostResponseModel, *util.Error)
	UserRepoGetById(id string) (*entity.User, *util.Error)
	UserRepoDeleteById(id string) (util.DeleteResponseType, *util.Error)
	UserRepositoryGetAll(filter util.Filter) (*entity.UserGetResponseModel, *util.Error)
	UserRepositoryFindByUsernameAndPassword(model entity.LoginRequestModel) (*entity.User, *util.Error)
}

func (u UserRepositoryType) UserRepoInsert(user entity.User) (*entity.UserPostResponseModel, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", user}}
	_, err := u.UserCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, util.UpsertFailed.ModifyApplicationName("user repository").ModifyErrorCode(4015)
	}
	return &entity.UserPostResponseModel{Id: user.Id}, nil
}
func (u UserRepositoryType) UserRepoGetById(id string) (*entity.User, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var user entity.User
	filter := bson.M{"_id": id}
	if err := u.UserCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, util.NotFound.ModifyApplicationName("user repository").ModifyErrorCode(4032)
	}
	return &user, nil
}
func (u UserRepositoryType) UserRepoDeleteById(id string) (util.DeleteResponseType, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	result, err := u.UserCollection.DeleteOne(ctx, filter)
	if err != nil || result.DeletedCount <= 0 {
		return util.DeleteResponseType{IsSuccess: false}, util.DeleteFailed.ModifyApplicationName("user repository").ModifyErrorCode(4033)
	}
	return util.DeleteResponseType{IsSuccess: true}, nil
}
func (u UserRepositoryType) UserRepositoryGetAll(filter util.Filter) (*entity.UserGetResponseModel, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	totalCount, err := u.UserCollection.CountDocuments(ctx, filter.Filters)
	if err != nil {
		return nil, util.CountGet.ModifyApplicationName("user repository").ModifyDescription(err.Error()).ModifyErrorCode(3002)
	}
	opts := options.Find().SetSkip(filter.Page).SetLimit(filter.PageSize)
	if filter.SortingField != "" && filter.SortingDirection != 0 {
		opts.SetSort(bson.D{{filter.SortingField, filter.SortingDirection}})
	}

	cur, err := u.UserCollection.Find(ctx, filter.Filters, opts)
	if err != nil {
		return nil, util.UnKnownError.ModifyApplicationName("user repository").ModifyOperation("GET").ModifyDescription(err.Error()).ModifyErrorCode(4044)
	}
	var users []entity.User
	err = cur.All(ctx, &users)
	if err != nil {
		return nil, util.UnKnownError.ModifyApplicationName("user repository").ModifyOperation("GET").ModifyDescription(err.Error()).ModifyErrorCode(4045)
	}
	return &entity.UserGetResponseModel{
		RowCount: totalCount,
		Users:    users,
	}, nil
}
func (u UserRepositoryType) UserRepositoryFindByUsernameAndPassword(model entity.LoginRequestModel) (*entity.User, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var user entity.User
	filter := map[string]interface{}{}
	filter["$and"] = bson.A{
		bson.M{"username": model.Username},
		bson.M{"password": strings.ToLower(*model.Password)},
	}
	if err := u.UserCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, util.NewError("user repository", "LOGIN", "There is no user with provided information.", http.StatusNotFound, 4056)
	}
	return &user, nil
}
