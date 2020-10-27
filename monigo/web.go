package monigo

import (
	"net/http"
	"strconv"
)

func StartWebInterface() {
	LogInfo("Web interface was listening on :%d", Config.Web.Port)
	_ = http.ListenAndServe(":" + strconv.Itoa(Config.Web.Port), nil)
}