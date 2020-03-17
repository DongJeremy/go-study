package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("hello world!")
	router := gin.Default()
	router.GET("/sleep", Sleep)
	router.Run(":1234")
}

func Sleep(c *gin.Context) {
	sec := c.DefaultQuery("sec", "1")
	secInt, _ := strconv.Atoi(sec)
	time.Sleep(time.Duration(secInt) * time.Second)
	c.String(200, "I sleep %s second.", sec)
}
