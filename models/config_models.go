package models

type Environment string

func (env Environment) IsValid() bool {
	switch env {
	case DevEnv, ProdEnv, UATEnv:
		return true
	}
	return false
}

const (
	DevEnv  Environment = "dev_env"
	ProdEnv Environment = "prod_env"
	UATEnv  Environment = "uat_env"
)

type AppConfig struct {
	APPBaseURL   string   `json:"app_base_url"`
	EtcdEndpoint []string `json:"etcd_endpoint"`
	DBConnURL    string   `json:"db_conn_url"`
	Service3Url  string   `json:"service_3_url"`
	// DBHost                          string    `json:"db_host"`
	// DBPort                          int       `json:"db_port"`
	// DBUser                          string    `json:"db_user"`
	// DBPassword                      string    `json:"db_password"`
	// DBName                          string    `json:"db_name"`
	AppServerPort                   string             `json:"app_server_port"`
	MinIOInfo                       MinIOInfo          `json:"minio_info"`
	TimeoutRequestRetryCount        int                `json:"timeout_request_retry_count"`
	TimeoutRequestRetryIntervalInMS int                `json:"timeout_request_retry_interval_in_ms"`
	Redis                           RedisInfo          `json:"redis"`
	Rigel                           RigelInfo          `json:"rigel"`
	KafkaConfig                     KafkaConfig        `json:"kafka_config"`
	KafkaBillingConfig              KafkaBillingConfig `json:"kafka_billing_config"`
	Service1Url                     string             `json:"service_1_url"`
	KafkaSxpFailureConfig           KafkaCacheConfig   `json:"kafka_sxp_failure_config"`
	KafkaSxpSuccessConfig           KafkaCacheConfig   `json:"kafka_sxp_success_config"`
	Serversage                      Serversage         `json:"serversage"`
	Service4Url                     string             `json:"service_4_url"`
	KafkaInvalidateCacheConfig      KafkaCacheConfig   `json:"kafka_invalidate_cache_config"`
	KafkaRtaReprocess               KafkaCacheConfig   `json:"kafka_rta_reprocess"`
}
type KafkaCacheConfig struct {
	IsKafkaOn    bool     `json:"is_kafka_on"`
	KafkaBrokers []string `json:"kafka_broker"`
	KafkaTopic   string   `json:"kafka_topic"`
	KafkaIndex   string   `json:"kafka_index"`
	Producer     struct {
		FlushFrequency int    `json:"flush_frequency"`
		Topic          string `json:"topic"`
	}
	Consumer struct {
		GroupID string   `json:"group_id"`
		Topic   []string `json:"topic"`
	}
}
type RigelInfo struct {
	AppName       string `json:"app_name"`
	ModuleName    string `json:"module_name"`
	VersionNumber int    `json:"version_number"`
	ConfigName    string `json:"config_name"`
}

type RedisInfo struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type KafkaConfig struct {
	IsKafkaOn    bool     `json:"is_kafka_on"`
	KafkaBrokers []string `json:"kafka_broker"`
	KafkaTopic   string   `json:"kafka_topic"`
}

type KafkaBillingConfig struct {
	IsKafkaOn    bool     `json:"is_kafka_on"`
	KafkaBrokers []string `json:"kafka_broker"`
	KafkaTopic   string   `json:"kafka_topic"`
}

type Serversage struct {
	IsInstrumentApplication bool   `json:"is_instrument_application"`
	OtelEndpoint            string `json:"otelEndpoint"`
	ServiceName             string `json:"serviceName"`
}
