package config

import "github.com/hashicorp/vault/api"

// Confs var
var Confs configs

type (
	// configs struct
	configs struct {
		Notifs struct {
			Host       string
			Path       string
			DebugAddr  string
			HTTPAddr   string
			GrpcAddr   string
			ThriftAddr string
		}
		Users struct {
			Host       string
			Path       string
			DebugAddr  string
			HTTPAddr   string
			GrpcAddr   string
			ThriftAddr string
		}
		Vault struct {
			Address string
			Token   string
			Logical *api.Logical
		}
		JWT struct {
			Secret  string
			RSecret string
			Path    string
		}
		Redis struct {
			Path     string
			Host     string
			Password string
			DB       string
		}
	}
)
