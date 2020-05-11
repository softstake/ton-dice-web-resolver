package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type TonWebResolverConfig struct {
	ContractAddr string `env:"CONTRACT_ADDR,required"`
	KeyFileBase  string `env:"PK_FILE_PATH" envDefault:"owner.pk"`
	StorageHost  string `env:"STORAGE_HOST,required"`
	StoragePort  int32  `env:"STORAGE_PORT" envDefault:"5300"`
	TonAPIHost   string `env:"TON_API_HOST,required"`
	TonAPIPort   int32  `env:"TON_API_PORT" envDefault:"5400"`
}

func GetConfig() TonWebResolverConfig {
	cfg := &TonWebResolverConfig{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal("Cannot parse initial ENV vars: ", err)
	}
	return *cfg
}
