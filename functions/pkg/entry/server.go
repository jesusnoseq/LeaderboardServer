package entry

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/persistence"
)

func GetEntryServer(salt string) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	dao := persistence.GetEntryDAO(persistence.Memory)

	api := NewEntryApi(dao, salt)
	api.SetupRoutes(router)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	return router
}
