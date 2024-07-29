package translate

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/edonyzpc/cogito/pkg/cogito"
	"github.com/edonyzpc/cogito/pkg/moonshot"
)

func readMarkdownFile(filePath string) ([]byte, error) {
	fileBuf, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	fileLength := len(string(fileBuf))

	if fileLength > 4096 {
		fmt.Println("file is larger than 4096, ", fileLength)
	}

	return fileBuf, nil
}

func writeMarkdownFile(filePath string, fileBuf []byte) error {
	err := os.WriteFile(filePath, fileBuf, 0644)
	if err != nil {
		return err
	}

	return nil
}

func renameEnMarkdownFile(filePath string) string {
	pathName := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)
	extName := filepath.Ext(fileName)
	newFileName := fileName[:len(fileName)-len(extName)] + "_en" + extName

	return path.Join(pathName, newFileName)
}

func Translate(filePath string) error {
	ctx := context.Background()

	client := moonshot.NewClient[*cogito.Cogito](&cogito.Cogito{
		URL:        "https://api.moonshot.cn/v1",
		APIKey:     os.Getenv("MOONSHOT_API_KEY"),
		HTTPClient: http.DefaultClient,
		Logger: func(ctx context.Context, caller string, request *http.Request, response *http.Response, elapse time.Duration) {
			log.Printf("[%s] %s %s", caller, request.URL, elapse)
		},
	})

	content, err := readMarkdownFile(filePath)
	if err != nil {
		return err
	}
	completion, err := client.CreateChatCompletion(ctx, &moonshot.ChatCompletionRequest{
		Messages: []*moonshot.Message{
			{
				Role:    moonshot.RoleSystem,
				Content: &moonshot.Content{Text: "你是一个精通将中文翻译成英文的专家，你能够将下面给出的markdown翻译成英文，要求：1.markdown语法相关内容不用翻译；2.http链接等固定字符串不用翻译；3.翻译后不要包含任何中文标点符号；4.翻译的英文要求准确地道符合英语母语使用习惯"},
			},
			{
				Role:    moonshot.RoleUser,
				Content: &moonshot.Content{Text: "```markdown\n" + string(content) + "\n```"},
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
	enMarkdownFile := renameEnMarkdownFile(filePath)
	err = writeMarkdownFile(enMarkdownFile, []byte(completion.GetMessageContent()))
	if err != nil {
		return err
	}
	return nil
}
