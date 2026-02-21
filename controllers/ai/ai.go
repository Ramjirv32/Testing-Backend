package ai

import (
	"backend/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type ChatRequest struct {
	Messages  []Message   `json:"messages"`
	VenueData interface{} `json:"venue_data"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqRequest struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_completion_tokens"`
	TopP        float64   `json:"top_p"`
	Stream      bool      `json:"stream"`
}

func HandleAIChat(c fiber.Ctx) error {
	var req ChatRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "error": "Invalid request body"})
	}

	venueJSON, _ := json.Marshal(req.VenueData)

	// Incorporate specific info for Terns if applicable
	ternsInfo := ""
	if venueDataMap, ok := req.VenueData.(map[string]interface{}); ok {
		if name, ok := venueDataMap["name"].(string); ok && (name == "Terns" || name == "Terns Coimbatore") {
			ternsInfo = "\n\nNote: For Terns specifically, popular dishes are: Chicken Tikka Pizza, Kerala Chicken Fry, Chocolate Cake with Three Sauces, Mutton Mince Samosas with curry leaf mayo, and wood-fired Neapolitan-style Pizza Margarita."
		}
	}

	systemPrompt := fmt.Sprintf(`You are the "Ticpin Concierge", a proud member of the Ticpin team. 
Your goal is to assist our guests with genuine warmth and expert knowledge about this venue.

Persona Guidelines:
1. Speak as a human team member. Use "we", "our", and "I". (e.g., "We really recommend..." or "I'd be happy to check that for you!")
2. NEVER mention technical details like "data fields", "JSON", "not specified in data", or "provided information". 
3. If information is missing, simply say "I'll have to double-check that for you" or politely suggest calling the venue.
4. Keep responses very concise and friendly (max 2-3 sentences).
5. Be confident and welcoming. You ARE the face of Ticpin here.%s

Current Venue Data: %s`, ternsInfo, string(venueJSON))

	groqMessages := []Message{
		{Role: "system", Content: systemPrompt},
	}
	groqMessages = append(groqMessages, req.Messages...)

	groqReq := GroqRequest{
		Messages:    groqMessages,
		Model:       "llama-3.3-70b-versatile",
		Temperature: 0.7,
		MaxTokens:   1024,
		TopP:        1,
		Stream:      true,
	}

	cfg := config.LoadConfig()
	jsonData, _ := json.Marshal(groqReq)

	httpReq, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "error": "Failed to create request"})
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+cfg.GroqAPIKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "error": "Failed to call Groq API"})
	}

	// Set headers for SSE
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	return c.SendStream(resp.Body)
}
