package logger

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func LogFormatter(param gin.LogFormatterParams) string {
	return fmt.Sprintf("%s | %d |\t %s | %s | %s\t\"%s\"\n",
		param.TimeStamp.Format("2006/01/02 15:04:05"),
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Request.Method,
		param.Request.URL.Path,
	)
}
