package controllers

import (
	"belajar-crud-mvc/models"
	"belajar-crud-mvc/services"
	"belajar-crud-mvc/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

// GET /api/v1/users
func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	page, _  := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search   := c.Query("search")

	users, total, totalPages, err := ctrl.service.GetAllUsers(page, limit, search)
	if err != nil {
		utils.InternalError(c, "Failed to fetch users", err.Error())
		return
	}

	utils.OK(c, "Users fetched successfully", gin.H{
		"users":       users,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

// GET /api/v1/users/:id
func (ctrl *UserController) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", "ID must be a number")
		return
	}

	user, err := ctrl.service.GetUserByID(uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, "Failed to fetch user", err.Error())
		return
	}

	utils.OK(c, "User fetched successfully", user)
}

// POST /api/v1/users
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var input models.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	user, err := ctrl.service.CreateUser(input)
	if err != nil {
		if err.Error() == "email already registered" {
			utils.BadRequest(c, err.Error(), "use a different email")
			return
		}
		utils.InternalError(c, "Failed to create user", err.Error())
		return
	}

	utils.Created(c, "User created successfully", user)
}

// PATCH /api/v1/users/:id
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", "ID must be a number")
		return
	}

	var input models.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	user, err := ctrl.service.UpdateUser(uint(id), input)
	if err != nil {
		if err.Error() == "user not found" {
			utils.NotFound(c, err.Error())
			return
		}
		if err.Error() == "no fields provided to update" {
			utils.BadRequest(c, err.Error(), "send at least one field")
			return
		}
		if err.Error() == "email already used by another user" {
			utils.BadRequest(c, err.Error(), "use a different email")
			return
		}
		utils.InternalError(c, "Failed to update user", err.Error())
		return
	}

	utils.OK(c, "User updated successfully", user)
}

// DELETE /api/v1/users/:id
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", "ID must be a number")
		return
	}

	if err := ctrl.service.DeleteUser(uint(id)); err != nil {
		if err.Error() == "user not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, "Failed to delete user", err.Error())
		return
	}

	utils.OK(c, "User deleted successfully", nil)
}
