package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	whois "github.com/undiabler/golang-whois"
)

type URLData struct {
	URL         string `json:"url"`
	Domain      string `json:"domain"`
	Scheme      string `json:"scheme"`
	TDL         string `json:"tld"`
	WhoisResult string `json:"whois_result"`
	Label       int    `json:"label"`
}

func main() {
	r := gin.Default()

	r.POST("/", func(c *gin.Context) {
		var inputData struct {
			URL string `form:"url" binding:"required"`
		}

		if err := c.ShouldBind(&inputData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		domain := inputData.URL

		// whoisの取得
		whoisResult, err := whois.GetWhois(domain)
		if err != nil {
			fmt.Println(err) // エラーメッセージを表示して問題を特定
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error in whois lookup"})
			return
		}

		// モデルに送信するデータを整形
		modelInput := map[string]interface{}{
			"url":          inputData.URL,
			"whois_result": whoisResult,
		}

		// JSONに変換
		modelInputJSON, err := json.Marshal(modelInput)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error in model prediction1"})
			return
		}

		// Pythonスクリプトを呼び出して予測結果を取得
		cmd := exec.Command("python3", "predict.py")
		cmd.Stderr = os.Stderr
		cmd.Stdin = strings.NewReader(string(modelInputJSON))

		// 実行
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error executing Python script:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error in model prediction2"})
			return
		}

		trimmedOutput := strings.TrimSpace(string(output))

		// Pythonスクリプトからの結果を数値として読み取る
		//ここのError in model prediction3が出る
		var modelOutput struct {
			Label int `json:"label"`
		}

		if err := json.Unmarshal([]byte(trimmedOutput), &modelOutput); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error in model prediction3"})
			return
		}

		// フロントに返すデータを整形
		urlData := URLData{
			URL:         inputData.URL,
			WhoisResult: whoisResult,
			Label:       modelOutput.Label,
		}

		c.JSON(http.StatusOK, urlData)

	})

	r.Run(":8080")
}
