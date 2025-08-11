package main

import (
	"flag"
	"log"
	"os"
	"signature/internal/adapter/inbound"
	"signature/internal/application/db"
	"signature/internal/constants"
	"signature/internal/entity/models"
	"signature/internal/service"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config YAML file")
	flag.Parse()
	data, err := os.ReadFile(*configPath)
	if err != nil {
		log.Panic("err loading yaml: ", err)
	}
	var config models.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Panic("err unmarshalling yaml: ", err)
	}

	constants.DBHOST = config.Database.Host
	constants.DBPORT = config.Database.Port
	constants.DBUSER = config.Database.User
	constants.DBNAME = config.Database.Name
	constants.DBPASS = config.Database.Pass

	constants.PORT = config.Port

	constants.PRIVATEKEY = config.Org.PrivateKey
	constants.SIGNMETHOD = config.Org.SignMethod
	constants.ORGTYPE = config.Org.OrgType
	constants.ORGID = config.Org.OrgId

	db := db.NewDB(constants.DBHOST, constants.DBPORT, constants.DBNAME, constants.DBUSER, constants.DBPASS)
	service := service.NewService(db)
	handler := inbound.NewHandler(service)
	done := make(chan any)
	defer close(done)
	go handler.StartWorker(done)
	router := gin.Default()
	router.POST(constants.ENDPOINT_REQPAY, handler.ReqPay)
	router.POST(constants.ENDPOINT_RESPPAY, handler.RespPay)
	switch constants.ORGTYPE {
	case constants.ORG_REGULATOR:
		router.POST(constants.ENDPOINT_ONBOARDBANK, handler.OnboardBank)
	case constants.ORG_BANK:
		router.POST(constants.ENDPOINT_ONBOARDREGULATOR, handler.OnboardRegulator)
	}
	if err := router.Run(":" + constants.PORT); err != nil {
		log.Panic(err)
	}
}
