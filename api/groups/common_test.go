package groups_test

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func startProxyServer(group data.GroupHandler, path string) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	routes := ws.Group(path)
	group.RegisterRoutes(routes, data.ApiRoutesConfig{}, func(_ *gin.Context) {}, func(_ *gin.Context) {}, func(_ *gin.Context) {})
	return ws
}

func loadResponse(rsp io.Reader, destination interface{}) {
	jsonParser := json.NewDecoder(rsp)
	err := jsonParser.Decode(destination)
	if err != nil {
		fmt.Println(err.Error())
	}
}
