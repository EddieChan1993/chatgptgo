package openai

import (
	"bytes"
	"encoding/base64"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/sashabaranov/go-openai"
	"image/png"
	"os"
)

/**
图片内容生成
*/
func CreateImgUrl(content string) {
	// Sample image by link
	reqUrl := gogpt.ImageRequest{
		Prompt:         content,
		Size:           openai.CreateImageSize512x512,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		N:              1,
	}

	respUrl, err := openAiIns.client.CreateImage(openAiIns.ctx, reqUrl)
	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		return
	}
	fmt.Println(respUrl.Data[0].URL)
}

func CreateImgPng(content string) {
	reqBase64 := gogpt.ImageRequest{
		Prompt:         content,
		Size:           openai.CreateImageSize512x512,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
	}

	respBase64, err := openAiIns.client.CreateImage(openAiIns.ctx, reqBase64)
	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		return
	}

	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		fmt.Printf("Base64 decode error: %v\n", err)
		return
	}

	r := bytes.NewReader(imgBytes)
	imgData, err := png.Decode(r)
	if err != nil {
		fmt.Printf("PNG decode error: %v\n", err)
		return
	}

	file, err := os.Create("image.png")
	if err != nil {
		fmt.Printf("File creation error: %v\n", err)
		return
	}
	defer file.Close()

	if err := png.Encode(file, imgData); err != nil {
		fmt.Printf("PNG encode error: %v\n", err)
		return
	}

	fmt.Println("The image was saved as example.png")
}
