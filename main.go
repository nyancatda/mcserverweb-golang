package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"encoding/json"
    "net/http"
    "github.com/gin-gonic/gin"
)

type Config struct {
	ServerName string `json:"servername"`
	IP string `json:"ip"`
	Port string `json:"port"`
	Introduced string `json:"introduced"`
	QQGrouplink string `json:"qqgrouplink"`
	Email string `json:"email"`
	ServerPort int `json:"serverport"`
}

type MotdBEJson struct {
    Status string `json:"status"`
    IP string `json:"ip"`
    Port string `json:"port"`
    Motd string `json:"motd"`
	Agreement string `json:"agreement"`
	Version string `json:"version"`
	Online string `json:"online"`
	Max string `json:"max"`
	Gamemode string `json:"gamemode"`
	Delay int `json:"delay"`
}

//读取配置文件
func getConfig()(Config){
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
func getMotdBE(ip string,port string)(MotdBEJson){
	url := "http://motdpe.blackbe.xyz/api.php?ip="+ip+"&port="+port
    res, err := http.Get(url)
    if err != nil {
        fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
        os.Exit(1)
    }
    body, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    if err != nil {
        fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
        os.Exit(1)
    }
    var config MotdBEJson
	json.Unmarshal([]byte(body), &config)
	fmt.Println(config)
    return config
}


func main() {
	gin.SetMode(gin.ReleaseMode)
    r := gin.Default()

	r.Static("/public", "./public")//定义静态资源目录
	r.LoadHTMLGlob("assets/**/*")
	ServerPort := fmt.Sprintf("%d", getConfig().ServerPort)
	fmt.Println("网站已运行在 "+ServerPort+" 端口")

    r.GET("/", func(c *gin.Context) {
		Config := getConfig()
		ServerInfo := getMotdBE(Config.IP,Config.Port)

		var Status string
		var Status_bool bool
		if ServerInfo.Status == "online"{
			Status = "在线"
			Status_bool = true
		} else {
			Status = "离线"
			Status_bool =false
		}

        c.HTML(http.StatusOK, "index/index.html", gin.H{
			"servername": Config.ServerName,
			"ip": Config.IP,
			"port": Config.Port,
			"introduced": Config.Introduced,
			"qqgrouplink": Config.QQGrouplink,
			"email": Config.Email,
			"status": Status,
			"status_bool": Status_bool,
			"online":ServerInfo.Online,
			"max":ServerInfo.Max,
			"delay":ServerInfo.Delay,
			"version":ServerInfo.Version})
    })

    r.Run(":"+ServerPort)
}