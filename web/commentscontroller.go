package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"net/http"
	"strings"
	"time"
)

type JsonTime time.Time

//实现它的json序列化方法
func (p JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(p).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

type Comments struct {
	Id       int64    `gorm:"column:id" json:"id"`
	PostId   int64    `gorm:"column:post_id" json:"post_id" binding:"required" `
	Author   string   `gorm:"column:author" json:"author" binding:"required"`
	Name     string   `gorm:"column:name" json:"name" binding:"required"`
	Content  string   `gorm:"column:content" json:"content" binding:"required"`
	Level    int64    `gorm:"column:level" json:"level"`
	Pid      int64    `gorm:"column:pid" json:"pid"`
	CreateAt JsonTime `gorm:"column:create_at" json:"create_at"`
}

func (c Comments) TableName() string {
	return "t_comments"
}

func QueryEscapeStr(s string) string {
	str := strings.Replace(s, "<", "%3c", -1)
	str = strings.Replace(str, ">", "%3e", -1)
	return str
}

func GetComments(engine *xorm.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		author := c.Query("author")
		postId := c.Query("postId")
		var list []Comments
		engine.Where("author = ? and post_id = ?", author, postId).Find(&list)
		c.JSON(200, gin.H{
			"content": list,
			"size":    len(list),
		})
	}
}

func PostComments(engine *xorm.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json Comments
		if err := c.ShouldBind(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    -1,
				"message": err.Error(),
			})
			return
		}

		json.CreateAt = JsonTime(time.Now())
		if json.Pid == 0 {
			json.Pid = -1
			json.Level = 1
		} else {
			parent := Comments{
				Id: json.Pid,
			}
			has, err := engine.Get(&parent)
			if !has {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    -2,
					"message": err.Error(),
				})
				return
			}
			json.Level = parent.Level + 1
			json.Pid = parent.Id
		}

		json.Name = QueryEscapeStr(json.Name)
		json.Content = QueryEscapeStr(json.Content)
		affected, err := engine.Insert(&json)
		if err != nil {
			c.JSON(200, gin.H{
				"code":    -1,
				"message": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"code":     0,
			"affected": affected,
		})
	}
}
