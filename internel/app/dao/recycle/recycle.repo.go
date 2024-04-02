package recycle

import (
	"gorm.io/gorm"
)

type RecycleRepo struct {
	DB *gorm.DB
}
