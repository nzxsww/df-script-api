# Eventos

Los eventos se registran dentro de `onEnable()` usando `events.on()`:

```js
events.on("NombreEvento", function(event) {
    // hacer algo con event
});
```

::: tip
Los nombres de eventos usan **PascalCase** sin sufijo: `"PlayerJoin"`, `"BlockBreak"`.
NO `"PlayerJoinEvent"`, NO `"player_join"`.
:::

## Eventos de jugador

### PlayerJoin

Se dispara cuando un jugador entra al servidor.

```js
events.on("PlayerJoin", function(event) {
    var player = event.getPlayer();
    player.sendMessage("¡Bienvenido!");

    // Personalizar el mensaje de entrada
    event.setJoinMessage("§e" + player.getName() + " entró al servidor.");

    // Cancelar el mensaje de entrada (no aparece en el chat)
    event.setCancelled(true);
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador que entró |
| `getJoinMessage()` | `string` | Mensaje de entrada actual |
| `setJoinMessage(msg)` | — | Cambiar el mensaje de entrada |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar/descancelar el evento |

---

### PlayerQuit

Se dispara cuando un jugador sale del servidor.

```js
events.on("PlayerQuit", function(event) {
    var player = event.getPlayer();
    event.setQuitMessage("§c" + player.getName() + " salió del servidor.");
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador que salió |
| `getQuitMessage()` | `string` | Mensaje de salida actual |
| `setQuitMessage(msg)` | — | Cambiar el mensaje de salida |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar el mensaje de salida |

---

### PlayerChat

Se dispara cuando un jugador envía un mensaje en el chat. El chat **siempre** es controlado manualmente — si no cancelás el evento, el mensaje se envía con el formato modificado.

```js
events.on("PlayerChat", function(event) {
    var player = event.getPlayer();
    var msg = event.getMessage();

    // Modificar el formato del chat
    event.setMessage("§7[§f" + player.getName() + "§7] §r" + msg);

    // Cancelar el mensaje (no se envía a nadie)
    event.setCancelled(true);
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador que escribió |
| `getMessage()` | `string` | El mensaje original |
| `setMessage(msg)` | — | Cambiar el mensaje |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar el envío del mensaje |

---

### PlayerMove

Se dispara cuando un jugador se mueve. **Cuidado:** se dispara muy frecuentemente.

```js
events.on("PlayerMove", function(event) {
    var y = event.getToY();

    // Evitar que el jugador caiga al vacío
    if (y < 0) {
        event.setCancelled(true);
    }
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador que se movió |
| `getFromX/Y/Z()` | `number` | Posición de origen |
| `getToX/Y/Z()` | `number` | Posición de destino |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar el movimiento |

---

### PlayerJump

Se dispara cuando un jugador salta.

```js
events.on("PlayerJump", function(event) {
    var player = event.getPlayer();
    // Cancelar para impedir saltar
    event.setCancelled(true);
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador que saltó |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar el salto |

---

### PlayerTeleport

Se dispara cuando un jugador es teleportado.

```js
events.on("PlayerTeleport", function(event) {
    // Redirigir la teleportación a una posición fija
    event.setPosition(0, 64, 0);
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador teleportado |
| `getX/Y/Z()` | `number` | Posición de destino |
| `setPosition(x, y, z)` | — | Cambiar la posición de destino |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar la teleportación |

---

### PlayerDeath

Se dispara cuando un jugador muere.

```js
events.on("PlayerDeath", function(event) {
    var player = event.getPlayer();

    // Hacer que el jugador conserve su inventario al morir
    event.setKeepInventory(true);

    player.sendMessage("§cMoriste, pero conservaste tu inventario.");
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador que murió |
| `getKeepInventory()` | `boolean` | Si conserva el inventario |
| `setKeepInventory(bool)` | — | Cambiar si conserva el inventario |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar la muerte |

---

### PlayerRespawn

Se dispara cuando un jugador reaparece tras morir.

```js
events.on("PlayerRespawn", function(event) {
    // Cambiar el punto de reaparición
    event.setPosition(0, 64, 0);
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador que reapareció |
| `getX/Y/Z()` | `number` | Posición de reaparición |
| `setPosition(x, y, z)` | — | Cambiar la posición de reaparición |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar la reaparición |

---

### PlayerHurt

Se dispara cuando un jugador recibe daño.

```js
events.on("PlayerHurt", function(event) {
    var player = event.getPlayer();

    // Reducir el daño a la mitad
    event.setDamage(event.getDamage() / 2);

    // O cancelar el daño completamente (modo god)
    event.setCancelled(true);
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador que recibió daño |
| `getDamage()` | `number` | Cantidad de daño |
| `setDamage(n)` | — | Cambiar la cantidad de daño |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar el daño |

---

### PlayerHeal

Se dispara cuando un jugador recupera vida.

```js
events.on("PlayerHeal", function(event) {
    // Duplicar la curación
    event.setHealth(event.getHealth() * 2);
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador que se curó |
| `getHealth()` | `number` | Cantidad de curación |
| `setHealth(n)` | — | Cambiar la cantidad de curación |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar la curación |

---

### PlayerFoodLoss

Se dispara cuando un jugador pierde hambre.

```js
events.on("PlayerFoodLoss", function(event) {
    // Evitar que el jugador pierda hambre
    event.setCancelled(true);
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador |
| `getFrom()` | `number` | Nivel de hambre anterior (0-20) |
| `getTo()` | `number` | Nuevo nivel de hambre |
| `setTo(n)` | — | Cambiar el nuevo nivel de hambre |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar la pérdida de hambre |

---

### PlayerExperienceGain

Se dispara cuando un jugador gana experiencia.

```js
events.on("PlayerExperienceGain", function(event) {
    // Duplicar la experiencia ganada
    event.setAmount(event.getAmount() * 2);
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador |
| `getAmount()` | `number` | Cantidad de experiencia |
| `setAmount(n)` | — | Cambiar la cantidad |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar la experiencia |

---

### PlayerToggleSprint

Se dispara cuando un jugador activa o desactiva el sprint.

```js
events.on("PlayerToggleSprint", function(event) {
    if (event.isSprinting()) {
        // El jugador empezó a correr
        event.setCancelled(true); // Impedir el sprint
    }
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador |
| `isSprinting()` | `boolean` | `true` si activó el sprint, `false` si lo desactivó |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar el cambio de sprint |

---

### PlayerToggleSneak

Se dispara cuando un jugador activa o desactiva el agachado.

```js
events.on("PlayerToggleSneak", function(event) {
    if (event.isSneaking()) {
        var player = event.getPlayer();
        player.sendMessage("§7Te agachaste.");
    }
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador |
| `isSneaking()` | `boolean` | `true` si se agachó, `false` si se levantó |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar el cambio |

---

### PlayerAttackEntity

Se dispara cuando un jugador ataca a una entidad.

```js
events.on("PlayerAttackEntity", function(event) {
    if (event.isCritical()) {
        var player = event.getPlayer();
        player.sendMessage("§c¡Golpe crítico!");
    }
    // Cancelar el ataque
    event.setCancelled(true);
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador que atacó |
| `getForce()` | `number` | Fuerza del golpe |
| `isCritical()` | `boolean` | Si fue un golpe crítico |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar el ataque |

---

### PlayerItemUse

Se dispara cuando un jugador usa un item (click derecho).

```js
events.on("PlayerItemUse", function(event) {
    event.setCancelled(true); // Impedir usar items
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar el uso del item |

---

### PlayerItemDrop

Se dispara cuando un jugador tira un item al suelo.

```js
events.on("PlayerItemDrop", function(event) {
    // Evitar tirar items
    event.setCancelled(true);
    event.getPlayer().sendMessage("§cNo podés tirar items aquí.");
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador |
| `getItemCount()` | `number` | Cantidad de items tirados |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar el tirar |

---

### PlayerItemPickup

Se dispara cuando un jugador recoge un item del suelo.

```js
events.on("PlayerItemPickup", function(event) {
    // Evitar recoger items
    event.setCancelled(true);
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador |
| `getItemCount()` | `number` | Cantidad de items recogidos |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar recoger |

---

## Eventos de bloque

### BlockBreak

Se dispara cuando un jugador rompe un bloque.

```js
events.on("BlockBreak", function(event) {
    var player = event.getPlayer();
    var y = event.getBlockY();

    if (y < 10) {
        event.setCancelled(true);
        player.sendMessage("§cNo podés romper bloques tan abajo.");
    }
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador que rompió |
| `getBlockX/Y/Z()` | `number` | Posición del bloque |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar la rotura |

---

### BlockPlace

Se dispara cuando un jugador coloca un bloque.

```js
events.on("BlockPlace", function(event) {
    var y = event.getBlockY();
    if (y > 200) {
        event.setCancelled(true);
        event.getPlayer().sendMessage("§cNo podés construir tan alto.");
    }
});
```

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayer()` | `player` | El jugador que colocó |
| `getBlockX/Y/Z()` | `number` | Posición del bloque |
| `isCancelled()` | `boolean` | Si el evento está cancelado |
| `setCancelled(bool)` | — | Cancelar la colocación |
