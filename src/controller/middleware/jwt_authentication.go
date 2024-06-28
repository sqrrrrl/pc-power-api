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
		Realm:      Realm,
		Key:        []byte(os.Getenv("JWT_SECRET")),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour * 24 * 31,

		Authenticator:   a.authenticator(),
		Unauthorized:    a.unauthorized(),
		PayloadFunc:     a.payloadFunc(),
		IdentityHandler: a.identityHandler(),
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
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

func (a *AuthenticationMiddleware) unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		if message == jwt.ErrEmptyCookieToken.Error() {
			message = "The token is invalid"
		}
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}

func (a *AuthenticationMiddleware) payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*JwtUser); ok {
			return jwt.MapClaims{
				jwt.IdentityKey: v,
			}
		}
		return jwt.MapClaims{}
	}
}

func (a *AuthenticationMiddleware) identityHandler() func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		var identity map[string]interface{}
		identity = claims["identity"].(map[string]interface{})
		return &JwtUser{
			ID:       identity["ID"].(string),
			Username: identity["Username"].(string),
		}
	}
}

func GetUserIdFromContext(c *gin.Context) string {
	if jwtUser, ok := c.Get(jwt.IdentityKey); ok {
		return jwtUser.(*JwtUser).ID
	}
	return ""
}
