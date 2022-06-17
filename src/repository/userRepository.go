package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ticketApp/src/type/entity"
	"time"
)

type UserRepositoryType struct {
	UserCollection *mongo.Collection
}

func NewUserRepository(userCollection *mongo.Collection) UserRepositoryType {
	return UserRepositoryType{UserCollection: userCollection}
}

type UserRepository interface {
	UserRepoInsert(user entity.User) (string, error)
	UserRepoGetById(id string) (*entity.User, error)
	UserRepoDeleteById(id string) (bool, error)
}

func (u UserRepositoryType) UserRepoInsert(user entity.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.Id = uuid.New().String()
	result, err := u.UserCollection.InsertOne(ctx, user)

	if result.InsertedID == nil || err != nil {
		errors.New("failed add")
		return "", err
	}
	return user.Id, nil
}

func (u UserRepositoryType) UserRepoGetById(id string) (*entity.User, error) {
	var user entity.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"baseentity.id": id}

	if err := u.UserCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u UserRepositoryType) UserRepoDeleteById(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := u.UserCollection.DeleteOne(ctx, bson.M{"baseentity.id": id})

	if err != nil || result.DeletedCount <= 0 {
		return false, err
	}
	return true, nil
}
