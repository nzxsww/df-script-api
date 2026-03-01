# Events

Events are registered inside `onEnable()` using `events.on()`:

```js
events.on("EventName", function(event) {
    // do something with event
});
```

::: tip
Event names use **PascalCase** without suffix: `"PlayerJoin"`, `"BlockBreak"`.
NOT `"PlayerJoinEvent"`, NOT `"player_join"`.
:::

## Player Events

### PlayerJoin

Fires when a player joins the server.

```js
events.on("PlayerJoin", function(event) {
    var player = event.getPlayer();
    player.sendMessage("Welcome!");
    event.setJoinMessage("§e" + player.getName() + " joined the server.");
});
```

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The player who joined |
| `getJoinMessage()` | `string` | Current join message |
| `setJoinMessage(msg)` | — | Change the join message |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel/uncancel the event |

---

### PlayerQuit

Fires when a player leaves the server.

```js
events.on("PlayerQuit", function(event) {
    var player = event.getPlayer();
    event.setQuitMessage("§c" + player.getName() + " left the server.");
});
```

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The player who left |
| `getQuitMessage()` | `string` | Current quit message |
| `setQuitMessage(msg)` | — | Change the quit message |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the quit message |

---

### PlayerChat

Fires when a player sends a chat message.

```js
events.on("PlayerChat", function(event) {
    var player = event.getPlayer();
    var msg = event.getMessage();
    event.setMessage("§7[§f" + player.getName() + "§7] §r" + msg);
});
```

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The player who wrote |
| `getMessage()` | `string` | The original message |
| `setMessage(msg)` | — | Change the message |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the message |

---

### PlayerMove

Fires when a player moves. **Warning:** fires very frequently.

```js
events.on("PlayerMove", function(event) {
    if (event.getToY() < 0) {
        event.setCancelled(true);
    }
});
```

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The player who moved |
| `getFromX/Y/Z()` | `number` | Origin position |
| `getToX/Y/Z()` | `number` | Destination position |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the movement |

---

### PlayerJump

Fires when a player jumps.

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The player who jumped |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the jump |

---

### PlayerTeleport

Fires when a player is teleported.

```js
events.on("PlayerTeleport", function(event) {
    event.setPosition(0, 64, 0);
});
```

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The teleported player |
| `getX/Y/Z()` | `number` | Destination position |
| `setPosition(x, y, z)` | — | Change destination |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the teleport |

---

### PlayerDeath

Fires when a player dies.

```js
events.on("PlayerDeath", function(event) {
    event.setKeepInventory(true);
    event.getPlayer().sendMessage("§cYou died but kept your inventory.");
});
```

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The player who died |
| `getKeepInventory()` | `boolean` | Whether inventory is kept |
| `setKeepInventory(bool)` | — | Set inventory keeping |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the death |

---

### PlayerRespawn

Fires when a player respawns after dying.

```js
events.on("PlayerRespawn", function(event) {
    event.setPosition(0, 64, 0);
});
```

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The respawned player |
| `getX/Y/Z()` | `number` | Respawn position |
| `setPosition(x, y, z)` | — | Change respawn position |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the respawn |

---

### PlayerHurt

Fires when a player takes damage.

```js
events.on("PlayerHurt", function(event) {
    event.setDamage(event.getDamage() / 2); // halve damage
    // or: event.setCancelled(true); // god mode
});
```

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The hurt player |
| `getDamage()` | `number` | Damage amount |
| `setDamage(n)` | — | Change damage amount |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the damage |

---

### PlayerHeal

Fires when a player heals.

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The healed player |
| `getHealth()` | `number` | Heal amount |
| `setHealth(n)` | — | Change heal amount |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the heal |

---

### PlayerFoodLoss

Fires when a player loses hunger.

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The player |
| `getFrom()` | `number` | Previous hunger level (0-20) |
| `getTo()` | `number` | New hunger level |
| `setTo(n)` | — | Change new hunger level |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the hunger loss |

---

### PlayerExperienceGain

Fires when a player gains experience.

```js
events.on("PlayerExperienceGain", function(event) {
    event.setAmount(event.getAmount() * 2); // double XP
});
```

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The player |
| `getAmount()` | `number` | XP amount |
| `setAmount(n)` | — | Change XP amount |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the XP gain |

---

### PlayerToggleSprint

Fires when a player toggles sprinting.

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The player |
| `isSprinting()` | `boolean` | `true` if started sprinting |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the toggle |

---

### PlayerToggleSneak

Fires when a player toggles sneaking.

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The player |
| `isSneaking()` | `boolean` | `true` if started sneaking |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the toggle |

---

### PlayerAttackEntity

Fires when a player attacks an entity.

```js
events.on("PlayerAttackEntity", function(event) {
    if (event.isCritical()) {
        event.getPlayer().sendMessage("§cCritical hit!");
    }
});
```

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The attacking player |
| `getForce()` | `number` | Attack force |
| `isCritical()` | `boolean` | Whether it was a critical hit |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the attack |

---

### PlayerItemUse

Fires when a player uses an item (right-click).

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The player |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel item use |

---

### PlayerItemDrop

Fires when a player drops an item.

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The player |
| `getItemCount()` | `number` | Number of items dropped |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the drop |

---

### PlayerItemPickup

Fires when a player picks up an item.

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The player |
| `getItemCount()` | `number` | Number of items picked up |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the pickup |

---

## Block Events

### BlockBreak

Fires when a player breaks a block.

```js
events.on("BlockBreak", function(event) {
    if (event.getBlockY() < 10) {
        event.setCancelled(true);
        event.getPlayer().sendMessage("§cCan't break blocks here.");
    }
});
```

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The player who broke |
| `getBlockX/Y/Z()` | `number` | Block position |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the break |

---

### BlockPlace

Fires when a player places a block.

| Method | Returns | Description |
|---|---|---|
| `getPlayer()` | `player` | The player who placed |
| `getBlockX/Y/Z()` | `number` | Block position |
| `isCancelled()` | `boolean` | Whether event is cancelled |
| `setCancelled(bool)` | — | Cancel the placement |
