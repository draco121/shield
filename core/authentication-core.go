package core

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/draco121/authenticationservice/repository"
	"github.com/draco121/common/jwt"
	"github.com/draco121/common/models"
	"github.com/draco121/common/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAuthenticationService interface {
	PasswordLogin(ctx context.Context, loginInput *models.LoginInput) (*models.LoginOutput, error)
	Authenticate(ctx context.Context, token string) (*models.JwtCustomClaims, error)
	RefreshLogin(ctx context.Context, refreshToken string) (*models.LoginOutput, error)
	Logout(ctx context.Context, token string) error
}

type authenticationService struct {
	IAuthenticationService
	authenticationRepository repository.IAuthenticationRepository
	userRepository           repository.IUserRepository
	client                   *mongo.Client
}

func NewAuthenticationService(client *mongo.Client, authenticationRepository repository.IAuthenticationRepository, userRepository repository.IUserRepository) IAuthenticationService {
	return &authenticationService{
		authenticationRepository: authenticationRepository,
		userRepository:           userRepository,
		client:                   client,
	}
}

func (s *authenticationService) PasswordLogin(ctx context.Context, loginInput *models.LoginInput) (*models.LoginOutput, error) {
	mongoSession, err := s.client.StartSession()
	if err != nil {
		utils.Logger.Error("failed to start mongo mongoSession", "error: ", err.Error())
		return nil, err
	}
	defer mongoSession.EndSession(ctx)
	err = mongoSession.StartTransaction()
	if err != nil {
		utils.Logger.Error("failed to start mongo transaction", "error: ", err.Error())
		return nil, err
	}
	user, err := s.userRepository.FindOneByEmail(ctx, loginInput.Email)
	if err != nil {
		utils.Logger.Error("failed to find user by email", "error: ", err.Error())
		_ = mongoSession.AbortTransaction(ctx)
		return nil, err
	} else {
		if utils.CheckPasswordHash(loginInput.Password, user.Password) {
			session := models.Session{
				UserId:    user.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				ID:        primitive.NewObjectID(),
			}
			id, err := s.authenticationRepository.InsertOne(ctx, &session)
			if err != nil {
				utils.Logger.Error("failed to insert mongoSession", "error: ", err.Error())
				return nil, err
			} else {
				claims := models.JwtCustomClaims{
					Email:     user.Email,
					UserId:    user.ID,
					Role:      user.Role,
					SessionId: id,
				}
				token, err := jwt.GenerateJWT(&claims)
				if err != nil {
					utils.Logger.Error("failed to generate JWT", "error: ", err.Error())
					return nil, err
				} else {
					refreshToken, err := jwt.GenerateRefreshToken(id)
					if err != nil {
						utils.Logger.Error("failed to generate refreshToken", "error: ", err.Error())
						return nil, err
					} else {
						_ = mongoSession.CommitTransaction(ctx)
						utils.Logger.Info("successfully authenticated")
						return &models.LoginOutput{
							Token:        token,
							RefreshToken: refreshToken,
						}, nil
					}
				}
			}
		} else {
			utils.Logger.Info("Invalid email or password")
			return nil, fmt.Errorf("invalid credentials")
		}
	}
}

func (s *authenticationService) Authenticate(ctx context.Context, token string) (*models.JwtCustomClaims, error) {
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
	claims, err := jwt.VerifyJwtToken(token)
	if err != nil {
		utils.Logger.Error("failed to verify token", "error: ", err.Error())
		return nil, err
	} else {
		_, err := s.authenticationRepository.FindOneById(ctx, claims.SessionId)
		if err != nil {
			utils.Logger.Error("failed to find user by id", "error: ", err.Error())
			return nil, err
		}
		_ = session.CommitTransaction(ctx)
		utils.Logger.Info("successfully authenticated")
		return &claims.JwtCustomClaims, nil
	}
}

func (s *authenticationService) RefreshLogin(ctx context.Context, refreshToken string) (*models.LoginOutput, error) {
	mongoSession, err := s.client.StartSession()
	if err != nil {
		utils.Logger.Error("failed to start mongo session", "error: ", err.Error())
		return nil, err
	}
	defer mongoSession.EndSession(ctx)
	err = mongoSession.StartTransaction()
	if err != nil {
		utils.Logger.Error("failed to start mongo transaction", "error: ", err.Error())
		return nil, err
	}
	claims, err := jwt.VerifyRefreshToken(refreshToken)
	if err != nil {
		if claims != nil {
			_, err = s.authenticationRepository.DeleteOneById(ctx, claims.SessionId)
			if err != nil {
				utils.Logger.Error("failed to delete user by id", "error: ", err.Error())
			} else {
				_ = mongoSession.CommitTransaction(ctx)
			}
		}
		return nil, err
	} else {
		session, err := s.authenticationRepository.FindOneById(ctx, claims.SessionId)
		if err != nil {
			utils.Logger.Error("failed to find user by id", "error: ", err.Error())
			return nil, err
		} else {
			session, err = s.authenticationRepository.UpdateOne(ctx, session)
			if err != nil {
				utils.Logger.Error("failed to update session", "error: ", err.Error())
				return nil, err
			} else {
				user, err := s.userRepository.FindOneById(ctx, session.UserId)
				if err != nil {
					utils.Logger.Error("failed to find user by id", "error: ", err.Error())
					return nil, err
				} else {
					claims := models.JwtCustomClaims{
						Email:     user.Email,
						Role:      user.Role,
						UserId:    user.ID,
						SessionId: session.ID,
					}
					newToken, err := jwt.GenerateJWT(&claims)
					if err != nil {
						utils.Logger.Error("failed to generate JWT token error: ", err.Error())
						return nil, err
					} else {
						_ = mongoSession.CommitTransaction(ctx)
						utils.Logger.Info("successfully refreshed tokens")
						return &models.LoginOutput{
							Token:        newToken,
							RefreshToken: refreshToken,
						}, nil
					}
				}
			}
		}
	}
}

func (s *authenticationService) Logout(ctx context.Context, token string) error {
	session, err := s.client.StartSession()
	if err != nil {
		utils.Logger.Error("failed to start mongo session", "error: ", err.Error())
		return err
	}
	defer session.EndSession(ctx)
	err = session.StartTransaction()
	if err != nil {
		utils.Logger.Error("failed to start mongo transaction", "error: ", err.Error())
		return err
	}
	claims, err := jwt.VerifyJwtToken(token)
	if claims != nil {
		_, err = s.authenticationRepository.DeleteOneById(ctx, claims.JwtCustomClaims.SessionId)
		if err != nil {
			utils.Logger.Error("failed to delete session error", err.Error())
			return err
		}
	}
	_ = session.CommitTransaction(ctx)
	utils.Logger.Info("logged out successfully")
	return nil

}
