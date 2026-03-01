package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/nzxsww/dragonfly-script-api/script/api"
	playerEvent "github.com/nzxsww/dragonfly-script-api/script/api/event/player"
	"github.com/nzxsww/dragonfly-script-api/script/api/loader"
	"github.com/nzxsww/dragonfly-script-api/script/api/plugin"
	"github.com/pelletier/go-toml/v2"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	chat.Global.Subscribe(chat.StdoutSubscriber{})

	conf, err := readConfig(slog.Default())
	if err != nil {
		panic(err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()

	pluginMgr := plugin.NewManager(srv)

	ldr := loader.NewLoader("plugins", srv, pluginMgr)
	loadedPlugins, err := ldr.LoadAll()
	if err != nil {
		slog.Error("Error loading plugins", "error", err)
	}

	for _, p := range loadedPlugins {
		fmt.Printf("Enabling plugin: %s\n", p.GetName())
		p.OnEnable()
	}

	srv.Listen()
	h := api.NewHandler(srv, pluginMgr)
	for p := range srv.Accept() {
		// Asignar el handler de eventos al jugador y disparar PlayerJoinEvent
		// Los plugins JS ya están suscritos directamente al sistema de eventos
		p.Handle(h)
		pluginMgr.CallEvent(playerEvent.NewPlayerJoinEvent(p))
	}

	// El servidor se cerró — notificar a todos los plugins para que hagan cleanup
	fmt.Println("Servidor cerrado. Desactivando plugins...")
	for _, p := range loadedPlugins {
		fmt.Printf("Disabling plugin: %s\n", p.GetName())
		p.OnDisable()
	}
}

func readConfig(log *slog.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}
		return c.Config(log)
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}
	return c.Config(log)
}
