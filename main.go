package main

import (
	"freewaypredict/controllers"
	_ "freewaypredict/routers"
	"github.com/astaxie/beego"
	"github.com/robfig/cron"
)

func main() {
	c := cron.New()
	c.AddFunc("0 0 03 * * *", controllers.GetSmoothData)
	c.Start()
	//controllers.GetSmoothData()

	beego.Run()

}
