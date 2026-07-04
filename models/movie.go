package models

type Movie struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"not null"`
	Director    string `json:"director" gorm:"not null"`
	Genre       string `json:"genre" gorm:"not null"`
	ReleaseYear int    `json:"release_year" gorm:"not null"`
	Description string `json:"description"`
}
