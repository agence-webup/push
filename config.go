package push

type StorageDriver string

const (
	MySQLStorageDriver  StorageDriver = "mysql"
	MemoryStorageDriver StorageDriver = "memory"
)

type RuntimeConfig struct {
	StorageDriver StorageDriver `toml:"storage_driver"`
	MySQL         *MySQLConfig
	APNS          *APNSConfig
	FCM           *FCMConfig
}

type MySQLConfig struct {
	Hostname string
	Port     string
	Database string
	Username string
	Password string
}

type APNSConfig struct {
	CertPath string `toml:"cert_path"`
	CertPass string `toml:"cert_passphrase"`
	Sandbox  bool
	Topic    string
}

type FCMConfig struct {
	URL       string
	ServerKey string `toml:"server_key"`
}
