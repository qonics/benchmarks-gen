package main

import (
	"fmt"
	"os"

	"cache-manager/config"
	"cache-manager/routes"
)

func main() {
	fmt.Println("Hello - GO-GIN")
	config.InitializeConfig()
	config.ConnectDb()
	defer config.SESSION.Close()
	defer config.DB.Close()
	// controller.SocketConnection()

	server := routes.InitRoutes()
	// ctx := context.Background()
	server.Run(":" + os.Getenv("server_port"))
}
