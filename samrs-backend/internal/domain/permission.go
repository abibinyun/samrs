package domain

type Permission struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	Name   string `gorm:"size:100;not null" json:"name"`
	Slug   string `gorm:"size:100;unique;not null" json:"slug"`
	Module string `gorm:"size:50" json:"module"`
}
