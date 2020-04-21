package main

import "github.com/astaxie/beego"
import _"FirstBeegoProject/models"
import _"FirstBeegoProject/routers"


func main() {
	beego.AddFuncMap("next",ShowNextPage)
	beego.AddFuncMap("pre",ShowPrePage)
	beego.Run()
}

func ShowNextPage(pageIndex int,pageCount int) int  {
	if pageIndex==pageCount {
		return pageCount
	}
	return pageIndex+1
}

func ShowPrePage(pageIndex int) int  {
	if pageIndex==1 {
		return 1
	}
	return pageIndex-1
}
