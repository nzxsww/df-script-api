# Getting Started

## Requirements

- **Go 1.21+** installed
- The server compiled (`go build -o server.exe .`) **or** downloaded from Releases:
  - https://github.com/nzxsww/df-script-api/releases
- A text editor to write JS

## Create Your First Plugin

### 1. Create the plugin folder

Inside the `plugins/` folder of the server, create a folder with your plugin name:

```
plugins/
└── my-first-plugin/    ← create this folder
```

### 2. Create `plugin.yml`

Inside that folder, create the `plugin.yml` file with your plugin information:

```yaml
name: MyFirstPlugin
version: 1.0.0
author: YourName
description: My first plugin for Dragonfly Script API
main: index.js
api-version: 1.0.0
```

| Field | Description |
|---|---|
| `name` | Unique plugin name. Appears in server logs. |
| `version` | Version in semver format (1.0.0) |
| `author` | Your name or nickname |
| `description` | Short description of what the plugin does |
| `main` | Main JS file (default: `index.js`) |
| `api-version` | API version (use `1.0.0`) |

### 3. Create `index.js`

Create the `index.js` file in the same folder:

```js
// This code runs when the plugin is loaded (before onEnable)
console.log("Loading: " + plugin.name + " v" + plugin.version);

function onEnable() {
    // Register your events and commands here
    console.log("Plugin enabled!");

    events.on("PlayerJoin", function(event) {
        var player = event.getPlayer();
        player.sendMessage("§aHello " + player.getName() + ", welcome!");
    });
}

function onDisable() {
    // Cleanup when the server closes
    console.log("Plugin disabled.");
}

// Export the lifecycle (required)
module = {
    onEnable: onEnable,
    onDisable: onDisable
};
```

### 4. Start the server

```bash
./server.exe
```

You should see in the console:

```
[Loader] Loading plugin: MyFirstPlugin v1.0.0 (author: YourName)
[MyFirstPlugin] Loading: MyFirstPlugin v1.0.0
[MyFirstPlugin] Plugin enabled!
```

### 5. Connect with Minecraft

Connect with **Minecraft Bedrock 1.21.130 - 1.21.132** to `localhost:19132`. When you join, you'll receive the welcome message.

::: tip
To connect from the same PC, use the address `127.0.0.1` port `19132`.
:::

## Resulting file structure

```
plugins/
└── my-first-plugin/
    ├── plugin.yml     ← plugin metadata
    ├── index.js       ← plugin logic
    └── config.yml     ← created automatically if you use config.save()
```

## Next steps

- [Detailed plugin structure →](./plugin-structure)
- [View all available events →](/en/api/events)
- [View all player methods →](/en/api/player)
- [Register commands →](/en/api/commands)
