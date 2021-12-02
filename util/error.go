package util

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func SendFailMessage(c *gin.Context, message string) {
        c.JSON(http.StatusBadRequest, gin.H{
            "fail": gin.H{
                "message": message,
            },
        })
}
