package models

import "gorm.io/gorm"

// T学生 与 demo1kratos 的 students 表结构一致。这里是文章服务、不拥有学生表，
// 保留这份镜像仅用于建文章前校验学生存在（两服务共用一个库）。
type T学生 struct {
	gorm.Model
	V名字 string `gorm:"column:name;type:varchar(255)" cnm:"V名字"`
	V年龄 int32  `gorm:"column:age;type:int" cnm:"V年龄"`
	V班级 string `gorm:"column:class_name;type:varchar(255)" cnm:"V班级"`
}

func (*T学生) TableName() string {
	return "students"
}
