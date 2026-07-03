package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OpenAIProvider 兼容 OpenAI 格式的 LLM 提供商
// 支持 OpenAI、DeepSeek、Claude (兼容模式)、Ollama 等
type OpenAIProvider struct {
	config        LLMConfig
	client        *http.Client
	promptBuilder *PromptBuilder
}

// OpenAI API 请求/响应结构
type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// NewOpenAIProvider 创建 OpenAI 兼容的 LLM 提供商
func NewOpenAIProvider(config LLMConfig) *OpenAIProvider {
	return &OpenAIProvider{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		promptBuilder: NewPromptBuilder(),
	}
}

// callLLM 调用 LLM API
func (p *OpenAIProvider) callLLM(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	reqBody := chatRequest{
		Model: p.config.Model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: p.config.Temperature,
		MaxTokens:   p.config.MaxTokens,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	// 构建请求 URL
	url := strings.TrimRight(p.config.BaseURL, "/") + "/chat/completions"

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if p.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w (body: %s)", err, string(body))
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("API 错误: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("API 返回空结果")
	}

	return strings.TrimSpace(chatResp.Choices[0].Message.Content), nil
}

// callLLMJSON 调用 LLM 并解析 JSON 响应
func (p *OpenAIProvider) callLLMJSON(ctx context.Context, systemPrompt, userPrompt string, target interface{}) error {
	systemPrompt += "\n\n请严格以 JSON 格式返回结果，不要包含任何其他文本或 markdown 标记。"

	result, err := p.callLLM(ctx, systemPrompt, userPrompt)
	if err != nil {
		return err
	}

	// 清理可能的 markdown 代码块标记
	result = cleanJSONResponse(result)

	if err := json.Unmarshal([]byte(result), target); err != nil {
		return fmt.Errorf("解析 JSON 失败: %w (raw: %s)", err, result)
	}

	return nil
}

// cleanJSONResponse 清理 LLM 返回的 JSON（去除 markdown 标记等）
func cleanJSONResponse(s string) string {
	s = strings.TrimSpace(s)
	// 去除 ```json ... ``` 标记
	if strings.HasPrefix(s, "```") {
		lines := strings.Split(s, "\n")
		if len(lines) > 2 {
			lines = lines[1 : len(lines)-1]
		}
		s = strings.Join(lines, "\n")
	}
	s = strings.TrimSpace(s)
	return s
}

// GenerateNarrative 生成场景描述
func (p *OpenAIProvider) GenerateNarrative(ctx context.Context, req NarrativeRequest) (string, error) {
	systemPrompt := "你是一个地下城探险游戏的叙事者。请用中文生成简短生动的场景描述（2-3句话，不超过100字）。风格：暗黑奇幻，略带幽默，像在写日志条目。"

	vars := map[string]string{
		"RoomType":   req.RoomType,
		"Depth":      fmt.Sprintf("%d", req.Depth),
		"RoomIndex":  fmt.Sprintf("%d", req.RoomIndex),
		"Context":    req.Context,
		"PlayerInfo": req.PlayerInfo,
	}
	if req.Extra != nil {
		for k, v := range req.Extra {
			vars[k] = v
		}
	}

	userPrompt, _ := p.promptBuilder.Build("narrative", vars)
	if userPrompt == "" {
		userPrompt = fmt.Sprintf("为一个地牢第%d层的%s房间生成描述。", req.Depth, req.RoomType)
	}

	return p.callLLM(ctx, systemPrompt, userPrompt)
}

// GenerateMonster 生成怪物
func (p *OpenAIProvider) GenerateMonster(ctx context.Context, req MonsterRequest) (*MonsterResponse, error) {
	systemPrompt := "你是一个地下城探险游戏的怪物设计师。请用中文设计一个有趣的怪物。"

	vars := map[string]string{
		"Depth":      fmt.Sprintf("%d", req.Depth),
		"RoomType":   req.RoomType,
		"IsBoss":     fmt.Sprintf("%v", req.IsBoss),
		"PlayerInfo": req.PlayerInfo,
	}

	userPrompt, _ := p.promptBuilder.Build("monster", vars)
	if userPrompt == "" {
		userPrompt = fmt.Sprintf("设计一个地牢第%d层的怪物，Boss=%v。用JSON格式返回name/description/traits/dialogue/strategy字段。", req.Depth, req.IsBoss)
	}

	var monster MonsterResponse
	if err := p.callLLMJSON(ctx, systemPrompt, userPrompt, &monster); err != nil {
		return nil, err
	}

	return &monster, nil
}

// GenerateEvent 生成随机事件
func (p *OpenAIProvider) GenerateEvent(ctx context.Context, req EventRequest) (*EventResponse, error) {
	systemPrompt := "你是一个地下城探险游戏的事件设计师。请用中文设计一个有趣的随机事件。"

	vars := map[string]string{
		"Depth":      fmt.Sprintf("%d", req.Depth),
		"PlayerInfo": req.PlayerInfo,
		"Previous":   req.Previous,
	}

	userPrompt, _ := p.promptBuilder.Build("event", vars)
	if userPrompt == "" {
		userPrompt = fmt.Sprintf("设计一个地牢第%d层的随机事件。用JSON格式返回title/description/choices字段，choices中每个选项包含text/description/outcome(good/bad/neutral)。", req.Depth)
	}

	var event EventResponse
	if err := p.callLLMJSON(ctx, systemPrompt, userPrompt, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

// GenerateDialogue 生成 NPC 对话
func (p *OpenAIProvider) GenerateDialogue(ctx context.Context, req DialogueRequest) (string, error) {
	systemPrompt := "你是一个地下城探险游戏的 NPC 扮演者。请用中文生成简短的对话（1-2句话），符合角色性格。直接输出对话内容。"

	vars := map[string]string{
		"NPCType":    req.NPCType,
		"Context":    req.Context,
		"PlayerInfo": req.PlayerInfo,
		"Topic":      req.Topic,
	}

	userPrompt, _ := p.promptBuilder.Build("dialogue", vars)
	if userPrompt == "" {
		userPrompt = fmt.Sprintf("一个%s类型的NPC，主题是%s。", req.NPCType, req.Topic)
	}

	return p.callLLM(ctx, systemPrompt, userPrompt)
}

// GenerateCombatAction 生成怪物战斗行为描述
func (p *OpenAIProvider) GenerateCombatAction(ctx context.Context, req CombatActionRequest) (string, error) {
	systemPrompt := "你是一个地下城战斗的旁白者。请用中文描述怪物的下一步行动（1-2句话，像战斗日志）。直接输出行动描述。"

	vars := map[string]string{
		"MonsterName":  req.MonsterName,
		"Traits":       strings.Join(req.Traits, "、"),
		"Strategy":     req.Strategy,
		"MonsterHP":    fmt.Sprintf("%d", req.MonsterHP),
		"MonsterMaxHP": fmt.Sprintf("%d", req.MonsterMaxHP),
		"PlayerHP":     fmt.Sprintf("%d", req.PlayerHP),
		"PlayerMaxHP":  fmt.Sprintf("%d", req.PlayerMaxHP),
		"Round":        fmt.Sprintf("%d", req.Round),
		"LastAction":   req.LastAction,
	}

	userPrompt, _ := p.promptBuilder.Build("combat_action", vars)
	if userPrompt == "" {
		userPrompt = fmt.Sprintf("怪物%s（%s）正在战斗，HP %d/%d，第%d回合。",
			req.MonsterName, strings.Join(req.Traits, "、"), req.MonsterHP, req.MonsterMaxHP, req.Round)
	}

	return p.callLLM(ctx, systemPrompt, userPrompt)
}
