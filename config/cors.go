package config

import "github.com/gin-contrib/cors"

func CORSConfig() cors.Config {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Authorization"}
	return config
}
