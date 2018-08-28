package tests

import (
	"flag"
	"github.com/primasio/wormhole/config"
	"github.com/primasio/wormhole/db"
	"log"
	"math/rand"
	"os"
	"time"
)

func InitTestEnv(configPath string) {
	environment := flag.String("e", "test", "")

	flag.Parse()

	// Init Config
	config.Init(*environment, &configPath)

	// Init Database
	if err := db.Init(); err != nil {
		log.Println("Database:", err)
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())
}
