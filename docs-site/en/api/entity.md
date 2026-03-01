# Entity Object

The `entity` object represents an individual entity in the world (item, TNT, arrow, mob, player, floating text, etc.). It is obtained through `world` methods or from entity events.

```js
// Get entities via world
var entities = world.getEntities();
var entity = entities[0];

// Get nearby entities
var nearby = world.getEntitiesInRadius(x, y, z, 10);
```

## Base methods (all entities)

These methods are available on **any** type of entity.

| Method | Returns | Description |
|---|---|---|
| `getUUID()` | `string` | Entity's unique UUID |
| `getType()` | `string` | Entity type (e.g. `"minecraft:item"`, `"minecraft:tnt"`) |
| `getX()` | `number` | X position |
| `getY()` | `number` | Y position |
| `getZ()` | `number` | Z position |
| `remove()` | — | Remove the entity from the world |
| `teleport(x, y, z)` | — | Teleport the entity (if supported) |
| `setVelocity(x, y, z)` | — | Change the entity's velocity (if supported) |

```js
var entities = world.getEntities();
for (var i = 0; i < entities.length; i++) {
    var e = entities[i];
    console.log(e.getType() + " UUID=" + e.getUUID());
    console.log("  pos: " + e.getX().toFixed(1) + "," + e.getY().toFixed(1) + "," + e.getZ().toFixed(1));
}
```

## Living entity methods (mobs)

Available when `typeof entity.getHealth === "function"`. Applies to mobs, animals and players.

| Method | Returns | Description |
|---|---|---|
| `getHealth()` | `number` | Current health |
| `getMaxHealth()` | `number` | Maximum health |
| `setMaxHealth(n)` | — | Change max health |
| `isDead()` | `boolean` | Whether the entity is dead |
| `hurt(damage)` | — | Apply attack damage |
| `heal(health)` | — | Heal the entity |
| `knockBack(x, y, z, force, height)` | — | Push the entity from a source position |
| `getSpeed()` | `number` | Movement speed |
| `addEffect(name, level, seconds)` | — | Apply a potion effect |
| `removeEffect(name)` | — | Remove an active effect |
| `clearEffects()` | — | Remove all active effects |

```js
var nearby = world.getEntitiesInRadius(x, y, z, 10);
for (var i = 0; i < nearby.length; i++) {
    var e = nearby[i];
    if (typeof e.getHealth === "function") {
        console.log(e.getType() + " health=" + e.getHealth().toFixed(1) + "/" + e.getMaxHealth().toFixed(1));
        e.addEffect("slowness", 1, 10);
        e.hurt(5);
    }
}
```

## Player-only methods

Available when `typeof entity.getName === "function"`. Players also have all living entity methods.

| Method | Returns | Description |
|---|---|---|
| `getName()` | `string` | Player name |
| `sendMessage(msg)` | — | Send a chat message |
| `sendTitle(text, subtitle)` | — | Show a title on screen |
| `disconnect(msg)` | — | Disconnect the player |

```js
var entities = world.getEntities();
for (var i = 0; i < entities.length; i++) {
    var e = entities[i];
    if (typeof e.getName === "function") {
        // It's a player
        e.sendMessage("§aHello " + e.getName() + "!");
        e.sendTitle("§aEvent", "§7Starting...");
    }
}
```

## Common entity types

| Type | Description | Living |
|---|---|---|
| `"minecraft:item"` | Item dropped on the ground | No |
| `"minecraft:tnt"` | Primed TNT | No |
| `"minecraft:arrow"` | Arrow | No |
| `"minecraft:xp_orb"` | Experience orb | No |
| `"minecraft:lightning_bolt"` | Lightning | No |
| `"minecraft:falling_block"` | Falling block | No |
| `"minecraft:fireworks_rocket"` | Firework | No |
| `"dragonfly:text"` | Floating text (created with `world.spawnEntity("text", ...)`) | No |

## Checking type before using optional methods

```js
var entities = world.getEntities();
for (var i = 0; i < entities.length; i++) {
    var e = entities[i];

    // Always available
    console.log(e.getType() + " at " + e.getX().toFixed(0) + "," + e.getY().toFixed(0) + "," + e.getZ().toFixed(0));

    // Only on mobs/players
    if (typeof e.getHealth === "function") {
        console.log("  health: " + e.getHealth().toFixed(1));
    }

    // Only on players
    if (typeof e.getName === "function") {
        console.log("  player: " + e.getName());
    }
}
```

## Removing entities

```js
// Remove by UUID (from world)
world.removeEntityByUUID("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx");

// Remove from the entity object directly
var entities = world.getEntities();
for (var i = 0; i < entities.length; i++) {
    if (entities[i].getType() === "minecraft:tnt") {
        entities[i].remove(); // Remove all TNT
    }
}

// Remove item entities within 5 blocks
var nearby = world.getEntitiesInRadius(x, y, z, 5);
for (var i = 0; i < nearby.length; i++) {
    if (nearby[i].getType() === "minecraft:item") {
        nearby[i].remove();
    }
}
```
