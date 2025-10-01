package util

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SetAuthCookie(c *gin.Context, userID int) {
    maxAge := 60 * 60 * 24 * 365 // 1 year
    c.SetCookie(
        "user_id",
        fmt.Sprintf("%d", userID),
        maxAge,
        "/",
        "",
        true,  // secure
        true,  // httpOnly
    )
}

func RemoveAuthCookie(c *gin.Context) {
    c.SetCookie(
        "user_id",
        "",
        -1,
        "/",
        "",
        true,  // secure
        true,  // httpOnly
    )
}
