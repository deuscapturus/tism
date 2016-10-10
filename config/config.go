// config provides configuration setting from the provided yaml file.
package config

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"flag"
)

type Configuration struct {
	Port            string      `yaml:"port,omitempty"`
	TLSCertFile     string   `yaml:"tls_cert_file,omitempty"`
	TLSKeyFile      string   `yaml:"tls_key_file,omitempty"`
	RevokedJWTs     []string `yaml:"revoked_api_keys,omitempty"`
	JWTsecret       string   `yaml:"token_secret,omitempty"`
	KeyRingFilePath string   `yaml:"keyring_path,omitempty"`
	GenAdminToken   bool
}

var Config *Configuration

// LoadConfiguration load athe configuration from the provided confioguration path.
func Load() {

	ConfigPath := flag.String("config", "config.yaml", "Configuration file path.")
	GenAdminToken := flag.Bool("t", false, "Generate a super admin token")

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

}
