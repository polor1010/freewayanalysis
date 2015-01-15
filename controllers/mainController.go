package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Predict() {

	GetSmoothData()

	fmt.Fprint(this.Ctx.ResponseWriter, "predict done")

}

func (this *MainController) Month() {
	this.Data["xsrfdata"] = template.HTML(this.XsrfFormHtml())

	date, _ := this.Ctx.Input.Params[":date"]
	locationID, _ := this.Ctx.Input.Params[":locationID"]

	fmt.Println(date)
	var location Location
	location.LocationID = locationID

	speedChartData := GetMonthByLocationID(date, location)

	this.Data["json"] = &speedChartData
	this.ServeJson()

	fmt.Println("(this *MainController) Month()")

}

func (this *MainController) Day() {

	this.Data["xsrfdata"] = template.HTML(this.XsrfFormHtml())

	date, _ := this.Ctx.Input.Params[":date"]

	fmt.Println("(this *MainController) Day()")
	locationID, _ := this.Ctx.Input.Params[":locationID"]

	fmt.Println(date)
	fmt.Println(locationID)
	var location Location
	location.LocationID = locationID

	speedChartData := GetDayByLocationID(date, location)

	this.Data["json"] = &speedChartData
	this.ServeJson()

}

func (this *MainController) All() {

	this.Data["xsrfdata"] = template.HTML(this.XsrfFormHtml())

	date, _ := this.Ctx.Input.Params[":date"]

	locations := GetAll(date)

	this.Data["json"] = &locations
	this.ServeJson()

}
