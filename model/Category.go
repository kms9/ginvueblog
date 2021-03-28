package model

import (
	"ginvueblog/upload"
	"ginvueblog/utils/errmsg"
	"gorm.io/gorm"
)

type Category struct {
	ID   uint   `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"type:varchar(20);not null" json:"name"`
}

// 查询分类是否存在
func CheckCategory(name string) (code int) {
	var cate Category
	upload.db.Select("id").Where("name = ?", name).First(&cate)
	if cate.ID > 0 {
		return errmsg.ERROR_CATENAME_USED //2001
	}
	return errmsg.SUCCSE
}

// 新增分类
func CreateCate(data *Category) int {
	err := upload.db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCSE
}

// 查询单个分类信息
func GetCateInfo(id int) (Category,int) {
	var cate Category
	upload.db.Where("id = ?",id).First(&cate)
	return cate,errmsg.SUCCSE
}

// 查询分类列表
func GetCate(pageSize int, pageNum int) ([]Category, int64) {
	var cate []Category
	var total int64
	upload.err = upload.db.Find(&cate).Limit(pageSize).Offset((pageNum - 1) * pageSize).Error
	upload.db.Model(&cate).Count(&total)
	if upload.err != nil && upload.err != gorm.ErrRecordNotFound {
		return nil, 0
	}
	return cate, total
}

// 编辑分类信息
func EditCate(id int, data *Category) int {
	var cate Category
	var maps = make(map[string]interface{})
	maps["name"] = data.Name

	upload.err = upload.db.Model(&cate).Where("id = ? ", id).Updates(maps).Error
	if upload.err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 删除分类
func DeleteCate(id int) int {
	var cate Category
	upload.err = upload.db.Where("id = ? ", id).Delete(&cate).Error
	if upload.err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
