package article

import (
	"github.com/1377195627/goblog/category"
	"github.com/jinzhu/gorm"
	"github.com/1377195627/goblog/database"
	"github.com/satori/go.uuid"
	"io/ioutil"
)

type Article struct {
	gorm.Model
	Title    string `gorm:"type:varchar(100);unique" json:"title" binding:"required"`
	Tag      string `json:"tag" binding:"required"`
	Uuid     string
	CategoryId uint `json:"category_id" binding:"required"`
}

func AddArticle(title, tag string, categoryId uint,context string) error {
	_,err := category.GetCategory(categoryId)
	if err != nil {
		return err
	}

	md_uuid := uuid.NewV1().String()

	article := Article{Title:title,Tag:tag,Uuid:md_uuid,CategoryId:categoryId}
	db := database.Instance()
	if err:= db.Create(&article).Error;err!=nil{
		return err
	}

	return ioutil.WriteFile("data/article/"+md_uuid+".md", []byte(context), 0644)
}
