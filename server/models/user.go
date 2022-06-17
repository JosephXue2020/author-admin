package models

import (
	"errors"
	"fmt"
	"goweb/author-admin/server/dao"
	"goweb/author-admin/server/pkg/e"
	"goweb/author-admin/server/pkg/setting"
	"goweb/author-admin/server/pkg/util"
)

// User grade
const (
	GUEST = iota
	NORMAL
	ADMIN
	SUPER
)

var GradeMap = map[int]string{
	GUEST:  "guest",
	NORMAL: "normal",
	ADMIN:  "admin",
	SUPER:  "super",
}

var RoleMap = func() map[string]int {
	m := make(map[string]int)
	for k, v := range GradeMap {
		m[v] = k
	}
	return m
}()

func RoleKeywords() []string {
	var slc []string
	for _, v := range GradeMap {
		slc = append(slc, v)
	}
	return slc
}

type User struct {
	// gorm.Model
	ID         int    `gorm:"primaryKey" json:"id" es:"keyword"`
	UUID       string `gorm:"unique" json:"uuid" es:"keyword"`
	Username   string `gorm:"unique" json:"username" es:"keyword"` // 用户名要求唯一
	Password   string `json:"password" es:"keyword"`
	Department string `json:"department" es:"text"`
	Role       string `json:"role" es:"keyword"`
	Creater    string `json:"creater" es:"keyword"`
	CreateOn   string `json:"createon" es:"keyword"`
}

// 达到阈值则允许
func (u *User) Permission(threshold int) bool {
	if threshold < GUEST || threshold > SUPER {
		return false
	}

	grade, ok := RoleMap[u.Role]
	if !ok || grade < threshold {
		return false
	}

	return true
}

// 低于阈值则允许
func (u *User) LessPermission(threshold int) bool {
	if threshold < GUEST || threshold > SUPER {
		return false
	}

	grade, ok := RoleMap[u.Role]
	if !ok || grade >= threshold {
		return false
	}

	return true
}

func CheckUser(username, password string) bool {
	var u User
	dao.DB.Select("id").Where(User{Username: username, Password: password}).First(&u)
	if u.ID > 0 {
		return true
	}
	return false
}

func SelectUserByID(id int) (User, error) {
	temp := User{}
	dao.DB.Where("id = ?", id).First(&temp)

	if temp.ID > 0 {
		return temp, nil
	}

	err := fmt.Errorf("User does not exist.")
	return temp, err
}

func SelectUserByUsername(username string) (User, error) {
	// 不验证唯一性
	temp := User{}
	dao.DB.Where("Username = ?", username).First(&temp)

	if temp.ID > 0 {
		return temp, nil
	}

	err := fmt.Errorf("User does not exist.")
	return temp, err
}

func SelectUserBatch(start, limit int, desc bool) []User {
	var users []User
	orderStr := "id"
	if desc {
		orderStr += " desc"
	}
	dao.DB.Order(orderStr).Limit(limit).Offset(start).Find(&users)
	return users
}

func SelectUserAll(desc bool) []User {
	var users []User
	orderStr := "id"
	if desc {
		orderStr += " desc"
	}
	dao.DB.Order(orderStr).Find(&users)
	return users
}

func CountUser() int {
	var count int
	dao.DB.Model(&User{}).Count(&count)
	return count
}

func AllowedRole(role string) bool {
	roles := RoleKeywords()
	return util.ContainStr(roles, role)
}

func GetGradeByName(username string) (int, error) {
	obj, err := SelectUserByUsername(username)
	if err != nil {
		return -1, err
	}

	grade, ok := RoleMap[obj.Role]
	if !ok {
		err = errors.New("Failed to get the grade of the user.")
		return -1, err
	}

	return grade, nil
}

func UserExist(username string) bool {
	temp := User{}
	dao.DB.Where("Username = ?", username).First(&temp)
	if temp.ID > 0 {
		return true
	}
	return false
}

func ValidateUserCreation(u User, creater string) (bool, int) {
	code := e.SUCCESS
	// 判断类型是否合法
	if !AllowedRole(u.Role) {
		code = e.ERROR_USER_INVALID
		return false, code
	}

	// 判断是否有足够权限：只能创建权限小于自己的用户
	createrGrade, err := GetGradeByName(creater)
	if err != nil {
		code = e.ERROR_USER_INVALID
		return false, code
	}
	if !u.LessPermission(createrGrade) {
		code = e.ERROR_USER_LACK_AUTHORITY
		return false, code
	}

	// 判断是否存在
	if UserExist(u.Username) {
		code = e.ERROR_USER_ALREADY_EXIST
		return false, code
	}

	return true, code
}

// 用户只能添加权限小于自己的用户
func AddUser(username, password, department, role, creater string) error {
	u := User{
		UUID:       util.GenUUID(),
		Username:   username,
		Password:   password,
		Department: department,
		Role:       role,
		Creater:    creater,
		CreateOn:   util.CurrentTimeStr(),
	}

	if valid, code := ValidateUserCreation(u, creater); !valid {
		err := fmt.Errorf(e.GetMsg(code))
		return err
	}

	dao.DB.Create(&u)
	return nil
}

// 自动添加super用户：一生二、二生三、三生万物
func addSuper() {
	name := setting.SuperUserName
	password := setting.SuperUserPassword

	// 判断是否存在。若存在，根据config更新密码
	u, err := SelectUserByUsername(name)
	if err == nil {
		if u.Password != password {
			u.Password = password
			dao.DB.Save(&u)
		}
		return
	}

	super := User{
		UUID:       util.GenUUID(),
		Username:   name,
		Password:   password,
		Department: "数据部",
		Role:       GradeMap[SUPER],
		Creater:    name,
		CreateOn:   util.CurrentTimeStr(),
	}
	dao.DB.Create(&super)
	return
}

func validateUserDeletion(id int, operator string) (bool, int) {
	code := e.SUCCESS

	userObj, err := SelectUserByID(id)
	if err != nil {
		code = e.ERROR_USER_NOT_EXIST
		return false, code
	}

	operatorGrade, err := GetGradeByName(operator)
	if err != nil {
		code = e.ERROR_USER_INVALID
		return false, code
	}
	if !userObj.LessPermission(operatorGrade) {
		code = e.ERROR_USER_LACK_AUTHORITY
		return false, code
	}

	return true, code
}

func DeleteUserByID(id int, operator string) error {
	if ok, code := validateUserDeletion(id, operator); !ok {
		err := fmt.Errorf(e.GetMsg(code))
		return err
	}

	dao.DB.Delete(&User{}, id)
	return nil
}

func DeleteUserByName(username, operator string) error {
	userObj := User{}
	dao.DB.Where("Username = ?", username).First(&userObj)

	if userObj.ID < 0 {
		err := fmt.Errorf("User to delete does not exist.")
		return err
	}

	return DeleteUserByID(userObj.ID, operator)
}

func ValidateUserUpdate(u User, operator string) (bool, int) {
	code := e.SUCCESS

	// 判断类型是否合法
	if !AllowedRole(u.Role) {
		code = e.ERROR_USER_INVALID
		return false, code
	}

	// 判断是否有足够权限：只能创建权限小于自己的用户
	operatorGrade, err := GetGradeByName(operator)
	if err != nil {
		code = e.ERROR_USER_INVALID
		return false, code
	}
	if !u.LessPermission(operatorGrade) {
		code = e.ERROR_USER_LACK_AUTHORITY
		return false, code
	}

	return true, code
}

func UpdateUser(id int, password, department, role, operator string) error {
	var code int
	u, err := SelectUserByID(id)
	if err != nil {
		code = e.ERROR_USER_NOT_EXIST
		err := fmt.Errorf(e.GetMsg(code))
		return err
	}

	u.Password = password
	u.Department = department
	u.Role = role

	if ok, code := ValidateUserUpdate(u, operator); !ok {
		err := fmt.Errorf(e.GetMsg(code))
		return err
	}

	dao.DB.Save(&u)
	return nil
}
