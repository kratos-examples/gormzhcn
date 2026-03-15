package models

import "gorm.io/gorm"

type T文章 struct {
	gorm.Model
	V标题 string `gorm:"column:title;type:varchar(255)" cnm:"V标题"`
	V内容 string `gorm:"column:content;type:text" cnm:"V内容"`
}

func (*T文章) TableName() string {
	return "articles"
}
