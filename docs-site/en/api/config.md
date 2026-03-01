# YAML Config

Each plugin has its own configuration in a YAML file at `plugins/plugin-name/config.yml`. The file is created automatically when the plugin calls `config.save()`.

## Recommended flow

```js
// 1. Define defaults at the start of the script (outside onEnable)
config.setDefaults({
    "prefix":      "§a[MyPlugin]§r",
    "max-players": 20,
    "debug-mode":  false,
    "multiplier":  1.5
});

function onEnable() {
    // 2. Read values (uses defaults if no config.yml exists)
    var prefix = config.getString("prefix", "§f[?]");
    console.log("Prefix: " + prefix);
}

function onDisable() {
    // 3. Save on close
    config.save();
}
```

## Available methods

| Method | Description |
|---|---|
| `config.setDefaults(obj)` | Define default values. Does NOT overwrite existing values. |
| `config.get(key, default)` | Get any value |
| `config.getString(key, default)` | Get value as string |
| `config.getInt(key, default)` | Get value as integer |
| `config.getBool(key, default)` | Get value as boolean |
| `config.getFloat(key, default)` | Get value as float |
| `config.set(key, value)` | Write a value in memory |
| `config.save()` | Save to disk (creates file if needed) |
| `config.reload()` | Reload from disk |

## Notes

- `config.yml` is only created when `config.save()` is called explicitly — typically in `onDisable()`.
- `setDefaults()` does not overwrite existing values — user changes in `config.yml` are preserved on restart.
- Call `config.reload()` if the file was edited manually while the server is running.
