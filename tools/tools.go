package tools

import (
	"ai-agent/db"
	"ai-agent/internal/logs"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"log"
)

func GetTools(ctx context.Context) ([]tool.BaseTool, []*schema.ToolInfo) {
	queryUserOrdersTool, err := utils.InferTool("query_user_orders", "查询用户订单信息：查询张三的订单", QueryOrdersUseUserNameFunc)
	queryUserInfoTool, err := utils.InferTool("query_user_info", "查询用户个人信息：查询张三的个人信息", QueryUserInfoFunc)
	if err != nil {
		panic("get query_user_orders tool failed, err=" + err.Error())
	}

	canUsableTools := []tool.BaseTool{
		queryUserOrdersTool,
		queryUserInfoTool,
	}

	infos := make([]*schema.ToolInfo, 0, len(canUsableTools))
	var info *schema.ToolInfo
	for _, canUsableTool := range canUsableTools {
		info, err = canUsableTool.Info(ctx)
		if err != nil {
			log.Fatalf("get ToolInfo failed, err=%v", err)
		}
		infos = append(infos, info)
	}

	return canUsableTools, infos
}

func GetUsefulTool() tool.InvokableTool {
	return &QueryUserTool{}
}

type QueryUserTool struct{}

func (t *QueryUserTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "query_user_info",
		Desc: "Query user info",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"name": {
				Type:     "string",
				Desc:     "name of user",
				Required: true,
			},
		}),
	}, nil
}

type QueryUserInfoRequest struct {
	Name string `json:"name"`
}

func (t *QueryUserTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 解析参数
	p := &QueryUserInfoRequest{}
	err := json.Unmarshal([]byte(argumentsInJSON), p)
	if err != nil {
		return "", err
	}

	// 请求后端服务
	users, err := db.QueryUserByName(p.Name)
	if err != nil {
		logs.Fatalf("查询异常: " + err.Error())
	}
	var userinfo []string
	for _, user := range users {
		item := fmt.Sprintf("%-5d %-20s %-25s %d\n", user.ID, user.Name, user.Email, user.Age)
		userinfo = append(userinfo, item)
	}
	// 序列化结果
	res, err := json.Marshal(userinfo)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
