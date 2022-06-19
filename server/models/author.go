package models

import (
	"fmt"
	"goweb/author-admin/server/dao"
	"goweb/author-admin/server/pkg/e"
	"goweb/author-admin/server/pkg/util"
	"log"
	"regexp"

	"github.com/astaxie/beego/validation"
	"github.com/jinzhu/gorm"
)

type Author struct {
	// gorm.Model
	Model `json:"model"`
	UUID  string `json:"uuid" es:"keyword"`

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

	// many to man foreign key
	Entries []Entry `gorm:"many2many:author_entry;" json:"entries" es:"object"`
}

type Entry struct {
	gorm.Model
	CDOI    string   `json:"cdoi" es:"keyword"`
	Authors []Author `gorm:"many2many:author_entry;" json:"authors" es:"object"`
	RawText string   `json:"rawtext" es:"text"`
}

func ValidateAuthor(a Author) bool {
	valid := validation.Validation{}
	valid.Required(a.Name, "name").Message("Name field of Author struct can not be empty.")
	valid.Match(a.Gender, regexp.MustCompile(`^[男|女]$`), "gender").Message("Gender field of Author struct must be 男 or 女.")

	if !valid.HasErrors() {
		return true
	}

	for _, err := range valid.Errors {
		log.Println("Author struct has error: ", err)
	}
	return false
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

func AuthorExist(name string) bool {
	temp := Author{}
	dao.DB.Where("Name = ?", name).First(&temp)
	if temp.ID > 0 {
		return true
	}
	return false
}

func ValidateAuthorCreation(a Author) (bool, int) {
	code := e.SUCCESS

	// 调用Author的Valid
	if !ValidateAuthor(a) {
		code := e.INVALID_PARAMS
		return false, code
	}

	// 判断是否存在
	if AuthorExist(a.Name) {
		code = e.ERROR_USER_ALREADY_EXIST
		return false, code
	}

	return true, code
}

func AddAuthor(a Author) error {
	a.UUID = util.GenUUID()

	if valid, code := ValidateAuthorCreation(a); !valid {
		err := fmt.Errorf(e.GetMsg(code))
		return err
	}

	dao.DB.Create(&a)
	return nil
}

func DeleteAuthorByID(id int) error {
	dao.DB.Delete(&Author{}, id)
	return nil
}

func UpdateAuthor() {

}
