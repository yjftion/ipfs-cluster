package ipfsproxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	ma "gx/ipfs/QmNTCey11oxhb1AxDnQBRHtdhap6Ctud872NjAYPYYXPuc/go-multiaddr"
	"gx/ipfs/QmYgGtLm9WJRgh6iuaZap8qVC1gqixFbZCNfhjLNBhWMCm/envconfig"

	"github.com/ipfs/ipfs-cluster/config"
)

const (
	configKey    = "ipfsproxy"
	envConfigKey = "cluster_ipfsproxy"
)

// Default values for Config.
const (
	DefaultListenAddr        = "/ip4/127.0.0.1/tcp/9095"
	DefaultNodeAddr          = "/ip4/127.0.0.1/tcp/5001"
	DefaultReadTimeout       = 0
	DefaultReadHeaderTimeout = 5 * time.Second
	DefaultWriteTimeout      = 0
	DefaultIdleTimeout       = 60 * time.Second
)

// Config allows to customize behaviour of IPFSProxy.
// It implements the config.ComponentConfig interface.
type Config struct {
	config.Saver

	// Listen parameters for the IPFS Proxy.
	ListenAddr ma.Multiaddr

	// Host/Port for the IPFS daemon.
	NodeAddr ma.Multiaddr

	// Maximum duration before timing out reading a full request
	ReadTimeout time.Duration

	// Maximum duration before timing out reading the headers of a request
	ReadHeaderTimeout time.Duration

	// Maximum duration before timing out write of the response
	WriteTimeout time.Duration

	// Server-side amount of time a Keep-Alive connection will be
	// kept idle before being reused
	IdleTimeout time.Duration
}

type jsonConfig struct {
	NodeMultiaddress string `json:"node_multiaddress"`

	ListenMultiaddress string `json:"listen_multiaddress"`
	ReadTimeout        string `json:"read_timeout"`
	ReadHeaderTimeout  string `json:"read_header_timeout"`
	WriteTimeout       string `json:"write_timeout"`
	IdleTimeout        string `json:"idle_timeout"`

	// Below fields are only here to maintain backward compatibility
	// They will be removed in future
	ProxyListenMultiaddress string `json:"proxy_listen_multiaddress,omitempty"`
	ProxyReadTimeout        string `json:"proxy_read_timeout,omitempty"`
	ProxyReadHeaderTimeout  string `json:"proxy_read_header_timeout,omitempty"`
	ProxyWriteTimeout       string `json:"proxy_write_timeout,omitempty"`
	ProxyIdleTimeout        string `json:"proxy_idle_timeout,omitempty"`
}

// toNewFields converts json config written in old style (fields starting with `proxy_`)
// to new style (without `proxy_`)
func (jcfg *jsonConfig) toNewFields() {
	if jcfg.ListenMultiaddress == "" {
		jcfg.ListenMultiaddress = jcfg.ProxyListenMultiaddress
	}

	if jcfg.ReadTimeout == "" {
		jcfg.ReadTimeout = jcfg.ProxyReadTimeout
	}

	if jcfg.ReadHeaderTimeout == "" {
		jcfg.ReadHeaderTimeout = jcfg.ProxyReadHeaderTimeout
	}

	if jcfg.WriteTimeout == "" {
		jcfg.WriteTimeout = jcfg.ProxyWriteTimeout
	}

	if jcfg.IdleTimeout == "" {
		jcfg.IdleTimeout = jcfg.ProxyIdleTimeout
	}
}

// ConfigKey provides a human-friendly identifier for this type of Config.
func (cfg *Config) ConfigKey() string {
	return configKey
}

// Default sets the fields of this Config to sensible default values.
func (cfg *Config) Default() error {
	proxy, err := ma.NewMultiaddr(DefaultListenAddr)
	if err != nil {
		return err
	}
	node, err := ma.NewMultiaddr(DefaultNodeAddr)
	if err != nil {
		return err
	}
	cfg.ListenAddr = proxy
	cfg.NodeAddr = node
	cfg.ReadTimeout = DefaultReadTimeout
	cfg.ReadHeaderTimeout = DefaultReadHeaderTimeout
	cfg.WriteTimeout = DefaultWriteTimeout
	cfg.IdleTimeout = DefaultIdleTimeout

	return nil
}

// Validate checks that the fields of this Config have sensible values,
// at least in appearance.
func (cfg *Config) Validate() error {
	var err error
	if cfg.ListenAddr == nil {
		err = errors.New("ipfsproxy.listen_multiaddress not set")
	}
	if cfg.NodeAddr == nil {
		err = errors.New("ipfsproxy.node_multiaddress not set")
	}

	if cfg.ReadTimeout < 0 {
		err = errors.New("ipfsproxy.read_timeout is invalid")
	}

	if cfg.ReadHeaderTimeout < 0 {
		err = errors.New("ipfsproxy.read_header_timeout is invalid")
	}

	if cfg.WriteTimeout < 0 {
		err = errors.New("ipfsproxy.write_timeout is invalid")
	}

	if cfg.IdleTimeout < 0 {
		err = errors.New("ipfsproxy.idle_timeout invalid")
	}

	return err

}

// LoadJSON parses a JSON representation of this Config as generated by ToJSON.
func (cfg *Config) LoadJSON(raw []byte) error {
	jcfg := &jsonConfig{}
	err := json.Unmarshal(raw, jcfg)
	if err != nil {
		logger.Error("Error unmarshaling ipfsproxy config")
		return err
	}

	// This is here only here to maintain backward compatibility
	// This won't be needed after removing old style fields(starting with `proxy_`)
	jcfg.toNewFields()

	err = cfg.Default()
	if err != nil {
		return fmt.Errorf("error setting config to default values: %s", err)
	}

	// override json config with env var
	err = envconfig.Process("cluster_ipfsproxy", jcfg)
	if err != nil {
		return err
	}

	proxyAddr, err := ma.NewMultiaddr(jcfg.ListenMultiaddress)
	if err != nil {
		return fmt.Errorf("error parsing proxy listen_multiaddress: %s", err)
	}
	nodeAddr, err := ma.NewMultiaddr(jcfg.NodeMultiaddress)
	if err != nil {
		return fmt.Errorf("error parsing ipfs node_multiaddress: %s", err)
	}

	cfg.ListenAddr = proxyAddr
	cfg.NodeAddr = nodeAddr

	err = config.ParseDurations(
		"ipfsproxy",
		&config.DurationOpt{Duration: jcfg.ReadTimeout, Dst: &cfg.ReadTimeout, Name: "read_timeout"},
		&config.DurationOpt{Duration: jcfg.ReadHeaderTimeout, Dst: &cfg.ReadHeaderTimeout, Name: "read_header_timeout"},
		&config.DurationOpt{Duration: jcfg.WriteTimeout, Dst: &cfg.WriteTimeout, Name: "write_timeout"},
		&config.DurationOpt{Duration: jcfg.IdleTimeout, Dst: &cfg.IdleTimeout, Name: "idle_timeout"},
	)
	if err != nil {
		return err
	}

	return cfg.Validate()
}

// ToJSON generates a human-friendly JSON representation of this Config.
func (cfg *Config) ToJSON() (raw []byte, err error) {
	// Multiaddress String() may panic
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	jcfg := &jsonConfig{}

	// Set all configuration fields
	jcfg.ListenMultiaddress = cfg.ListenAddr.String()
	jcfg.NodeMultiaddress = cfg.NodeAddr.String()
	jcfg.ReadTimeout = cfg.ReadTimeout.String()
	jcfg.ReadHeaderTimeout = cfg.ReadHeaderTimeout.String()
	jcfg.WriteTimeout = cfg.WriteTimeout.String()
	jcfg.IdleTimeout = cfg.IdleTimeout.String()

	raw, err = config.DefaultJSONMarshal(jcfg)
	return
}
