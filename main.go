package main

import (
	"freewayanalysis/controllers"
	_ "freewayanalysis/routers"
	"github.com/astaxie/beego"
	"github.com/robfig/cron"
)

func main() {
	c := cron.New()
	c.AddFunc("0 0 03 * * *", controllers.GetSmoothData)
	c.Start()

	beego.Run()

}
