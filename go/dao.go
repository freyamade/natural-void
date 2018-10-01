package naturalvoid

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"sync"
)

type DAO struct {
	DB *gorm.DB
}

var daoInstance *DAO
var daoOnce sync.Once

func GetDAO() (*DAO, error) {
	var err error
	daoOnce.Do(func() {
		daoInstance = &DAO{}
		err = daoInstance.new()
	})
	if err != nil {
		return nil, err
	}
	return daoInstance, nil
}

func (dao *DAO) new() error {
	// Initialize a new DAO object
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable")
	if err != nil {
		return err
	}
	dao.DB = db
	return nil
}

func (dao *DAO) Close() {
	dao.DB.Close()
}
