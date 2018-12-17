package remote_signer

import (
	"os"
	"strconv"
	"strings"
)

var SyslogServer string
var SyslogFacility string
var PrivateKeyFolder string
var KeyPrefix string
var SKSServer string
var HttpPort int
var MaxKeyRingCache int
var EnableRethinkSKS bool
var RethinkDBHost string
var RethinkDBPort int
var RethinkDBUsername string
var RethinkDBPassword string
var RethinkDBPoolSize int
var DatabaseName string
var MasterGPGKeyPath string
var MasterGPGKeyPasswordPath string
var MasterGPGKeyBase64Encoded bool
var KeysBase64Encoded bool
var IgnoreKubernetesCA bool

func init() {
	// Pre init
	MaxKeyRingCache = -1
	HttpPort = -1
	RethinkDBPort = -1
	RethinkDBPoolSize = -1

	// Load envvars
	SyslogServer = os.Getenv("SYSLOG_IP")
	SyslogFacility = os.Getenv("SYSLOG_FACILITY")
	PrivateKeyFolder = os.Getenv("PRIVATE_KEY_FOLDER")
	SKSServer = os.Getenv("SKS_SERVER")
	KeyPrefix = os.Getenv("KEY_PREFIX")

	var maxKeyRingCache = os.Getenv("MAX_KEYRING_CACHE_SIZE")
	if maxKeyRingCache != "" {
		i, err := strconv.ParseInt(maxKeyRingCache, 10, 32)
		if err != nil {
			MaxKeyRingCache = int(i)
		}
	}

	var hp = os.Getenv("HTTP_PORT")
	if hp != "" {
		i, err := strconv.ParseInt(hp, 10, 32)
		if err != nil {
			HttpPort = int(i)
		}
	}

	EnableRethinkSKS = strings.ToLower(os.Getenv("ENABLE_RETHINK_SKS")) == "true"

	RethinkDBHost = os.Getenv("RETHINKDB_HOST")
	RethinkDBUsername = os.Getenv("RETHINKDB_USERNAME")
	RethinkDBPassword = os.Getenv("RETHINKDB_PASSWORD")

	var rdbport = os.Getenv("RETHINKDB_PORT")
	if rdbport != "" {
		i, err := strconv.ParseInt(rdbport, 10, 32)
		if err != nil {
			RethinkDBPort = int(i)
		}
	}

	var poolSize = os.Getenv("RETHINKDB_POOL_SIZE")
	if poolSize != "" {
		i, err := strconv.ParseInt(poolSize, 10, 32)
		if err != nil {
			RethinkDBPoolSize = int(i)
		}
	}

	DatabaseName = os.Getenv("DATABASE_NAME")
	MasterGPGKeyPath = os.Getenv("MASTER_GPG_KEY_PATH")
	MasterGPGKeyPasswordPath = os.Getenv("MASTER_GPG_KEY_PASSWORD_PATH")
	MasterGPGKeyBase64Encoded = strings.ToLower(os.Getenv("MASTER_GPG_KEY_BASE64_ENCODED")) == "true"

	KeysBase64Encoded = strings.ToLower(os.Getenv("KEYS_BASE64_ENCODED")) == "true"
	IgnoreKubernetesCA = strings.ToLower(os.Getenv("IGNORE_KUBERNETES_CA")) == "true"

	// Set defaults if not defined
	if SyslogServer == "" {
		SyslogServer = "127.0.0.1"
	}

	if SyslogFacility == "" {
		SyslogFacility = "LOG_USER"
	}

	if PrivateKeyFolder == "" {
		PrivateKeyFolder = "./keys"
	}

	if MaxKeyRingCache == -1 {
		MaxKeyRingCache = 1000
	}

	if HttpPort == -1 {
		HttpPort = 5100
	}

	if RethinkDBHost == "" {
		RethinkDBHost = "127.0.0.1"
	}

	if RethinkDBUsername == "" {
		RethinkDBUsername = "admin"
	}

	if RethinkDBPort == -1 {
		RethinkDBPort = 28015
	}

	if RethinkDBPoolSize == -1 {
		RethinkDBPoolSize = 10
	}

	if DatabaseName == "" {
		DatabaseName = "remote_signer"
	}

	// Other stuff
	_ = os.Mkdir(PrivateKeyFolder, 0770)
}
