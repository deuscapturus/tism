// config provides configuration setting from the provided yaml file.
package config

import (
	"flag"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
)

type Configuration struct {
	Port            string   `yaml:"port,omitempty"`
	TLSCertFile     string   `yaml:"tls_cert_file,omitempty"`
	TLSKeyFile      string   `yaml:"tls_key_file,omitempty"`
	TLSDir          string   `yaml:"tls_directory,omitempty"`
	RevokedJWTs     []string `yaml:"revoked_tokens,omitempty"`
	JWTsecret       string   `yaml:"token_secret,omitempty"`
	KeyRingFilePath string   `yaml:"keyring_path,omitempty"`
	GenAdminToken   bool
	GenCert         bool
}

var Config *Configuration

// Load configuration at startup
func init() {
	Load()
}

// LoadConfiguration load athe configuration from the provided confioguration path.
func Load() {

	// Default configuration values
	Config = &Configuration{
		Port:            "8080",
		KeyRingFilePath: "gpgkeys/secring.gpg",
		TLSCertFile:     "cert.crt",
		TLSKeyFile:      "cert.key",
		TLSDir:          "./cert/",
	}

	// Command line configuration values
	ConfigPath := flag.String("config", "config.yaml", "Configuration file path.")

	GenAdminToken := flag.Bool("t", false, "Generate a super admin token")
	GenCert := flag.Bool("c", false, "Generate a new TLS certificate. WARNING; WILL OVERWRITE")

	flag.Parse()

	ConfigFileBytes, err := ioutil.ReadFile(*ConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(ConfigFileBytes, &Config)
	if err != nil {
		log.Fatal(err)
	}

	Config.GenAdminToken = *GenAdminToken
	Config.GenCert = *GenCert
}
