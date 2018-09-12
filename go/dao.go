package naturalvoid

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DAO struct {
	DB *gorm.DB
}

func (dao *DAO) New() (error) {
	// Initialize a new DAO object
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		return err
	}
	dao.DB = db
    return nil
}

func (dao *DAO) Close() {
	dao.DB.Close()
}
