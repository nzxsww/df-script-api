# World Object

The `world` object is available globally in every plugin. It allows interacting with the Minecraft world: placing blocks, spawning entities, managing particles, and working with online players.

::: warning Important — world transactions
The server uses a **transaction system** to access the world, similar to a database. Each command and each event have their own active transaction. If you try to open a second transaction from inside a first one, the server freezes (deadlock).

Methods that **read world state** (`getEntities`, `getEntitiesInRadius`, `getBlock`, etc.) **only work from inside an event or a command**:

```js
// ✅ Correct — inside an event
events.on("PlayerJoin", function(event) {
    var entities = world.getEntities(); // OK
});

// ✅ Correct — inside a command
commands.register("entities", "List entities", function(player, args) {
    var entities = world.getEntities(); // OK
});

// ❌ Incorrect — in onEnable() there is no active transaction
function onEnable() {
    var entities = world.getEntities(); // Deadlock / returns empty
}
```

Methods that **write** like `spawnEntity` or `setBlock` do work from `onEnable()` since they open their own transaction internally.
:::

## Blocks

| Method | Returns | Description |
|---|---|---|
| `setBlock(x, y, z, name)` | — | Place a block at the given position |
| `getBlock(x, y, z)` | `string` | Get the block name at the given position |
| `getHighestBlock(x, z)` | `number` | Get the Y of the highest non-air block at X,Z |

```js
world.setBlock(100, 64, 100, "minecraft:stone");
world.setBlock(100, 65, 100, "minecraft:grass_block");

var block = world.getBlock(100, 64, 100); // "minecraft:stone"
var y = world.getHighestBlock(100, 100);  // e.g. 65
```

## Particles

| Method | Description |
|---|---|
| `spawnParticle(x, y, z, name)` | Spawn a particle effect at the given position |

**Available particles:**

| Name | Description |
|---|---|
| `"flame"` | Fire/torch flame |
| `"lava"` | Lava bubble |
| `"water_drip"` | Water drip |
| `"lava_drip"` | Lava drip |
| `"explosion"` | Huge explosion |
| `"bone_meal"` | Bone meal usage |
| `"evaporate"` | Water evaporation |
| `"snowball"` | Snowball poof |
| `"egg_smash"` | Egg smash |
| `"entity_flame"` | Entity on fire |

```js
world.spawnParticle(100, 65, 100, "flame");
world.spawnParticle(100, 65, 100, "explosion");
```

## Entities

| Method | Description |
|---|---|
| `spawnEntity(type, x, y, z)` | Spawn an entity at the given position |
| `spawnEntity(type, x, y, z, options)` | Spawn an entity with additional options |

**Available types and options:**

| Type | Options | Description |
|---|---|---|
| `"lightning"` | — | Lightning bolt |
| `"tnt"` | `{ fuse: 4 }` | Primed TNT. `fuse` = seconds until explosion (default: 4) |
| `"text"` | `{ text: "§aHello" }` | Floating text |
| `"experience_orb"` | `{ amount: 50 }` | Experience orb with XP amount (default: 1) |
| `"item"` | `{ item: "minecraft:diamond", count: 1 }` | Item on the ground |

```js
world.spawnEntity("lightning", 100, 64, 100);
world.spawnEntity("tnt", 100, 65, 100, { fuse: 4 });
world.spawnEntity("text", 0, 70, 0, { text: "§a§lWelcome to the server!" });
world.spawnEntity("experience_orb", 100, 64, 100, { amount: 50 });
world.spawnEntity("item", 100, 64, 100, { item: "minecraft:diamond", count: 5 });
```

## Entities

### Spawning entities

| Method | Description |
|---|---|
| `spawnEntity(type, x, y, z)` | Spawn an entity at the given position |
| `spawnEntity(type, x, y, z, options)` | Spawn an entity with additional options |

**Available types and options:**

| Type | Options | Description |
|---|---|---|
| `"lightning"` | — | Lightning bolt |
| `"tnt"` | `{ fuse: 4 }` | Primed TNT. `fuse` = seconds until explosion (default: 4) |
| `"text"` | `{ text: "§aHello" }` | Floating text |
| `"experience_orb"` | `{ amount: 50 }` | Experience orb with XP amount (default: 1) |
| `"item"` | `{ item: "minecraft:diamond", count: 1 }` | Item on the ground |

```js
world.spawnEntity("lightning", x, y, z);
world.spawnEntity("tnt", x, y, z, { fuse: 4 });
world.spawnEntity("text", x, y, z, { text: "§a§lWelcome!" });
world.spawnEntity("experience_orb", x, y, z, { amount: 100 });
world.spawnEntity("item", x, y, z, { item: "minecraft:diamond", count: 5 });
```

### Getting and removing entities

| Method | Returns | Description |
|---|---|---|
| `getEntities()` | `entity[]` | Get all entities in the world |
| `getEntitiesInRadius(x, y, z, radius)` | `entity[]` | Entities within the given spherical radius (in blocks) |
| `removeEntityByUUID(uuid)` | `boolean` | Remove an entity by UUID. Returns `true` if found and removed |

Each returned entity is an object with its own methods. See the full reference at [Entity Object](/en/api/entity).

```js
// List all entities in the world
var entities = world.getEntities();
for (var i = 0; i < entities.length; i++) {
    var e = entities[i];
    console.log(e.getType() + " at " + e.getX() + "," + e.getY() + "," + e.getZ());
}

// Entities within 10 blocks of the player
events.on("PlayerJoin", function(event) {
    var p = event.getPlayer();
    var nearby = world.getEntitiesInRadius(p.getX(), p.getY(), p.getZ(), 10);
    p.sendMessage("§eNearby entities: §f" + nearby.length);
    for (var i = 0; i < nearby.length; i++) {
        var e = nearby[i];
        if (typeof e.getHealth === "function") {
            p.sendMessage("§7- §f" + e.getType() + " §7health=" + e.getHealth().toFixed(1));
        } else {
            p.sendMessage("§7- §f" + e.getType());
        }
    }
});

// Remove entity by UUID
world.removeEntityByUUID("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx");

// Remove all TNT entities
var entities = world.getEntities();
for (var i = 0; i < entities.length; i++) {
    if (entities[i].getType() === "minecraft:tnt") {
        entities[i].remove();
    }
}
```

**Common entity types:**

| Type | Description |
|---|---|
| `"minecraft:item"` | Item dropped on the ground |
| `"minecraft:tnt"` | Primed TNT |
| `"minecraft:arrow"` | Arrow |
| `"minecraft:xp_orb"` | Experience orb |
| `"minecraft:lightning_bolt"` | Lightning |
| `"minecraft:falling_block"` | Falling block |
| `"minecraft:fireworks_rocket"` | Firework |
| `"dragonfly:text"` | Floating text (created with `world.spawnEntity("text", ...)`) |

## Players

| Method | Returns | Description |
|---|---|---|
| `getPlayers()` | `player[]` | Get all online players as player objects |
| `getPlayerCount()` | `number` | Number of players currently online |
| `broadcast(msg)` | — | Send a message to all online players |

```js
// Broadcast a message to everyone
world.broadcast("§aThe server will restart in 5 minutes!");

// Iterate over all players
var players = world.getPlayers();
for (var i = 0; i < players.length; i++) {
    players[i].sendTitle("§aEvent!", "§7Starting now");
    players[i].addEffect("speed", 2, 60);
}

// Player count
var count = world.getPlayerCount();
world.broadcast("§7There are §f" + count + " §7players online.");
```
