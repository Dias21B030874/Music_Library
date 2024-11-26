package routes

import (
	"music_library/songs"

	"github.com/gin-gonic/gin"
)

func RegisterSongRoutes(router *gin.Engine) {
	group := router.Group("/songs")
	{
		group.GET("", songs.GetSongs)
		group.GET("/:id/verses", songs.GetSongVerses)
		group.POST("", songs.CreateSong)
		group.PUT("/:id", songs.UpdateSong)
		group.DELETE("/:id", songs.DeleteSong)
	}
}
