package agent

import (
	"ai-agent/tools"
	"context"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/fatih/color"
	"log"
	"strings"
)

type Agent struct {
	runnable compose.Runnable[[]*schema.Message, []*schema.Message]
}

func NewAgent(ctx context.Context) (*Agent, error) {
	model := createOpenAIChatModel(ctx)
	canUsableTools, infos := tools.GetTools(ctx)
	err := model.BindTools(infos)

	if err != nil {
		log.Fatalf("绑定Tools失败, 错误信息为:%v", err)
		return nil, err
	}
	toolsNode, err := compose.NewToolNode(context.Background(), &compose.ToolsNodeConfig{
		Tools: canUsableTools,
	})

	if err != nil {
		log.Fatalf("NewToolNode failed, err=%v", err)
		return nil, err
	}

	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	chain.
		AppendChatModel(model, compose.WithNodeName("chat_model")).
		AppendToolsNode(toolsNode, compose.WithNodeName("tools"))

	// 编译并运行 chain
	agent, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("chain.Compile failed, err=%v", err)
		return nil, err
	}
	return &Agent{
		runnable: agent,
	}, nil
}

func (a *Agent) Invoke(ctx context.Context, input string) (string, error) {
	_, err := a.runnable.Invoke(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: input,
		},
	})
	if err != nil {
		log.Println()
		if strings.Contains(err.Error(), "no tool call found in input") {
			color.Red("抱歉，我还不会处理这种业务：%s", input)
		}
		return "", err
	}

	// 输出结果
	//for idx, msg := range resp {
	//	log.Println()
	//	log.Fatalf("message %d: %s: %s", idx, msg.Role, msg.Content)
	//}
	return "", nil
}
