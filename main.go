package main

import (
	"belajar-crud-mvc/config"
	"belajar-crud-mvc/di"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	// 2. Connect database
	config.ConnectDB()

	// 3. aktivasi container dan mengeluarkan isi container
	container := di.BuildContainer()
	err := container.Invoke(func(r *gin.Engine) {
		if err := r.Run(":8080"); err != nil {
			log.Fatal("cannot connect to server: ", err)
		}
	})

	if err != nil {
		log.Fatal("cannot initialize application: ", err)
	}
}
