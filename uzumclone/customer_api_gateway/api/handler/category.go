package handler

import (
	ct "customer-api-gateway/genproto/catalog_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Router /category [post]
// @Summary Create Category
// @Description API for creating a category
// @Tags category
// @Accept       json
// @Produce      json
// @Param        category body catalog_service.CreateCategory true "Category"
// @Success      201 {object} models.ResponseSuccess
// @Failure      404 {object} models.ResponseError
// @Failure      500 {object} models.ResponseError
func (h *handler) CreateCategory(c *gin.Context) {
	category := &ct.CreateCategory{}
	if err := c.ShouldBindJSON(category); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while reading request body")
		return
	}

	resp, err := h.grpcClient.CategoryService().Create(c.Request.Context(), category)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while creating category")
		return
	}

	c.JSON(http.StatusCreated, resp)
}
