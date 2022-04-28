package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Rank struct {
	gorm.Model
	ID           uint `gorm:"primarykey"`
	Name         string
	Milliseconds uint `gorm:"index:seconds_idx,sort:desc"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RankEntry struct {
	Hash         string `json:"hash" binding:"required,alphanum,len=64"`
	Name         string `json:"name" binding:"required,alphanum,min=3,max=32"`
	Milliseconds uint   `json:"milliseconds" binding:"required,gte=0,lte=9999000"`
}

func main() {
	viper.SetConfigFile("app.env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Invalid config file", err)
	}

	dsn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%s",
		viper.Get("DB_HOST"),
		viper.Get("DB_NAME"),
		viper.Get("DB_USER"),
		viper.Get("DB_PASS"),
		viper.Get("DB_PORT"))
	db := getDB(dsn)

	salt := viper.Get("SALT").(string)
	port := viper.Get("HTTP_PORT").(string)
	startServer(db, port, salt)
}

func getDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Rank{})

	return db
}

func startServer(db *gorm.DB, port string, salt string) {
	api := Api{
		db:   db,
		salt: salt,
	}
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/ranks", api.GetRanks)
	router.POST("/ranks", api.PostRanks)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	err := router.Run(":" + port)
	if err != err {
		log.Fatal("error runing http server", err)
	}
}

type Api struct {
	db   *gorm.DB
	salt string
}

func (a Api) GetRanks(context *gin.Context) {
	var ranks []Rank
	err := a.db.Order("milliseconds desc").Limit(50).Find(&ranks).Error
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	rankEntries := make([]RankEntry, len(ranks))
	for i, r := range ranks {
		rankEntries[i] = RankEntry{
			Name:         r.Name,
			Milliseconds: r.Milliseconds,
		}
	}
	context.JSON(http.StatusOK, rankEntries)
}

func (a Api) PostRanks(context *gin.Context) {
	body := RankEntry{}

	if err := context.ShouldBindJSON(&body); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errMsg := map[string]string{}
			for _, fe := range ve {
				errMsg[fe.Field()] = fe.Error()
			}
			context.AbortWithStatusJSON(http.StatusBadRequest, errMsg)
		}
		return
	}

	if err := checkValidEntry(body, a.salt); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"hash": err.Error()})
		return
	}

	rank := Rank{
		Name:         body.Name,
		Milliseconds: body.Milliseconds,
	}

	if err := a.db.Create(&rank).Error; err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	context.JSON(http.StatusCreated, &body)
}

func checkValidEntry(r RankEntry, salt string) error {
	sum := sha256.Sum256([]byte(r.Name + strconv.Itoa(int(r.Milliseconds)) + salt))
	if hex.EncodeToString(sum[:]) != r.Hash {
		return fmt.Errorf("invalid hash... Don't be evil!")
	}
	return nil
}
