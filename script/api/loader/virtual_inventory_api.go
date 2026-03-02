package loader

import "github.com/dop251/goja"

// registerVirtualInventories expone un API JS para inventarios virtuales.
func (l *Loader) registerVirtualInventories(vm *goja.Runtime, p *ScriptPlugin) {
	vm.Set("inventory", map[string]interface{}{
		"createMenu": func(opts map[string]interface{}) map[string]interface{} {
			title, _ := opts["title"].(string)
			menuType, _ := opts["type"].(string)
			sizeVal, _ := opts["size"].(int64)
			if title == "" {
				title = p.name
			}
			if menuType == "" {
				menuType = "chest"
			}
			mw := buildMenuWrapper(vm, p, title, menuType, int(sizeVal))
			return menuWrapperToJS(vm, mw)
		},
	})
}
