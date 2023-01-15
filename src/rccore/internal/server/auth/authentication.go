package auth

import (
	"rccore/internal/common"

	"github.com/authorizerdev/authorizer-go"
	"github.com/gin-gonic/gin"
)

func AuthorizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defaultHeaders := map[string]string{}
		authorizerClient, err := authorizer.NewAuthorizerClient(common.BuildSecret.AuthClientID, common.BuildSecret.AuthURL, "", defaultHeaders)
		if err != nil {
			// unauthorized
			c.AbortWithStatusJSON(401, "unauthorized - unable to create authorizer client")
			return
		}
		profile, err := authorizerClient.GetProfile(map[string]string{
			"Authorization": c.Request.Header.Get("Authorization"),
		})
		if err != nil {
			// unauthorized
			c.AbortWithStatusJSON(401, "unauthorized - unable to get profile")
			common.Logger.Error(err)
			return
		}
		common.Logger.Info(profile.Roles)
		c.Next()
	}
}
