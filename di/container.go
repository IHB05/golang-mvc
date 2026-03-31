package di

import (
	"belajar-crud-mvc/config"
	"belajar-crud-mvc/controllers"
	"belajar-crud-mvc/repositories"
	"belajar-crud-mvc/services"
	"belajar-crud-mvc/routes"
	"log"

	"go.uber.org/dig"
	"gorm.io/gorm"
)

func BuildContainer() *dig.Container {
	c := dig.New() // -> bikin container baru

	mustProvide(c, provideDB) // memasukan akses database ke dalam container

	mustProvide(c, repositories.NewProductRepository) // memasukan contructor yang berupa metode, bukan fungsi agar container dapat info cara buatnya
	mustProvide(c, repositories.NewTransactionRepository)
	mustProvide(c, repositories.NewUserRepository)

	mustProvide(c, services.NewProductService)
	mustProvide(c, services.NewTransactionService)
	mustProvide(c, services.NewUserService)

	mustProvide(c, controllers.NewProductController)
	mustProvide(c, controllers.NewTransactionController)
	mustProvide(c, controllers.NewUserController)

	mustProvide(c, routes.NewRouter)

	return c
}

func mustProvide(c *dig.Container, constructor interface{}) { // pake interface biar bisa memuat banyak tipe data constructor
	if err := c.Provide(constructor); err != nil { // check kalau fitur Provide error atau tidak
		log.Fatalf("failed to provide dependency %v", err) // kalau error ada isinya maka send log.Fatalf (lgsg terminate)
	}
}

func provideDB() *gorm.DB {
	return config.DB
}
