package hello

import (
	"github.com/astaxie/beego"
)

type HelloController struct {
	beego.Controller
}

// @Title Get
// @Description print the string
// @Success 200
// @Failure 403
// @router / [Get]
func (h *HelloController) Get() {
	h.Ctx.WriteString("hello world")
}
