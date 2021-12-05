package handler

import (
    "net/http"
    "os"
    "time"

    "github.com/backy4rd/zootube-media/util"

    "github.com/disintegration/imaging"
    "github.com/gin-gonic/gin"
)

const iconHeight = 64
var acceptedPhotoMimetypes = []string { "image/jpeg", "image/png" }
var acceptedPhotoExtensions = []string { "jpg", "png" }

func RemovePhotoHandler(c *gin.Context) {
    filename := c.Param("filename")
    os.Remove("./static/photos/" + filename)

    c.Writer.WriteHeader(204)
}

func UploadPhotoHandler(c *gin.Context) {
    _photo, _ := c.FormFile("photo")

    if _photo == nil {
        util.SendFailMessage(c, "missing photo")
        return
    }
    if !util.IsStringInArray(acceptedPhotoMimetypes, _photo.Header.Get("Content-Type")) {
        util.SendFailMessage(c, "photo type is not accepted");
        return
    }
    photoExtension, err := util.GetFileExtension(_photo.Filename)
    if err != nil {
        util.SendFailMessage(c, "photo name is not valid");
        return
    }
    if !util.IsStringInArray(acceptedPhotoExtensions, photoExtension) {
        util.SendFailMessage(c, "photo type is not accepted");
        return
    }

    photoFilename := util.GenerateRandomString(32) + "." + photoExtension
    photoPath := "./temp/photos/" + photoFilename
    c.SaveUploadedFile(_photo, photoPath)
    time.AfterFunc(time.Hour, func() {
        os.Remove(photoPath)
    })

    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "photo": photoFilename,
        },
    })
}

func ProcessBannerHandler(c *gin.Context) {
    _bannerFilename := c.Param("filename")

    _bannerPath := "./temp/photos/" + _bannerFilename
    bannerPath := "./static/photos/" + _bannerFilename
    if !util.IsFileExist(_bannerPath) {
        util.SendFailMessage(c, "banner not found")
        return

    }

    util.MoveFile(_bannerPath, bannerPath)

    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "bannerPath": "/photos/" + _bannerFilename,
        },
    })
}

func ProcessAvatarHandler(c *gin.Context) {
    _avatarFilename := c.Param("filename")

    avatarPath := "./static/photos/" + _avatarFilename
    _avatarPath := "./temp/photos/" + _avatarFilename
    if !util.IsFileExist(_avatarPath) {
        util.SendFailMessage(c, "avatar not found")
        return
    }

    iconFilename := util.GenerateRandomString(32) + ".jpg"
    iconPath := "./static/photos/" + iconFilename
    icon, err := imaging.Open(_avatarPath)
    if err != nil {
        util.SendFailMessage(c, "error occur while opening avatar")
        return
    }
    icon = imaging.Fill(icon, iconHeight, iconHeight, imaging.Center, imaging.Lanczos)
    err = imaging.Save(icon, iconPath)
    if err != nil {
        util.SendFailMessage(c, "error occur while cropping avatar")
        return
    }
    util.MoveFile(_avatarPath, avatarPath)

    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "avatarPath": "/photos/" + _avatarFilename,
            "iconPath": "/photos/" + iconFilename,
        },
    })
}

