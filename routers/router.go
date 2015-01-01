package routers

import (
	"freewaypredict/controllers"
	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/all/:date", &controllers.MainController{}, "get:All")
	beego.Router("/month/:date/:locationID", &controllers.MainController{}, "get:Month")
	beego.Router("/day/:date/:locationID", &controllers.MainController{}, "get:Day")
	beego.Router("/predict/", &controllers.MainController{}, "get:Predict")

	//beego.Router("/month/:locationID",&controllers.MainController{})

	beego.SetStaticPath("/images", "images")
	beego.SetStaticPath("/css", "css")
	beego.SetStaticPath("/js", "js")

}
