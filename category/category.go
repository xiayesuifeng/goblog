package category

import (
	"github.com/jinzhu/gorm"
	"errors"
	"github.com/1377195627/goblog/database"
)

type Category struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);unique" json:"name" binding:"required"`
}

func AddCategory(name string) error {
	db := database.Instance()
	return db.Create(&Category{Name:name}).Error
}

func GetCategorys() []Category {
	db := database.Instance()
	categorys := make([]Category,0)
	db.Find(&categorys)
	return categorys
}

func GetCategory(id uint) (Category,error) {
	db := database.Instance()
	category := Category{}
	if db.Find(&category,id).RecordNotFound() {
		return category, errors.New("category not found")
	}

	return category,nil
}

func DeleteCategory(id int) error {
	db := database.Instance()
	return db.Unscoped().Delete(&Category{},id).Error
}

func (c *Category) SetName(name string) error {
	db := database.Instance()
	return db.Model(&c).Update("name",name).Error
}