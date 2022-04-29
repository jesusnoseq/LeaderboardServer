// package rdb

// import (
// 	"time"

// func getDB(dsn string) *gorm.DB {
// 	db, err := gorm.Open(postgres.New(postgres.Config{
// 		DSN:                  dsn,
// 		PreferSimpleProtocol: true,
// 	}), &gorm.Config{})
// 	if err != nil {
// 		panic(err)
// 	}

// 	db.AutoMigrate(&Rank{})

// 	return db
// }

// 	"gorm.io/gorm"
// )

// type Rank struct {
// 	gorm.Model
// 	ID           uint `gorm:"primarykey"`
// 	Name         string
// 	Milliseconds uint `gorm:"index:seconds_idx,sort:desc"`
// 	CreatedAt    time.Time
// 	UpdatedAt    time.Time
// }

// dsn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%s",
// viper.Get("DB_HOST"),
// viper.Get("DB_NAME"),
// viper.Get("DB_USER"),
// viper.Get("DB_PASS"),
// viper.Get("DB_PORT"))
// db := getDB(dsn)