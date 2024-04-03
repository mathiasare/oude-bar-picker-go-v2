package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Bar struct {
	gorm.Model
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	LocationUrl string         `json:"locationUrl"`
	ImageUrl    string         `json:"imageUrl"`
	Tags        []string       `gorm:"serializer:json;default:'[]'" json:"tags"`
}

type Vote struct {
	ID           string `gorm:"primarykey" json:"voteCode"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Participants []Participant  `gorm:"foreignKey:VoteId" json:"participants"`
	WinnerId     *uint          `json:"winner_id"`
	Winner       Bar            `gorm:"foreignKey:WinnerId" json:"winner"`
}

type Participant struct {
	gorm.Model
	Name   string `json:"name"`
	BarId  *uint  `json:"bar_id"`
	VoteId string `json:"vote_id"`
	Bar    Bar    `gorm:"foreignKey:BarId" json:"bar"`
	Vote   Vote   `gorm:"foreignKey:VoteId" json:"vote"`
}
