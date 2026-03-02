# Player Object

The `player` object is available in all events (via `event.getPlayer()`) and in command callbacks. It has 50+ methods to interact with the player.

```js
events.on("PlayerJoin", function(event) {
    var player = event.getPlayer();
    player.sendMessage("Hello " + player.getName());
});
```

## Identity

| Method | Returns | Description |
|---|---|---|
| `getName()` | `string` | Player name (gamertag) |
| `getUUID()` | `string` | Unique player UUID |
| `getXUID()` | `string` | Xbox Live XUID |
| `getNameTag()` | `string` | Tag visible above head |
| `setNameTag(tag)` | — | Change visible tag |

## Messages

| Method | Description |
|---|---|
| `sendMessage(msg)` | Send message to player chat |
| `sendPopup(msg)` | Show popup in HUD (above item bar) |
| `sendTip(msg)` | Show tip in HUD (above health bar) |
| `sendToast(title, msg)` | Show toast notification |
| `sendJukeboxPopup(msg)` | Jukebox-style popup |

## Connection

| Method | Returns | Description |
|---|---|---|
| `disconnect(msg)` | — | Kick the player with a message |
| `transfer(address)` | — | Transfer to another server (`"ip:port"`) |
| `getLatency()` | `number` | Latency in milliseconds |

## Position & Movement

| Method | Returns | Description |
|---|---|---|
| `getX()` | `number` | X coordinate |
| `getY()` | `number` | Y coordinate |
| `getZ()` | `number` | Z coordinate |
| `teleport(x, y, z)` | — | Teleport the player |
| `setVelocity(x, y, z)` | — | Change velocity/impulse |

## Physical State

| Method | Returns | Description |
|---|---|---|
| `getHealth()` | `number` | Current health |
| `getMaxHealth()` | `number` | Max health |
| `setMaxHealth(n)` | — | Change max health |
| `getFoodLevel()` | `number` | Hunger level (0-20) |
| `setFoodLevel(n)` | — | Change hunger |
| `isOnGround()` | `boolean` | Whether on ground |
| `isSneaking()` | `boolean` | Whether sneaking |
| `isSprinting()` | `boolean` | Whether sprinting |
| `isFlying()` | `boolean` | Whether flying |
| `isSwimming()` | `boolean` | Whether swimming |
| `isDead()` | `boolean` | Whether dead |
| `isImmobile()` | `boolean` | Whether immobile |

## Experience

| Method | Returns | Description |
|---|---|---|
| `getExperience()` | `number` | Total experience points |
| `getExperienceLevel()` | `number` | Experience level |
| `addExperience(n)` | — | Add experience points |
| `setExperienceLevel(n)` | — | Set experience level |

## Game Mode

| Method | Returns | Description |
|---|---|---|
| `getGameMode()` | `string` | Current mode (`"survival"`, `"creative"`, `"adventure"`, `"spectator"`) |
| `setGameMode(mode)` | — | Change game mode |

## Flight

::: warning
`startFlying()` and `stopFlying()` only work in **creative** or **spectator** mode. In survival they are silently ignored. If you need flight in survival, change the gamemode first.
:::

| Method | Returns | Description |
|---|---|---|
| `isFlying()` | `boolean` | Whether flying |
| `startFlying()` | — | Enable flight |
| `stopFlying()` | — | Disable flight |

## Visual Effects

| Method | Returns | Description |
|---|---|---|
| `isInvisible()` | `boolean` | Whether invisible |
| `setInvisible()` | — | Make invisible |
| `setVisible()` | — | Make visible |

## Speed

| Method | Returns | Description |
|---|---|---|
| `getSpeed()` | `number` | Current movement speed |
| `setSpeed(n)` | — | Change speed (default: `0.1`) |

## Inventory

| Method | Returns | Description |
|---|---|---|
| `giveItem(name, count)` | `boolean` | Give item. Returns `false` if unknown or inventory full |
| `clearInventory()` | — | Clear the entire inventory |
| `getItemCount(name)` | `number` | How many of that item the player has |

```js
player.giveItem("minecraft:diamond", 10);
var diamonds = player.getItemCount("minecraft:diamond");
player.clearInventory();
```

## Sounds

| Method | Description |
|---|---|
| `playSound(name)` | Play a sound to the player |

**Available sounds:**

| Name | When to use |
|---|---|
| `"click"` | Interaction, selection |
| `"levelup"` | Achievement, level up |
| `"pop"` | Item received, confirmation |
| `"burp"` | Eating |
| `"deny"` | Error, action denied |
| `"door_open"` | Open door |
| `"door_close"` | Close door |
| `"chest_open"` | Open chest |
| `"chest_close"` | Close chest |
| `"anvil_land"` | Heavy object, impact |
| `"bow_shoot"` | Shoot |
| `"arrow_hit"` | Arrow impact |

## Titles

| Method | Description |
|---|---|
| `sendTitle(text, subtitle)` | Show a large title in the center of the screen with an optional subtitle |

```js
player.sendTitle("§aWelcome!", "§7Have a great time");
```

## Potion Effects

| Method | Description |
|---|---|
| `addEffect(name, level, seconds)` | Apply a potion effect. `level` starts at 1 (level I) |
| `removeEffect(name)` | Remove a specific effect |
| `clearEffects()` | Remove all active effects |

**Available effects:**

| Name | Description |
|---|---|
| `"speed"` | Move faster |
| `"slowness"` | Move slower |
| `"haste"` | Mine faster |
| `"mining_fatigue"` | Mine slower |
| `"strength"` | Deal more damage |
| `"jump_boost"` | Jump higher |
| `"nausea"` | Distorted vision |
| `"regeneration"` | Regenerate health |
| `"resistance"` | Take less damage |
| `"fire_resistance"` | Immune to fire |
| `"water_breathing"` | Breathe underwater |
| `"invisibility"` | Invisible to others |
| `"blindness"` | Reduced vision |
| `"night_vision"` | See in the dark |
| `"hunger"` | Lose hunger faster |
| `"weakness"` | Deal less damage |
| `"poison"` | Lose health over time |
| `"wither"` | Lose health (bypasses armour) |
| `"health_boost"` | Increase max health |
| `"absorption"` | Extra health points |
| `"saturation"` | Restore food instantly |
| `"levitation"` | Float upward |
| `"slow_falling"` | Fall slowly |
| `"conduit_power"` | Underwater haste + vision |
| `"darkness"` | Reduced vision (like blindness) |

```js
player.addEffect("speed", 2, 30);      // Speed II for 30 seconds
player.addEffect("regeneration", 1, 10); // Regeneration I for 10 seconds
player.removeEffect("speed");
player.clearEffects();
```

## Armour

| Method | Returns | Description |
|---|---|---|
| `setArmour(slot, itemName)` | — | Equip an armour piece. Slot: `0`=helmet, `1`=chestplate, `2`=leggings, `3`=boots |
| `getArmour(slot)` | `string` | Get the item name in that armour slot. Returns `""` if empty |
| `clearArmour()` | — | Remove all equipped armour |

```js
player.setArmour(0, "minecraft:diamond_helmet");
player.setArmour(1, "minecraft:diamond_chestplate");
player.setArmour(2, "minecraft:diamond_leggings");
player.setArmour(3, "minecraft:diamond_boots");

var helmet = player.getArmour(0); // "minecraft:diamond_helmet"
player.clearArmour();
```

## Scoreboard

| Method | Description |
|---|---|
| `sendScoreboard(sb)` | Send a scoreboard to the player (created with `scoreboard.create()`) |
| `removeScoreboard()` | Remove the scoreboard from the player's screen |

```js
var sb = scoreboard.create("§6My Server");
sb.setLine(0, "§7Players: 5");
sb.setLine(1, "§7Map: Lobby");
player.sendScoreboard(sb);

// Later...
player.removeScoreboard();
```

See the [Scoreboard API](/en/api/scoreboard) for full documentation including ScoreboardManager and Live Scoreboard.

## Commands

| Method | Description |
|---|---|
| `executeCommand(cmd)` | Execute a command as the player |
