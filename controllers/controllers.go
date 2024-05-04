package controllers

import (
	"github.com/draco121/common/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"

	"github.com/draco121/authenticationservice/core"
	"github.com/draco121/common/models"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	authenticationService core.IAuthenticationService
	userService           core.IUserService
}

func NewControllers(authenticationService core.IAuthenticationService, userService core.IUserService) Controllers {
	c := Controllers{
		authenticationService: authenticationService,
		userService:           userService,
	}
	return c
}

func (s *Controllers) Login(c *gin.Context) {
	var loginInput models.LoginInput
	if err := c.ShouldBind(&loginInput); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
	} else {
		res, err := s.authenticationService.PasswordLogin(c, &loginInput)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, res)
		}
	}
}

func (s *Controllers) RefreshLogin(c *gin.Context) {
	refreshToken := c.GetHeader("refreshToken")
	if refreshToken != "" {
		result, err := s.authenticationService.RefreshLogin(c, refreshToken)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

func (s *Controllers) Logout(c *gin.Context) {
	token := c.GetHeader("Authentication")
	err := s.authenticationService.Logout(c, token)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		c.Status(http.StatusNoContent)
	}
}

func (s *Controllers) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
	} else {
		user.Role = constants.Tenant
		res, err := s.userService.CreateUser(c, &user)
		if err != nil {
			c.JSON(409, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(201, res)
		}
	}
}

func (s *Controllers) GetUserProfile(c *gin.Context) {
	result, err := s.userService.GetUserById(c)
	if err != nil {
		c.JSON(404, gin.H{
			"message": err.Error(),
		})
	} else {
		c.JSON(200, result)
	}
}

func (s *Controllers) UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
	} else {
		res, err := s.userService.UpdateUser(c, &user)
		if err != nil {
			c.JSON(404, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(201, gin.H{
				"result": res,
			})
		}
	}
}

func (s *Controllers) DeleteUser(c *gin.Context) {
	id := c.Query("id")
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, err)
	}
	if id != "" {
		_, err := s.userService.DeleteUser(c, userId)
		if err != nil {
			c.JSON(404, gin.H{
				"message": err.Error(),
			})
		} else {
			c.Status(204)
		}
	} else {
		c.JSON(400, gin.H{
			"message": "user id not provided",
		})
	}

}
