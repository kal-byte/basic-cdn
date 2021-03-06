package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type _utils struct{}

func (u *_utils) getShortenedName(name string, length int) string {
	if len(name) > 14 {
		return name[len(name)-14:]
	}

	return name
}

var utils = &_utils{}

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("There was an error initializing the environment variables. Error: %s", err)
		return
	}

	router := gin.New()

	router.MaxMultipartMemory = 50 << 20 // 50 MiB
	router.Static("/static", "./static")

	router.Use(gin.Logger())
	router.Use(gin.CustomRecovery(func(ctx *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}))

	router.GET("/", indexEndPoint)
	router.POST("/upload", downloadSentFile)

	router.Run(":7070")
}

func indexEndPoint(ctx *gin.Context) {
	redirectURL := os.Getenv("REDIRECT_URL")
	ctx.Redirect(http.StatusPermanentRedirect, redirectURL)
}

func downloadSentFile(ctx *gin.Context) {
	token := ctx.Query("token")

	if token != os.Getenv("TOKEN") {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("You provided an inproper token."))
		return
	}

	file, err := ctx.FormFile("file")

	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	filename := filepath.Base(file.Filename)
	path := ""

	path = utils.getShortenedName(filename, 14)

	if err := ctx.SaveUploadedFile(file, "static/"+path); err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	baseURL := os.Getenv("BASE_URL")
	ctx.JSON(http.StatusOK, gin.H{"path": baseURL + "static/" + path})
}
