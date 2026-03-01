# Objeto Player

El objeto `player` está disponible en todos los eventos (via `event.getPlayer()`) y en los callbacks de comandos. Tiene más de 50 métodos para interactuar con el jugador.

```js
events.on("PlayerJoin", function(event) {
    var player = event.getPlayer();
    player.sendMessage("Hola " + player.getName());
});
```

## Identidad

| Método | Retorna | Descripción |
|---|---|---|
| `getName()` | `string` | Nombre del jugador (gamertag) |
| `getUUID()` | `string` | UUID único del jugador |
| `getXUID()` | `string` | XUID de Xbox Live |
| `getNameTag()` | `string` | Tag visible sobre la cabeza |
| `setNameTag(tag)` | — | Cambiar el tag visible |

```js
player.setNameTag("§c[Admin] §f" + player.getName());
```

## Mensajes

| Método | Descripción |
|---|---|
| `sendMessage(msg)` | Enviar mensaje al chat del jugador |
| `sendPopup(msg)` | Mostrar popup en el HUD (sobre la barra de items) |
| `sendTip(msg)` | Mostrar tip en el HUD (sobre la barra de vida) |
| `sendToast(titulo, msg)` | Mostrar notificación tipo toast en la esquina |
| `sendJukeboxPopup(msg)` | Popup estilo jukebox |

```js
player.sendMessage("§aMensaje en el chat");
player.sendPopup("§eTexto en el HUD");
player.sendTip("§bTip rápido");
player.sendToast("§l¡Logro!", "§rCompletaste una misión");
```

## Conexión

| Método | Retorna | Descripción |
|---|---|---|
| `disconnect(msg)` | — | Kickear al jugador con un mensaje |
| `transfer(dirección)` | — | Transferir a otro servidor (`"ip:puerto"`) |
| `getLatency()` | `number` | Latencia en milisegundos |

```js
player.disconnect("§cFuiste baneado.");
player.transfer("otro-servidor.com:19132");
console.log("Latencia: " + player.getLatency() + "ms");
```

## Posición y movimiento

| Método | Retorna | Descripción |
|---|---|---|
| `getX()` | `number` | Coordenada X |
| `getY()` | `number` | Coordenada Y |
| `getZ()` | `number` | Coordenada Z |
| `teleport(x, y, z)` | — | Teleportar al jugador |
| `setVelocity(x, y, z)` | — | Cambiar velocidad/impulso |

```js
var x = player.getX();
var y = player.getY();
var z = player.getZ();
player.sendMessage("Estás en " + Math.floor(x) + ", " + Math.floor(y) + ", " + Math.floor(z));

player.teleport(0, 64, 0); // teleportar al spawn
player.setVelocity(0, 1, 0); // lanzar al jugador hacia arriba
```

## Estado físico

| Método | Retorna | Descripción |
|---|---|---|
| `getHealth()` | `number` | Vida actual |
| `getMaxHealth()` | `number` | Vida máxima |
| `setMaxHealth(n)` | — | Cambiar la vida máxima |
| `getFoodLevel()` | `number` | Nivel de hambre (0-20) |
| `setFoodLevel(n)` | — | Cambiar el hambre |
| `isOnGround()` | `boolean` | Si está tocando el suelo |
| `isSneaking()` | `boolean` | Si está agachado |
| `isSprinting()` | `boolean` | Si está corriendo |
| `isFlying()` | `boolean` | Si está volando |
| `isSwimming()` | `boolean` | Si está nadando |
| `isDead()` | `boolean` | Si está muerto |
| `isImmobile()` | `boolean` | Si está inmóvil |

```js
var hp = player.getHealth();
var maxHp = player.getMaxHealth();
player.sendMessage("§cVida: §f" + Math.floor(hp) + "/" + Math.floor(maxHp));
player.setFoodLevel(20); // llenar el hambre
```

## Experiencia

| Método | Retorna | Descripción |
|---|---|---|
| `getExperience()` | `number` | Puntos de experiencia totales |
| `getExperienceLevel()` | `number` | Nivel de experiencia |
| `addExperience(n)` | — | Agregar puntos de experiencia |
| `setExperienceLevel(n)` | — | Establecer el nivel de experiencia |

```js
player.addExperience(100);
player.setExperienceLevel(30);
console.log("Nivel: " + player.getExperienceLevel());
```

## Modo de juego

| Método | Retorna | Descripción |
|---|---|---|
| `getGameMode()` | `string` | Modo actual (`"survival"`, `"creative"`, `"adventure"`, `"spectator"`) |
| `setGameMode(modo)` | — | Cambiar el modo de juego |

```js
if (player.getGameMode() === "survival") {
    player.setGameMode("creative");
    player.sendMessage("§aModo creativo activado.");
}
```

## Vuelo

::: warning
`startFlying()` y `stopFlying()` solo funcionan en modos **creative** o **spectator**. En survival son ignorados silenciosamente. Si necesitás vuelo en survival, cambiá el gamemode primero.
:::

| Método | Retorna | Descripción |
|---|---|---|
| `isFlying()` | `boolean` | Si está volando |
| `startFlying()` | — | Activar vuelo |
| `stopFlying()` | — | Desactivar vuelo |

```js
// Activar vuelo en cualquier modo
if (player.getGameMode() === "survival") {
    player.setGameMode("creative");
}
player.startFlying();
```

## Efectos visuales

| Método | Retorna | Descripción |
|---|---|---|
| `isInvisible()` | `boolean` | Si está invisible |
| `setInvisible()` | — | Hacer invisible |
| `setVisible()` | — | Hacer visible |

```js
player.setInvisible();
setTimeout(function() {
    player.setVisible();
}, 5000); // visible de nuevo después de 5 segundos
```

## Velocidad

| Método | Retorna | Descripción |
|---|---|---|
| `getSpeed()` | `number` | Velocidad de movimiento actual |
| `setSpeed(n)` | — | Cambiar velocidad (default: `0.1`) |

```js
player.setSpeed(0.2); // doble de velocidad
player.setSpeed(0.1); // velocidad normal
```

## Inventario

| Método | Retorna | Descripción |
|---|---|---|
| `giveItem(nombre, cantidad)` | `boolean` | Dar un item. Retorna `false` si el item no existe o el inventario está lleno |
| `clearInventory()` | — | Limpiar todo el inventario |
| `getItemCount(nombre)` | `number` | Cuántos items de ese tipo tiene |

```js
// Los nombres de items usan el formato de Minecraft
player.giveItem("minecraft:diamond", 10);
player.giveItem("minecraft:stone", 64);
player.giveItem("minecraft:diamond_sword", 1);

var diamonds = player.getItemCount("minecraft:diamond");
player.sendMessage("Tenés §b" + diamonds + " §fdiamantes.");

player.clearInventory();
```

## Sonidos

Reproducir un sonido al jugador en su posición actual.

| Método | Descripción |
|---|---|
| `playSound(nombre)` | Reproducir un sonido |

**Sonidos disponibles:**

| Nombre | Cuándo usarlo |
|---|---|
| `"click"` | Interacción, selección |
| `"levelup"` | Logro, subida de nivel |
| `"pop"` | Item recibido, confirmación |
| `"burp"` | Comer |
| `"deny"` | Error, acción denegada |
| `"door_open"` | Abrir puerta |
| `"door_close"` | Cerrar puerta |
| `"chest_open"` | Abrir cofre |
| `"chest_close"` | Cerrar cofre |
| `"anvil_land"` | Objeto pesado, impacto |
| `"bow_shoot"` | Disparar |
| `"arrow_hit"` | Impacto de flecha |

```js
player.playSound("levelup");
player.playSound("deny");
```

## Comandos

| Método | Descripción |
|---|---|
| `executeCommand(cmd)` | Ejecutar un comando como si lo escribiera el jugador |

```js
player.executeCommand("/spawn");
```
