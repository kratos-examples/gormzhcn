package models

import "gorm.io/gorm"

type T学生 struct {
	gorm.Model
	V名字 string `gorm:"column:name;type:varchar(255)" cnm:"V名字"`
}

func (*T学生) TableName() string {
	return "students"
}
