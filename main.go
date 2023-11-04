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
    err = os.MkdirAll("./temp/videos", 0777)
    if err != nil {
        panic("make static folders getting error")
    }
    err = os.MkdirAll("./temp/photos", 0777)
    if err != nil {
        panic("make static folders getting error")
    }
}

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())
    makeStaticFolders()
    port := os.Getenv("PORT")

    router := gin.Default()

    router.Use(CORSMiddleware())

    router.POST("/photos", handler.UploadPhotoHandler)
    router.POST("/videos", handler.UploadVideoHandler)

    router.PATCH("/thumbnails/:filename", handler.ProcessThumbnailHandler)
    router.PATCH("/videos/:filename", handler.ProcessVideoHandler)
    router.PATCH("/avatars/:filename", handler.ProcessAvatarHandler)
    router.PATCH("/banners/:filename", handler.ProcessBannerHandler)

    router.DELETE("/photos/:filename", handler.RemovePhotoHandler)
    router.DELETE("/thumbnails/:filename", handler.RemoveThumbnailHandler)
    router.DELETE("/videos/:filename", handler.RemoveVideoHandler)

    router.Static("/", "./static");

    if (port == "") {
        port = "8080"
    }

    handler.CompressQueue.Run()
    router.Run(":" + port)
}
