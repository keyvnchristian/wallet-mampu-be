package config

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func GetPostgresDSN() string {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString(`database.host`),
		viper.GetString(`database.user`),
		viper.GetString(`database.password`),
		viper.GetString(`database.name`))
	return dsn
}

func Conn() (*sqlx.DB, error) {
	conn, err := sqlx.Open(`postgres`, GetPostgresDSN())
	if err != nil {
		fmt.Println("Database connection error:", err)
		return conn, err
	}

	return conn, err
}
