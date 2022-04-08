package main

import (
	"encoding/json"
	"fmt"
	"go-vksave/models"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

const (
	TOKEN     = "3e9649476875cd14154ec4ad0891e41589ae282163dd2e3032958588ad51c72c16f8f113908abec966143"
	MethodURL = "https://api.vk.com/method/"
	API_ID    = "2685278"
	AUTHORIZE = "https://oauth.vk.com/authorize?" +
		"client_id=" + API_ID +
		"&display=popup" +
		"&redirect_uri=https://oauth.vk.com/blank.html" +
		"&scope=messages,offline" +
		"&response_type=token" +
		"&v=5.131" +
		"&state=123456"
)

func download(url string) {
	fileName := "img/" + url[strings.LastIndex(url, "/")+1:strings.LastIndex(url, "/")+16]
	output, err := os.Create(fileName)
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()
	io.Copy(output, response.Body)
}

func main() {
	fmt.Println("hello world")
	ir := getImages("")
	sortPhotoBySizes(&ir)
	for _, item := range ir.Response.Items {
		download(item.Attachment.Photo.Sizes[0].URL)
	}
}

func sortPhotoBySizes(ir *models.ImageResponse) {
	for _, item := range ir.Response.Items {
		sort.Sort(sort.Reverse(models.ByHeight(item.Attachment.Photo.Sizes)))
	}

}

func getImages(startWith string) models.ImageResponse {
	resp, err := http.Get(MethodURL + "messages.getHistoryAttachments?v=5.131&access_token=" + TOKEN +
		"&media_type=photo" +
		"&peer_id=2000000199" +
		"&count=200")
	body, err := ioutil.ReadAll(resp.Body)
	//err := exec.Command("rundll32", "url.dll,FileProtocolHandler", AUTHORIZE).Start()
	if err != nil {
		panic(err)
	}
	var imageResp models.ImageResponse
	if err := json.Unmarshal(body, &imageResp); err != nil {
		panic(err)
	}
	return imageResp
}
