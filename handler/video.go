package handler

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/backy4rd/zootube-static-server/util"

	"github.com/gin-gonic/gin"
)

func RemoveVideoHandler(c *gin.Context) {
    filename := c.Param("filename")
    err := os.Remove("./static/videos/" + filename)

    if (err == nil) {
        c.Writer.WriteHeader(200)
    } else {
        c.Writer.WriteHeader(400)
    }
}

func UploadVideoHandler(c *gin.Context) {
    _video, _ := c.FormFile("video")
    _seek := c.PostForm("seek")

    seek, err := strconv.Atoi(_seek)

    if _video == nil || err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "message": "invalid or missing fields",
            },
        })
        return
    }
    if (!strings.HasPrefix(_video.Header.Get("Content-Type"), "video")) {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "message": "video is not an video",
            },
        })
        return
    }
    videoExtension, err := util.GetFileExtension(_video.Filename)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "message": "video name is not valid",
            },
        })
        return
    }

    videoFilename := util.GenerateRandomString(32) + "." + videoExtension
    videoPath := "./static/videos/" + videoFilename
    thumbnailFilename := util.GenerateRandomString(32) + ".jpg"
    thumbnailPath := "./static/thumbnails/" + thumbnailFilename

    c.SaveUploadedFile(_video, videoPath)
    util.ExtractFrame(videoPath, seek, thumbnailHeight, thumbnailPath)

    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "videoPath": "/videos/" + videoFilename,
            "thumbnailPath": "/thumbnails/" + thumbnailFilename,
        },
    })
}
