# Console

The `console` object lets you print messages to the server console. All messages automatically include the `[PluginName]` prefix.

## Methods

| Method | When to use |
|---|---|
| `console.log(msg)` | General informational messages |
| `console.info(msg)` | Important information |
| `console.warn(msg)` | Warnings |
| `console.error(msg)` | Errors |

```js
console.log("Plugin loaded");
// Output: [MyPlugin] Plugin loaded

console.warn("Config not found, using defaults");
// Output: [MyPlugin] [WARN] Config not found, using defaults

console.error("Could not save config");
// Output: [MyPlugin] [ERROR] Could not save config
```
