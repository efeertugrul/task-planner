package utility

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	DB              *gorm.DB
	sqlDB           *sql.DB
	testModels      []interface{}
	functionNameMap = map[string]func(string) string{}

	driverFunctionMap = map[string]string{}
	registeredDrivers = map[string]struct{}{}
)

func GetTestDB() *gorm.DB {
	if DB == nil {
		DB = NewTestDB()
	}

	return DB
}

func NewTestDB() *gorm.DB {
	if DB == nil {
		Initialize()
	}

	return DB
}

// Initialize opens a connection to the database.
func Initialize() {
	var (
		err      error
		log      = logger.Info
		connPool *sql.Conn
	)
	driverName := "sqlite3"
	if len(functionNameMap) > 0 {
		driverName = "sqlite3_extended"
		RegisterDBFunctions()
	}

	if sqlDB, err = sql.Open(driverName, "file::memory:?cache=shared"); err != nil {
		fmt.Printf("%s", err.Error())
	}

	if connPool, err = sqlDB.Conn(context.Background()); err != nil {
		fmt.Printf("%s", err.Error())
	}

	sqliteDialector := sqlite.Open("file::memory:?cache=shared").(*sqlite.Dialector)
	sqliteDialector.DriverName = driverName
	sqliteDialector.Conn = connPool

	if DB, err = gorm.Open(sqliteDialector, &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
			NameReplacer:  nil,
			NoLowerCase:   false,
		},
		FullSaveAssociations:                     false,
		Logger:                                   logger.Default.LogMode(log),
		DisableForeignKeyConstraintWhenMigrating: true,
		IgnoreRelationshipsWhenMigrating:         true,
	}); err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
}

// CloseTestDB closes the database connection.
func CloseTestDB() {
	ClearDatabase()

	if sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			fmt.Printf("%s", err.Error())
		}

		return
	}

	if db, err := DB.DB(); err != nil {
		fmt.Printf("%s", err.Error())
	} else {
		if err = db.Close(); err != nil {
			fmt.Printf("%s", err.Error())
		}
	}
}

func ClearTables() {
	for _, model := range testModels {
		if err := DB.Where("1=1").Unscoped().Delete(model).Error; err != nil {
			fmt.Printf("%s", err.Error())
		}
	}
}

func ClearDatabase() {
	for _, model := range testModels {
		if err := DB.Migrator().DropTable(model); err != nil {
			fmt.Printf("%s", err.Error())
		}
	}
}

// AutoMigrate auto-migrates the models to create the tables.
func AutoMigrate(models ...interface{}) {
	currentModels := testModels
	if len(currentModels) > 0 {
		ClearDatabase()
	}

	testModels = models

	for i := range testModels {
		if err := DB.Migrator().AutoMigrate(models[i]); err != nil {
			fmt.Printf("%s", err.Error())
		}
	}
}

func RegisterDBFunctions() {
	for driverName, functionName := range driverFunctionMap {
		if _, driverRegistered := registeredDrivers[driverName]; !driverRegistered {
			if function, ok := functionNameMap[functionName]; ok {
				sql.Register(driverName, &sqlite3.SQLiteDriver{
					ConnectHook: func(conn *sqlite3.SQLiteConn) error {
						return conn.RegisterFunc(functionName, function, true)
					},
				})
			}

			registeredDrivers[driverName] = struct{}{}
		}

	}
}
