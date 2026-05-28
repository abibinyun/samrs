package handlers

import (
	assetusecase "samrs-backend/internal/usecase/asset"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
)

type PublicAssetHandler struct {
	usecase assetusecase.PublicAssetUsecase
}

func NewPublicAssetHandler(u assetusecase.PublicAssetUsecase) *PublicAssetHandler {
	return &PublicAssetHandler{usecase: u}
}

func (h *PublicAssetHandler) GetByTenantAndCode(c *gin.Context) {
	tenantSlug := c.Param("slug")
	assetCode := c.Param("code")

	asset, err := h.usecase.GetPublicAsset(tenantSlug, assetCode)
	if err != nil {
		util.ErrorResponseFromErr(c, "Aset tidak ditemukan", err)
		return
	}

	util.SuccessResponse(c, "Aset ditemukan", asset)
}
