package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbDriver         string `mapstructure:"DB_DRIVER"`
	DbSource         string `mapstructure:"DB_SOURCE"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDb       string `mapstructure:"POSTGRES_DB"`
	ServerAddress    string `mapstructure:"SERVER_ADDRESS"`
	JwtSecret        string `mapstructure:"JWT_SECRET"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("util/config") // Menetapkan folder tempat file config berada
	viper.SetConfigName("config")      // Nama file config tanpa ekstensi
	viper.SetConfigType("json")        // Menggunakan format JSON

	viper.AutomaticEnv()

	err = viper.ReadInConfig() // Membaca file config
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config) // Mengurai isi file ke dalam struct Config
	return
}
