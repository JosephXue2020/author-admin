package models

import (
	"fmt"
	"goweb/author-admin/server/dao"

	"github.com/jinzhu/gorm"
)

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
	Subject     string `json:"subject"`
	Publication string `json:"publication"`
	DBKPosition string `json:"dbkposition"`
}

type Author struct {
	gorm.Model
	UUID string `json:"uuid"`

	AuthorNature
	// Name    string `json:"name"`
	// Gender  string `json:"gender"`
	// Nation  string `json:"nation"`
	// BornIn  string `json:"bornin"`
	// BornAt  string `json:"bornat"`
	// DeathAt string `json:"deathat"`

	AuthorResume
	// GraduateAt string `json:"graduateat"`
	// Company    string `json:"company"`
	// Position   string `json:"position"`
	// JobTitle   string `json:"jobtitle"`
	// Honor      string `json:"honor"`
	// Telephone  string `json:"telephone"`
	// Cellphone  string `json:"cellphone"`
	// Email      string `json:"email"`
	// PostalAddr string `json:"postaladdr"`
	// Desc       string `json:"desc"`

	AuthorStudy
	// Subject     string `json:"subject"`
	// Publication string `json:"publication"`
	// DBKPosition string `json:"dbkposition"`

	Entries []Entry `gorm:"many2many:author_entry;" json:"entries"`
}

type Entry struct {
	gorm.Model
	CDOI    string   `json:"cdoi"`
	Authors []Author `gorm:"many2many:author_entry;" json:"authors"`
	XMLText string   `json:"xmltext"`
}

func SelectAuthorByID(id int) (Author, error) {
	temp := Author{}
	dao.DB.Where("id = ?", id).First(&temp)

	if temp.ID > 0 {
		return temp, nil
	}

	err := fmt.Errorf("Author does not exist.")
	return temp, err
}

func SelectAuthorBatch(start, limit int, desc bool) []Author {
	var authors []Author
	orderStr := "id"
	if desc {
		orderStr += " desc"
	}
	dao.DB.Order(orderStr).Limit(limit).Offset(start).Find(&authors)
	return authors
}

func SelectAuthorAll(desc bool) []Author {
	var authors []Author
	orderStr := "id"
	if desc {
		orderStr += " desc"
	}
	dao.DB.Order(orderStr).Find(&authors)
	return authors
}

func CountAuthor() int {
	var count int
	dao.DB.Model(&Author{}).Count(&count)
	return count
}
