package dal

import (
	"github.com/crazygit/hpv-notification/config"
	"github.com/crazygit/hpv-notification/internal/dal/model"
	"gorm.io/driver/sqlite"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbInstance *gorm.DB

func getDatabaseLogger() logger.Interface {
	dbLogLevel := logger.Silent
	if config.AppConfig.Debug {
		dbLogLevel = logger.Info
	}
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  dbLogLevel,  // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
}

func configConnectionPool(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to config connection pool: %v", err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func InitDBInstance() {
	db, err := gorm.Open(sqlite.Open(config.AppConfig.Database.Dsn), &gorm.Config{
		Logger:          getDatabaseLogger(),
		CreateBatchSize: 1000,
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %s, err: %s", config.AppConfig.Database.Dsn, err)
	}
	configConnectionPool(db)

	err = db.AutoMigrate(&model.Place{})
	if err != nil {
		log.Fatalf("Failed migrate dal, err: %s", err)
	}
	dbInstance = db
}

// GetInstance get database session
func GetInstance() *gorm.DB {
	return dbInstance
}
