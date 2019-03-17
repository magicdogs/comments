package main

import (
	"comments/constant"
	"comments/utils"
	"comments/web"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	engine, err := xorm.NewEngine("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	engine.ShowSQL(true)
	defer engine.Close()
	r := gin.Default()
	r.Use(utils.CORSMiddleware())
	r.Use(gin.Logger())
	r.GET(constant.RestCommentsUrl, web.GetComments(engine))
	r.POST(constant.RestCommentsUrl, web.PostComments(engine))
	r.Static("/static", "./static")
	r.Run(":8888")
}
