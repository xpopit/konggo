package main

import (
	"net/http"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

// const PluginName = "header-validation"
// const Version = "0.1.0"
// const Priority = 1000

// const FailedResponse = `{"error": "%s is required"}`

//	type Config struct {
//		HeaderKey string `json:"header_key"`
//	}
//
// Version ของ custom plugin ของเรานั้นเอง
const Version = "0.1"

// Priority เป็นการบอกว่า custom plugin ของเราจะถูกทำ
// ในลำดับที่เท่าไหร่ของ Plugin ท้ังหมดที่เปิดการใช้งาน หากค่าของ Priority
// มีค่าสูงสุด Custom Plugin นี้จะถูกทำก่อน ถ้ามีค่าต่ำสุดก็จะถูกทำที่หลัง
const Priority = 1

type Config struct {
	//ใช้สำหรับกำหนด endpoint ของ Auth Server เนื่องจากเราจะไม่ hardcode ลงไป
	//แต่เราจะไปใส่ค่าบน Kong manager แทน
	Endpoint string
}

// main ตรงนี้ไม่มีไรเป็น start server custom plugin
func main() {
	server.StartServer(New, Version, Priority)
}

func New() interface{} {
	return &Config{}
}

//	func main() {
//		err := server.StartServer(New, Version, Priority)
//		if err != nil {
//			log.Fatalf("Failed start %s plugin", PluginName)
//		}
//	}
//
// Auth server แบบหลอกๆ
func callStubAuthServer(endpoint string, token string) (statusCode int, body []byte) {
	if len(token) == 0 {
		return 401, []byte("Authorization required")
	}
	if token == "TEST_TOKEN" {
		return 200, []byte("ok")
	}
	return
}

func (conf *Config) Access(kong *pdk.PDK) {
	auth, err := kong.Request.GetHeader("Authorization")
	if err != nil {
		kong.Log.Err(err)
		return
	}

	// ทำการตรวจสอบ Access Token ว่าถูกต้องหรือไม่
	sc, body := callStubAuthServer(conf.Endpoint, auth)

	//เมื่อส่ง Access Token ไปตรวจที่ Auth Server แล้วพบว่าไม่ผ่านก็ทำการดีด Request
	//นี้กลับไปยัง Client ทันที
	if sc != http.StatusOK {
		kong.Response.Exit(sc, body, nil)
		return
	}

	//ต่อจากนี้ Kong จะทำการส่ง request ต่อไปยัง API Services ที่กำหนดไว้ใน
	//Gateway Service , Routing
}
