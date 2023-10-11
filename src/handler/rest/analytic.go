package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get Dashboard Widget
// @Description Get Dashboard Widget by UMKM ID
// @Security BearerAuth
// @Tags Analytic
// @Param umkm_id path integer true "umkm id"
// @Produce json
// @Success 200 {object} entity.Response{data=entity.WidgetDashboardResult}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/umkm/{umkm_id}/analytic/dashboard-widget [GET]
func (r *rest) GetDashboardWidget(ctx *gin.Context) {
	var param entity.AnalyticParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := r.uc.Analytic.GetDashboardWidget(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get dashboard widget", result)
}

// @Summary Get All Dashboard Widget
// @Description Get All Dashboard Widget
// @Security BearerAuth
// @Tags Analytic
// @Produce json
// @Success 200 {object} entity.Response{data=entity.WidgetDashboardResult}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/admin/analytic/dashboard-widget [GET]
func (r *rest) GetAllDashboardWidget(ctx *gin.Context) {
	result, err := r.uc.Analytic.GetAllDashboardWidget(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get all dashboard widget", result)
}
