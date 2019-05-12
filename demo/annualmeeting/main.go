package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

/*
	导入用户数据抽奖的模型:
		导入所有数据，抽取其中一个人
*/

var userList []string
var mu sync.Mutex

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&LotteryController{})
	return app
}

func main() {
	app := newApp()
	userList = make([]string, 0)
	mu = sync.Mutex{}

	app.Run(iris.Addr(":8080"))
}

type LotteryController struct {
	Ctx iris.Context
}

func (c *LotteryController) Get() string {
	count := len(userList)
	return fmt.Sprintf("当前总共参与抽奖的用户数: %d\n", count)
}

func (c *LotteryController) PostImport() string {
	strUsers := c.Ctx.FormValue("users")
	users := strings.Split(strUsers, ",")
	mu.Lock()
	defer mu.Unlock()
	count1 := len(userList)
	for _, u := range users {
		u = strings.TrimSpace(u)
		if len(u) > 0 {
			userList = append(userList, u)
		}
	}
	count2 := len(userList)
	return fmt.Sprintf("当前总共参与抽奖的用户数: %d, 成功导入的用户数是: %d\n", count2, (count2 - count1))
}

func (c *LotteryController) GetLucky() string {
	mu.Lock()
	defer mu.Unlock()
	count := len(userList)
	if count > 1 {
		seed := time.Now().UnixNano()
		index := rand.New(rand.NewSource(seed)).Int31n(int32(count))
		user := userList[index]
		userList = append(userList[0:index], userList[index+1:]...)
		return fmt.Sprintf("当前中奖用户: %s, 剩余用户数: %d\n", user, count-1)
	} else if count == 1 {
		user := userList[0]
		userList = userList[0:0]
		return fmt.Sprintf("当前中奖用户: %s, 剩余用户数: %d\n", user, count-1)
	} else {
		return fmt.Sprintf("当前参与用户数为0，请导入用户数据\n")
	}
}
