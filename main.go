package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
	// Ginのルーターを初期化
	r := gin.Default()

	//URLを受け取る
	r.POST("/process-url", func(c *gin.Context) {
		var inputData struct {
			URL string `form:"url" binding:"required"`
		}

		if err := c.ShouldBind(&inputData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		// inputData.URL にフォームから送信されたURLが格納されているので、ここで適切な処理を行う
		// 受け取ったURLをAIに送信して結果を取得
		aiResponse, err := sendURLtoAI(inputData.URL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending URL to AI"})
			return
		}

		// AIからのレスポンス（0か1）をフロントに返す
		c.JSON(http.StatusOK, aiResponse)
	})

	// サーバーを起動
	r.Run(":8080")
}

// AIにURLを送信して結果を取得する関数
// できててほしい
func sendURLtoAI(url string) (int, error) {
	data := map[string]string{"url": url}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	aiURL := "http://localhost:5000/predict"

	req, err := http.NewRequest("POST", aiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var responseData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return 0, err
	}

	// responseDataから審議判定を取得して整数値に変換する
	// この例では、"result" フィールドが 0 または 1 として返されることを仮定しています
	result, ok := responseData["result"].(float64)
	if !ok {
		return 0, fmt.Errorf("Invalid response format")
	}

	return int(result), nil
}
