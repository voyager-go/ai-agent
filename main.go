package main

import (
	"ai-agent/agent"
	"ai-agent/db"
	"bufio"
	"context"
	"fmt"
	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)

func main() {
	color.Green("欢迎进入指令式智能体交互系统!")
	color.Green("尝试输入 'help' 查询更多有用的信息或者输入 'exit' 退出系统.")

	reader := bufio.NewReader(os.Stdin)

	ctx := context.Background()
	aiAgent, err := agent.NewAgent(ctx)
	if err != nil {
		panic(err)
	}
	for {
		fmt.Print("\n> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input. Please try again.")
			continue
		}

		input = strings.TrimSpace(input)

		switch strings.ToLower(input) {
		case "help":
			showHelp()
		case "exit":
			color.Red("下次再见哦 (#^.^#)")
			return
		default:
			color.Yellow("处理中，请稍候~")
			resp, _ := aiAgent.Invoke(ctx, input)
			//if err != nil {
			//	panic(err)
			//}

			fmt.Println(resp)
		}
	}
}

func init() {
	db.InitDbConnect()
}

// 显示帮助信息
func showHelp() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("  help       - Show this help message")
	fmt.Println("  exit       - Exit the system")
}
