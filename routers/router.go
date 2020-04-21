package routers

import (
	"FirstBeegoProject/controllers"
    "github.com/astaxie/beego/context"
    "github.com/astaxie/beego"
)

func init() {
    beego.InsertFilter("/Article/*",beego.BeforeRouter,FilterFunc)
    beego.Router("/", &controllers.LoginController{},"get:ShowLogin;post:HandleLogin")
    beego.Router("/regist",&controllers.RegistController{},"get:ShowRg;post:HandleRg")
    beego.Router("/Article/ShowArticle",&controllers.ArticleController{},"get:ShowArticleList;post:HandleSelect")
    beego.Router("/Article/AddArticle",&controllers.ArticleController{},"get:ShowAddArticleList;post:HandleAddArticle")
    beego.Router("/Article/ArticleContent",&controllers.ArticleController{},"get:ShowArticleContent")
    beego.Router("/Article/UpdateArticle",&controllers.ArticleController{},"get:ShowUpdateArticle;post:HandleUpdateArticle")
    beego.Router("/Article/DeleteArticle",&controllers.ArticleController{},"get:ShowDeleteArticle")
    beego.Router("/Article/AddArticleType",&controllers.ArticleController{},"get:ShowAddArticleType;post:HandleAddType")
    beego.Router("/Article/Logout",&controllers.LoginController{},"get:LogOut")
    beego.Router("/Article/deleteType",&controllers.ArticleController{},"get:DeleteType")
}
var FilterFunc = func(ctx *context.Context) {
    username := ctx.Input.Session("userName")
    if username==nil {
        ctx.Redirect(302,"/")
    }
}
