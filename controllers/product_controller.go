package controllers

import (
	"belajar-crud-mvc/models"
	"belajar-crud-mvc/services"
	"belajar-crud-mvc/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service services.ProductService
}

func NewProductController(service services.ProductService) *ProductController {
	return &ProductController{service: service}
}

func (ctrl *ProductController) GetAllProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	category := c.Query("category")

	products, total, totalPages, err := ctrl.service.GetAllProducts(page, limit, search, category)
	if err != nil {
		utils.InternalError(c, "Failed to fetch products", err.Error())
		return
	}

	utils.OK(c, "Products fetched successfully", gin.H{
		"products":    products,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

func (ctrl *ProductController) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid product ID", "ID must be a number")
		return
	}

	product, err := ctrl.service.GetProductByID(uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, "Failed to fetch product", err.Error())
		return
	}

	utils.OK(c, "Product fetched successfully", product)
}

func (ctrl *ProductController) CreateProduct(c *gin.Context) {
	var input models.CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	product, err := ctrl.service.CreateProduct(input)
	if err != nil {
		utils.InternalError(c, "Failed to create product", err.Error())
		return
	}

	utils.Created(c, "Product created successfully", product)
}

func (ctrl *ProductController) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid product ID", "ID must be a number")
		return
	}

	var input models.UpdateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	product, err := ctrl.service.UpdateProduct(uint(id), input)
	if err != nil {
		if err.Error() == "product not found" {
			utils.NotFound(c, err.Error())
			return
		}
		if err.Error() == "no fields provided to update" {
			utils.BadRequest(c, err.Error(), "send at least one field")
			return
		}
		utils.InternalError(c, "Failed to update product", err.Error())
		return
	}

	utils.OK(c, "Product updated successfully", product)
}

func (ctrl *ProductController) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid product ID", "ID must be a number")
		return
	}

	if err := ctrl.service.DeleteProduct(uint(id)); err != nil {
		if err.Error() == "product not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, "Failed to delete product", err.Error())
		return
	}

	utils.OK(c, "Product deleted successfully", nil)
}
