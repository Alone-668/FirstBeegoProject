package controllers

import (
	"FirstBeegoProject/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"path"
	"strconv"
	"time"
)

type ArticleController struct {
	beego.Controller
}

func (this *ArticleController)ShowArticleList()  {

	typeName := this.GetString("select")
	o:=orm.NewOrm()
	var articles[] models.Article
	qs := o.QueryTable("Article")
	var count int64
	var err error
	if typeName==""{
		count, err = qs.RelatedSel("ArticleType").Count()
		if err != nil {
			beego.Info("查询错误")
			return
		}
	}else {
		count, err = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
		if err != nil {
			beego.Info("查询错误")
			return
		}
	}
	this.Data["typeName"] = typeName

	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}
	pageSize := 1
	start:=pageSize*(pageIndex-1)
	if typeName ==""{
		_, err = qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)
		if err != nil {
			beego.Info("查询错误")
			return
		}
	}else {
		_, err = qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
		if err != nil {
			beego.Info("查询错误")
			return
		}
	}

	pageCount := int(math.Ceil(float64(count) / float64(pageSize)))
	FirstPage:=false
	if pageIndex == 1 {
		FirstPage = true
	}
	FinalPage:=false
	if pageIndex ==pageCount {
		FinalPage =true
	}
	var articleType[] models.ArticleType
	o.QueryTable("ArticleType").All(&articleType)
	username := this.GetSession("userName")
	this.Data["username"]=username.(string)

	this.Data["articleType"] = articleType

	this.Data["FirstPage"]= FirstPage
	this.Data["FinalPage"]= FinalPage
	this.Data["pageIndex"]= pageIndex
	this.Data["count"]=count
	this.Data["pageCount"]=pageCount
	this.Data["articles"]=articles
	this.Layout = "layout.html"
	this.LayoutSections=make(map[string]string)
	this.LayoutSections["Script"] = "script.html"
	this.TplName="index.html"
}
func (this*ArticleController)HandleSelect()  {
	typeName := this.GetString("select")
	if typeName == "" {
		beego.Info("下拉框传递数据失败")
		return
	}
	var articles models.Article
	o:=orm.NewOrm()
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
	beego.Info(articles)
}
func (this *ArticleController)ShowAddArticleList()  {
	o:=orm.NewOrm()
	var articleType[] models.ArticleType
	o.QueryTable("ArticleType").All(&articleType)
	username := this.GetSession("userName")
	this.Data["username"]=username.(string)
	this.Data["articleType"] = articleType
	this.Layout = "layout.html"
	this.TplName="add.html"

}
func (this *ArticleController)HandleAddArticle()  {
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	articleType := this.GetString("select")
	file, header, err := this.GetFile("uploadname")
	defer file.Close()
	if err != nil {
		beego.Info("文件上传错误")
		return
	}

	ext := path.Ext(header.Filename)
	filename:= time.Now().Format("2006-01-02 15:04:05")

	err = this.SaveToFile("uploadname", "static/img/"+filename+ext)
	if err!=nil {
		beego.Info("保存图片失败")
		return
	}

	beego.Info(articleName+" "+content)
	beego.Info(filename+ext)

	//插入数据
	o:=orm.NewOrm()
	article:=models.Article{}
	article.Atiname = articleName
	article.Acontent = content
	var artiType models.ArticleType
	artiType.TypeName = articleType
	err = o.Read(&artiType, "TypeName")
	if err != nil {
		beego.Info("类型读取错误")
		return
	}
	article.ArticleType = &artiType
	article.Aimg = "./static/img/"+filename+ext


	_, err = o.Insert(&article)
	if err != nil {
		beego.Info("插入数据失败")
		return
	}
	this.Redirect("/Article/ShowArticle",302)


}
func (this*ArticleController)ShowArticleContent()  {
	//获取GET请求中的id
	id, err := this.GetInt("id")
	if err != nil {
		beego.Info("请求路径出错")
		this.Layout = "layout.html"
		this.TplName="index.html"
	}
	o:=orm.NewOrm()
	var article models.Article
	article.Id=id
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("Id",id).All(&article)
	beego.Info(article.ArticleType.TypeName)
	article.Acount+=1
	m2M := o.QueryM2M(&article, "Users")
	username := this.GetSession("userName")
	user:=models.User{Username:username.(string) }
	o.Read(&user,"Username")
	_, err = m2M.Add(&user)
	if err != nil {
		beego.Info("多对多插入失败")
		return
	}
	o.Update(&article)



	//多对多查询
	var users[] models.User
	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)


	this.Data["username"]=username.(string)
	this.Data["users"]=users
	this.Data["article"] = article
	this.Layout = "layout.html"
	this.TplName="content.html"

}
func (this*ArticleController)ShowUpdateArticle()  {
	id, err := this.GetInt("id")
	if err != nil {
		beego.Info("修改错误")
		this.Redirect("/Article/ShowArticle",302)
		return
	}
	o:=orm.NewOrm()
	var article models.Article
	article.Id = id
	err = o.Read(&article)
	if err != nil {
		beego.Info("数据查询错误")
		return
	}
	username := this.GetSession("userName")
	this.Data["username"]=username.(string)
	this.Data["article"] = article
	this.Layout = "layout.html"
	this.TplName="update.html"

}

func (this *ArticleController)ShowDeleteArticle()  {
	id, err := this.GetInt("id")
	if err != nil {
		beego.Info("查询错误")
		this.Redirect("/Article/ShowArticle",302)
		return
	}
	article:=models.Article{}
	article.Id=id
	o:=orm.NewOrm()
	o.Delete(&article)
	this.Redirect("/Article/ShowArticle",302)

}
func (this*ArticleController)HandleUpdateArticle()  {
	id, _ := this.GetInt("id")
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	file, header, err := this.GetFile("uploadname")
	defer file.Close()
	if err != nil {
		beego.Info("文件上传错误")
		return
	}
	//校验数据，如果数据出错,返回当前编辑页面
	if err !=nil || articleName =="" || content == "" {
		beego.Info("编辑数据不完整")
		this.Redirect("/Article/UpdateArticle?id="+strconv.Itoa(id),302)
		return
	}

	ext := path.Ext(header.Filename)
	filename:= time.Now().Format("2006-01-02 15:04:05")

	err = this.SaveToFile("uploadname", "static/img/"+filename+ext)
	if err!=nil {
		beego.Info("保存图片失败")
		return
	}

	o:=orm.NewOrm()
	article:=models.Article{Id: id}
	err = o.Read(&article)
	if err != nil {
		beego.Info("传递的文章id错误")
		this.Redirect("/Article/UpdateArticle?id="+strconv.Itoa(id),302)
		return
	}
	article.Atiname=articleName
	article.Acontent=content
	article.Aimg="static/img/"+filename+ext
	o.Update(&article)
	this.Redirect("/Article/ShowArticle",302)
}
func (this *ArticleController)ShowAddArticleType()  {
	o:=orm.NewOrm()
	var articleTypes[]models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	username := this.GetSession("userName")
	this.Data["username"]=username.(string)
	this.Data["articleTypes"] = articleTypes
	this.Layout = "layout.html"
	this.TplName = "addType.html"
}
func (this *ArticleController) HandleAddType() {
	typeName := this.GetString("typeName")
	if typeName=="" {
		beego.Info("添加类型不能为空")
		return
	}
	beego.Info(typeName)
	o:=orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName=typeName
	_, err := o.Insert(&articleType)
	if err != nil {
		beego.Info("数据插入失败")
		return
	}
	this.Redirect("/Article/AddArticleType",302)

}
func (this*ArticleController)DeleteType()  {
	id, _ := this.GetInt("id")
	o:=orm.NewOrm()
	articleType:=models.ArticleType{Id: id}
	_, err := o.Delete(&articleType)
	if err != nil {
		beego.Info("删除失败")
		return
	}
	this.Redirect("/Article/AddArticleType",302)
}