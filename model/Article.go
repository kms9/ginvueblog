package model

import (
	"ginvueblog/setup"
	"ginvueblog/upload"
	"ginvueblog/utils/errmsg"
	"gorm.io/gorm"
)

type Article struct {
	Category Category `gorm:"foreignkey:Cid"`
	gorm.Model
	Title        string `gorm:"type:varchar(100);not null" json:"title"`
	Cid          int    `gorm:"type:int;not null" json:"cid"`
	Desc         string `gorm:"type:varchar(200)" json:"desc"`
	Content      string `gorm:"type:longtext" json:"content"`
	Img          string `gorm:"type:varchar(100)" json:"img"`
	CommentCount int    `gorm:"type:int;not null;default:0" json:"comment_count"`
	ReadCount    int    `gorm:"type:int;not null;default:0" json:"read_count"`
}

// 新增文章
func CreateArt(data *Article) int {
	err := setup.MysqlDB.Create(&data).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCSE
}

//  查询分类下的所有文章
func GetCateArt(id int, pageSize int, pageNum int) ([]Article, int, int64) {
	var cateArtList []Article
	var total int64
	
	err := setup.MysqlDB.Select("article.id,title, `cid`, img, created_at, `desc`, comment_count, read_count").Preload("Category").Limit(pageSize).Offset((pageNum-1)*pageSize).Where(
		"cid =?", id).Find(&cateArtList).Error
	
	setup.MysqlDB.Model(&cateArtList).Where(
		"cid =?", id).Count(&total)
	if err != nil {
		return nil, errmsg.ERROR_CATE_NOT_EXIST, 0
	}
	return cateArtList, errmsg.SUCCSE, total
}

//  查询单个文章
func GetArtInfo(id int) (Article, int) {
	var art Article
	err := setup.MysqlDB.Where("id = ?", id).Preload("Category").First(&art).Error
	setup.MysqlDB.Model(&art).Where("id = ?", id).UpdateColumn("read_count", gorm.Expr("read_count + ?", 1))
	if err != nil {
		return art, errmsg.ERROR_ART_NOT_EXIST
	}
	return art, errmsg.SUCCSE
}

//  查询文章列表
func GetArt(title string, pageSize int, pageNum int) ([]Article, int, int64) {
	var articleList []Article
	var err error
	var total int64
	
	if title == "" {
		err = setup.MysqlDB.Debug().Select("article.id, title, img, created_at, updated_at, `desc`, comment_count, read_count, Category.name").Preload("Category").Joins("Category").Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("Created_At DESC").Find(&articleList).Error
		// 单独计数
		setup.MysqlDB.Model(&articleList).Count(&total)
		if err != nil {
			return nil, errmsg.ERROR, 0
		}
		return articleList, errmsg.SUCCSE, total
	}else{
		err = setup.MysqlDB.Select("article.id,title, img, created_at, updated_at, `desc`, comment_count, read_count, category.name").Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Preload("Category").Where("title LIKE ?",
			title+"%",
		).Find(&articleList).Error
		// 单独计数
		setup.MysqlDB.Model(&articleList).Where("title LIKE ?",
			title+"%",
		).Count(&total)

		if err != nil {
			return nil, errmsg.ERROR, 0
		}
		return articleList, errmsg.SUCCSE, total
	}

}

// 编辑文章
func EditArt(id int, data *Article) int {
	var art Article
	var maps = make(map[string]interface{})
	maps["title"] = data.Title
	maps["cid"] = data.Cid
	maps["desc"] = data.Desc
	maps["content"] = data.Content
	maps["img"] = data.Img
	
	upload.err = setup.MysqlDB.Model(&art).Where("id = ? ", id).Updates(&maps).Error
	if upload.2+err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 删除文章
func DeleteArt(id int) int {
	var art Article
	upload.err = setup.MysqlDB.Where("id = ? ", id).Delete(&art).Error
	if upload.err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
