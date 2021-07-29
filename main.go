package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/backy4rd/zootube-static-server/handler"

	"github.com/gin-gonic/gin"
)

func main() {
    rand.Seed(time.Now().UnixNano())
    router := gin.New()

    router.POST("/avatars", handler.UploadAvatarHandler)
    router.POST("/banners", handler.UploadBannerHandler)
    router.POST("/thumbnails", handler.UploadThumbnailHandler)
    router.POST("/videos", handler.UploadVideoHandler)

    router.DELETE("/photos/:filename", handler.RemovePhotoHandler)
    router.DELETE("/thumbnails/:filename", handler.RemoveThumbnailHandler)
    router.DELETE("/videos/:filename", handler.RemoveVideoHandler)

    router.Static("/", "./static");

    port := os.Getenv("PORT")
    if (port == "") {
        port = "8080"
    }

    router.Run(":" + port)
}
