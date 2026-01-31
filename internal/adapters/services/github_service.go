package services

import (
    "context"
    "encoding/base64"
    "fmt"
    "log"

    github "github.com/google/go-github/v62/github"
    "golang.org/x/oauth2"
    "ai-content-creator-agent/internal/infrastructure/config"
)

type GitHubService struct {
    client *github.Client
    repoOwner string
    repoName  string
}

func NewGitHubService(cfg *config.Config) *GitHubService {
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: cfg.GitHubToken},
    )
    tc := oauth2.NewClient(ctx, ts)
    
    return &GitHubService{
        client:    github.NewClient(tc),
        repoOwner: "DKTJONATHAN",
        repoName:  "TheReport",
    }
}

func (g *GitHubService) PublishArticle(title, content string) (string, error) {
    // Generate slug from title
    slug := generateSlug(title)
    filePath := fmt.Sprintf("src/content/posts/%s.md", slug)
    
    // Build commit message per your spec
    commitMsg := fmt.Sprintf("content: publish article %s", title)
    
    // Encode content for GitHub API
    fileContent := base64.StdEncoding.EncodeToString([]byte(content))
    
    // Create file directly on main
    fileOpts := &github.CreateFileOptions{
        Message: &commitMsg,
        Content: []byte(fileContent),
    }
    
    _, _, err := g.client.Repositories.CreateFile(
        context.Background(),
        g.repoOwner, g.repoName,
        filePath,
        fileOpts,
    )
    
    if err != nil {
        return "", fmt.Errorf("github publish failed: %v", err)
    }
    
    log.Printf("✅ Published %s → %s", title, filePath)
    return fmt.Sprintf("https://jonathanmwaniki.co.ke/posts/%s", slug), nil
}

func generateSlug(title string) string {
    slug := strings.ToLower(title)
    slug = strings.ReplaceAll(slug, " ", "-")
    slug = strings.ReplaceAll(slug, ".", "")
    slug = strings.ReplaceAll(slug, ",", "")
    slug = strings.ReplaceAll(slug, "?", "")
    if len(slug) > 60 {
        slug = slug[:60]
    }
    return slug
}