package naturalvoid

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// DB Models
// Model representing a User in the DB
// Users come from LDAP so this just stores some information for mapping stories to their DMs
type User struct {
	gorm.Model
	Stories  []Story
	UserName string
}

// A story is a collection of episodes
// ASTTW Chapter 1 is a story, as well as Delve Into Delirium
type Story struct {
	gorm.Model
	Description pq.StringArray `gorm:"type:text[]"`
	Episodes    []Episode
	Name        string
	ShortName   string
	Slug        string
	UserID      uint
}

// An Episode relates to a recorded file.
// Episodes are mapped to a Story for obvious reasons.
type Episode struct {
	gorm.Model
	Description pq.StringArray `gorm:"type:text[]"`
	Name        string
	Number      int
	StoryID     uint
}
