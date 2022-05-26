package models

import "github.com/jinzhu/gorm"

type AuthorNature struct {
	Name    string `json:"name"`
	Gender  string `json:"gender"`
	Nation  string `json:"nation"`
	BornIn  string `json:"bornin"`
	BornAt  string `json:"bornat"`
	DeathAt string `json:"deathat"`
}

type AuthorResume struct {
	GraduateAt string `json:"graduateat"`
	Company    string `json:"company"`
	Position   string `json:"position"`
	JobTitle   string `json:"jobtitle"`
	Honor      string `json:"honor"`
	Telephone  string `json:"telephone"`
	Cellphone  string `json:"cellphone"`
	Email      string `json:"email"`
	PostalAddr string `json:"postaladdr"`
	Desc       string `json:"desc"`
}

type AuthorStudy struct {
	Subject     string `json:""`
	Publication string `json:""`
	DBKPosition string `json:""`
}

type Author struct {
	gorm.Model
	UUID string `json:"uuid"`
	AuthorNature
	AuthorResume
	AuthorStudy
	Entries []Entry `gorm:"many2many:author_entry;" json:"entries"`
}

type Entry struct {
	gorm.Model
	CDOI    string   `json:"cdoi"`
	Authors []Author `gorm:"many2many:author_entry;" json:"authors"`
	XMLText string   `json:"xmltext"`
}
