package entry

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/models"
	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/persistence"
)

type EntryApi struct {
	dao  persistence.EntryDAO
	salt string
}

func NewEntryApi(dao persistence.EntryDAO, salt string) EntryApi {
	return EntryApi{dao, salt}
}

func (a EntryApi) SetupRoutes(router *gin.Engine) {
	router.GET("/entry", a.GetEntries)
	router.POST("/entry", a.PostEntry)
	router.PUT("/entry", a.PostEntry)
}

func (a EntryApi) GetEntries(context *gin.Context) {
	entries, err := a.dao.GetEntries(context)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	rankEntries := make([]models.Entry, len(entries))
	for i, r := range entries {
		rankEntries[i] = models.Entry{
			Name:  r.Name,
			Score: r.Score,
		}
	}
	context.JSON(http.StatusOK, rankEntries)
}

func (a EntryApi) PostEntry(context *gin.Context) {
	body := models.Entry{}

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

	if err := CheckValidEntry(body, a.salt); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"hash": err.Error()})
		return
	}

	entry := models.Entry{
		Name:  body.Name,
		Score: body.Score,
	}

	result, err := a.dao.CreateEntry(context, entry)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	log.Printf("Result: %+v", result)
	context.JSON(http.StatusCreated, &body)
}
