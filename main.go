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
			agent.NewReactAgent(ctx).Invoke(ctx, input)
			//agent.NewAgent(ctx).Invoke(ctx, input)
		}
	}
}

func init() {
	db.InitDbConnect()
	//db.CreateTables()
	//db.GenerateMockData()
}

// 显示帮助信息
func showHelp() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("  help       - Show this help message")
	fmt.Println("  exit       - Exit the system")
}
