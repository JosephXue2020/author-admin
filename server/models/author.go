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
	UUID string `json:"uuid" es:"keyword"`

	// AuthorNature
	Name    string `json:"name" es:"keyword"`
	Gender  string `json:"gender" es:"keyword"`
	Nation  string `json:"nation" es:"keyword"`
	BornIn  string `json:"bornin" es:"text"`
	BornAt  string `json:"bornat" es:"date"`
	DeathAt string `json:"deathat" es:"date"`

	// AuthorResume
	GraduateAt string `json:"graduateat" es:"keyword"`
	Company    string `json:"company" es:"keyword"`
	Position   string `json:"position" es:"keyword"`
	JobTitle   string `json:"jobtitle" es:"keyword"`
	Honor      string `json:"honor" es:"keyword"`
	Telephone  string `json:"telephone" es:"keyword"`
	Cellphone  string `json:"cellphone" es:"keyword"`
	Email      string `json:"email" es:"keyword"`
	PostalAddr string `json:"postaladdr" es:"text"`
	Desc       string `json:"desc" es:"text"`

	// AuthorStudy
	Subject     string `json:"subject" es:"keyword"`
	Publication string `json:"publication" es:"text"`
	DBKPosition string `json:"dbkposition" es:"keyword"`

	Entries []Entry `gorm:"many2many:author_entry;" json:"entries" es:"object"`
}

type Entry struct {
	gorm.Model
	CDOI    string   `json:"cdoi" es:"keyword"`
	Authors []Author `gorm:"many2many:author_entry;" json:"authors" es:"object"`
	XMLText string   `json:"xmltext" es:"text"`
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
