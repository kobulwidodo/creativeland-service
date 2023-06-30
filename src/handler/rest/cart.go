package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Add menu to cart
// @Description Add menu to cart
// @Security BearerAuth
// @Tags Cart
// @Param cart body entity.CreateCartParam true "cart info"
// @Produce json
// @Success 200 {object} entity.Response{data=entity.Cart{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/cart/create [POST]
func (r *rest) AddMenuToCart(ctx *gin.Context) {
	var param entity.CreateCartParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	cart, err := r.uc.Cart.Create(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully add item to cart", cart)
}

// @Summary Decrease Item
// @Description Decrease Item from cart
// @Security BearerAuth
// @Tags Cart
// @Param cart_id path integer true "cart id"
// @Produce json
// @Success 200 {object} entity.Response{}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/cart/{cart_id}/decrease [PUT]
func (r *rest) DecreaseItem(ctx *gin.Context) {
	var param entity.CartParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	err := r.uc.Cart.DecreaseItem(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully decrease item to cart", nil)
}

// @Summary Get List Cart
// @Description Get List Cart by User Logged in
// @Security BearerAuth
// @Tags Cart
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.Cart{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/cart [GET]
func (r *rest) GetListCartByUser(ctx *gin.Context) {
	carts, err := r.uc.Cart.GetListByUser(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get list items on cart", carts)
}

// @Summary Delete Cart
// @Description Get an item on cart
// @Security BearerAuth
// @Tags Cart
// @Produce json
// @Param cart_id path integer true "cart id"
// @Success 200 {object} entity.Response{}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/cart/{cart_id} [DELETE]
func (r *rest) DeleteItemCart(ctx *gin.Context) {
	var selectParam entity.CartParam
	if err := ctx.ShouldBindUri(&selectParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	err := r.uc.Cart.Delete(ctx.Request.Context(), selectParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully delete an item on cart", nil)
}
