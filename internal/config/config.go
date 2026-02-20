package config

// Config holds the application configuration
type Config struct {
	DBDriver      string
	DBSource      string
	ServerAddress string
}

// LoadConfig returns the hardcoded configuration for now
// In a real app, this might load from env vars or a file
func LoadConfig() Config {
	return Config{
		DBDriver:      "mysql",
		// Matches connection.php: localhost, root, empty pass, college_db, port 3307
		DBSource:      "root:@tcp(127.0.0.1:3307)/college_db?parseTime=true",
		ServerAddress: ":8080",
	}
}
