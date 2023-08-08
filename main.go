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
	Label       string `json:"label"`
}

func main() {
	r := gin.Default()

	r.POST("/process-url", func(c *gin.Context) {
		var inputData struct {
			URL string `form: "url" binding: "required"`
		}

		if err := c.ShouldBind(&inputData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		//whoisの取得
		whoisResult, err := whois.GetWhois(inputData.URL)
		if err != nil {
			fmt.Println(err)
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error in model prediction"})
			return
		}

		// Pythonスクリプトを呼び出して予測結果を取得
		cmd := exec.Command(".venv/bin/python3", "src/predict.py")
		cmd.Stderr = os.Stderr
		out, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(out))
		cmd.Stdin = strings.NewReader(string(modelInputJSON))

		// 標準出力を取得
		output, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error in model prediction"})
			return
		}

		//Pythonスクリプトからの結果をパースしてフロントに返す
		var modelOutput struct {
			Label string `json:"label"`
		}

		if err := json.Unmarshal(output, &modelOutput); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error in model prediction"})
			return
		}

		//フロントに返すデータを整形
		urlData := URLData{
			URL:         inputData.URL,
			WhoisResult: whoisResult,
			Label:       modelOutput.Label,
		}

		c.JSON(http.StatusOK, urlData)

	})

	r.Run(":8080")
}

//できていてほしい
