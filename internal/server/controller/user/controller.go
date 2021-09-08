package user

import (
	"fmt"
	"github.com/Viverov/guideliner/internal/cradle"
	us "github.com/Viverov/guideliner/internal/domains/user/service"
	"github.com/Viverov/guideliner/internal/server/controller/utils"
	"github.com/Viverov/guideliner/internal/server/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	httpResponder utils.HttpResponder
}

func NewUserController(responder utils.HttpResponder) *Controller {
	return &Controller{
		httpResponder: responder,
	}
}

func (uc *Controller) Init(router *gin.Engine, cradle *cradle.Cradle, prefix string) {
	router.POST(prefix+"/users/registration", createRegisterHandler(cradle, uc.httpResponder))
	router.POST(prefix+"/users/login", createLoginHandler(cradle, uc.httpResponder))
	router.GET(
		prefix+"/users/me",
		middleware.CreateAuthMiddleware(cradle, uc.httpResponder),
		createMeHandler(cradle, uc.httpResponder),
	)
	router.POST(
		prefix+"/users/change_password",
		middleware.CreateAuthMiddleware(cradle, uc.httpResponder),
		createChangePasswordHandler(cradle, uc.httpResponder),
	)

}

type registerBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerResponse struct {
	UserID uint `json:"user_id"`
}

func createRegisterHandler(cradle *cradle.Cradle, er utils.HttpResponder) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		body := &registerBody{}
		if err := ctx.ShouldBindJSON(body); err != nil {
			er.Response(ctx, http.StatusBadRequest, "Validation error", err.Error(), "")
			return
		}

		dto, err := cradle.GetUserService().Register(body.Email, body.Password)
		if err != nil {
			switch err.(type) {
			case *us.EmailAlreadyExistError:
				er.Response(ctx, http.StatusBadRequest, "Email already taken", "", "")
				return
			default:
				er.InternalError(ctx, err.Error())
				return
			}
		}

		ctx.JSON(http.StatusOK, registerResponse{UserID: dto.ID()})
	}
}

type loginBody struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func createLoginHandler(cradle *cradle.Cradle, er utils.HttpResponder) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		body := &loginBody{}
		if err := ctx.ShouldBindJSON(body); err != nil {
			er.Response(ctx, http.StatusBadRequest, "Validation error", err.Error(), "")
			return
		}

		dto, err := cradle.GetUserService().FindByEmail(body.Email)
		if err != nil {
			er.InternalError(ctx, err.Error())
			return
		}
		if dto == nil {
			er.Response(ctx, http.StatusBadRequest, "Email not found", "", "")
			return
		}

		isValid, err := cradle.GetUserService().ValidateCredentials(body.Email, body.Password)
		if err != nil {
			er.InternalError(ctx, err.Error())
			return
		}
		if !isValid {
			er.Response(ctx, http.StatusBadRequest, "Invalid email/password", "", "")
			return
		}

		token, err := cradle.GetUserService().GetToken(dto.ID())
		if err != nil {
			er.InternalError(ctx, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, loginResponse{Token: token})
	}
}

type meResponse struct {
	ID    uint   `json:"id"  binding:"required"`
	Email string `json:"email"  binding:"required"`
}

func createMeHandler(cradle *cradle.Cradle, hr utils.HttpResponder) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		iUserId, _ := ctx.Get("user_id")
		userId, ok := iUserId.(uint)
		if !ok {
			hr.InternalError(ctx, fmt.Sprintf("Can't convert user_id into uint. userId value: %d", userId))
			return
		}

		dto, err := cradle.GetUserService().FindById(userId)
		if err != nil || dto == nil {
			hr.InternalError(ctx, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, meResponse{
			ID:    dto.ID(),
			Email: dto.Email(),
		})
	}
}

type changePasswordBody struct {
	OldPassword string `json:"old_password"  binding:"required"`
	NewPassword string `json:"new_password"  binding:"required"`
}

type changePasswordResponse struct{}

func createChangePasswordHandler(cradle *cradle.Cradle, hr utils.HttpResponder) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		iUserId, _ := ctx.Get("user_id")
		userId := iUserId.(uint)

		iEmail, _ := ctx.Get("user_email")
		email := iEmail.(string)

		body := &changePasswordBody{}
		if err := ctx.ShouldBindJSON(body); err != nil {
			hr.Response(ctx, http.StatusBadRequest, "Validation error", err.Error(), "")
			return
		}

		isValid, err := cradle.GetUserService().ValidateCredentials(email, body.OldPassword)
		if err != nil {
			hr.InternalError(ctx, err.Error())
			return
		}
		if !isValid {
			hr.Response(ctx, http.StatusBadRequest, "Invalid old password", "", "")
			return
		}

		err = cradle.GetUserService().ChangePassword(userId, body.NewPassword)
		if err != nil {
			hr.InternalError(ctx, err.Error())
		}

		ctx.JSON(http.StatusOK, changePasswordResponse{})
	}
}
