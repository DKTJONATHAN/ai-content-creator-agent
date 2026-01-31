package config

import (
    "encoding/json"
    "fmt"
    "github.com/joho/godotenv"
    "log"
    "os"
)

type Config struct {
\t// Existing fields
\tProjectID           string `json:"project_id"`
\tGeminiAPIKey        string `json:"gemini_api_key"`
\tNewsAPIKey          string `json:"news_api_key"`
\tServiceAccountKeyPath string `json:"service_account_key_path"`
\tDB_NAME             string `json:"db_name"`
\t
\t// NEW: Jonathan's super agent config
\tGitHubToken         string `json:"github_token"`
\tSuperAgentConfig    SuperAgentConfig `json:"jonathan_config"`
}

type SuperAgentConfig struct {
\tMission            string                 `json:"mission_statement"`
\tPillars            map[string]interface{} `json:"content_pillars"`
\tStyleRules         map[string]interface{} `json:"CRITICAL_STYLE_RULES"`
\tVoiceSignatures    []string               `json:"VOICE_SIGNATURES"`
\tResearchRules      map[string]interface{} `json:"RESEARCH_PROTOCOL"`
\tFrontmatterTemplate map[string]interface{} `json:"EXACT_ARTICLE_FORMAT"`
\tGitConfig          map[string]interface{} `json:"GIT_WORKFLOW"`
\tPublishing         map[string]interface{} `json:"publishing"`
}

func LoadConfig() (*Config, error) {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    config := &Config{
\t\tProjectID:           os.Getenv("PROJECT_ID"),
\t\tGeminiAPIKey:        os.Getenv("GEMINI_API_KEY"),
\t\tNewsAPIKey:          os.Getenv("NEWS_API_KEY"),
\t\tServiceAccountKeyPath: os.Getenv("SERVICE_ACCOUNT_KEY_PATH"),
\t\tDB_NAME:             os.Getenv("DB_NAME"),
\t\tGitHubToken:         os.Getenv("GITHUB_TOKEN"),
    }

    // NEW: Load your super JSON
    superPath := os.Getenv("SUPER_CONFIG_PATH")
    if superPath == "" {
        superPath = "./super-agent-config-v3.json"
    }
    
    superData, err := os.ReadFile(superPath)
    if err != nil {
        log.Printf("Warning: Super config not found at %s. Download it first.", superPath)
    } else {
        // Parse your super JSON
        var superConfig SuperAgentConfig
        if err := json.Unmarshal(superData, &superConfig); err != nil {
            log.Printf("Failed to parse super config: %v", err)
        } else {
            config.SuperAgentConfig = superConfig
            log.Println("✅ Jonathan's super config loaded!")
        }
    }

    return config, nil
}