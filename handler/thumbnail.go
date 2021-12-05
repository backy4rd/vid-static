package handler

import (
    "net/http"
    "os"

    "github.com/backy4rd/zootube-media/util"
    "github.com/disintegration/imaging"
    "github.com/gin-gonic/gin"
)

const thumbnailHeight = 160

func RemoveThumbnailHandler(c *gin.Context) {
    filename := c.Param("filename")
    os.Remove("./static/thumbnails/" + filename)

    c.Writer.WriteHeader(204)
}

func ProcessThumbnailHandler(c *gin.Context) {
    _thumbnailFileName := c.Param("filename")

    thumbnailPath := "./temp/photos/" + _thumbnailFileName
    resizedThumbnailPath := "./static/thumbnails/" + _thumbnailFileName
    if !util.IsFileExist(thumbnailPath) {
        util.SendFailMessage(c, "thumbnail not found")
        return
    }

    thumbnail, err := imaging.Open(thumbnailPath)
    if err != nil {
        util.SendFailMessage(c, "error occur while opening thumbnail")
        return
    }
    thumbnail = imaging.Resize(thumbnail, 0, thumbnailHeight, imaging.Lanczos)
    err = imaging.Save(thumbnail, resizedThumbnailPath)
    if err != nil {
        util.SendFailMessage(c, "error occur while cropping thumbnail")
        return
    }
    os.Remove(thumbnailPath)

    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "thumbnailPath": "/thumbnails/" + _thumbnailFileName,
        },
    })
}
