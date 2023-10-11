package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Create Order
// @Description Create New Order
// @Security BearerAuth
// @Tags Transaction
// @Param transaction body entity.CreateTransactionParam true "transaction info"
// @Produce json
// @Success 200 {object} entity.Response{data=int}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/transaction/create [POST]
func (r *rest) CreateOrder(ctx *gin.Context) {
	var inputParam entity.CreateTransactionParam
	if err := ctx.ShouldBindJSON(&inputParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	id, err := r.uc.Transaction.Create(ctx.Request.Context(), inputParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "successfully created new order", gin.H{"id": id})
}

// @Summary Get Order
// @Description Get Detail Order
// @Security BearerAuth
// @Tags Transaction
// @Param transaction_id path integer true "transaction id"
// @Produce json
// @Success 200 {object} entity.Response{data=entity.TransactionDetailResponse}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/transaction/{transaction_id} [GET]
func (r *rest) GetOrderDetail(ctx *gin.Context) {
	var param entity.TransactionParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := r.uc.Transaction.GetOrderDetail(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get order detail", result)
}

// @Summary Get Transaction List by UMKM
// @Description Get Transaction List by UMKM ID
// @Security BearerAuth
// @Tags Transaction
// @Param umkm_id path integer true "umkm id"
// @Param status query string false "status" Enums(in_cart, unpaid, paid, done)
// @Param order_id query string false "order_id"
// @Produce json
// @Success 200 {object} entity.Response{data=entity.TransactionDetailResponse}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/umkm/{umkm_id}/transactions [GET]
func (r *rest) GetTransactionListUmkm(ctx *gin.Context) {
	var param entity.TransactionParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := ctx.ShouldBindQuery(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := r.uc.Transaction.GetTransactionListByUmkm(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get transactions list", result)
}

// @Summary Get Transaction List
// @Description Get Transaction List
// @Security BearerAuth
// @Tags Transaction
// @Param order_id query string false "order_id"
// @Param limit query int true "limit"
// @Param page query int true "page"
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.TransactionDetailResponse}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/admin/transactions [GET]
func (r *rest) GetTransactionList(ctx *gin.Context) {
	var param entity.TransactionParam
	if err := ctx.ShouldBindQuery(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := r.uc.Transaction.GetTransactionList(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get transactions list", result)
}

// @Summary Get Recap Transaction
// @Description Get Recap Transaction List
// @Security BearerAuth
// @Tags Transaction
// @Param date query string false "date"
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.SalesRecapResponse{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/admin/transactions/recap [GET]
func (r *rest) GetRecapSalesList(ctx *gin.Context) {
	var param entity.TransactionParam
	if err := ctx.ShouldBindQuery(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := r.uc.Transaction.GetRecapSalesList(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get recap transaction list", result)
}

// @Summary Get My Transaction
// @Description Get My Transactions
// @Security BearerAuth
// @Tags Transaction
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.TransactionDetailResponse}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/transaction/me [GET]
func (r *rest) GetMyTransaction(ctx *gin.Context) {
	var param entity.TransactionParam

	if err := ctx.ShouldBindQuery(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := r.uc.Transaction.GetMyTransaction(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get transactions list", result)
}

// @Summary Complete Orders
// @Description Mark order as done
// @Security BearerAuth
// @Tags Transaction
// @Param umkm_id path integer true "umkm id"
// @Param transaction_id path integer true "transaction id"
// @Produce json
// @Success 200 {object} entity.Response{}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/umkm/{umkm_id}/transaction/{transaction_id}/mark-as-done [PUT]
func (r *rest) CompleteOrder(ctx *gin.Context) {
	var param entity.TransactionParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	err := r.uc.Transaction.CompleteOrder(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully mark as done", nil)
}
