package controller

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/goblog/core"
	"os"
)

type Admin struct {
}

func (a *Admin) Login(ctx *gin.Context) {
	type Data struct {
		Password string `json:"password" form:"password" binding:"required"`
	}

	data := Data{}

	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": "password is null",
		})
		return
	}

	session := sessions.Default(ctx)

	md5Data := md5.Sum([]byte(data.Password))
	sha1Data := sha1.Sum([]byte(md5Data[:]))
	passwd := hex.EncodeToString(sha1Data[:])

	if passwd == core.Conf.Password {
		session.Set("login", true)
		session.Save()

		ctx.JSON(200, gin.H{
			"code": 0,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": "password errors",
		})
	}
}

func (a *Admin) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)

	login := session.Get("login")
	if login != nil {
		session.Set("login", nil)
		session.Save()
		ctx.JSON(200, gin.H{
			"code": 0,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": "no login",
		})
	}
}

func (a *Admin) GetInfo(ctx *gin.Context) {
	logo := "/api/logo"
	_, err := os.Stat(core.Conf.DataDir + "/logo")
	if err != nil {
		if os.IsNotExist(err) {
			logo = "none"
		}
	}
	ctx.JSON(200, gin.H{
		"name":        core.Conf.Name,
		"useCategory": core.Conf.UseCategory,
		"logo":        logo,
	})
}

func (a *Admin) PatchInfo(ctx *gin.Context) {
	type Data struct {
		Name        string `json:"name"`
		UseCategory *bool  `json:"useCategory"`
	}

	data := Data{}

	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": err.Error(),
		})
	} else if data.Name == "" && data.UseCategory == nil {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": "need name or useCategory",
		})
	} else {
		if data.Name != "" {
			core.Conf.Name = data.Name
		}

		if data.UseCategory != nil {
			core.Conf.UseCategory = *data.UseCategory
		}

		if err := core.SaveConf(); err != nil {
			ctx.JSON(200, gin.H{
				"code":    100,
				"message": err.Error(),
			})
		} else {
			ctx.JSON(200, gin.H{
				"code": 0,
			})
		}
	}
}

func (a *Admin) GetLogo(ctx *gin.Context) {
	ctx.File(core.Conf.DataDir + "/logo")
}

func (a *Admin) PutLogo(ctx *gin.Context) {
}
