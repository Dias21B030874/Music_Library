package interfaces

type NewSongRequest struct {
	Group string `json:"group" binding:"required"`
	Song  string `json:"song" binding:"required"`
}

type SongResponse struct {
	ID          uint   `json:"id"`
	Group       string `json:"group"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type PaginatedSongsResponse struct {
	Songs      []SongResponse `json:"songs"`
	TotalCount int            `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
}

type SongVersesResponse struct {
	Verse   string `json:"verse"`
	VerseID int    `json:"verse_id"`
}
