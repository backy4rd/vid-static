package handler

import (
	"net/http"
	"os"
	"strings"

	"github.com/backy4rd/zootube-media/util"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

const iconHeight = 64

func RemovePhotoHandler(c *gin.Context) {
    filename := c.Param("filename")
    err := os.Remove("./static/photos/" + filename)

    if (err == nil) {
        c.Writer.WriteHeader(200)
    } else {
        c.Writer.WriteHeader(400)
    }
}

func UploadAvatarHandler(c *gin.Context) {
    _avatar, _ := c.FormFile("avatar")

    if _avatar == nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "message": "missing avatar",
            },
        })
        return
    }
    if !strings.HasPrefix(_avatar.Header.Get("Content-Type"), "image") {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "message": "avatar is not an image",
            },
        })
        return
    }
    avatarExtension, err := util.GetFileExtension(_avatar.Filename)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "message": "avatar name is not valid",
            },
        })
        return
    }

    avatarFilename := util.GenerateRandomString(32) + "." + avatarExtension
    iconFilename := util.GenerateRandomString(32) + ".jpg"
    avatarPath := "./static/photos/" + avatarFilename
    iconPath := "./static/photos/" + iconFilename

    c.SaveUploadedFile(_avatar, avatarPath)

    icon, err := imaging.Open(avatarPath)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "message": "error occur while opening avatar",
            },
        })
        return
    }
    icon = imaging.Fill(icon, iconHeight, iconHeight, imaging.Center, imaging.Lanczos)
    err = imaging.Save(icon, iconPath)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "message": "error occur while cropping avatar",
            },
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "avatarPath": "/photos/" + avatarFilename,
            "iconPath": "/photos/" + iconFilename,
        },
    })
}

func UploadBannerHandler(c *gin.Context) {
    _banner, _ := c.FormFile("banner")

    if (_banner == nil) {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "message": "missing banner",
            },
        })
        return
    }
    if (!strings.HasPrefix(_banner.Header.Get("Content-Type"), "image")) {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "message": "banner is not an image",
            },
        })
        return
    }
    bannerExtension, err := util.GetFileExtension(_banner.Filename)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "message": "banner name is not valid",
            },
        })
        return
    }

    bannerFilename := util.GenerateRandomString(32) + "." + bannerExtension
    bannerPath := "./static/photos/" + bannerFilename

    c.SaveUploadedFile(_banner, bannerPath)

    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "bannerPath": "/photos/" + bannerFilename,
        },
    })
}
