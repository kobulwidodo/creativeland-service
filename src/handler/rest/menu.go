package rest

import (
	"go-clean/src/business/entity"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Summary Create Menu
// @Description Create New Menu
// @Security BearerAuth
// @Tags Menu
// @Param umkm_id path integer true "umkm id"
// @Param menu body entity.CreateMenuParam true "menu info"
// @Produce json
// @Success 200 {object} entity.Response{data=entity.Menu{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/umkm/{umkm_id}/menu/create [POST]
func (r *rest) CreateMenu(ctx *gin.Context) {
	var menuParam entity.MenuParam
	if err := ctx.ShouldBindUri(&menuParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	var menuInput entity.CreateMenuParam
	if err := ctx.ShouldBindJSON(&menuInput); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	menu, err := r.uc.Menu.Create(menuInput, menuParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "successfully created new menu", menu)
}

// @Summary Get Menu
// @Description Get Menu By ID
// @Security BearerAuth
// @Tags Menu
// @Produce json
// @Param menu_id path integer true "menu id"
// @Success 200 {object} entity.Response{data=entity.Menu{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/menu/{menu_id} [GET]
func (r *rest) GetMenuByID(ctx *gin.Context) {
	var menuParam entity.MenuParam
	if err := ctx.ShouldBindUri(&menuParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	menu, err := r.uc.Menu.GetById(menuParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get menu", menu)
}

// @Summary Get Menu List
// @Description Get a list of Menu
// @Security BearerAuth
// @Tags Menu
// @Produce json
// @Param umkm_id query integer false "umkm id"
// @Param name query string false "name"
// @Success 200 {object} entity.Response{data=[]entity.Menu{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/menu [GET]
func (r *rest) GetMenuList(ctx *gin.Context) {
	var menuParam entity.MenuParam
	if err := ctx.ShouldBindWith(&menuParam, binding.Query); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	menus, err := r.uc.Menu.GetAll(menuParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get menu list", menus)
}

// @Summary Update Menu
// @Description Update a Menu
// @Security BearerAuth
// @Tags Menu
// @Produce json
// @Param menu_id path integer true "menu id"
// @Param menu body entity.UpdateMenuParam true "menu info"
// @Success 200 {object} entity.Response{}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/menu/{menu_id} [PUT]
func (r *rest) UpdateMenu(ctx *gin.Context) {
	var updateParam entity.UpdateMenuParam
	if err := ctx.ShouldBindJSON(&updateParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	log.Printf("%s\n", reflect.TypeOf(updateParam.Price))

	var selectParam entity.MenuParam
	if err := ctx.ShouldBindUri(&selectParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	err := r.uc.Menu.Update(selectParam, updateParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully update menu", nil)
}

// @Summary Delete Menu
// @Description Delete a Menu
// @Security BearerAuth
// @Tags Menu
// @Produce json
// @Param menu_id path integer true "menu id"
// @Success 200 {object} entity.Response{}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/menu/{menu_id} [DELETE]
func (r *rest) DeleteMenu(ctx *gin.Context) {
	var selectParam entity.MenuParam
	if err := ctx.ShouldBindUri(&selectParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	err := r.uc.Menu.Delete(selectParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully delete menu", nil)
}
