package models

type Song struct {
	ID          uint   `gorm:"primaryKey"`
	Group       string `gorm:"not null"`
	Title       string `gorm:"not null"`
	ReleaseDate string `gorm:"not null"`
	Text        string `gorm:"type:text"`
	Link        string `gorm:"not null"`
	CreatedAt   string `gorm:"autoCreateTime"`
	UpdatedAt   string `gorm:"autoUpdateTime"`
}
