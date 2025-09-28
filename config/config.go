package config

type Config struct {
    Port      string
    MongoURI  string
    JWTSecret string
    RedisURL  string
    NatsURL   string
}

func Load() Config {
    // TODO: Read .env and return Config
    return Config{}
}
