package core

import (
	"context"
	"github.com/draco121/horizon/constants"
	"github.com/draco121/horizon/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"shield/repository"

	"github.com/draco121/horizon/utils"
)

type IUserService interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserById(ctx context.Context) (*models.User, error)
}

type userService struct {
	IUserService
	repo   repository.IUserRepository
	client *mongo.Client
}

func NewUserService(client *mongo.Client, repository repository.IUserRepository) IUserService {
	return &userService{
		repo:   repository,
		client: client,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	session, err := s.client.StartSession()
	if err != nil {
		utils.Logger.Error("failed to start mongo session", "error: ", err.Error())
		return nil, err
	}
	defer session.EndSession(ctx)
	err = session.StartTransaction()
	if err != nil {
		utils.Logger.Error("failed to start mongo transaction", "error: ", err.Error())
		return nil, err
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.Logger.Error("failed to hash password", "error: ", err.Error())
		return nil, err
	} else {
		user.Password = hashedPassword
		user.Role = constants.Tenant
		user, err = s.repo.InsertOne(ctx, user)
		if err != nil {
			utils.Logger.Error("failed to insert user", "error: ", err.Error())
			return nil, err
		} else {
			_ = session.CommitTransaction(ctx)
			utils.Logger.Info("inserted user")
			return user, nil
		}
	}
}

func (s *userService) GetUserById(ctx context.Context) (*models.User, error) {
	session, err := s.client.StartSession()
	if err != nil {
		utils.Logger.Error("failed to start mongo session", "error: ", err.Error())
		return nil, err
	}
	defer session.EndSession(ctx)
	err = session.StartTransaction()
	if err != nil {
		utils.Logger.Error("failed to start mongo transaction", "error: ", err.Error())
		return nil, err
	}
	userId := ctx.Value("UserId").(primitive.ObjectID)
	user, err := s.repo.FindOneById(ctx, userId)
	if err != nil {
		utils.Logger.Error("failed to find user", "error: ", err.Error())
		return nil, err
	} else {
		_ = session.CommitTransaction(ctx)
		utils.Logger.Info("fetched user")
		return user, nil
	}
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	session, err := s.client.StartSession()
	if err != nil {
		utils.Logger.Error("failed to start mongo session", "error: ", err.Error())
		return nil, err
	}
	defer session.EndSession(ctx)
	err = session.StartTransaction()
	if err != nil {
		utils.Logger.Error("failed to start mongo transaction", "error: ", err.Error())
		return nil, err
	}
	user, err := s.repo.FindOneByEmail(ctx, email)
	if err != nil {
		utils.Logger.Error("failed to find user", "error: ", err.Error())
		return nil, err
	} else {
		_ = session.CommitTransaction(ctx)
		utils.Logger.Info("fetched user")
		return user, nil
	}
}

func (s *userService) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	session, err := s.client.StartSession()
	if err != nil {
		utils.Logger.Error("failed to start mongo session", "error: ", err.Error())
		return nil, err
	}
	defer session.EndSession(ctx)
	err = session.StartTransaction()
	if err != nil {
		utils.Logger.Error("failed to start mongo transaction", "error: ", err.Error())
		return nil, err
	}
	newPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.Logger.Error("failed to hash password", "error: ", err.Error())
		return nil, err
	} else {
		user.Password = newPassword
		user, err := s.repo.UpdateOne(ctx, user)
		if err != nil {
			utils.Logger.Error("failed to update user", "error: ", err.Error())
			return nil, err
		} else {
			_ = session.CommitTransaction(ctx)
			utils.Logger.Info("updated user")
			return user, nil
		}
	}
}

func (s *userService) DeleteUser(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	session, err := s.client.StartSession()
	if err != nil {
		utils.Logger.Error("failed to start mongo session", "error: ", err.Error())
		return nil, err
	}
	defer session.EndSession(ctx)
	err = session.StartTransaction()
	if err != nil {
		utils.Logger.Error("failed to start mongo transaction", "error: ", err.Error())
		return nil, err
	}
	user, err := s.repo.DeleteOneById(ctx, id)
	if err != nil {
		utils.Logger.Error("failed to delete user", "error: ", err.Error())
		return nil, err
	} else {
		_ = session.CommitTransaction(ctx)
		utils.Logger.Info("deleted user")
		return user, nil
	}
}
