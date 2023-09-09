package main

import (
	"easylist/routes"
	util "easylist/utility"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	if config.GinMode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	routes.Init()
}
