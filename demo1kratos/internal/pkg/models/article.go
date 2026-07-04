package models

import "gorm.io/gorm"

// T文章 与 demo2kratos 的 articles 表结构一致。这里是学生服务、不拥有文章表，
// 保留这份镜像仅用于删学生时顺带删掉他名下的文章（两服务共用一个库）。
type T文章 struct {
	gorm.Model
	V标题   string `gorm:"column:title;type:varchar(255)" cnm:"V标题"`
	V内容   string `gorm:"column:content;type:text" cnm:"V内容"`
	V学生ID int64  `gorm:"column:student_id;index" cnm:"V学生ID"`
}

func (*T文章) TableName() string {
	return "articles"
}
