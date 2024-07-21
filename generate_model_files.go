//go:build generate_models_file
// +build generate_models_file

package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

var (
	baseUrl = flag.String("baseUrl", os.Getenv("MOONSHOT_BASE_URL"), "base_url for moonshot service")
	key     = flag.String("key", os.Getenv("MOONSHOT_API_KEY"), "api_key for moonshot service, which would be included in an `Authorization` HTTP header")
	file    = flag.String("file", "models.go", "output file which contains all moonshot models")
)

func main() {
	if err := generate(); err != nil {
		fmt.Println(err)
		return
	}
}

func generate() error {
	flag.Parse()
	ctx := context.Background()
	if *baseUrl == "" {
		*baseUrl = "https://api.moonshot.cn/v1"
	}
	client := NewClient(caller{})
	models, err := client.ListModels(ctx)
	if err != nil {
		defer CloseErrorResponseBody(err)
		if mshErr := ParseError(err); mshErr != nil {
			return errors.New("moonshot: " + mshErr.Message)
		}
		return err
	}
	var code bytes.Buffer
	fmt.Fprintf(&code, "// Code generated by generate_models_file.go, DO NOT EDIT.\n\n")
	fmt.Fprintf(&code, "package %s\n\n", os.Getenv("GOPACKAGE"))
	fmt.Fprintf(&code, "//go:generate go run -tags=generate_models_file . %s\n\n", strings.Join(os.Args[1:], " "))
	fmt.Fprintf(&code, `
const (
	RoleSystem    = "system"
	RoleAssistant = "assistant"
	RoleUser      = "user"
	RoleTool      = "tool"
	RoleCache     = "cache"
)

const (
	ToolTypeFunction = "function"
)

const (
	FinishReasonToolCalls = "tool_calls"
	FinishReasonStop      = "stop"
)

const (
	ContentPartTypeText     = "text"
	ContentPartTypeImageUrl = "image_url"
)

const (
	ImageUrlDetailLow  = "low"
	ImageUrlDetailHigh = "high"
	ImageUrlDetailAuto = "auto"
)

const (
	ResponseFormatJSONObject = "json_object"
	ResponseFormatText = "text"
)

`)
	fmt.Fprintf(&code, "const (\n")
	sort.SliceStable(models.Data, func(i, j int) bool {
		return models.Data[i].ID < models.Data[j].ID
	})
	for _, model := range models.Data {
		fmt.Fprintf(&code, "Model%s = \"%s\"\n", toCamel(replace(model.ID)), model.ID)
	}
	fmt.Fprintf(&code, ")")
	source, err := format.Source(code.Bytes())
	if err != nil {
		return err
	}
	if err = os.WriteFile(*file, source, 0644); err != nil {
		return err
	}
	fmt.Printf("Successfully generated the '%s' file", *file)
	return nil
}

type caller struct{}

func (caller) BaseUrl() string { return strings.TrimSuffix(*baseUrl, "/") }
func (caller) Key() string     { return *key }

func (caller) Log(
	_ context.Context,
	_ string,
	request *http.Request,
	response *http.Response,
	elapse time.Duration,
) {
	var status string
	if response != nil {
		status = response.Status
	}
	fmt.Printf("http: %s %s %s %s", request.Method, request.URL, status, elapse)
}

var replacements = map[string]string{
	"v1": "",
}

func replace(name string) string {
	for before, after := range replacements {
		name = strings.ReplaceAll(name, before, after)
	}
	return name
}

func toCamel(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	n := strings.Builder{}
	n.Grow(len(s))
	capNext := true
	for _, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if capNext {
			if vIsLow {
				v += 'A'
				v -= 'a'
			}
		}
		if vIsCap || vIsLow {
			n.WriteByte(v)
			capNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}
