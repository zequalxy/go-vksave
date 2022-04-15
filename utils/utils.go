package utils

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
)

func Download(url string) {
	fileName := "img/" + generateName() + ".jpg"
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

func generateName() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
