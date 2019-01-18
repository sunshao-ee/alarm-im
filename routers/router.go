package routers

import (
	"alarm-im/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/im", &controllers.ImController{})
}
