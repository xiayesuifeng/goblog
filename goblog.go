package goblog

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"log"
	"os"
)

func InstallRouter(context *gin.Context) {
	//if context.Request.Method==http.MethodGet {
	//	context.HTML(http.StatusOK, "install.html", gin.H{})
	//	return
	//}

	err := context.Bind(Conf)
	if err != nil {
		log.Fatal(err)
	}

	if err = Install(); err != nil {
		log.Fatal(err)
		context.String(http.StatusOK, "安装失败")
	}

	context.String(http.StatusOK, "安装完成")

}

func HomeRouter(context *gin.Context){
	if _, err := os.Stat("goblog.lock"); err != nil {
		if os.IsNotExist(err) {
			//context.Redirect(http.StatusMovedPermanently, "/install")
			context.String(http.StatusOK,"博客未安装")
			return
		}
	}
	context.String(http.StatusOK,"GoBlog")
}