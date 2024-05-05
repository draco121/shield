package repository

import (
	"context"
	"errors"
	"time"

	"github.com/draco121/horizon/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IAuthenticationRepository interface {
	InsertOne(ctx context.Context, session *models.Session) (primitive.ObjectID, error)
	UpdateOne(ctx context.Context, session *models.Session) (*models.Session, error)
	FindOneById(ctx context.Context, id primitive.ObjectID) (*models.Session, error)
	DeleteOneById(ctx context.Context, id primitive.ObjectID) (*models.Session, error)
}

type authenticationRepository struct {
	IAuthenticationRepository
	db *mongo.Database
}

func NewAuthenticationRepository(database *mongo.Database) IAuthenticationRepository {
	return &authenticationRepository{
		db: database,
	}
}

func (ur *authenticationRepository) InsertOne(ctx context.Context, session *models.Session) (primitive.ObjectID, error) {

	result, err := ur.db.Collection("sessions").InsertOne(ctx, session)
	if err != nil {
		return primitive.NilObjectID, err
	} else {
		id := result.InsertedID.(primitive.ObjectID)
		return id, nil
	}
}

func (ur *authenticationRepository) UpdateOne(ctx context.Context, session *models.Session) (*models.Session, error) {
	filter := bson.M{"_id": session.ID}
	update := bson.M{"$set": bson.M{
		"updatedAt": time.Now(),
	}}
	result := models.Session{}
	err := ur.db.Collection("sessions").FindOneAndUpdate(ctx, filter, update).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	} else {
		return &result, nil
	}
}

func (ur *authenticationRepository) FindOneById(ctx context.Context, id primitive.ObjectID) (*models.Session, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	result := models.Session{}
	err := ur.db.Collection("sessions").FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	} else {
		return &result, nil
	}

}

func (ur *authenticationRepository) DeleteOneById(ctx context.Context, id primitive.ObjectID) (*models.Session, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	result := models.Session{}
	err := ur.db.Collection("users").FindOneAndDelete(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	} else {
		return &result, nil
	}

}
