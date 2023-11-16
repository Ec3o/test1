package main

import (
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var html = `
<!DOCTYPE html>
<html>
<head>
    <title>File Upload</title>
</head>
<body>
    <h1>Upload a File</h1>
    <form method="POST" action="/unzip" enctype="multipart/form-data">
        <input type="file" name="file" required>
        <br>
        <input type="submit" value="Upload">
    </form>
</body>
</html>
`

func init() {
	// disable the output of gin
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if err := os.MkdirAll("uploads", 0777); err != nil {
		panic(err)
	} //创建一个叫'upload'的路径
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Check the attachment for next step.")
	})
	r.GET("/unzip", func(c *gin.Context) {
		// 渲染HTML表单
		c.Data(200, "text/html", []byte(html))
	})
	r.POST("/unzip", func(c *gin.Context) {
		fh, err := c.FormFile("file") //获取一个名叫file变量的值
		if err != nil {
			log.Println(err)
			c.String(500, "读取文件失败")
		}
		savepath := filepath.Clean(filepath.Join("uploads", strings.ReplaceAll(fh.Filename, "../", "")))

		file, err := fh.Open()
		if err != nil {
			log.Println(err)
			c.String(500, "读取文件失败")
		}
		defer file.Close()

		saveto, err := os.OpenFile(savepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			log.Println(err)
			c.String(500, "保存文件失败")
		}
		defer saveto.Close()

		_, err = io.Copy(saveto, file)
		if err != nil {
			print(savepath)
			log.Println(err)
			c.String(500, "保存文件失败")
		}

		res, err := Unzip(savepath)
		if err != nil {
			log.Println(err)
			print(savepath)
			c.String(500, "解压失败")
		}

		// show unzip log
		c.String(200, string(res))
	})

	r.Run(":8180")
}

func Unzip(zipfile string) ([]byte, error) {
	unzipdir := filepath.Join("output", strconv.Itoa(rand.Int()))
	print(unzipdir)
	output, err := exec.Command("unzip", "-n", "-v", zipfile, "-d", unzipdir).CombinedOutput()
	if err != nil {
		return nil, err
	}

	return output, nil
}
