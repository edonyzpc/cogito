//go:build !generate_models_file && !context_cache
// +build !generate_models_file,!context_cache

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/edonyzpc/cogito/pkg/cogito"
	"github.com/edonyzpc/cogito/pkg/moonshot"
	"github.com/edonyzpc/cogito/pkg/translate"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cogito",
		Short: "A brief description of your application",
		Long:  `A longer description that spans multiple lines.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	translateCmd = &cobra.Command{
		Use:   "translate",
		Short: "translate file",
		Long:  `translate file`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := translate.Translate(translateFilePath); err != nil {
				log.Fatalln(err)
			}
		},
	}
	translateFilePath string
	apiKey            string
)

func init() {
	rootCmd.Flags().StringVar(&apiKey, "api", "", "api key to access the LLMs")

	translateCmd.Flags().StringVarP(&translateFilePath, "filepath", "f", "", "full file path name to be translated by cogito")
	translateCmd.MarkFlagRequired("filepath")

	rootCmd.AddCommand(translateCmd)
}

func RunDemo() error {
	ctx := context.Background()

	client := moonshot.NewClient[*cogito.Cogito](&cogito.Cogito{
		URL:        "https://api.moonshot.cn/v1",
		APIKey:     os.Getenv("MOONSHOT_API_KEY"),
		HTTPClient: http.DefaultClient,
		Logger: func(ctx context.Context, caller string, request *http.Request, response *http.Response, elapse time.Duration) {
			log.Printf("[%s] %s %s", caller, request.URL, elapse)
		},
	})

	estimateTokenCount, err := client.EstimateTokenCount(ctx, &moonshot.EstimateTokenCountRequest{
		Messages: []*moonshot.Message{
			{
				Role:    moonshot.RoleSystem,
				Content: &moonshot.Content{Text: "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。"},
			},
			{
				Role:    moonshot.RoleUser,
				Content: &moonshot.Content{Text: "你好，我叫李雷，1+1等于多少？"},
			},
		},
		Model:       moonshot.ModelMoonshot8K,
		MaxTokens:   4096,
		N:           1,
		Temperature: "0.3",
	})

	if err != nil {
		return err
	}

	log.Printf("total_tokens=%d", estimateTokenCount.Data.TotalTokens)

	balance, err := client.CheckBalance(ctx)
	if err != nil {
		return err
	}

	log.Printf("balance=%s", balance.Data.AvailableBalance)

	completion, err := client.CreateChatCompletion(ctx, &moonshot.ChatCompletionRequest{
		Messages: []*moonshot.Message{
			{
				Role:    moonshot.RoleSystem,
				Content: &moonshot.Content{Text: "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。"},
			},
			{
				Role:    moonshot.RoleUser,
				Content: &moonshot.Content{Text: "你好，我叫李雷，1+1等于多少？"},
			},
		},
		Model:       moonshot.ModelMoonshot8K,
		MaxTokens:   4096,
		N:           1,
		Temperature: "0.3",
	})

	if err != nil {
		return err
	}

	fmt.Println(completion.GetMessageContent())

	stream, err := client.CreateChatCompletionStream(ctx, &moonshot.ChatCompletionStreamRequest{
		Messages: []*moonshot.Message{
			{
				Role:    moonshot.RoleSystem,
				Content: &moonshot.Content{Text: "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。"},
			},
			{
				Role:    moonshot.RoleUser,
				Content: &moonshot.Content{Text: "写一个小故事，讲的是一个叫“龙猫”的勇士积极抵抗魔族入侵，保卫 Kimi 女神。"},
			},
		},
		Model:       moonshot.ModelMoonshot8K,
		MaxTokens:   4096,
		N:           1,
		Temperature: "0.3",
	})

	if err != nil {
		return err
	}

	defer stream.Close()
	for chunk := range stream.C {
		fmt.Printf("%s", chunk.GetDeltaContent())
	}
	fmt.Println("")

	if err = stream.Err(); err != nil {
		return err
	}

	stream, err = client.CreateChatCompletionStream(ctx, &moonshot.ChatCompletionStreamRequest{
		Messages: []*moonshot.Message{
			{
				Role:    moonshot.RoleSystem,
				Content: &moonshot.Content{Text: "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。"},
			},
			{
				Role:    moonshot.RoleUser,
				Content: &moonshot.Content{Text: "写一个小故事，讲的是有一个叫“龙猫”的人，每天会在各个群聊里游荡，挑选一些感兴趣的话题回复，每个群都以得到龙猫老师的回复为荣，请写一个跌宕起伏的剧情，讲述“龙猫”与各个群聊的爱恨情仇。"},
			},
		},
		Model:       moonshot.ModelMoonshot8K,
		MaxTokens:   4096,
		N:           1,
		Temperature: "0.3",
	})

	if err != nil {
		return err
	}

	defer stream.Close()
	message := stream.CollectMessage()
	fmt.Println(message.Content.Text)

	if err = stream.Err(); err != nil {
		return err
	}

	pdf, err := os.Open("file-test.pdf")
	if err != nil {
		return err
	}

	defer pdf.Close()

	file, err := client.UploadFile(ctx, &moonshot.UploadFileRequest{
		File:    pdf,
		Purpose: "file-extract",
	})

	if err != nil {
		return err
	}

	log.Printf("file_id=%q; status=%s", file.ID, file.Status)

	content, err := client.RetrieveFileContent(ctx, file.ID)
	if err != nil {
		return err
	}

	fmt.Println(string(content))

	return nil
}

func main() {
	/*
		if err := RunDemo(); err != nil {
			log.Fatalln(err)
		}
	*/
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
