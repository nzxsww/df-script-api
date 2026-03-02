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

## Títulos

| Método | Descripción |
|---|---|
| `sendTitle(texto, subtitulo)` | Mostrar un título grande en el centro de la pantalla con subtítulo opcional |

```js
player.sendTitle("§a¡Bienvenido!", "§7Que la pases bien");
```

## Efectos de poción

| Método | Descripción |
|---|---|
| `addEffect(nombre, nivel, segundos)` | Aplicar un efecto de poción. `nivel` empieza en 1 (nivel I) |
| `removeEffect(nombre)` | Quitar un efecto específico |
| `clearEffects()` | Quitar todos los efectos activos |

**Efectos disponibles:**

| Nombre | Descripción |
|---|---|
| `"speed"` | Moverse más rápido |
| `"slowness"` | Moverse más lento |
| `"haste"` | Minar más rápido |
| `"mining_fatigue"` | Minar más lento |
| `"strength"` | Hacer más daño |
| `"jump_boost"` | Saltar más alto |
| `"nausea"` | Visión distorsionada |
| `"regeneration"` | Regenerar vida |
| `"resistance"` | Recibir menos daño |
| `"fire_resistance"` | Inmune al fuego |
| `"water_breathing"` | Respirar bajo el agua |
| `"invisibility"` | Invisible para otros |
| `"blindness"` | Visión reducida |
| `"night_vision"` | Ver en la oscuridad |
| `"hunger"` | Perder hambre más rápido |
| `"weakness"` | Hacer menos daño |
| `"poison"` | Perder vida gradualmente |
| `"wither"` | Perder vida (ignora armadura) |
| `"health_boost"` | Aumentar vida máxima |
| `"absorption"` | Puntos de vida extra |
| `"saturation"` | Restaurar comida instantáneamente |
| `"levitation"` | Flotar hacia arriba |
| `"slow_falling"` | Caer lentamente |
| `"conduit_power"` | Haste + visión bajo el agua |
| `"darkness"` | Visión muy reducida |

```js
player.addEffect("speed", 2, 30);       // Velocidad II por 30 segundos
player.addEffect("regeneration", 1, 10); // Regeneración I por 10 segundos
player.removeEffect("speed");
player.clearEffects();
```

## Armadura

| Método | Retorna | Descripción |
|---|---|---|
| `setArmour(slot, nombre)` | — | Equipar una pieza de armadura. Slot: `0`=casco, `1`=pechera, `2`=pantalones, `3`=botas |
| `getArmour(slot)` | `string` | Obtener el nombre del item en ese slot. Retorna `""` si está vacío |
| `clearArmour()` | — | Quitar toda la armadura |

```js
player.setArmour(0, "minecraft:diamond_helmet");
player.setArmour(1, "minecraft:diamond_chestplate");
player.setArmour(2, "minecraft:diamond_leggings");
player.setArmour(3, "minecraft:diamond_boots");

var casco = player.getArmour(0); // "minecraft:diamond_helmet"
player.clearArmour();
```

## Scoreboard

| Método | Descripción |
|---|---|
| `sendScoreboard(sb)` | Envía un scoreboard al jugador (creado con `scoreboard.create()`) |
| `removeScoreboard()` | Quita el scoreboard de la pantalla del jugador |

```js
var sb = scoreboard.create("§6Mi Servidor");
sb.setLine(0, "§7Jugadores: 5");
sb.setLine(1, "§7Mapa: Lobby");
player.sendScoreboard(sb);

// Más tarde...
player.removeScoreboard();
```

Ver la [Scoreboard API](/api/scoreboard) para documentación completa incluyendo ScoreboardManager y Live Scoreboard.

## Comandos

| Método | Descripción |
|---|---|
| `executeCommand(cmd)` | Ejecutar un comando como si lo escribiera el jugador |

```js
player.executeCommand("/spawn");
```
