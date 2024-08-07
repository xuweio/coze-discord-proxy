package discord

import (
	"coze-discord-proxy/common"
	"coze-discord-proxy/model"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"strings"
	"time"
)

// processMessage 提取并处理消息内容及其嵌入元素
func processMessageUpdate(m *discordgo.MessageUpdate) model.ReplyResp {
	var embedUrls []string
	for _, embed := range m.Embeds {
		if embed.Image != nil {
			embedUrls = append(embedUrls, embed.Image.URL)
		}
	}

	return model.ReplyResp{
		Content:   m.Content,
		EmbedUrls: embedUrls,
	}
}

func processMessageUpdateForOpenAI(m *discordgo.MessageUpdate) model.OpenAIChatCompletionResponse {

	if len(m.Embeds) != 0 {
		for _, embed := range m.Embeds {
			if embed.Image != nil && !strings.Contains(m.Content, embed.Image.URL) {
				if m.Content != "" {
					m.Content += "\n"
				}
				m.Content += fmt.Sprintf("%s\n![Image](%s)", embed.Image.URL, embed.Image.URL)
			}
		}
	}

	promptTokens := common.CountTokens(m.ReferencedMessage.Content)
	completionTokens := common.CountTokens(m.Content)

	return model.OpenAIChatCompletionResponse{
		ID:      m.ID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "Coze-Model",
		Choices: []model.OpenAIChoice{
			{
				Index: 0,
				Message: model.OpenAIMessage{
					Role:    "assistant",
					Content: m.Content,
				},
			},
		},
		Usage: model.OpenAIUsage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		},
	}
}

func processMessageUpdateForOpenAIImage(m *discordgo.MessageUpdate) model.OpenAIImagesGenerationResponse {
	var response model.OpenAIImagesGenerationResponse

	if common.SliceContains(common.CozeDailyLimitErrorMessages, m.Content) {
		response.Data = append(response.Data, &model.OpenAIImagesGenerationDataResponse{
			RevisedPrompt: m.Content,
		})
		res := model.OpenAIImagesGenerationResponse{
			Created: time.Now().Unix(),
			Data:    response.Data,
		}
		return res
	}

	re := regexp.MustCompile(`\]\((https?://[^\s\)]+)\)`)
	subMatches := re.FindAllStringSubmatch(m.Content, -1)

	if len(subMatches) == 0 {

		if len(m.Embeds) != 0 {
			for _, embed := range m.Embeds {
				if embed.Image != nil && !strings.Contains(m.Content, embed.Image.URL) {
					//	if m.Content != "" {
					//		m.Content += "\n"
					//	}
					response.Data = append(response.Data, &model.OpenAIImagesGenerationDataResponse{
						URL:           embed.Image.URL,
						RevisedPrompt: m.Content,
					})
				}
			}
		}

		if len(m.Attachments) != 0 {
			for _, attachment := range m.Attachments {
				if attachment.ProxyURL != "" && !strings.Contains(m.Content, attachment.ProxyURL) {
					//	if m.Content != "" {
					//		m.Content += "\n"
					//	}
					response.Data = append(response.Data, &model.OpenAIImagesGenerationDataResponse{
						URL:           attachment.ProxyURL,
						RevisedPrompt: m.Content,
					})
				}
			}
		}

		if len(m.Message.Components) > 0 && len(m.Embeds) == 0 && len(m.Attachments) == 0 {
			response.Data = append(response.Data, &model.OpenAIImagesGenerationDataResponse{
				RevisedPrompt: m.Content,
			})
		}
	}

	for _, match := range subMatches {
		response.Data = append(response.Data, &model.OpenAIImagesGenerationDataResponse{
			URL:           match[1],
			RevisedPrompt: m.Content,
		})
	}

	return model.OpenAIImagesGenerationResponse{
		Created: time.Now().Unix(),
		Data:    response.Data,
	}
}

// processMessage 提取并处理消息内容及其嵌入元素
func processMessageCreate(m *discordgo.MessageCreate) model.ReplyResp {
	var embedUrls []string
	for _, embed := range m.Embeds {
		if embed.Image != nil {
			embedUrls = append(embedUrls, embed.Image.URL)
		}
	}

	return model.ReplyResp{
		Content:   m.Content,
		EmbedUrls: embedUrls,
	}
}

func processMessageCreateForOpenAI(m *discordgo.MessageCreate) model.OpenAIChatCompletionResponse {

	if len(m.Embeds) != 0 {
		for _, embed := range m.Embeds {
			if embed.Image != nil && !strings.Contains(m.Content, embed.Image.URL) {
				if m.Content != "" {
					m.Content += "\n"
				}
				m.Content += fmt.Sprintf("%s\n![Image](%s)", embed.Image.URL, embed.Image.URL)
			}
		}
	}

	promptTokens := common.CountTokens(m.ReferencedMessage.Content)
	completionTokens := common.CountTokens(m.Content)

	return model.OpenAIChatCompletionResponse{
		ID:      m.ID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "Coze-Model",
		Choices: []model.OpenAIChoice{
			{
				Index: 0,
				Message: model.OpenAIMessage{
					Role:    "assistant",
					Content: m.Content,
				},
			},
		},
		Usage: model.OpenAIUsage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		},
	}
}

func processMessageCreateForOpenAIImage(m *discordgo.MessageCreate) model.OpenAIImagesGenerationResponse {
	var response model.OpenAIImagesGenerationResponse

	if common.SliceContains(common.CozeDailyLimitErrorMessages, m.Content) {
		response.Data = append(response.Data, &model.OpenAIImagesGenerationDataResponse{
			RevisedPrompt: m.Content,
		})
		res := model.OpenAIImagesGenerationResponse{
			Created: time.Now().Unix(),
			Data:    response.Data,
		}
		return res
	}

	re := regexp.MustCompile(`\]\((https?://[^\s\)]+)\)`)
	subMatches := re.FindAllStringSubmatch(m.Content, -1)

	if len(subMatches) == 0 {

		if len(m.Embeds) != 0 {
			for _, embed := range m.Embeds {
				if embed.Image != nil && !strings.Contains(m.Content, embed.Image.URL) {
					//	if m.Content != "" {
					//		m.Content += "\n"
					//	}
					response.Data = append(response.Data, &model.OpenAIImagesGenerationDataResponse{
						URL:           embed.Image.URL,
						RevisedPrompt: m.Content,
					})
				}
			}
		}

		if len(m.Attachments) != 0 {
			for _, attachment := range m.Attachments {
				if attachment.ProxyURL != "" && !strings.Contains(m.Content, attachment.ProxyURL) {
					//	if m.Content != "" {
					//		m.Content += "\n"
					//	}
					response.Data = append(response.Data, &model.OpenAIImagesGenerationDataResponse{
						URL:           attachment.ProxyURL,
						RevisedPrompt: m.Content,
					})
				}
			}
		}

		if len(m.Message.Components) > 0 && len(m.Embeds) == 0 && len(m.Attachments) == 0 {
			response.Data = append(response.Data, &model.OpenAIImagesGenerationDataResponse{
				RevisedPrompt: m.Content,
			})
		}
	}

	for _, match := range subMatches {
		response.Data = append(response.Data, &model.OpenAIImagesGenerationDataResponse{
			URL:           match[1],
			RevisedPrompt: m.Content,
		})
	}

	return model.OpenAIImagesGenerationResponse{
		Created: time.Now().Unix(),
		Data:    response.Data,
	}
}
