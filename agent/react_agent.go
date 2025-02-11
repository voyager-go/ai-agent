package agent

import (
	"ai-agent/tools"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"io"
	"log"
)

type ReactAgent struct {
	*react.Agent
}

func NewReactAgent(ctx context.Context) *ReactAgent {
	persona := `你是企业管理系统的智能助手，可以执行已加载的方法来操作管理系统。`
	model := createOpenAIChatModel(ctx)

	canUsableTools, infos := tools.GetTools(ctx)
	err := model.BindTools(infos)

	if err != nil {
		log.Fatalln(err)
	}

	reagent, err := react.NewAgent(ctx, &react.AgentConfig{
		Model: model,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: canUsableTools,
		},
		MessageModifier: react.NewPersonaModifier(persona),
	})

	if err != nil {
		log.Fatalln(err)
	}

	return &ReactAgent{reagent}

	//sr, err := reagent.Stream(ctx, []*schema.Message{
	//	{
	//		Role:    schema.User,
	//		Content: "从系统中查询张三的爸爸的信息",
	//	},
	//}, agent.WithComposeOptions(compose.WithCallbacks(&LoggerCallback{})))
	//if err != nil {
	//	logs.Errorf("failed to stream: %v", err)
	//	return
	//}
	//
	//defer sr.Close() // remember to close the stream
	//
	//logs.Infof("\n\n===== start streaming =====\n\n")
	//
	//for {
	//	msg, err := sr.Recv()
	//	if err != nil {
	//		if errors.Is(err, io.EOF) {
	//			// finish
	//			break
	//		}
	//		// error
	//		logs.Infof("failed to recv: %v", err)
	//		return
	//	}
	//
	//	// 打字机打印
	//	logs.Tokenf("%v", msg.Content)
	//}
	//
	//logs.Infof("\n\n===== finished =====\n")
	//select {}
}

func (r *ReactAgent) Invoke(ctx context.Context, input string) {
	var outMessage *schema.Message
	outMessage, err := r.Agent.Generate(ctx, []*schema.Message{
		schema.UserMessage(input),
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(outMessage.Content)
}

type LoggerCallback struct {
	callbacks.HandlerBuilder // 可以用 callbacks.HandlerBuilder 来辅助实现 callback
}

func (cb *LoggerCallback) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	fmt.Println("==================")
	inputStr, _ := json.MarshalIndent(input, "", "  ") // nolint: byted_s_returned_err_check
	fmt.Printf("[OnStart] %s\n", string(inputStr))
	return ctx
}

func (cb *LoggerCallback) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	fmt.Println("=========[OnEnd]=========")
	outputStr, _ := json.MarshalIndent(output, "", "  ") // nolint: byted_s_returned_err_check
	fmt.Println(string(outputStr))
	return ctx
}

func (cb *LoggerCallback) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	fmt.Println("=========[OnError]=========")
	fmt.Println(err)
	return ctx
}

func (cb *LoggerCallback) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo,
	output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {

	var graphInfoName = "PregelGraph"

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("[OnEndStream] panic err:", err)
			}
		}()

		defer output.Close() // remember to close the stream in defer

		fmt.Println("=========[OnEndStream]=========")
		for {
			frame, err := output.Recv()
			if errors.Is(err, io.EOF) {
				// finish
				break
			}
			if err != nil {
				fmt.Printf("internal error: %s\n", err)
				return
			}

			s, err := json.Marshal(frame)
			if err != nil {
				fmt.Printf("internal error: %s\n", err)
				return
			}

			if info.Name == graphInfoName { // 仅打印 graph 的输出, 否则每个 stream 节点的输出都会打印一遍
				fmt.Printf("%s: %s\n", info.Name, string(s))
			}
		}

	}()
	return ctx
}

func (cb *LoggerCallback) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo,
	input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	defer input.Close()
	return ctx
}
