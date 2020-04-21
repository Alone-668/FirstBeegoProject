package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)
import (
	_"github.com/go-sql-driver/mysql"
)

type User struct {
	Id int
	Username string
	Passwd string
	Articles[] *Article `orm:"rel(m2m)"`
}
type Article struct {
	Id int `orm:"pk;auto"`
	Atiname string `orm:"size(20)"`
	Acontent string `orm:"size(500)"`
	Aimg string	`orm:"size(50);null"`
	Atime time.Time	`orm:"type(datetime);auto_now_add"`
	Acount int	`orm:"default(0)"`
	ArticleType *ArticleType `orm:"rel(fk)"`
	Users[] *User`orm:"reverse(many)"`
}
type ArticleType struct {
	Id int
	TypeName string `orm:"size(20)"`
	Articles[] *Article`orm:"reverse(many)"`
}
func init()  {
	orm.RegisterDataBase("default","mysql","root:root@tcp(127.0.0.1:3306)/newsWeb?charset=utf8")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	orm.RunSyncdb("default",false,true)
}
