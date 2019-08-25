package conf

import (
	"flag"
	uuid "github.com/satori/go.uuid"
	"strconv"
)

var (
	//Ref = "/join/"
	Port = 8100
	PublicVersion = "dcrvnwww"
)

const (
	RefSession = "RefSession"
)
func PortToServe() string {
	return ":" + strconv.Itoa(Port)
}

func Init() {
	u4 := uuid.NewV4()
	PublicVersion = u4.String()
	flag.IntVar(&Port, "port", 8100, "http port")
	flag.Parse()
}
