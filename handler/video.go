package handler

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/backy4rd/zootube-media/util"

	"github.com/gin-gonic/gin"
)

var acceptedVideoExtensions = []string { "mp4", "mkv", "webm" }

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

    if _video == nil {
        util.SendFailMessage(c, "missing video");
        return
    }
    if !strings.HasPrefix(_video.Header.Get("Content-Type"), "video") {
        util.SendFailMessage(c, "video is not an video");
        return
    }
    videoExtension, err := util.GetFileExtension(_video.Filename)
    if err != nil {
        util.SendFailMessage(c, "video name is not valid");
        return
    }
    if !util.IsStringInArray(acceptedVideoExtensions, videoExtension) {
        util.SendFailMessage(c, "video type is not accepted");
        return
    }

    videoFilename := util.GenerateRandomString(32) + "." + videoExtension
    videoPath := "./temp/" + videoFilename
    c.SaveUploadedFile(_video, videoPath)
    time.AfterFunc(time.Hour, func() {
        os.Remove(videoPath)
    })

    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "video": videoFilename,
        },
    })
}

func ProcessVideoHandler(c *gin.Context) {
    _video := c.Param("filename")
    _seek := c.PostForm("seek")

    seek, err := strconv.Atoi(_seek)
    if _video == "" || err != nil {
        util.SendFailMessage(c, "missing parameters");
        return
    }
    duration, err := util.GetVideoDuration("./temp/" + _video);
    if err != nil {
        util.SendFailMessage(c, "video not found");
        return
    }
    if seek > duration {
        util.SendFailMessage(c, "seek cannot greater than video duration");
        return
    }
    err = util.MoveFile("./temp/" + _video, "./static/videos/" + _video);
    if err != nil {
        util.SendFailMessage(c, "video not found");
        return
    }

    thumbnailFilename := util.GenerateRandomString(32) + ".jpg"
    thumbnailPath := "./static/thumbnails/" + thumbnailFilename
    util.ExtractFrame("./static/videos/" + _video, seek, thumbnailHeight, thumbnailPath)


    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "videoPath": "/videos/" + _video,
            "thumbnailPath": "/thumbnails/" + thumbnailFilename,
            "duration": duration,
        },
    })
}
