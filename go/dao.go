package naturalvoid

import (
	"bytes"
	"text/template"
	"sync"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
	confString := "host={{.DBHost}} port={{.DBPort}} user={{.DBUser}} dbname={{.DBName}} password={{.DBPass}} sslmode={{.DBSecure}}"
	conf := GetConf()
	t, err := template.New("conf").Parse(confString)
	var buffer bytes.Buffer
	err = t.Execute(&buffer, conf)
	if err != nil {
		panic(err)
	}
	// Initialize a new DAO object
	db, err := gorm.Open("postgres", buffer.String())
	if err != nil {
		return err
	}
	dao.DB = db
	return nil
}

func (dao *DAO) Close() {
	dao.DB.Close()
}
