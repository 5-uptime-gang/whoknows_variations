package util

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SetAuthCookie(c *gin.Context, userID int) {
    c.SetCookie(
        "user_id",
        fmt.Sprintf("%d", userID),
        0, // maxAge in seconds; 0 means session cookie
        "/",
        "",
        false,  // secure
        true,  // httpOnly
    )
}

func RemoveAuthCookie(c *gin.Context) {
    c.SetCookie(
        "user_id",
        "",
        -1, // maxAge in seconds; -1 means delete cookie
        "/",
        "",
        false,  // secure
        true,  // httpOnly
    )
}
