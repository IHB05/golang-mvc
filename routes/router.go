package routes

import (
	"belajar-crud-mvc/controllers"
	"belajar-crud-mvc/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

type RouterParams struct {
	dig.In // membuat parameter dari variabel tanpa harus memasukan secara manual

	ProductController     *controllers.ProductController // memasukan alamat ke ProductController biar dia sendiri bisa akses kedalamnya
	UserController        *controllers.UserController
	TransactionController *controllers.TransactionController
}

func NewRouter(p RouterParams) *gin.Engine { // harus NewRouter, tidak boleh newRouter karena kalo kecil cuma bisa local doang
	r := gin.Default() // memanggil gin.default (constructor gin.engine) yang dimasukan ke dalam r, gin.default() bertipe gin.engine

	r.Use(middleware.CORS())                // connect ke middleware, fungsi use untuk memasang middleware pada engine
	r.GET("/health", func(c *gin.Context) { // path : /health, handler func(c *gin.Context) dan isinya, bisa diisi yg laen selain func (method struct)
		c.JSON(http.StatusOK, gin.H{ // http.statusok = 200 (ngambil status), gin.H (isi map string gin) (tipe parameter ini adalah interface bisa diisi tipe data lain (slice or struct))
			"status":    "ok",
			"message":   "server is running",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	v1 := r.Group("/api/v1")
	{ // untuk menunjukan hirarki
		products := v1.Group("/products") // path = /api/v1/products
		{
			products.GET("", p.ProductController.GetAllProducts) // contoh handler dengan method
			products.GET("/:id", p.ProductController.GetProduct)
			products.POST("", p.ProductController.CreateProduct)
			products.PATCH("/:id", p.ProductController.UpdateProduct)
			products.DELETE("/:id", p.ProductController.DeleteProduct)
		}

		users := v1.Group("/users")
		{
			users.GET("", p.UserController.GetAllUsers)
			users.GET("/:id", p.UserController.GetUser)
			users.POST("", p.UserController.CreateUser)
			users.PATCH("/:id/", p.UserController.UpdateUser)
			users.DELETE("/:id", p.UserController.DeleteUser)
		}

		transactions := v1.Group("/transactions")
		{
			transactions.GET("", p.TransactionController.GetAllTransactions)
			transactions.GET("/:id", p.TransactionController.GetTransaction)
			transactions.POST("", p.TransactionController.CreateTransaction)
			transactions.PATCH("/:id/status", p.TransactionController.UpdateTransactionStatus)
			transactions.DELETE("/:id", p.TransactionController.DeleteTransaction)
		}
	}

	return r // bertipe gin.engine, ini yang dicari di main.go saat invoke yaitu siapa yang dapat membuat gin.engine

}

// type Controllers struct {
// 	Product     *controllers.ProductController
// 	User        *controllers.UserController
// 	Transaction *controllers.TransactionController
// }

// func SetupRouter(ctrl Controllers) *gin.Engine {
// 	r := gin.Default()
// 	r.Use(middleware.CORS())

// 	// Health check
// 	r.GET("/health", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"status":    "ok",
// 			"message":   "Server is running",
// 			"timestamp": time.Now().Format(time.RFC3339),
// 		})
// 	})

// 	v1 := r.Group("/api/v1")
// 	{
// 		// Product routes
// 		products := v1.Group("/products")
// 		{
// 			products.GET("", ctrl.Product.GetAllProducts)
// 			products.GET("/:id", ctrl.Product.GetProduct)
// 			products.POST("", ctrl.Product.CreateProduct)
// 			products.PATCH("/:id", ctrl.Product.UpdateProduct)
// 			products.DELETE("/:id", ctrl.Product.DeleteProduct)
// 		}

// 		// User routes
// 		users := v1.Group("/users")
// 		{
// 			users.GET("", ctrl.User.GetAllUsers)
// 			users.GET("/:id", ctrl.User.GetUser)
// 			users.POST("", ctrl.User.CreateUser)
// 			users.PATCH("/:id", ctrl.User.UpdateUser)
// 			users.DELETE("/:id", ctrl.User.DeleteUser)
// 		}

// 		// Transaction routes
// 		transactions := v1.Group("/transactions")
// 		{
// 			transactions.GET("", ctrl.Transaction.GetAllTransactions)
// 			transactions.GET("/:id", ctrl.Transaction.GetTransaction)
// 			transactions.POST("", ctrl.Transaction.CreateTransaction)
// 			transactions.PATCH("/:id/status", ctrl.Transaction.UpdateTransactionStatus)
// 			transactions.DELETE("/:id", ctrl.Transaction.DeleteTransaction)
// 		}
// 	}

// 	return r
// }
