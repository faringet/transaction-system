package config

type Config struct {
	PostgresDB PostgresDB `mapstructure:"POSTGRES_DB"`
	Kafka      Kafka      `mapstructure:"KAFKA"`
	Logger     Logger     `mapstructure:"LOGGER"`
	Scheduler  Scheduler  `mapstructure:"SCHEDULER"`
	LocalURL   string     `mapstructure:"THIS_APP_URL"`
}

type PostgresDB struct {
	Addr     string `mapstructure:"ADDR"`
	User     string `mapstructure:"USER"`
	Password string `mapstructure:"PASSWORD"`
	Database string `mapstructure:"DATABASE"`
}

type Kafka struct {
	Brokers []string `mapstructure:"BROKERS"`
	Topic   string   `mapstructure:"TOPIC"`
	GroupID string   `mapstructure:"GROUP_ID"`
}

type Logger struct {
	Production  string `mapstructure:"PRODUCTION"`
	Development string `mapstructure:"DEVELOPMENT"`
}

type Scheduler struct {
	Update int `mapstructure:"UPDATE"`
}
