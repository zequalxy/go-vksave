package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"go-vksave/models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var tkn, chatId string
var flag = Flag{false}

const (
	ENDMESSAGE = "LastURLs"
	TOKEN      = "3e9649476875cd14154ec4ad0891e41589ae282163dd2e3032958588ad51c72c16f8f113908abec966143"
	MethodURL  = "https://api.vk.com/method/messages.getHistoryAttachments?v=5.131&access_token="
	API_ID     = "2685278"
	AUTHORIZE  = "https://oauth.vk.com/authorize?" +
		"client_id=" + API_ID +
		"&display=popup" +
		"&redirect_uri=https://oauth.vk.com/blank.html" +
		"&scope=messages,offline" +
		"&response_type=token" +
		"&v=5.131" +
		"&state=123456"
)

func parseToken(token string) {
	t := strings.Split(token, "access_token=")[1]
	tkn = strings.Split(t, "&")[0]
}

func parseChatId(chatUrl string) {
	chatId = strings.Split(chatUrl, "sel=")[1]
	if chatId[0] == 'c' {
		i, _ := strconv.Atoi(chatId[1:])
		i += 2000000000
		chatId = strconv.Itoa(i)
	}
}

func startDownload(c *gin.Context) {
	parseToken(c.Request.FormValue("token"))
	parseChatId(c.Request.FormValue("chatId"))
	start()
	c.Redirect(http.StatusFound, "/")
}

func generator(out chan string) {
	start := ""
	for {
		ir := getImages(start)
		if len(ir.Response.Items) == 0 {
			out <- ENDMESSAGE
			return
		}
		sortPhotoBySizes(&ir)
		for _, item := range ir.Response.Items {
			out <- item.Attachment.Photo.Sizes[0].URL
		}
		start = ir.Response.NextFrom
	}
}

func sortPhotoBySizes(ir *models.ImageResponse) {
	for _, item := range ir.Response.Items {
		sort.Sort(sort.Reverse(models.ByHeight(item.Attachment.Photo.Sizes)))
	}

}

func getImages(startWith string) models.ImageResponse {
	resp, err := http.Get(MethodURL + tkn +
		"&media_type=photo" +
		"&peer_id=" + chatId + // 2000000114
		"&count=200" +
		"&start_from=" + startWith)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var imageResp models.ImageResponse
	if err := json.Unmarshal(body, &imageResp); err != nil {
		panic(err)
	}
	return imageResp
}

func auth(c *gin.Context) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", AUTHORIZE).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", AUTHORIZE).Start()
	case "darwin":
		err = exec.Command("open", AUTHORIZE).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
	c.Redirect(http.StatusFound, "/")
}

func start() {
	links := make(chan string)
	quit := make(chan bool)
	b := new(models.Balancer)
	b.Init(links)

	keys := make(chan os.Signal, 1)
	signal.Notify(keys, os.Interrupt)

	go b.Balance(quit)
	go generator(links)

	fmt.Println("Начинаем загрузку изображений")

	for {
		select {
		case <-keys:
			fmt.Println("CTRL-C: Ожидаю завершения активных загрузок")
			quit <- true
		case <-quit:
			fmt.Println("Загрузки завершены!")
			err := zipImg()
			if err != nil {
				return
			}
			flag.Flag = true
			return

		}
	}
}

func zipImg() error {
	source := "img"
	zipfile, err := os.Create("./assets/img.zip")
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

type Flag struct {
	Flag bool
}

func main() {
	fmt.Println("hello, friend")
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./assets")
	//router.StaticFile("/assets/images", "./assets/images")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", flag)
	})
	router.GET("/auth", auth)
	router.POST("/download", startDownload)

	err := router.Run(":8080")
	if err != nil {
		return
	}

}
