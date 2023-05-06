package rest

import (
	"errors"
	"fmt"
	"go-clean/src/business/entity"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (r *rest) httpRespSuccess(ctx *gin.Context, code int, message string, data interface{}) {
	resp := entity.Response{
		Meta: entity.Meta{
			Message: message,
			Code:    code,
			IsError: false,
		},
		Data: data,
	}
	ctx.JSON(code, resp)
}

func (r *rest) httpRespError(ctx *gin.Context, code int, err error) {
	resp := entity.Response{
		Meta: entity.Meta{
			Message: err.Error(),
			Code:    code,
			IsError: true,
		},
		Data: nil,
	}
	log.Default().Println(err)
	ctx.AbortWithStatusJSON(code, resp)
}

func (r *rest) VerifyUser(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		r.httpRespError(ctx, http.StatusUnauthorized, errors.New("empty token"))
		return
	}

	var tokenString string
	_, err := fmt.Sscanf(authHeader, "Bearer %v", &tokenString)
	if err != nil {
		r.httpRespError(ctx, http.StatusUnauthorized, errors.New("invalid token"))
		return
	}

	token, err := r.ValidateToken(tokenString)
	if err != nil {
		r.httpRespError(ctx, http.StatusUnauthorized, err)
		return
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		r.httpRespError(ctx, http.StatusUnauthorized, errors.New("failed to claim token"))
		return
	}

	c := ctx.Request.Context()
	user := entity.User{}

	if claim["is_guest"].(bool) {
		user = entity.User{
			GuestID: claim["guest_id"].(string),
		}
	} else {
		user, err = r.uc.User.GetById(uint(claim["id"].(float64)))
		if err != nil {
			r.httpRespError(ctx, http.StatusUnauthorized, errors.New("error while getting user"))
			return
		}
	}

	c = r.auth.SetUserAuthInfo(c, user.ConvertToAuthUser(), tokenString)
	ctx.Request = ctx.Request.WithContext(c)

	ctx.Next()
}

func (r *rest) VerifyAdmin(ctx *gin.Context) {
	user, err := r.auth.GetUserAuthInfo(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, http.StatusUnauthorized, err)
		return
	}

	if !user.User.IsAdmin {
		r.httpRespError(ctx, http.StatusUnauthorized, errors.New("dont have access"))
		return
	}

	ctx.Next()
}

func (r *rest) VerifyCart(ctx *gin.Context) {
	user, err := r.auth.GetUserAuthInfo(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, http.StatusUnauthorized, err)
		return
	}

	var selectParam entity.CartParam
	if err := ctx.ShouldBindUri(&selectParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := r.uc.Cart.ValidateCart(ctx, selectParam.ID, user.User.GuestID); err != nil {
		r.httpRespError(ctx, http.StatusUnauthorized, err)
		return
	}

	ctx.Next()
}

func (r *rest) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("token invalid")
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
