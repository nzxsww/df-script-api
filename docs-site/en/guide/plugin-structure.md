# Plugin Structure

## Lifecycle

Every plugin has three phases:

```
Server starts
      ↓
 Root script code runs       ← setDefaults(), global variables
      ↓
   onEnable() is called      ← register events and commands
      ↓
   [Server is running]
      ↓
   onDisable() is called     ← save config, cleanup
      ↓
 Server closes
```

::: warning Important
`events.on()` and `commands.register()` **must be called inside `onEnable()`**, never at the root level of the script. The root level runs before the system is ready to receive registrations.
:::

## Available global variables

| Variable | Type | Description |
|---|---|---|
| `plugin` | object | Plugin metadata |
| `console` | object | Server logger |
| `events` | object | Register event listeners |
| `commands` | object | Register commands |
| `config` | object | Plugin YAML config |
| `setTimeout` | function | Run code after a delay |
| `setInterval` | function | Run code repeatedly |
| `clearInterval` | function | Cancel an interval |

### `plugin`

```js
plugin.name        // "MyPlugin"
plugin.version     // "1.0.0"
plugin.author      // "YourName"
plugin.dataFolder  // "plugins/my-plugin"
```

## Full template

```js
// =========================================
// My Plugin — Dragonfly Script API
// =========================================

console.log("Loading: " + plugin.name + " v" + plugin.version);

// Default configuration
// config.yml is created at plugins/my-plugin/config.yml
config.setDefaults({
    "prefix":       "§a[MyPlugin]§r",
    "welcome":      "Welcome to the server!",
    "max-players":  20
});

// =========================================
// Lifecycle
// =========================================

function onEnable() {
    var prefix  = config.getString("prefix", "§f[?]");
    var welcome = config.getString("welcome", "Welcome");

    events.on("PlayerJoin", function(event) {
        var player = event.getPlayer();
        player.sendMessage(prefix + " " + welcome);
    });

    commands.register("hello", "Say hello", function(player, args) {
        player.sendMessage(prefix + " Hello, " + player.getName() + "!");
    });

    console.log("Plugin enabled with prefix: " + prefix);
}

function onDisable() {
    config.save();
    console.log("Plugin disabled.");
}

// Export (required)
module = {
    onEnable:  onEnable,
    onDisable: onDisable
};
```

## Colors in messages

Minecraft Bedrock uses `§` followed by a code to apply colors:

| Code | Color/Format |
|---|---|
| `§a` | Green |
| `§c` | Red |
| `§e` | Yellow |
| `§f` | White |
| `§7` | Gray |
| `§b` | Aqua |
| `§d` | Magenta |
| `§r` | Reset |
| `§l` | Bold |
| `§o` | Italic |

```js
player.sendMessage("§aGreen §cRed §eYellow §r Normal");
```
