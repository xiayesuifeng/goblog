package article

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"gitlab.com/xiayesuifeng/goblog/category"
	"gitlab.com/xiayesuifeng/goblog/conf"
	"gitlab.com/xiayesuifeng/goblog/database"
	"io/ioutil"
	"os"
)

type Article struct {
	gorm.Model
	Title      string `gorm:"type:varchar(100);unique" json:"title" binding:"required"`
	Tag        string `json:"tag" binding:"required"`
	Uuid       string
	CategoryId uint `json:"category_id"`
	Private    bool `gorm:"default:0" json:"private"`
}

func AddArticle(title, tag string, categoryId uint, private bool, context string) error {
	_, err := category.GetCategory(categoryId)
	if err != nil {
		return err
	}

	var tmp uuid.UUID
	if tmp, err = uuid.NewV4(); err != nil {
		return err
	}
	md_uuid := tmp.String()

	article := Article{Title: title, Tag: tag, Uuid: md_uuid, CategoryId: categoryId, Private: private}
	db := database.Instance()
	if err := db.Create(&article).Error; err != nil {
		return err
	}

	return ioutil.WriteFile(conf.Conf.DataDir+"/article/"+md_uuid+".md", []byte(context), 0644)
}

func EditArticle(id, categoryId uint, title, tag, context string, private bool) error {
	article := Article{}

	db := database.Instance()
	if db.First(&article, id).RecordNotFound() {
		return errors.New("article not found")
	}

	err := ioutil.WriteFile(conf.Conf.DataDir+"/article/"+article.Uuid+".md", []byte(context), 0644)
	if err != nil {
		return err
	}

	return db.Model(&article).Updates(map[string]interface{}{
		"category_id": categoryId,
		"title": title,
		"tag": tag,
		"private": private,
	}).Error
}

func DeleteArticle(id int) error {
	db := database.Instance()

	article := Article{}
	if db.First(&article, id).RecordNotFound() {
		return errors.New("article not found")
	}

	os.Remove(conf.Conf.DataDir + "/article/" + article.Uuid + ".md")
	return db.Unscoped().Delete(&article).Error
}
