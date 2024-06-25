package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/pc-power-api/src/api"
	"github.com/pc-power-api/src/infra/repo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

const IdentityKey = "id"
const Realm = "PcPowerApi"

type JwtUser struct {
	ID       string
	Username string
}

type AuthenticationMiddleware struct {
	userRepository *repo.UserRepository
}

func NewAuthenticationMiddleware(userRepository *repo.UserRepository) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		userRepository: userRepository,
	}
}

func (a *AuthenticationMiddleware) AuthMiddleware() (gin.HandlerFunc, *jwt.GinJWTMiddleware) {
	authMiddleware, err := jwt.New(a.initAuthSecurity())
	if err != nil {
		log.Fatal(err.Error())
	}
	return func(context *gin.Context) {
		errInit := authMiddleware.MiddlewareInit()
		if errInit != nil {
			log.Fatal(errInit.Error())
		}
	}, authMiddleware
}

func (a *AuthenticationMiddleware) initAuthSecurity() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:       Realm,
		Key:         []byte(os.Getenv("JWT_SECRET")),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour * 24 * 31,
		IdentityKey: IdentityKey,

		Authenticator: a.authenticator(),
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}

func (a *AuthenticationMiddleware) authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var credentials *api.Credentials
		if err := c.ShouldBind(&credentials); err != nil {
			return "", jwt.ErrMissingLoginValues
		}

		user, aerr := a.userRepository.GetByUsername(credentials.Username)

		if (aerr == nil) && (bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)) == nil) {
			return &JwtUser{
				ID:       user.ID,
				Username: user.Username,
			}, nil
		}
		return nil, jwt.ErrFailedAuthentication
	}
}
