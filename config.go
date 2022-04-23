package main

import (
	"fmt"
	"strings"
        "encoding/json"
	"os"
)



// config
type Configuration struct {
    User    []string
    Password   []string
    Server   []string
    Port   []string
    Db   []string
}

var sqlcred = getconfig()

func getconfig() string {
file, _ := os.Open("conf.json")
defer file.Close()
decoder := json.NewDecoder(file)
configuration := Configuration{}
errcode := decoder.Decode(&configuration)
if errcode != nil {
  fmt.Println("error:", errcode)
}

return ``+strings.Join(configuration.User," ")+`:`+strings.Join(configuration.Password," ")+`@tcp(`+strings.Join(configuration.Server," ")+`:`+strings.Join(configuration.Port," ")+`)/`+strings.Join(configuration.Db," ")+``
}
