package main

import (
\t"ai-content-creator-agent/internal/adapters/controllers"
\t"ai-content-creator-agent/internal/adapters/repositories"
\t"ai-content-creator-agent/internal/adapters/services"
\t"ai-content-creator-agent/internal/domain/interfaces"
\t"ai-content-creator-agent/internal/domain/usecases"
\t"ai-content-creator-agent/internal/infrastructure/api"
\t"ai-content-creator-agent/internal/infrastructure/config"

\t"context"
\t"log"

\t"github.com/gin-gonic/gin"
)

func main() {
\t// Load configuration
\tconfiguration, err := config.LoadConfig()
\tif err != nil {
\t\tlog.Fatal(err.Error())
\t}

\t// Initialize services
\tnlpService := services.NewNLPService(configuration.GeminiAPIKey)
\tscheduleService := services.NewScheduleService()

\t// Initialize api services
\tnewsApi := api.NewNewsAPI(configuration.NewsAPIKey)

\t// Initialize repositories
\tcompanyRepo := repositories.NewCompanyRepository(context.Background(), configuration.ProjectID)
\tcontentRepo := repositories.NewContentRepository(context.Background(), configuration.ProjectID)

\t// Initialize usecases
\tresearchUsecase := usecases.NewResearchUsecase(newsApi, nlpService)
\tcompanyUsecase := usecases.NewCompanyUsecase(companyRepo, scheduleService)
\tcontentUsecase := usecases.NewContentUsecase(contentRepo, researchUsecase, nlpService)
\tsocialMediaUsecase := usecases.NewSocialMediaMgtUsecase(contentUsecase, companyUsecase)

\t// NEW: Jonathan's GitHub service
\tgithubService := services.NewGitHubService(configuration)

\t// Initialize Controllers
\tcompanyController := controllers.NewCompanyController(companyUsecase)
\tsocialMediaController := controllers.NewSocialMediaMgtController(socialMediaUsecase)

\t// Initialize router
\trouter := setupRouter(companyController, socialMediaController, githubService)

\t// Start server
\terr = router.Run(":8080")
\tif err != nil {
\t\tlog.Fatal(err)
\t}
}

func setupRouter(companyController interfaces.CompanyController, socialMediaController controllers.SocialMediaMgtController, githubService *services.GitHubService) *gin.Engine {
\trouter := gin.Default()

\t// Existing routes
\trouter.POST("/post-content/:id", socialMediaController.PostOnFacebook)
\trouter.POST("/reply/:id", socialMediaController.ReplyToComments)
\trouter.POST("/company/register", companyController.RegisterCompany)
\trouter.GET("/company/:id", companyController.GetCompany)

\t// NEW: Jonathan's publishing endpoint
\trouter.POST("/jonathan-publish", func(c *gin.Context) {
\t\ttitle := "Test Jonathan Post"
\t\tcontent := "---
title: "Test"
date: "2026-01-31"
---

# Test Article"
\t\t
\t\turl, err := githubService.PublishArticle(title, content)
\t\tif err != nil {
\t\t\tc.JSON(500, gin.H{"error": err.Error()})
\t\t\treturn
\t\t}
\t\tc.JSON(200, gin.H{"published": true, "url": url})
\t})

\treturn router
}