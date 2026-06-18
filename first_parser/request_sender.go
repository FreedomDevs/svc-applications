package first_parser

import (
	"context"
	_ "embed"
	"encoding/json"
	"regexp"

	//	"log"
	"strings"

	"google.golang.org/genai"
)

//go:embed prompt.txt
var promptTemplate string

func marshalText(input string) string {
	b, _ := json.Marshal(input)
	s := string(b)
	return s
}

func SendRequest(nickname string, age string, about string, invited_by string, ctx context.Context, client *genai.Client) (*genai.GenerateContentResponse, error) {
	prompt := strings.NewReplacer(
		"%nickname%", marshalText(nickname),
		"%age%", marshalText(age),
		"%about%", marshalText(about),
		"%invited_by%", marshalText(invited_by),
	).Replace(promptTemplate)

	contents := []*genai.Content{genai.NewContentFromText(prompt, genai.RoleUser)}

	return client.Models.GenerateContent(ctx, "gemma-3-27b-it", contents, nil)
}

type UserProfile struct {
	Friends          []string       `json:"friends"`
	MentionedServers []string       `json:"mentioned_servers"`
	Categories       map[string]any `json:"categories"`
}

func ParseSomeDataData(response string) UserProfile {
	lines := strings.Split(response, "\n")
	profile := UserProfile{
		Categories: make(map[string]any),
	}

	// Регулярки для парсинга ключей и значений
	reRoot := regexp.MustCompile(`^(\w+):\s*(.*)$`)
	reCategory := regexp.MustCompile(`^\*\s+([\w\s]+):\s*(.*)$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "categories:" {
			continue
		}

		if strings.HasPrefix(line, "*") {
			// Парсим вложенные категории
			matches := reCategory.FindStringSubmatch(line)
			if len(matches) == 3 {
				key := strings.ToLower(strings.ReplaceAll(matches[1], " ", "_"))
				val := matches[2]

				// Если в значении есть запятая, превращаем в слайс (например, Playstyle)
				if strings.Contains(val, ",") {
					parts := strings.Split(val, ",")
					for i := range parts {
						parts[i] = strings.TrimSpace(parts[i])
					}
					profile.Categories[key] = parts
				} else {
					profile.Categories[key] = val
				}
			}
		} else {
			// Парсим корневые поля (friends, mentioned_servers)
			matches := reRoot.FindStringSubmatch(line)
			if len(matches) == 3 {
				key := matches[1]
				val := matches[2]

				var list []string
				if val != "None" && val != "" {
					parts := strings.Split(val, ",")
					for _, p := range parts {
						list = append(list, strings.TrimSpace(p))
					}
				}

				if key == "friends" {
					profile.Friends = list
				} else if key == "mentioned_servers" {
					profile.MentionedServers = list
				}
			}
		}
	}

	return profile
}
