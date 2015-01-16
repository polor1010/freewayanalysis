package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"html/template"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {

	this.XSRFExpire = 7200
	this.Data["xsrfdata"] = template.HTML(this.XsrfFormHtml())
	//ctx.Output.Header("Access-Control-Allow-Origin", "*")

	//this.Ctx.SetHeader("Access-Control-Allow-Origin", "*", true)

}

func (this *MainController) Predict() {

	GetSmoothData()

	fmt.Fprint(this.Ctx.ResponseWriter, "predict done")

}

func (this *MainController) Month() {

	date, _ := this.Ctx.Input.Params[":date"]
	locationID, _ := this.Ctx.Input.Params[":locationID"]

	fmt.Println(date)
	var location Location
	location.LocationID = locationID

	speedChartData := GetMonthByLocationID(date, location)

	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	this.Data["json"] = &speedChartData
	this.ServeJson()

	fmt.Println("(this *MainController) Month()")

}

func (this *MainController) Day() {

	date, _ := this.Ctx.Input.Params[":date"]

	fmt.Println("(this *MainController) Day()")
	locationID, _ := this.Ctx.Input.Params[":locationID"]

	fmt.Println(date)
	fmt.Println(locationID)
	var location Location
	location.LocationID = locationID

	speedChartData := GetDayByLocationID(date, location)

	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	this.Data["json"] = &speedChartData
	this.ServeJson()

}

func (this *MainController) All() {

	fmt.Println("all")
	date, _ := this.Ctx.Input.Params[":date"]

	locations := GetAll(date)

	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	this.Data["json"] = &locations
	this.ServeJson()

}
