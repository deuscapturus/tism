// config provides configuration setting from the provided yaml file.
package config

import (
	"flag"
	"github.com/deuscapturus/tism/randid"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

type Configuration struct {
	Port            string   `yaml:"port,omitempty"`
	Address         string   `yaml:"address,omitempty"`
	TLSCertFile     string   `yaml:"tls_cert_file,omitempty"`
	TLSKeyFile      string   `yaml:"tls_key_file,omitempty"`
	TLSDir          string   `yaml:"tls_directory,omitempty"`
	RevokedJWTs     []string `yaml:"revoked_tokens,omitempty"`
	JWTsecret       string   `yaml:"token_secret,omitempty"`
	KeyRingFilePath string   `yaml:"keyring_path,omitempty"`
	GenAdminToken   bool
	GenCert         bool
	ConfigPath      string
}

// Config default configuration values
var Config = &Configuration{
	KeyRingFilePath: "gpgkeys/secring.gpg",
	TLSCertFile:     "cert.crt",
	TLSKeyFile:      "cert.key",
	TLSDir:          "./cert/",
}

// Load configuration at startup
func init() {
	parseFlags()
	renderConfigTemplate()
	Load()
}

// ParseFlags Command line configuration values
func parseFlags() {


	ConfigPath := flag.String("config", "config.yaml", "Configuration file path.")
	GenAdminToken := flag.Bool("t", false, "Generate a super admin token")
	GenCert := flag.Bool("c", false, "Generate a new TLS certificate. WARNING; WILL OVERWRITE")
	Port := flag.String("p", "8080", "Port to listen on")
	Address := flag.String("a", "0.0.0.0", "Address to listen on")

	flag.Parse()

	Config.GenAdminToken = *GenAdminToken
	Config.GenCert = *GenCert
	Config.ConfigPath = *ConfigPath
	Config.Port = *Port
	Config.Address = *Address
}

// renderConfigTemplate
func renderConfigTemplate() {

	configTemplate, err := template.ParseFiles(Config.ConfigPath)
	if err != nil {
		log.Println(err)
	}

	type ConfigValues struct {
		TokenSecret string
	}
	configTemplate.Option("missingkey=error")
	//TODO Do not touch the file unless template variables are found.
	configFile, err := os.Create(Config.ConfigPath)
	if err != nil {
		log.Fatal("Unable to create configuration file")
	}

	err = configTemplate.Execute(configFile, ConfigValues{generateTokenSecret()})
	if err != nil {
		log.Fatal("Unable to write configuration file")
	}
}

// Load configuration from the provided confioguration path.
func Load() {

	ConfigFileBytes, err := ioutil.ReadFile(Config.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(ConfigFileBytes, &Config)
	if err != nil {
		log.Fatal(err)
	}

	if Config.JWTsecret == "" {
		log.Fatal("TokenSecret is undefined")
	}
}

// generateTokenSecret generates a random 256 character string
func generateTokenSecret() string {

	return randid.GenerateSecret(256)
}
