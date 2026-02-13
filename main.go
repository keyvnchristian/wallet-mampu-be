package main

import (
	"log"
	"mampu/config"
	"mampu/repository"
	"mampu/transport"
	"mampu/usecase"
	"net/http"

	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

func init() {
	viper.SetConfigType(`json`)
	viper.AddConfigPath(`./config`)
	viper.SetConfigName(`config`)
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Fatal error config file: %s\n", err)
	}
}

func main() {
	db, err := config.Conn()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	log.Println("database connected")

	repos := repository.NewRepositories(db)

	walletUC := usecase.NewWalletService(db, repos)

	healthHandler := transport.NewHealthHandler(db)

	walletHandler := transport.NewWalletHandler(walletUC)

	http.HandleFunc("/ping", healthHandler.Ping)
	http.HandleFunc("/withdraw", walletHandler.Withdraw)
	http.HandleFunc("/wallet", walletHandler.GetWallet)
	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}
	log.Println("server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
