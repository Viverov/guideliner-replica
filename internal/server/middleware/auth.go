package middleware

import (
	"github.com/Viverov/guideliner/internal/cradle"
	tokens "github.com/Viverov/guideliner/internal/domains/user/token_provider"
	"github.com/Viverov/guideliner/internal/server/controller/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateAuthMiddleware create auth middleware. Middleware set "user_id" and "user_email" field in context on success
func CreateAuthMiddleware(cradle *cradle.Cradle, hr utils.HttpResponder) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("api-token")
		if token == "" {
			hr.Response(ctx, http.StatusUnauthorized, "api-token header must be set", "", "")
			ctx.Abort()
			return
		}

		dto, err := cradle.GetUserService().GetUserByToken(token)
		if err != nil {
			switch err.(type) {
			case *tokens.NotTokenError, *tokens.ExpiredTokenError:
				hr.Response(ctx, http.StatusUnauthorized, err.Error(), "", "")
				ctx.Abort()
				return
			default:
				hr.InternalError(ctx, err.Error())
				ctx.Abort()
				return
			}
		}

		ctx.Set("user_id", dto.ID())
		ctx.Set("user_email", dto.Email())

		ctx.Next()
	}
}
