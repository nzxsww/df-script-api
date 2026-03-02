package loader

import (
	"fmt"

	dfitem "github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/dop251/goja"
)

// registerItemAPI expone el objeto global `item` para crear itemStacks desde JS.
func (l *Loader) registerItemAPI(vm *goja.Runtime, p *ScriptPlugin) {
	vm.Set("item", map[string]interface{}{
		"create": func(name string, count int) map[string]interface{} {
			it, ok := world.ItemByName(name, 0)
			if !ok {
				fmt.Printf("[%s] item.create: item desconocido '%s'\n", p.name, name)
				return newItemWrapper(dfitem.Stack{})
			}
			if count < 1 {
				count = 1
			}
			return newItemWrapper(dfitem.NewStack(it, count))
		},
	})
}
