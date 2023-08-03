package main

import (
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

	//
	r.POST("/process-url", func(c *gin.Context) {
		var inputData struct {
			URL string `form:"url" binding:"required"`
		}

		if err := c.ShouldBind(&inputData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		// inputData.URL にフォームから送信されたURLが格納されているので、ここで適切な処理を行う

		c.JSON(http.StatusOK, gin.H{"message": "URL received successfully"})
	})

	//受け取ったurlをAiに渡す

	//AIから帰ってきた審議判定をフロントに返す

	// サーバーを起動
	r.Run(":8080")
}
