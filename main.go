package main

import (
    "strings"
    "os"
    "github.com/gin-gonic/gin"
)

func postPhotoHandler(c *gin.Context) {
    file, _ := c.FormFile("file")

    if (file == nil) {
        c.Writer.WriteHeader(400)
        return
    }
    if (!strings.HasPrefix(file.Header.Get("Content-Type"), "image")) {
        c.Writer.WriteHeader(400)
        return
    }

    // Upload the file to specific dst.
    c.SaveUploadedFile(file, "./static/photos/" + file.Filename)

    c.Writer.WriteHeader(200)
}

func postThumbnailHandler(c *gin.Context) {
    file, _ := c.FormFile("file")

    if (file == nil) {
        c.Writer.WriteHeader(400)
        return
    }
    if (!strings.HasPrefix(file.Header.Get("Content-Type"), "image")) {
        c.Writer.WriteHeader(400)
        return
    }

    // Upload the file to specific dst.
    c.SaveUploadedFile(file, "./static/thumbnails/" + file.Filename)

    c.Writer.WriteHeader(200)
}

func postVideoHandler(c *gin.Context) {
    file, _ := c.FormFile("file")

    if (file == nil) {
        c.Writer.WriteHeader(400)
        return
    }
    if (!strings.HasPrefix(file.Header.Get("Content-Type"), "video")) {
        c.Writer.WriteHeader(400)
        return
    }

    // Upload the file to specific dst.
    c.SaveUploadedFile(file, "./static/videos/" + file.Filename)

    c.Writer.WriteHeader(200)
}

func removePhotoHandler(c *gin.Context) {
    filename := c.Param("filename")
    err := os.Remove("./static/photos/" + filename)

    if (err == nil) {
        c.Writer.WriteHeader(200)
    } else {
        c.Writer.WriteHeader(400)
    }
}

func removeVideoHandler(c *gin.Context) {
    filename := c.Param("filename")
    err := os.Remove("./static/videos/" + filename)

    if (err == nil) {
        c.Writer.WriteHeader(200)
    } else {
        c.Writer.WriteHeader(400)
    }
}

func removeThumbnailHandler(c *gin.Context) {
    filename := c.Param("filename")
    err := os.Remove("./static/thumbnails/" + filename)

    if (err == nil) {
        c.Writer.WriteHeader(200)
    } else {
        c.Writer.WriteHeader(400)
    }
}

func main() {
    router := gin.Default()

    router.POST("/photos", postPhotoHandler)
    router.POST("/thumbnails", postThumbnailHandler)
    router.POST("/videos", postVideoHandler)
    router.DELETE("/photos/:filename", removePhotoHandler)
    router.DELETE("/thumbnails/:filename", removeThumbnailHandler)
    router.DELETE("/videos/:filename", removeVideoHandler)

    router.Static("/", "./static");

    port := os.Getenv("PORT")
    if (port == "") {
        port = "8080"
    }

    router.Run(":" + port)
}
