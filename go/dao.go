package naturalvoid

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"sync"
)

type DAO struct {
	DB *gorm.DB
}

var instance *DAO
var once sync.Once

func GetDAO() *DAO {
	once.Do(func() {
		instance = &DAO{}
		err := instance.new()
		if err != nil {
			panic(err)
		}
	})
	return instance
}

func (dao *DAO) new() error {
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
