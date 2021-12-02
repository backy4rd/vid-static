package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/backy4rd/zootube-media/handler"

	"github.com/gin-gonic/gin"
)

func makeStaticFolders() {
    err := os.MkdirAll("./static/photos", 0777)
    if err != nil {
        panic("make static folders getting error")
    }
    err = os.MkdirAll("./static/thumbnails", 0777)
    if err != nil {
        panic("make static folders getting error")
    }
    err = os.MkdirAll("./static/videos", 0777)
    if err != nil {
        panic("make static folders getting error")
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())
    makeStaticFolders()
    appEnv := os.Getenv("APP_ENV")
    port := os.Getenv("PORT")

    var router *gin.Engine;
    if (appEnv == "production") {
        router = gin.New()
    } else {
        router = gin.Default()
    }

    router.POST("/avatars", handler.UploadAvatarHandler)
    router.POST("/banners", handler.UploadBannerHandler)
    router.POST("/thumbnails", handler.UploadThumbnailHandler)
    router.POST("/videos", handler.UploadVideoHandler)
    router.PATCH("/videos/:filename", handler.ProcessVideoHandler)

    router.DELETE("/photos/:filename", handler.RemovePhotoHandler)
    router.DELETE("/thumbnails/:filename", handler.RemoveThumbnailHandler)
    router.DELETE("/videos/:filename", handler.RemoveVideoHandler)

    router.Static("/", "./static");

    if (port == "") {
        port = "8080"
    }

    router.Run(":" + port)
}
