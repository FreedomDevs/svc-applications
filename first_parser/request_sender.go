package first_parser

import (
	"context"
	_ "embed"
	//	"log"
	"strings"

	"google.golang.org/genai"
)

//go:embed prompt.txt
var promptTemplate string

func SendRequest(nickname string, age string, about string, invited_by string, ctx context.Context, client *genai.Client) (*genai.GenerateContentResponse, error) {
	prompt := strings.NewReplacer(
		"%nickname%", nickname,
		"%age%", age,
		"%about%", about,
		"%invited_by%", invited_by,
	).Replace(promptTemplate)

	//log.Println(prompt)

	contents := []*genai.Content{genai.NewContentFromText(prompt, genai.RoleUser)}

	return client.Models.GenerateContent(ctx, "gemma-3-27b-it", contents, nil)
}
