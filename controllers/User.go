package controllers

import (
	"FirstBeegoProject/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type RegistController struct {
	beego.Controller
}

func (this*RegistController) ShowRg() {
	this.TplName="register.html"

}

func (this*RegistController) HandleRg()  {
	username := this.GetString("userName")
	passwd := this.GetString("password")
	//beego.Info(username,passwd)
	if username=="" || passwd==""{
		beego.Info("请输入完整数据")
		this.TplName="register.html"
		return
	}
	o:=orm.NewOrm()
	user:=models.User{}
	user.Username=username
	user.Passwd=passwd
	_, err := o.Insert(&user)
	if err != nil {
		beego.Info("注册失败")
		this.TplName="register.html"
		return
	}
	this.Redirect("/",302)

}

type LoginController struct {
	beego.Controller
}

func (this*LoginController)ShowLogin()  {
	userName := this.Ctx.GetCookie("userName")
	if userName =="" {
		this.Data["username"]=""
		this.Data["check"]=""
	}else {
		this.Data["username"]=userName
		this.Data["check"]="checked"
	}
	this.TplName="login.html"
}
func (this*LoginController)HandleLogin()  {
	username := this.GetString("userName")
	passwd := this.GetString("password")
	if username==""||passwd=="" {
		beego.Info("清输入数据")
		this.TplName="login.html"
		return
	}
	o:=orm.NewOrm()
	user:=models.User{}
	user.Username=username
	err := o.Read(&user, "Username")
	if err != nil {
		beego.Info("用户名错误")
		this.TplName="login.html"
		return
	}
	if user.Passwd !=passwd {
		beego.Info("密码错误")
		this.TplName="login.html"
		return
	}
	cookie := this.GetString("remember")
	if cookie == "on"{
		this.Ctx.SetCookie("userName",username,time.Second*3600)
	}else {
		this.Ctx.SetCookie("userName","sss",-1)
	}
	this.SetSession("userName",username)
	this.Redirect("/Article/ShowArticle",302)
}
func (this *LoginController	)LogOut()  {
	this.DelSession("userName")
	this.Redirect("/",302)
}

