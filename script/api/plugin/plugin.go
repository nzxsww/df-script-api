package plugin

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/nzxsww/dragonfly-script-api/script/api/config"
)

type Plugin interface {
	OnEnable()
	OnDisable()
	GetName() string
	GetDataFolder() string
	GetServer() *server.Server
	GetConfig() *config.Config
}

type BasePlugin struct {
	name       string
	dataFolder string
	srv        *server.Server
	cfg        *config.Config
}

func NewBasePlugin(name, dataFolder string, srv *server.Server) *BasePlugin {
	return &BasePlugin{
		name:       name,
		dataFolder: dataFolder,
		srv:        srv,
		cfg:        nil,
	}
}

func (p *BasePlugin) OnEnable()                 {}
func (p *BasePlugin) OnDisable()                {}
func (p *BasePlugin) GetName() string           { return p.name }
func (p *BasePlugin) GetDataFolder() string     { return p.dataFolder }
func (p *BasePlugin) GetServer() *server.Server { return p.srv }

func (p *BasePlugin) GetConfig() *config.Config {
	if p.cfg == nil {
		p.cfg = config.New(p.dataFolder+"/config.yml", p)
		p.cfg.SetDefaults(p.getDefaultConfig())
		p.cfg.Load()
	}
	return p.cfg
}

func (p *BasePlugin) getDefaultConfig() map[string]interface{} {
	return map[string]interface{}{}
}
