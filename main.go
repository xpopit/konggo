package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

type Config struct {
	URL            string
	ConnectTimeout int
	SendTimeout    int
	ReadTimeout    int
}

// Version of the plugin
const Version = "0.1.0"
const Priority = 1000
const PluginName = "konggo"

func New() interface{} {
	return &Config{}
}

//	func callStubAuthServer(endpoint string, token string) (statusCode int, body []byte) {
//		if len(token) == 0 {
//			return 401, []byte("Authorization required")
//		}
//		if token == "TEST_TOKEN" {
//			return 200, []byte("ok")
//		}
//		return
//	}
func (conf Config) Access(kong *pdk.PDK) {
	// auth, err := kong.Request.GetHeader("Authorization")
	// if err != nil {
	// 	kong.Log.Err(err)
	// 	return
	// }
	// log.Println("Authorization: ", auth)

	// http.NewRequest(kong.Request.GetMethod(), conf.Endpoint, nil)
	path, err := kong.Request.GetPath()
	if err != nil {
		kong.Log.Err("Error getting path: " + err.Error())
		kong.Response.Exit(500, []byte("Internal Server Error"), nil)
		return
	}

	query, err := kong.Request.GetRawQuery()
	if err != nil {
		kong.Log.Err("Error getting query: " + err.Error())
		kong.Response.Exit(500, []byte("Internal Server Error"), nil)
		return
	}

	// Construct the URL
	url := conf.URL + path
	if query != "" {
		url += "?" + query
	}
	method, err := kong.Request.GetMethod()
	if err != nil {
		kong.Log.Err("Failed to create request: " + err.Error())
		kong.Response.Exit(500, []byte("Internal Server Error"), nil)
		return
	}
	// Create a new HTTP request
	req, err := http.NewRequest(method, url, nil) // Assuming the method is GET
	if err != nil {
		kong.Log.Err("Failed to create request: " + err.Error())
		kong.Response.Exit(500, []byte("Internal Server Error"), nil)
		return
	}

	// headerResponse := make(map[string][]string, 0)
	// headerResponse["Content-Type"] = []string{"application/json"}

	// if headerKey == "" {
	// 	kong.Response.Exit(400, fmt.Sprintf(FailedResponse, conf.HeaderKey), headerResponse)
	// }
	client := &http.Client{
		Timeout: time.Duration(conf.ConnectTimeout+conf.SendTimeout+conf.ReadTimeout) * time.Millisecond,
	}

	headers, err := kong.Request.GetHeaders(-1)
	if err != nil {
		kong.Log.Err(err.Error())
		kong.Response.Exit(500, []byte("Internal Server Error"), nil)
		return
	}

	for k, v := range headers {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		kong.Log.Err(err.Error())
		kong.Response.Exit(500, []byte("Internal Server Error"), nil)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		kong.Response.Exit(401, []byte("Unauthorized"), nil)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	server.StartServer(New, Version, Priority)
}
