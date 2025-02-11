package tools

import (
	"ai-agent/db"
	"context"
	"fmt"
	"strings"
)

type QueryUserOrdersCondition struct {
	UserName string `json:"name" jsonschema:"description=name of user'"`
}

type QueryUserInfoCondition struct {
	UserName string `json:"name" jsonschema:"description=name of user'"`
}

func QueryOrdersUseUserNameFunc(_ context.Context, req *QueryUserOrdersCondition) (string, error) {
	result, err := db.QueryOrdersByUserName(req.UserName)
	if err != nil {
		panic(err)
	}

	for user, orders := range result {
		fmt.Printf("用户ID: %d, 姓名: %s, 邮箱: %s\n", user.ID, user.Name, user.Email)
		if len(orders) == 0 {
			fmt.Println("  没有订单")
		} else {
			fmt.Println("  订单信息:")
			for _, order := range orders {
				fmt.Printf("    订单ID: %d, 金额: %.2f, 下单时间: %s\n",
					order.ID, order.Amount, order.OrderDate.Format("2006-01-02 15:04:05"))
			}
		}
		fmt.Println()
	}
	return `{"msg": "查询用户订单信息成功"}`, nil
}

func QueryUserInfoFunc(_ context.Context, req *QueryUserInfoCondition) (string, error) {
	users, err := db.QueryUserByName(req.UserName)
	if err != nil {
		panic("查询异常: " + err.Error())
	}
	if len(users) == 0 {
		return "未查询到" + req.UserName + "的信息", nil
	}
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%-5s %-20s %-25s %s\n", "工号", "姓名", "邮箱", "年龄")
	fmt.Println(strings.Repeat("-", 80))
	for _, user := range users {
		fmt.Printf("%-5d %-20s %-25s %d\n", user.ID, user.Name, user.Email, user.Age)
	}
	fmt.Println(strings.Repeat("=", 80))
	return `{"msg": "查询用户信息成功"}`, nil
}
