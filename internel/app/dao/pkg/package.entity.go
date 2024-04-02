package pkg

import "gorm.io/gorm"

type Pkg struct {
	Id    int8
	Size  int8
	Name  string
	Price float32
}

func GetPkgDB(old *gorm.DB) *gorm.DB {
	return old.Model(&Pkg{})
}
