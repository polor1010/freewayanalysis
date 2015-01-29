package main

import (
	"freewayanalysis/controllers"
	_ "freewayanalysis/routers"
	"github.com/astaxie/beego"
	"github.com/robfig/cron"
)

func main() {
	c := cron.New()
	c.AddFunc("0 0 1 * * *", controllers.GetSmoothData)
	c.Start()

	//controllers.GetSmoothData()

	beego.Run()

}
