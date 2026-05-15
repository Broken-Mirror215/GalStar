package api

import (
	"Gal-Finder/internal/response"
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type VNDBapi struct{}

func NewVNDBApi() *VNDBapi {
	return &VNDBapi{}
}

type VNDBresponse struct {
	Results []VNDBItem `json:"results"`
	More    bool       `json:"more"`
}

type VNDBItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	//防止不存在的字段用指针
	AltTitle *string    `json:"alttitle"`
	Released *string    `json:"released"`
	Rating   *float64   `json:"rating"`
	Image    *VNDBImage `json:"image"`
}

type VNDBImage struct {
	URL       string  `json:"url"`
	Thumbnail string  `json:"thumbnail"`
	Sexual    float64 `json:"sexual"`
	Violence  float64 `json:"violence"`
}

func (a *VNDBapi) Search(c *gin.Context) {
	keyword := c.Query("q") //gin提供的查询字符串参数
	if keyword == "" {
		response.Fail(c, 400, 400, "not get the keyword q")
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		response.Fail(c, 400, 400, "not get the page")
		return
	}

	//vndb query format
	body := map[string]interface{}{
		"filters": []interface{}{"search", "=", keyword},
		"fields":  "title,alttitle,released,rating,image{url,thumbnail,sexual,violence}",
		"sort":    "searchrank",
		"results": 20,
		"page":    page,
	}

	//这个是把body变成Json类型的
	raw, err := json.Marshal(body)
	if err != nil {
		response.Fail(c, 500, 500, "json marshal error")
		return
	}

	client := http.Client{
		Timeout: time.Second * 3,
	}

	//这个函数的第三个参数又在干什么
	//感觉像是在调用vndb的api去查询
	req, err := http.NewRequest("POST", "https://api.vndb.org/kana/vn", bytes.NewReader(raw))
	if err != nil {
		response.Fail(c, 500, 500, "create vndb request error")
		return
	}

	//这个是发给VNDB的请求头，但是数据得查
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		response.Fail(c, 500, 500, "vndb request error")
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		response.Fail(c, 502, 502, "vndb request error")
		return
	}

	var result VNDBresponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		response.Fail(c, 500, 500, "vndb response error")
		return
	}

	response.Success(c, gin.H{
		"list":    result.Results,
		"page":    page,
		"hasMore": result.More,
	})
}
