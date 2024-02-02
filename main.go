package main

import (
	"fmt"
	"os"

	"benchmarks-gin/config"
	"benchmarks-gin/routes"
)

func main() {
	fmt.Println("Hello - GO-GIN")
	config.InitializeConfig()
	// controller.SocketConnection()

	server := routes.InitRoutes()
	// ctx := context.Background()
	// for i := 0; i < 10; i++ {
	// 	uuid := gocql.TimeUUID().String()
	// 	fmt.Println(uuid)
	// }
	server.Run(":" + os.Getenv("server_port"))
}
