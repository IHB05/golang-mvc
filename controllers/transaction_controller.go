package controllers

import (
	"belajar-crud-mvc/models"
	"belajar-crud-mvc/services"
	"belajar-crud-mvc/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	service services.TransactionService
}

func NewTransactionController(service services.TransactionService) *TransactionController {
	return &TransactionController{service: service}
}

// GET /api/v1/transactions
func (ctrl *TransactionController) GetAllTransactions(c *gin.Context) {
	page, _  := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Filter by user_id kalau ada
	var userID uint
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		id, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			utils.BadRequest(c, "Invalid user_id", "user_id must be a number")
			return
		}
		userID = uint(id)
	}

	transactions, total, totalPages, err := ctrl.service.GetAllTransactions(page, limit, userID)
	if err != nil {
		utils.InternalError(c, "Failed to fetch transactions", err.Error())
		return
	}

	utils.OK(c, "Transactions fetched successfully", gin.H{
		"transactions": transactions,
		"total":        total,
		"page":         page,
		"limit":        limit,
		"total_pages":  totalPages,
	})
}

// GET /api/v1/transactions/:id
func (ctrl *TransactionController) GetTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid transaction ID", "ID must be a number")
		return
	}

	transaction, err := ctrl.service.GetTransactionByID(uint(id))
	if err != nil {
		if err.Error() == "transaction not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, "Failed to fetch transaction", err.Error())
		return
	}

	utils.OK(c, "Transaction fetched successfully", transaction)
}

// POST /api/v1/transactions
func (ctrl *TransactionController) CreateTransaction(c *gin.Context) {
	var input models.CreateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	transaction, err := ctrl.service.CreateTransaction(input)
	if err != nil {
		errMsg := err.Error()

		// ✅ pakai strings.Contains — lebih aman dari slice string
		switch {
		case errMsg == "user not found":
			utils.NotFound(c, errMsg)
		case errMsg == "product not found":
			utils.NotFound(c, errMsg)
		case strings.Contains(errMsg, "insufficient stock"):
			utils.BadRequest(c, errMsg, "reduce quantity or choose another product")
		default:
			utils.InternalError(c, "Failed to create transaction", errMsg)
		}
		return
	}

	utils.Created(c, "Transaction created successfully", transaction)
}

// PATCH /api/v1/transactions/:id/status
func (ctrl *TransactionController) UpdateTransactionStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid transaction ID", "ID must be a number")
		return
	}

	var input models.UpdateTransactionStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	transaction, err := ctrl.service.UpdateTransactionStatus(uint(id), input)
	if err != nil {
		errMsg := err.Error()
		switch {
		case errMsg == "transaction not found":
			utils.NotFound(c, errMsg)
		case errMsg == "cancelled transaction cannot be updated":
			utils.BadRequest(c, errMsg, "transaction is already cancelled")
		default:
			utils.InternalError(c, "Failed to update transaction status", errMsg)
		}
		return
	}

	utils.OK(c, "Transaction status updated successfully", transaction)
}

// DELETE /api/v1/transactions/:id
func (ctrl *TransactionController) DeleteTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid transaction ID", "ID must be a number")
		return
	}

	if err := ctrl.service.DeleteTransaction(uint(id)); err != nil {
		if err.Error() == "transaction not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, "Failed to delete transaction", err.Error())
		return
	}

	utils.OK(c, "Transaction deleted successfully", nil)
}
