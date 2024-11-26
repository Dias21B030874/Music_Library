package songs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"music_library/database"
	"music_library/interfaces"
	"music_library/models"

	"github.com/gin-gonic/gin"
)

// GetSongs godoc
// @Summary Get all songs
// @Description Get a list of songs with pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of songs per page" default(10)
// @Success 200 {object} interfaces.PaginatedSongsResponse
// @Router /songs [get]
func GetSongs(c *gin.Context) {
	db := database.GetDB()

	var songs []models.Song
	var total int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	db.Model(&models.Song{}).Count(&total)
	db.Offset(offset).Limit(pageSize).Find(&songs)

	c.JSON(http.StatusOK, interfaces.PaginatedSongsResponse{
		Songs:      convertToSongResponses(songs),
		TotalCount: int(total),
		Page:       page,
		PageSize:   pageSize,
	})
}

func GetSongVerses(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 1

	var song models.Song
	if err := db.First(&song, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
		return
	}

	verses := paginateVerses(song.Text, page, pageSize)
	if len(verses) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no verses found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"verses": verses})
}

func fetchSongDetails(group, song string) (*interfaces.SongResponse, error) {
	// Формируем URL для запроса к внешнему API
	url := fmt.Sprintf("http://localhost:8081/info?group=%s&song=%s", group, song)

	// Отправляем GET запрос
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call external API: %v", err)
	}
	defer resp.Body.Close()

	// Если ответ не успешный, выводим ошибку
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned error: %s", resp.Status)
	}

	// Читаем тело ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Десериализуем JSON в структуру SongResponse
	var songDetails interfaces.SongResponse
	if err := json.Unmarshal(body, &songDetails); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %v", err)
	}

	// Возвращаем успешно полученные данные
	return &songDetails, nil
}

func CreateSong(c *gin.Context) {
	db := database.GetDB()

	var req interfaces.NewSongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	songDetails, err := fetchSongDetails(req.Group, req.Song)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newSong := models.Song{
		Group:       req.Group,
		Title:       req.Song,
		ReleaseDate: songDetails.ReleaseDate,
		Text:        songDetails.Text,
		Link:        songDetails.Link,
	}

	db.Create(&newSong)
	c.JSON(http.StatusCreated, newSong)
}

func UpdateSong(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")

	var song models.Song
	if err := db.First(&song, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
		return
	}

	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&song)
	c.JSON(http.StatusOK, song)
}

func DeleteSong(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")

	if err := db.Delete(&models.Song{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

func convertToSongResponses(songs []models.Song) []interfaces.SongResponse {
	var responses []interfaces.SongResponse
	for _, song := range songs {
		responses = append(responses, interfaces.SongResponse{
			ID:          song.ID,
			Group:       song.Group,
			Title:       song.Title,
			ReleaseDate: song.ReleaseDate,
			Link:        song.Link,
		})
	}
	return responses
}

func paginateVerses(text string, page, pageSize int) []string {
	verses := splitTextByVerses(text)
	start := (page - 1) * pageSize
	if start >= len(verses) {
		return nil
	}

	end := start + pageSize
	if end > len(verses) {
		end = len(verses)
	}

	return verses[start:end]
}

func splitTextByVerses(text string) []string {
	return strings.Split(text, "\n\n")
}
