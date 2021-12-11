package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/backy4rd/zootube-media/util"

	"github.com/gin-gonic/gin"
)

var acceptedVideoExtensions = []string { "mp4", "mkv", "webm" }
var acceptedVideoMimetypes = []string { "video/x-matroska", "video/mp4", "video/webm" }

var CompressQueue = util.NewTaskQueue()
var API_SERVER_ENDPOINT = os.Getenv("API_SERVER_ENDPOINT")

func RemoveVideoHandler(c *gin.Context) {
    filename := c.Param("filename")
    os.Remove("./static/videos/" + filename)

    c.Writer.WriteHeader(204)
}

func UploadVideoHandler(c *gin.Context) {
    _video, _ := c.FormFile("video")

    if _video == nil {
        util.SendFailMessage(c, "missing video");
        return
    }
    if !util.IsStringInArray(acceptedVideoMimetypes, _video.Header.Get("Content-Type")) {
        util.SendFailMessage(c, "video type is not accepted");
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
    videoPath := "./temp/videos/" + videoFilename
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
    _id := c.PostForm("video_id")

    if _video == "" || _id == "" {
        util.SendFailMessage(c, "missing parameters");
        return
    }
    duration, err := util.GetVideoDuration("./temp/videos/" + _video);
    if err != nil {
        util.SendFailMessage(c, "video not found");
        return
    }
    seek, err := strconv.Atoi(_seek)
    if _seek != "" && err != nil {
        seek = duration / 2
    }
    if seek > duration {
        util.SendFailMessage(c, "seek cannot greater than video duration");
        return
    }
    err = util.MoveFile("./temp/videos/" + _video, "./static/videos/" + _video);
    if err != nil {
        util.SendFailMessage(c, "video not found");
        return
    }
    quality, err := util.GetVideoQuality("./static/videos/" + _video);
    if err != nil {
        util.SendFailMessage(c, "video not found");
        return
    }

    thumbnailFilename := util.GenerateRandomString(32) + ".jpg"
    thumbnailPath := "./static/thumbnails/" + thumbnailFilename
    util.ExtractFrame("./static/videos/" + _video, seek, thumbnailHeight, thumbnailPath)

    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "video" + strconv.Itoa(quality) + "Path": "/videos/" + _video,
            "thumbnailPath": "/thumbnails/" + thumbnailFilename,
            "duration": duration,
        },
    })

    go CompressQueue.Push(func() {
        compressedFilename := util.GenerateRandomString(32) + ".mp4"
        compressedPath := "./static/videos/" + compressedFilename
        err := util.Compress360p("./static/videos/" + _video, compressedPath)
        if err != nil {
            fmt.Println(err)
            return
        }

        body := url.Values{}
        body.Add("video360Path", "/videos/" + compressedFilename)
        endpoint := API_SERVER_ENDPOINT + "/videos/" + _id + "/qualities"

        _, err = http.PostForm(endpoint, body);
        if err != nil {
            os.Remove(compressedPath);
        }
    })
}
