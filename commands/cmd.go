package main

import (
	_ "teamdo/middleware"
	_ "teamdo/models"

	"github.com/kfchen81/beego/orm"
)

func main() {
	orm.Debug = true

	orm.RunCommand()
}
