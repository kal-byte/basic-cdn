package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("There was an error initializing the environment variables. Error: %s", err)
		return
	}

	router := gin.New()

	router.MaxMultipartMemory = 50 << 20 // 8 MiB
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
	ctx.Redirect(http.StatusPermanentRedirect, "https://www.kal-byte.co.uk/")
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

	if len(filename) > 14 {
		path = filename[len(filename)-14:]
	} else {
		path = filename
	}

	if err := ctx.SaveUploadedFile(file, "static/"+path); err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	baseURL := os.Getenv("BASE_URL")
	ctx.JSON(http.StatusOK, gin.H{"path": baseURL + "static/" + path})
}
