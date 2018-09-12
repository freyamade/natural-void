package naturalvoid

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DAO struct {
	DB *gorm.DB
}

func (dao *DAO) New() {
	// Initialize a new DAO object
	db, err := gorm.Open("sqlite3", "./db")
	if err != nil {
		panic(err)
	}
	dao.DB = db
}

func (dao *DAO) Close() {
	dao.DB.Close()
}
