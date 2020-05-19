package main

import (
	"github.com/gin-gonic/gin"
	"mdup/api"
)

/**
 * user: ZY
 * Date: 2020/4/7 17:26
 */



func main(){
  r:=gin.Default()
  r.Static(api.StaticViewPath,"."+api.StaticViewPath)
  r.POST("/upload",api.DoFileUpload)
  r.GET("/upload",api.FileUpload)
  r.Run(":8081")
}

