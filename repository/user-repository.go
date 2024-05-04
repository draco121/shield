package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/draco121/common/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserRepository interface {
	InsertOne(ctx context.Context, user *models.User) (*models.User, error)
	UpdateOne(ctx context.Context, user *models.User) (*models.User, error)
	FindOneById(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	FindOneByEmail(ctx context.Context, email string) (*models.User, error)
	DeleteOneById(ctx context.Context, id primitive.ObjectID) (*models.User, error)
}

type userRepository struct {
	IUserRepository
	db *mongo.Database
}

func NewUserRepository(database *mongo.Database) IUserRepository {
	return &userRepository{
		db: database,
	}
}

func (ur *userRepository) InsertOne(ctx context.Context, user *models.User) (*models.User, error) {
	result, _ := ur.FindOneByEmail(ctx, user.Email)
	if result != nil {
		return nil, fmt.Errorf("record exists")
	} else {
		user.ID = primitive.NewObjectID()
		_, err := ur.db.Collection("users").InsertOne(ctx, user)
		if err != nil {
			return nil, err
		} else {
			return user, nil
		}
	}
}

func (ur *userRepository) UpdateOne(ctx context.Context, user *models.User) (*models.User, error) {
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{
		"password": user.Password,
	}}
	result := models.User{}
	err := ur.db.Collection("users").FindOneAndUpdate(ctx, filter, update).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	} else {
		return &result, nil
	}
}

func (ur *userRepository) FindOneById(ctx context.Context, id primitive.ObjectID) (*models.User, error) {

	filter := bson.D{{Key: "_id", Value: id}}
	result := models.User{}
	err := ur.db.Collection("users").FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	} else {
		return &result, nil
	}

}

func (ur *userRepository) FindOneByEmail(ctx context.Context, email string) (*models.User, error) {
	filter := bson.D{{Key: "email", Value: email}}
	result := models.User{}
	err := ur.db.Collection("users").FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	} else {
		return &result, nil
	}
}

func (ur *userRepository) DeleteOneById(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	result := models.User{}
	err := ur.db.Collection("users").FindOneAndDelete(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	} else {
		return &result, nil
	}

}
