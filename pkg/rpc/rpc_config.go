package rpc

import (
	"github.com/nspcc-dev/neo-go/pkg/encoding/fixedn"
)

type (
	// Config is an RPC service configuration information.
	Config struct {
		Address              string `yaml:"Address"`
		Enabled              bool   `yaml:"Enabled"`
		EnableCORSWorkaround bool   `yaml:"EnableCORSWorkaround"`
		// MaxGasInvoke is the maximum amount of GAS which
		// can be spent during an RPC call.
		MaxGasInvoke           fixedn.Fixed8 `yaml:"MaxGasInvoke"`
		MaxIteratorResultItems int           `yaml:"MaxIteratorResultItems"`
		MaxFindResultItems     int           `yaml:"MaxFindResultItems"`
		MaxNEP11Tokens         int           `yaml:"MaxNEP11Tokens"`
		Port                   uint16        `yaml:"Port"`
		SessionEnabled         bool          `yaml:"SessionEnabled"`
		SessionExpirationTime  int           `yaml:"SessionExpirationTime"`
		SessionBackedByMPT     bool          `yaml:"SessionBackedByMPT"`
		SessionPoolSize        int           `yaml:"SessionPoolSize"`
		StartWhenSynchronized  bool          `yaml:"StartWhenSynchronized"`
		TLSConfig              TLSConfig     `yaml:"TLSConfig"`
	}

	// TLSConfig describes SSL/TLS configuration.
	TLSConfig struct {
		Address  string `yaml:"Address"`
		CertFile string `yaml:"CertFile"`
		Enabled  bool   `yaml:"Enabled"`
		Port     uint16 `yaml:"Port"`
		KeyFile  string `yaml:"KeyFile"`
	}
)
