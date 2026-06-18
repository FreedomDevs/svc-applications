package final_reviewer

import (
	"context"
	_ "embed"
	"encoding/json"
	"svc-applications/first_parser"

	"google.golang.org/genai"
)

//go:embed prompt.txt
var promptTemplate string

//go:embed schema.json
var schema []byte

func SendRequest(classification first_parser.UserProfile, nickname string, age string, about string, invited_by string, ctx context.Context, client *genai.Client) (*genai.GenerateContentResponse, error) {
	var schemaJson any
	_ = json.Unmarshal(schema, &schemaJson)

	config := &genai.GenerateContentConfig{
		ResponseMIMEType:   "application/json",
		SystemInstruction:  genai.NewContentFromText(promptTemplate, genai.RoleUser),
		ResponseJsonSchema: schemaJson,
		ThinkingConfig: &genai.ThinkingConfig{
			IncludeThoughts: true,
		},
	}
	prompt, _ := json.Marshal(map[string]any{
		"name":                nickname,
		"age":                 age,
		"about":               about,
		"invited_by":          invited_by,
		"auto_classification": classification,
	})

	return client.Models.GenerateContent(ctx, "gemini-3.1-flash-lite-preview", genai.Text(string(prompt)), config)
}
