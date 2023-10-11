package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get Withdraw List
// @Description Get Withdraw List
// @Security BearerAuth
// @Tags Withdraw
// @Param date query string false "date"
// @Param limit query int true "limit"
// @Param page query int true "page"
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.TransactionDetailResponse}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/admin/withdraw [GET]
func (r *rest) GetWithdrawList(ctx *gin.Context) {
	var param entity.WithdrawParam
	if err := ctx.ShouldBindQuery(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := r.uc.Withdraw.GetList(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get withdraw list", result)
}

// @Summary Create Withdraw
// @Description Create new Withdraw
// @Security BearerAuth
// @Tags Withdraw
// @Param withdraw body entity.CreateWithdrawParam true "withdraw info"
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.TransactionDetailResponse}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/admin/withdraw [POST]
func (r *rest) CreateWithdraw(ctx *gin.Context) {
	var withdrawInput entity.CreateWithdrawParam
	if err := ctx.ShouldBindJSON(&withdrawInput); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	withdraw, err := r.uc.Withdraw.Create(ctx.Request.Context(), withdrawInput)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "successfully created new withdraw", withdraw)
}
