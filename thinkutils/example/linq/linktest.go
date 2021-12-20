package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

type User struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
}

type UserCopy struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
}

func main() {
	//log.Info("FXXK")

	szJson := `[{"Name":"name-a","Age":1},{"Name":"name-b","Age":2},{"Name":"name-c","Age":3}]`
	var lstUser []User
	err := thinkutils.JSONUtils.FromJson(szJson, &lstUser)
	if err != nil {
		fmt.Println(err)
	}

	var names []string
	linq.From(lstUser).WhereT(func(user User) bool {
		return user.Age > 1
	}).SelectT(func(user User) string {
		return user.Name
	}).ToSlice(&names)

	fmt.Println(names)

	name := linq.From(lstUser).WhereT(func(user User) bool {
		return user.Age > 0
	}).WhereT(func(user User) bool {
		return user.Age > 1
	}).SelectT(func(user User) string {
		return user.Name
	}).First()

	fmt.Println(name)

	var users []UserCopy
	linq.From(lstUser).WhereT(func(user User) bool {
		return user.Age > 0
	}).SelectT(func(user User) UserCopy {
		u := UserCopy{
			Name: user.Name,
			Age:  user.Age,
		}
		log.Info("%p", &u)
		return u
	}).ToSlice(&users)

	for i := 0; i < len(users); i++ {
		u := users[i]
		log.Info("%p", &u)
	}
	fmt.Println(users)
}
