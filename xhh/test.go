package xhh

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"xhhrobot/ai"
)

func RunTest() {
	fmt.Println("=== 测试模式 ===")

	Init()
	if Info.Cookie == "" {
		fmt.Println("未检测到 cookie，请先运行 -mode login")
		return
	}

	fmt.Print("输入帖子 ID：")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimRight(input, "\r\n")

	linkID, _ := strconv.Atoi(input)
	if linkID == 0 {
		fmt.Println("帖子 ID 不合法")
		return
	}

	fmt.Print("输入 @ 你的用户 UID：")
	uidInput, _ := reader.ReadString('\n')
	uidInput = strings.TrimRight(uidInput, "\r\n")
	uid, _ := strconv.Atoi(uidInput)
	if !Check(uid) {
		fmt.Printf("UID被过滤\n")
		return
	}

	fmt.Print("输入 @ 你的消息内容：")
	msg, _ := reader.ReadString('\n')
	msg = strings.TrimRight(msg, "\r\n")

	contents, topics, tags := GetLinkInfo(linkID, 0)
	if len(contents) <= 1 {
		fmt.Println("获取帖子信息失败，请检查 link_id 是否正确")
		return
	}

	reply := ai.GetAiReply(contents, msg, topics, tags)
	if reply == "" {
		fmt.Println("AI 返回为空，请检查 AI 配置")
		return
	}

	fmt.Println("\n=== AI 回复 ===")
	fmt.Println(reply)
}
func NoAiTest() {
	fmt.Println("=== 测试模式 ===")

	Init()
	if Info.Cookie == "" {
		fmt.Println("未检测到 cookie，请先运行 -mode login")
		return
	}

	fmt.Print("输入帖子 ID：")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimRight(input, "\r\n")

	linkID, _ := strconv.Atoi(input)
	if linkID == 0 {
		fmt.Println("帖子 ID 不合法")
		return
	}
	Reply("test", input, "", "", "0")
}
