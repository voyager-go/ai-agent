package tools

import (
	"context"
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
