package database

import (
	"github.com/kevinchapron/BasicLogger/Logging"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

const DB_FILENAME = "FSHK.db"

var singletonDB *FSHKDB

func GetDatabase() *FSHKDB {
	if singletonDB == nil {
		var _db = FSHKDB{}
		singletonDB = &_db
	}
	return singletonDB
}

type FSHKDB struct {
	innerDB    *gorm.DB
	existingDB bool
}

func (db *FSHKDB) checkIfExists() {
	info, err := os.Stat(DB_FILENAME)
	if os.IsNotExist(err) {
		db.existingDB = false
		return
	}
	db.existingDB = info.IsDir()
}

func (db *FSHKDB) defaultPopulation() {
	// populate
}

func (db *FSHKDB) Connect() {
	db.checkIfExists()

	if db.innerDB == nil {
		_db, err := gorm.Open(sqlite.Open(DB_FILENAME), &gorm.Config{})
		if err != nil {
			Logging.Error(err)
		}

		for _, item := range migrations() {
			_db.AutoMigrate(item)
		}

		db.innerDB = _db
	}

	if !db.existingDB {
		db.defaultPopulation()
	}
}
