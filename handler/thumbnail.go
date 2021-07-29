package handler

import (
	"net/http"
	"os"
	"strings"

	"github.com/backy4rd/zootube-static-server/util"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

const thumbnailHeight = 160

func RemoveThumbnailHandler(c *gin.Context) {
    filename := c.Param("filename")
    err := os.Remove("./static/thumbnails/" + filename)

    if (err == nil) {
        c.Writer.WriteHeader(200)
    } else {
        c.Writer.WriteHeader(400)
    }
}

func UploadThumbnailHandler(c *gin.Context) {
    _thumbnail, _ := c.FormFile("thumbnail")

    if _thumbnail == nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "message": "missing thumbnail",
            },
        })
        return
    }
    if !strings.HasPrefix(_thumbnail.Header.Get("Content-Type"), "image") {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "message": "thumbnail is not an image",
            },
        })
        return
    }

    tempFilename := util.GenerateRandomString(32) + _thumbnail.Filename
    thumbnailFilename := util.GenerateRandomString(32) + ".jpg"
    tempPath := "./temp/" + tempFilename
    thumbnailPath := "./static/thumbnails/" + thumbnailFilename

    c.SaveUploadedFile(_thumbnail, tempPath)
    defer os.Remove("./temp/" + tempFilename);

    thumbnail, err := imaging.Open(tempPath)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "message": "error occur while opening thumbnail",
            },
        })
        return
    }
    thumbnail = imaging.Resize(thumbnail, 0, thumbnailHeight, imaging.Lanczos)
    err = imaging.Save(thumbnail, thumbnailPath)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "message": "error occur while cropping thumbnail",
            },
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "thumbnailPath": "/thumbnails/" + thumbnailFilename,
        },
    })
}
