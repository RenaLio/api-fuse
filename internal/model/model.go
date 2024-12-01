package model

import "time"

type BaseModel struct {
	ID        uint64    `gorm:"primary_key;auto_increment;comment:主键ID"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
	//DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"`
}
