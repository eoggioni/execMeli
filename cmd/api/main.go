package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routes(r)
	r.Run()
}
