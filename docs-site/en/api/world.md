# World Object

The `world` object is available globally in every plugin. It allows interacting with the Minecraft world: placing blocks, spawning entities, managing particles, and working with online players.

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
| `spawnLightning(x, y, z)` | Strike lightning at the given position |
| `spawnTNT(x, y, z, fuse)` | Spawn a primed TNT. `fuse` is seconds until explosion |
| `spawnText(x, y, z, text)` | Spawn a floating text entity |
| `spawnExperienceOrb(x, y, z, amount)` | Spawn an experience orb with the given XP amount |

```js
// Strike lightning at coordinates
world.spawnLightning(100, 64, 100);

// Spawn TNT with 4 second fuse
world.spawnTNT(100, 65, 100, 4);

// Floating text (supports color codes)
world.spawnText(0, 70, 0, "§a§lWelcome to the server!");

// Drop experience orbs
world.spawnExperienceOrb(100, 64, 100, 50);
```

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
