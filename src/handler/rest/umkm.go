package rest

import (
	"fmt"
	"go-clean/src/business/entity"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Summary Create Umkm
// @Description Create New Umkm
// @Security BearerAuth
// @Tags Umkm
// @Param umkm body entity.CreateUmkmParam true "umkm info"
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.Umkm{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/umkm/create [POST]
func (r *rest) CreateUmkm(ctx *gin.Context) {
	var umkmInput entity.CreateUmkmParam
	if err := ctx.ShouldBindJSON(&umkmInput); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	umkm, err := r.uc.Umkm.Create(umkmInput)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "successfully created new umkm", umkm)
}

// @Summary Get Umkm
// @Description Get Umkm By ID
// @Security BearerAuth
// @Tags Umkm
// @Produce json
// @Param umkm_id path integer true "umkm id"
// @Success 200 {object} entity.Response{data=entity.Umkm{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/umkm/{umkm_id} [GET]
func (r *rest) GetUmkmByID(ctx *gin.Context) {
	var umkmParam entity.UmkmParam
	if err := ctx.ShouldBindUri(&umkmParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	umkm, err := r.uc.Umkm.Get(umkmParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get umkm", umkm)
}

// @Summary Get Umkm List
// @Description Get a list of UMKM
// @Security BearerAuth
// @Tags Umkm
// @Produce json
// @Param name query string false "name param"
// @Param status query string false "status umkm" Enums(open, close)
// @Success 200 {object} entity.Response{data=entity.Umkm{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/umkm [GET]
func (r *rest) GetUmkmList(ctx *gin.Context) {
	var umkmParam entity.UmkmParam
	if err := ctx.ShouldBindWith(&umkmParam, binding.Query); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	umkms, err := r.uc.Umkm.GetList(umkmParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get umkm list", umkms)
}

// @Summary Update Umkm
// @Description Update a UMKM
// @Security BearerAuth
// @Tags Umkm
// @Produce json
// @Param umkm_id path integer true "umkm id"
// @Param umkm body entity.UpdateUmkmParam true "umkm info"
// @Success 200 {object} entity.Response{}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/umkm/{umkm_id} [PUT]
func (r *rest) UpdateUmkm(ctx *gin.Context) {
	var updateParam entity.UpdateUmkmParam
	if err := ctx.ShouldBindJSON(&updateParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	var selectParam entity.UmkmParam
	if err := ctx.ShouldBindUri(&selectParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	err := r.uc.Umkm.Update(selectParam, updateParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully update umkm", nil)
}

// @Summary Delete Umkm
// @Description Delete a UMKM
// @Security BearerAuth
// @Tags Umkm
// @Produce json
// @Param umkm_id path integer true "umkm id"
// @Success 200 {object} entity.Response{}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/umkm/{umkm_id} [DELETE]
func (r *rest) DeleteUmkm(ctx *gin.Context) {
	var selectParam entity.UmkmParam
	if err := ctx.ShouldBindUri(&selectParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	err := r.uc.Umkm.Delete(selectParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully delete umkm", nil)
}

// @Summary Upload Umkm Image
// @Description Upload Umkm Image
// @Security BearerAuth
// @Tags Umkm
// @Accept multipart/form-data
// @Produce json
// @Param umkm_id path integer true "umkm id"
// @Param file formData file true "Upload file"
// @Success 200 {object} entity.Response{}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/umkm/{umkm_id}/upload-image [POST]
func (r *rest) UploadImageUmkm(ctx *gin.Context) {
	var selectParam entity.UmkmParam
	if err := ctx.ShouldBindUri(&selectParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	path := fmt.Sprintf("public/assets/umkm/%d-%s-%d", selectParam.ID, file.Filename, time.Now().Unix())

	if err := ctx.SaveUploadedFile(file, path); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	err = r.uc.Umkm.SaveImage(ctx.Request.Context(), selectParam, path)
	if err != nil {
		os.Remove(path)
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully update tenant's image", nil)
}
