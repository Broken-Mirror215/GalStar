package api

import (
	"Gal-Finder/internal/response"
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type VNDBapi struct{}

func NewVNDBApi() *VNDBapi {
	return &VNDBapi{}
}

func (a *VNDBapi) Search(c *gin.Context) {
	keyword := c.Query("q") //gin提供的查询字符串参数
	if keyword == "" {
		response.Fail(c, 401, 401, "not get the keyword q")
		return
	}
	page := c.DefaultQuery("page", "1")

	//这个body的设计是根据vndb介绍的哪个接口去设计的？
	body := map[string]interface{}{
		"filter":   []interface{}{"search", "=", keyword},
		"fields":   "title,alttitle,released,rating,image{url,thumbnail,sexual,violence}",
		"sort":     "searchrank",
		"rescults": 20,
		"page":     page,
	}

	//这个是什么？
	raw, _ := json.Marshal(body)

	client := http.Client{
		Timeout: time.Second * 3,
	}

	//这个函数的第三个参数又在干什么

	req, err := http.NewRequest("POST", "https://api.vndb.org/kana/vn", bytes.NewReader(raw))
	if err != nil {
		response.Fail(c, 500, 500, "create vndb request error")
		return
	}

	//这个又是什么？？
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		response.Fail(c, 500, 500, "vndb request error")
		return
	}
	defer resp.Body.Close()

	var rescult map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rescult); err != nil {
		response.Fail(c, 500, 500, "vndb response error")
		return
	}
	response.Success(c, rescult)
}
