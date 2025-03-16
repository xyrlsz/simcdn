package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"simcdn/config"
	"simcdn/logger"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	nodes     = []string{} // 其他节点地址
	nodeIndex = 0
	nodeMutex = &sync.Mutex{}
)

func runSever(conf *config.Config) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.LoggerWithFormatter(logger.LogFormatter))
	nodes = conf.NodeHosts

	// 处理动态资源请求
	r.GET(fmt.Sprintf("/%s/*filename", conf.RelativePath), func(c *gin.Context) {
		filename := c.Param("filename")[1:] // 去掉前导 "/"

		// 检查本地文件
		for _, staticDir := range conf.LocalAssetsPaths {
			filePath := filepath.Join(staticDir, filename)
			if _, err := os.Stat(filePath); err == nil {
				c.File(filePath)
				return
			}

		}

		// 若 nodes 为空，则返回 404
		if len(nodes) == 0 {
			c.String(http.StatusNotFound, "File not found")
			return
		}

		// 从其他节点获取资源
		nodeMutex.Lock()
		nodeURL := nodes[nodeIndex]
		nodeIndex = (nodeIndex + 1) % len(nodes) // 轮询选择节点
		nodeMutex.Unlock()

		target, err := url.Parse(nodeURL)
		if err != nil {
			c.String(http.StatusInternalServerError, "Invalid node URL")
			return
		}

		// 反向代理到其他 CDN 节点
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// 启动服务器

	log.Default().Printf("CDN server has started on  %s:%d.\n", conf.ListenOn, conf.Port)
	if err := r.Run(fmt.Sprintf("%s:%d", conf.ListenOn, conf.Port)); err != nil {
		log.Default().Fatal(fmt.Sprint("Server failed to start:", err))
		os.Exit(1)
	}

}

func main() {
	conf := config.GetConfig()
	log.Default().Println("Press Ctrl+C to exit.")

	if conf != nil {
		go func() {
			runSever(conf)
		}()
		select {}
	} else {
		for {
		}
	}

}
