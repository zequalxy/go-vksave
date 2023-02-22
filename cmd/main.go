package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zequalxy/go-vksave/models"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
)

const (
	ENDMESSAGE = "LastURLs"
	MethodURL  = "https://api.vk.com/method/messages.getHistoryAttachments?v=5.131&access_token="
)

var (
	TOKEN   = ""
	IMGDIR  = "img"
	CHAT_ID = ""
)

func init() {
	flag.StringVar(&TOKEN, "t", TOKEN, "токен")
	flag.StringVar(&CHAT_ID, "c", CHAT_ID, "id чата")
	flag.StringVar(&IMGDIR, "d", IMGDIR, "директория для картинок")
}

func parseToken(token string) {
	t := strings.Split(token, "access_token=")[1]
	TOKEN = strings.Split(t, "&")[0]
}

func parseChatId(chatUrl string) {
	CHAT_ID = strings.Split(chatUrl, "sel=")[1]
	if CHAT_ID[0] == 'c' {
		i, _ := strconv.Atoi(CHAT_ID[1:])
		i += 2000000000
		CHAT_ID = strconv.Itoa(i)
	}
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
	resp, err := http.Get(MethodURL + TOKEN +
		"&media_type=photo" +
		"&peer_id=" + CHAT_ID +
		"&count=200" +
		"&start_from=" + startWith)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var imageResp models.ImageResponse
	if err := json.Unmarshal(body, &imageResp); err != nil {
		panic(err)
	}
	return imageResp
}

func main() {
	flag.Parse()
	parseToken(TOKEN)
	parseChatId(CHAT_ID)
	if err := os.MkdirAll(IMGDIR, 666); err != nil {
		panic(err)
	}
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
			return

		}
	}
}
