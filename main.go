package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Config struct {
	ServerName  string `json:"servername"`
	IP          string `json:"ip"`
	Port        string `json:"port"`
	Introduced  string `json:"introduced"`
	QQGrouplink string `json:"qqgrouplink"`
	Email       string `json:"email"`
	ServerPort  int    `json:"serverport"`
}

type MotdBEJson struct {
	Status     string `json:"status"`
	Host       string `json:"host"`
	Motd       string `json:"motd"`
	Agreement  int    `json:"agreement"`
	Version    string `json:"version"`
	Online     int    `json:"online"`
	Max        int    `json:"max"`
	Level_name string `json:"level_name"`
	Gamemode   string `json:"gamemode"`
	Delay      int    `json:"delay"`
}

//读取配置文件
func getConfig() Config {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	json.Unmarshal([]byte(byteValue), &config)
	return config
}

//请求MotdPE API
//https://wiki.blackbe.xyz/api/motd.html
func getMotdBE(ip string, port string) MotdBEJson {
	url := "https://motdbe.blackbe.xyz/api?host=" + ip + ":" + port
	client := http.Client{Timeout: 10 * time.Second} //设置10秒超时
	res, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
	}
	var config MotdBEJson
	json.Unmarshal([]byte(body), &config)
	return config
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Static("/assets", "./assets") //定义静态资源目录
	r.LoadHTMLGlob("assets/template/**/*")
	ServerPort := fmt.Sprintf("%d", getConfig().ServerPort)
	fmt.Println("网站已运行在 " + ServerPort + " 端口")

	r.GET("/", func(c *gin.Context) {
		Config := getConfig()
		ServerInfo := getMotdBE(Config.IP, Config.Port)

		var Status string
		var Status_bool bool
		if ServerInfo.Status == "online" {
			Status = "在线"
			Status_bool = true
		} else {
			Status = "离线"
			Status_bool = false
		}

		c.HTML(http.StatusOK, "index/index.html", gin.H{
			"servername":  Config.ServerName,
			"ip":          Config.IP,
			"port":        Config.Port,
			"introduced":  Config.Introduced,
			"qqgrouplink": Config.QQGrouplink,
			"email":       Config.Email,
			"status":      Status,
			"status_bool": Status_bool,
			"online":      ServerInfo.Online,
			"max":         ServerInfo.Max,
			"delay":       ServerInfo.Delay,
			"version":     ServerInfo.Version})
	})

	r.Run(":" + ServerPort)
}
