package main

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

// OllamaClient represents an Ollama API client
type OllamaClient struct {
        baseURL string
        model   string
        client  *http.Client
}

// GenerateRequest represents the request to Ollama
type GenerateRequest struct {
        Model  string `json:"model"`
        Prompt string `json:"prompt"`
        Stream bool   `json:"stream"`
}

// GenerateResponse represents the response from Ollama
type GenerateResponse struct {
        Model         string `json:"model"`
        Response      string `json:"response"`
        Done          bool   `json:"done"`
        Context       []int  `json:"context,omitempty"`
        TotalDuration int64  `json:"total_duration,omitempty"`
        Error         string `json:"error,omitempty"`
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(baseURL, model string) *OllamaClient {
        if baseURL == "" {
                baseURL = "http://localhost:11434"
        }
        if model == "" {
                model = "llama2"
        }

        return &OllamaClient{
                baseURL: baseURL,
                model:   model,
                client: &http.Client{
                        Timeout: 300 * time.Second,
                },
        }
}

// Generate sends a prompt to Ollama and returns the response
func (c *OllamaClient) Generate(prompt string) (string, error) {
        return c.GenerateWithCancel(prompt, nil)
}

// GenerateWithCancel sends a prompt to Ollama with cancellation support
func (c *OllamaClient) GenerateWithCancel(prompt string, cancel <-chan bool) (string, error) {
        if prompt == "" {
                return "", fmt.Errorf("empty prompt")
        }

        reqBody := GenerateRequest{
                Model:  c.model,
                Prompt: prompt,
                Stream: false,
        }

        jsonData, err := json.Marshal(reqBody)
        if err != nil {
                return "", fmt.Errorf("failed to marshal request: %w", err)
        }

        url := fmt.Sprintf("%s/api/generate", c.baseURL)

        // Create context with cancellation
        ctx, ctxCancel := context.WithCancel(context.Background())
        defer ctxCancel()

        // Monitor cancel channel in goroutine
        if cancel != nil {
                go func() {
                        select {
                        case <-cancel:
                                ctxCancel()
                        case <-ctx.Done():
                        }
                }()
        }

        req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
        if err != nil {
                return "", fmt.Errorf("failed to create request: %w", err)
        }
        req.Header.Set("Content-Type", "application/json")

        resp, err := c.client.Do(req)
        if err != nil {
                if ctx.Err() == context.Canceled {
                        return "", fmt.Errorf("cancelled")
                }
                // Check if it's a connection error
                if strings.Contains(err.Error(), "connection refused") {
                        return "", fmt.Errorf("cannot connect to Ollama (is it running?)")
                }
                return "", fmt.Errorf("request failed: %w", err)
        }
        defer resp.Body.Close()

        // Check status code
        if resp.StatusCode != http.StatusOK {
                body, _ := io.ReadAll(resp.Body)
                bodyStr := string(body)
                
                // Try to parse error from JSON
                var errResp GenerateResponse
                if json.Unmarshal(body, &errResp) == nil && errResp.Error != "" {
                        return "", fmt.Errorf("ollama error: %s", errResp.Error)
                }
                
                if len(bodyStr) > 200 {
                        bodyStr = bodyStr[:200] + "..."
                }
                return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, bodyStr)
        }

        body, err := io.ReadAll(resp.Body)
        if err != nil {
                if ctx.Err() == context.Canceled {
                        return "", fmt.Errorf("cancelled")
                }
                return "", fmt.Errorf("failed to read response: %w", err)
        }

        if len(body) == 0 {
                return "", fmt.Errorf("empty response from Ollama")
        }

        var genResp GenerateResponse
        if err := json.Unmarshal(body, &genResp); err != nil {
                return "", fmt.Errorf("failed to parse response: %w", err)
        }

        // Check for error in response
        if genResp.Error != "" {
                return "", fmt.Errorf("ollama error: %s", genResp.Error)
        }

        // Clean up the response - just trim whitespace, don't filter content
        response := strings.TrimSpace(genResp.Response)
        
        if response == "" {
                return "", fmt.Errorf("received empty response from model")
        }

        return response, nil
}

// IsAvailable checks if Ollama is running
func (c *OllamaClient) IsAvailable() bool {
        url := fmt.Sprintf("%s/api/tags", c.baseURL)
        client := &http.Client{Timeout: 2 * time.Second}
        
        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        defer cancel()
        
        req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
        if err != nil {
                return false
        }
        
        resp, err := client.Do(req)
        if err != nil {
                return false
        }
        defer resp.Body.Close()
        
        return resp.StatusCode == http.StatusOK
}

// CheckModel verifies if a model is available
func (c *OllamaClient) CheckModel() error {
        url := fmt.Sprintf("%s/api/tags", c.baseURL)
        
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        
        req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
        if err != nil {
                return fmt.Errorf("failed to create request: %w", err)
        }
        
        resp, err := c.client.Do(req)
        if err != nil {
                return fmt.Errorf("cannot connect to Ollama: %w", err)
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != http.StatusOK {
                return fmt.Errorf("ollama returned status %d", resp.StatusCode)
        }
        
        body, err := io.ReadAll(resp.Body)
        if err != nil {
                return fmt.Errorf("failed to read response: %w", err)
        }
        
        var result struct {
                Models []struct {
                        Name string `json:"name"`
                } `json:"models"`
        }
        
        if err := json.Unmarshal(body, &result); err != nil {
                return fmt.Errorf("failed to parse response: %w", err)
        }
        
        // Check if requested model exists
        modelFound := false
        for _, model := range result.Models {
                if strings.HasPrefix(model.Name, c.model) {
                        modelFound = true
                        break
                }
        }
        
        if !modelFound {
                availableModels := make([]string, 0, len(result.Models))
                for _, model := range result.Models {
                        availableModels = append(availableModels, model.Name)
                }
                return fmt.Errorf("model '%s' not found. Available: %s", 
                        c.model, strings.Join(availableModels, ", "))
        }
        
        return nil
}
