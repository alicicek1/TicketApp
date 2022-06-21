package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
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
	UserRepoInsert(user *entity.User) (map[string]interface{}, *util.Error)
	UserRepoGetById(id string) (*entity.User, *util.Error)
	UserRepoDeleteById(id string) (bool, *util.Error)
}

func (u UserRepositoryType) UserRepoInsert(user *entity.User) (map[string]interface{}, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if user != nil && user.Id != nil {
		var userUpdate entity.User
		filter := bson.M{"_id": user.Id}
		if err := u.UserCollection.FindOne(ctx, filter).Decode(&userUpdate); err != nil {
			return nil, &util.Error{
				ApplicationName: "user repository",
				Operation:       "GET",
				Description:     "There is no user by provided id. Update failed.",
				StatusCode:      http.StatusNotFound,
				ErrorCode:       4002,
			}
		}

		if user.Type == 0 {
			user.Type = userUpdate.Type
		}
		if user.Username == "" {
			user.Username = userUpdate.Username
		}
		if user.Password == "" {
			user.Password = userUpdate.Password
		}
		if user.Email == "" {
			user.Email = userUpdate.Email
		}
		if user.CreatedAt.IsZero() {
			user.CreatedAt = userUpdate.CreatedAt
		} else {
			user.CreatedAt = userUpdate.CreatedAt
		}

		update := bson.M{"$set": bson.M{"" +
			"type": user.Type,
			"username":  user.Username,
			"password":  user.Password,
			"email":     user.Email,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
			"_id":       user.Id,
		}}

		res, err := u.UserCollection.UpdateOne(context.TODO(), filter, update)
		if res != nil && (res.ModifiedCount < 1 || err != nil) {
			return nil, &util.Error{
				ApplicationName: "user repository",
				Operation:       "PUT",
				Description:     "failed add",
				StatusCode:      http.StatusBadRequest,
				ErrorCode:       4000,
			}
		}
		return map[string]interface{}{
			"id": user.Id.Hex(),
		}, nil
	} else {
		res, err := u.UserCollection.InsertOne(ctx, user)

		if res != nil && (res.InsertedID == nil || err != nil) {
			return nil, &util.Error{
				ApplicationName: "user repository",
				Operation:       "POST",
				Description:     "failed add",
				StatusCode:      http.StatusBadRequest,
				ErrorCode:       4000,
			}
		}

		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			return map[string]interface{}{
				"id": oid.Hex(),
			}, nil
		}
	}

	return nil, &util.Error{
		ApplicationName: "user repository",
		Operation:       "POST",
		Description:     "model added. object id returning error",
		StatusCode:      http.StatusCreated,
		ErrorCode:       4001,
	}
}

func (u UserRepositoryType) UserRepoGetById(id string) (*entity.User, *util.Error) {
	var user entity.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &util.Error{
			ApplicationName: "user repository",
			Operation:       "GET",
			Description:     "failed parsing hex to object id",
			StatusCode:      http.StatusBadRequest,
			ErrorCode:       4003,
		}
	}
	filter := bson.M{"_id": objId}
	if err := u.UserCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, &util.Error{
			ApplicationName: "user repository",
			Operation:       "GET",
			Description:     "failed get",
			StatusCode:      http.StatusNotFound,
			ErrorCode:       4002,
		}
	}
	return &user, nil
}

func (u UserRepositoryType) UserRepoDeleteById(id string) (bool, *util.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, &util.Error{
			ApplicationName: "user repository",
			Operation:       "GET",
			Description:     "failed parsing hex to object id",
			StatusCode:      http.StatusBadRequest,
			ErrorCode:       4003,
		}
	}

	filter := bson.M{"_id": objId}
	result, err := u.UserCollection.DeleteOne(ctx, filter)

	if err != nil || result.DeletedCount <= 0 {
		return false, &util.Error{
			ApplicationName: "user repository",
			Operation:       "DELETE",
			Description:     "failed delete",
			StatusCode:      http.StatusBadRequest,
			ErrorCode:       4003,
		}
	}
	return true, nil
}
