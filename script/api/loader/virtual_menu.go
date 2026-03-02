package loader

import (
	"fmt"
	"reflect"

	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/item"
	dfplayer "github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/dop251/goja"
)

type menuWrapper struct {
	menu      inv.Menu
	vm        *goja.Runtime
	plugin    *ScriptPlugin
	onClick   goja.Callable
	onClose   goja.Callable
	container inv.Container
	customInv bool
}

type menuSubmittable struct {
	wrapper *menuWrapper
}

func (m menuSubmittable) Submit(p *dfplayer.Player, it item.Stack) {
	if m.wrapper.onClick == nil {
		return
	}
	name, _ := it.Item().EncodeItem()
	jsItem := map[string]interface{}{
		"name":  name,
		"count": it.Count(),
	}
	clickType := getLastClickType(p)
	if _, err := m.wrapper.onClick(goja.Undefined(), m.wrapper.vm.ToValue(newPlayerWrapper(p)), m.wrapper.vm.ToValue(jsItem), m.wrapper.vm.ToValue(clickType)); err != nil {
		fmt.Printf("[%s] Error en menu.onClick: %v\n", m.wrapper.plugin.name, err)
	}
}

func (m menuSubmittable) Close(p *dfplayer.Player) {
	if m.wrapper.onClose == nil {
		return
	}
	if _, err := m.wrapper.onClose(goja.Undefined(), m.wrapper.vm.ToValue(newPlayerWrapper(p))); err != nil {
		fmt.Printf("[%s] Error en menu.onClose: %v\n", m.wrapper.plugin.name, err)
	}
}

func menuContainerFromType(menuType string, size int) inv.Container {
	switch menuType {
	case "hopper":
		return inv.ContainerHopper{}
	case "barrel":
		return inv.ContainerBarrel{}
	case "dropper":
		return inv.ContainerDropper{}
	case "ender_chest":
		return inv.ContainerEnderChest{}
	case "chest":
		if size == 54 {
			return inv.ContainerChest{DoubleChest: true}
		}
		return inv.ContainerChest{DoubleChest: false}
	default:
		// fallback a chest
		if size == 54 {
			return inv.ContainerChest{DoubleChest: true}
		}
		return inv.ContainerChest{DoubleChest: false}
	}
}

func menuSizeFromType(menuType string) int {
	switch menuType {
	case "hopper":
		return 5
	case "dropper":
		return 9
	case "barrel":
		return 27
	case "ender_chest":
		return 27
	default:
		return 27
	}
}

func buildMenuWrapper(vm *goja.Runtime, plugin *ScriptPlugin, title, menuType string, size int) *menuWrapper {
	if size == 0 {
		size = menuSizeFromType(menuType)
	}
	container := menuContainerFromType(menuType, size)
	mw := &menuWrapper{vm: vm, plugin: plugin, container: container}
	mw.menu = inv.NewMenu(menuSubmittable{wrapper: mw}, title, container)
	return mw
}

func menuWrapperToJS(vm *goja.Runtime, mw *menuWrapper) map[string]interface{} {
	return map[string]interface{}{
		"setItems": func(items []interface{}) {
			stacks := make([]item.Stack, mw.container.Size())
			for _, entry := range items {
				m, ok := entry.(map[string]interface{})
				if !ok {
					continue
				}
				slot, _ := m["slot"].(int64)
				name, _ := m["name"].(string)
				countVal, _ := m["count"].(int64)
				if name == "" {
					continue
				}
				it, ok := world.ItemByName(name, 0)
				if !ok {
					fmt.Printf("[%s] menu.setItems: item desconocido '%s'\n", mw.plugin.name, name)
					continue
				}
				if int(slot) >= len(stacks) || int(slot) < 0 {
					continue
				}
				count := int(countVal)
				if count <= 0 {
					count = 1
				}
				stacks[int(slot)] = item.NewStack(it, count)
			}
			mw.menu = mw.menu.WithStacks(stacks...)
		},
		"pattern": func(rows []string, key map[string]map[string]interface{}) {
			stacks := make([]item.Stack, mw.container.Size())
			for r, row := range rows {
				for c := 0; c < len(row); c++ {
					ch := string(row[c])
					entry, ok := key[ch]
					if !ok {
						continue
					}
					name, _ := entry["name"].(string)
					if name == "" {
						continue
					}
					countVal, _ := entry["count"].(int64)
					count := int(countVal)
					if count <= 0 {
						count = 1
					}
					it, ok := world.ItemByName(name, 0)
					if !ok {
						fmt.Printf("[%s] menu.pattern: item desconocido '%s'\n", mw.plugin.name, name)
						continue
					}
					slot := c + r*9
					if slot < 0 || slot >= len(stacks) {
						continue
					}
					stacks[slot] = item.NewStack(it, count)
				}
			}
			mw.menu = mw.menu.WithStacks(stacks...)
		},
		"onClick": func(fn goja.Callable) {
			mw.onClick = fn
		},
		"onClose": func(fn goja.Callable) {
			mw.onClose = fn
		},
		"open": func(target goja.Value) {
			mw.openMenu(target, false)
		},
		"update": func(target goja.Value) {
			mw.openMenu(target, true)
		},
		"close": func(target goja.Value) {
			mw.closeMenu(target)
		},
	}
}

func (mw *menuWrapper) resolvePlayerName(target goja.Value) string {
	if target == nil || goja.IsUndefined(target) || goja.IsNull(target) {
		return ""
	}
	if target.ExportType().Kind() == reflect.String {
		return target.String()
	}
	if obj, ok := target.(*goja.Object); ok {
		if fnVal := obj.Get("getName"); fnVal != nil {
			if fn, ok := goja.AssertFunction(fnVal); ok {
				nameVal, err := fn(goja.Undefined(), target)
				if err == nil {
					return nameVal.String()
				}
			}
		}
	}
	return ""
}

func (mw *menuWrapper) openMenu(target goja.Value, update bool) {
	name := mw.resolvePlayerName(target)
	if name == "" || mw.plugin.srv == nil {
		return
	}
	go func() {
		handle, ok := mw.plugin.srv.PlayerByName(name)
		if !ok {
			return
		}
		handle.ExecWorld(func(tx *world.Tx, e world.Entity) {
			if pl, ok := e.(*dfplayer.Player); ok {
				if update {
					inv.UpdateMenu(pl, mw.menu)
				} else {
					inv.SendMenu(pl, mw.menu)
				}
			}
		})
	}()
}

func (mw *menuWrapper) closeMenu(target goja.Value) {
	name := mw.resolvePlayerName(target)
	if name == "" || mw.plugin.srv == nil {
		return
	}
	go func() {
		handle, ok := mw.plugin.srv.PlayerByName(name)
		if !ok {
			return
		}
		handle.ExecWorld(func(tx *world.Tx, e world.Entity) {
			if pl, ok := e.(*dfplayer.Player); ok {
				inv.CloseContainer(pl)
			}
		})
	}()
}

