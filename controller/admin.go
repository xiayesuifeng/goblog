package controller

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"gitlab.com/xiayesuifeng/goblog/article"
	"gitlab.com/xiayesuifeng/goblog/category"
	"gitlab.com/xiayesuifeng/goblog/conf"
	"gitlab.com/xiayesuifeng/goblog/core"
	"gitlab.com/xiayesuifeng/goblog/database"
	"io"
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
		ctx.JSON(200, core.FailResult("password is null"))
		return
	}

	session := sessions.Default(ctx)

	md5Data := md5.Sum([]byte(data.Password))
	sha1Data := sha1.Sum([]byte(md5Data[:]))
	passwd := hex.EncodeToString(sha1Data[:])

	if passwd == conf.Conf.Password {
		session.Set("login", true)
		session.Save()

		ctx.JSON(200, core.SuccessResult())
	} else {
		ctx.JSON(200, core.FailResult("password errors"))
	}
}

func (a *Admin) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)

	login := session.Get("login")
	if login != nil {
		session.Set("login", nil)
		session.Save()
		ctx.JSON(200, core.SuccessResult())
	} else {
		ctx.JSON(200, core.Result(core.ResultUnauthorizedCode, "no login"))
	}
}

func (a *Admin) GetInfo(ctx *gin.Context) {
	logo := "/api/logo"
	_, err := os.Stat(conf.Conf.DataDir + "/logo")
	if err != nil {
		if os.IsNotExist(err) {
			logo = "none"
		}
	}
	ctx.JSON(200, gin.H{
		"name":        conf.Conf.Name,
		"useCategory": conf.Conf.UseCategory,
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
		ctx.JSON(200, core.FailResult(err.Error()))
	} else if data.Name == "" && data.UseCategory == nil {
		ctx.JSON(200, core.FailResult("need name or useCategory"))
	} else {
		if data.Name != "" {
			conf.Conf.Name = data.Name
		}

		if data.UseCategory != nil {
			conf.Conf.UseCategory = *data.UseCategory
			if !conf.Conf.UseCategory {
				db := database.Instance()

				tmp := category.Category{Name: "other"}
				if db.Where(&tmp).First(&tmp).RecordNotFound() {
					if err := db.Create(&tmp).Error; err != nil {
						ctx.JSON(200, core.FailResult(err.Error()))
						return
					}
				}

				db.Model(&article.Article{}).Updates(article.Article{CategoryId: tmp.ID})
				conf.Conf.OtherCategoryId = tmp.ID
			}
		}

		if err := conf.SaveConf(); err != nil {
			ctx.JSON(200, core.FailResult(err.Error()))
		} else {
			ctx.JSON(200, core.SuccessResult())
		}
	}
}

func (a *Admin) GetLogo(ctx *gin.Context) {
	ctx.File(conf.Conf.DataDir + "/logo")
}

func (a *Admin) PutLogo(ctx *gin.Context) {
	logo, _, err := ctx.Request.FormFile("logo")
	if err != nil {
		ctx.JSON(200, core.FailResult(err.Error()))
	} else {
		file, err := os.Create(conf.Conf.DataDir + "/logo")
		if err != nil {
			ctx.JSON(200, core.FailResult(err.Error()))
			return
		}

		defer file.Close()

		io.Copy(file, logo)
		ctx.JSON(200, core.SuccessResult())
	}
}

func (a *Admin) GetAssets(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	ctx.File(conf.Conf.DataDir + "/assets/" + uuid)
}

func (a *Admin) PutAssets(ctx *gin.Context) {
	assets, _, err := ctx.Request.FormFile("assets")
	if err != nil {
		ctx.JSON(200, core.FailResult(err.Error()))
	} else {
		uuid := a.GetAssetsUuid()
		file, err := os.Create(conf.Conf.DataDir + "/assets/" + uuid)
		if err != nil {
			ctx.JSON(200, core.FailResult(err.Error()))
			return
		}

		defer file.Close()

		io.Copy(file, assets)
		ctx.JSON(200, core.SuccessDataResult("path", "/api/assets/"+uuid))
	}
}

func (a *Admin) GetAssetsUuid() string {
	if uuid, err := uuid.NewV4(); err == nil {
		if _, err := os.Stat(conf.Conf.DataDir + "/assets/" + uuid.String()); os.IsNotExist(err) {
			return uuid.String()
		}
	}

	return a.GetAssetsUuid()
}
