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

## Entities

| Method | Returns | Description |
|---|---|---|
| `getEntities()` | `entity[]` | Get all entities in the world |
| `getEntitiesInRadius(x, y, z, radius)` | `entity[]` | Entities within the given spherical radius (in blocks) |
| `removeEntityByUUID(uuid)` | `boolean` | Remove an entity by UUID. Returns `true` if found and removed |

**Entity object returned:**

| Method | Available on | Description |
|---|---|---|
| `getUUID()` | All | Entity's unique UUID |
| `getType()` | All | Entity type (e.g. `"minecraft:item"`, `"minecraft:tnt"`, `"minecraft:arrow"`) |
| `getX/Y/Z()` | All | Position in the world |
| `remove()` | All | Remove entity from the world |
| `teleport(x, y, z)` | All (if supported) | Teleport the entity |
| `setVelocity(x, y, z)` | All (if supported) | Change velocity |
| `getHealth()` | Living (mobs) | Current health |
| `getMaxHealth()` | Living (mobs) | Maximum health |
| `setMaxHealth(n)` | Living (mobs) | Change max health |
| `isDead()` | Living (mobs) | Whether the entity is dead |
| `hurt(damage)` | Living (mobs) | Apply damage |
| `heal(health)` | Living (mobs) | Heal |
| `knockBack(x, y, z, force, height)` | Living (mobs) | Push back |
| `addEffect(name, level, seconds)` | Living (mobs) | Apply potion effect |
| `removeEffect(name)` | Living (mobs) | Remove effect |
| `clearEffects()` | Living (mobs) | Remove all effects |
| `getSpeed()` | Living (mobs) | Movement speed |
| `getName()` | Players | Player name |
| `sendMessage(msg)` | Players | Send message |
| `sendTitle(text, subtitle)` | Players | Show title |
| `disconnect(msg)` | Players | Disconnect |

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
| `"dragonfly:text"` | Floating text (created with `spawnText`) |

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
