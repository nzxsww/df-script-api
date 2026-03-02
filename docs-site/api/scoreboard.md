# Scoreboard API

La Scoreboard API permite mostrar información en el panel lateral derecho de la pantalla del jugador (el scoreboard de Minecraft Bedrock). Soporta tres modos: **estático**, **gestionado por jugador** (ScoreboardManager) y **auto-actualizable** (Live Scoreboard).

## Scoreboard estático

El modo más simple. Creás un scoreboard, le ponés líneas y lo enviás a un jugador.

```js
var sb = scoreboard.create("§6§lMi Servidor");
sb.setLine(0, "§7Jugadores: 5");
sb.setLine(1, "§7Mapa: Lobby");
sb.setLine(2, "");
sb.setLine(3, "§fplay.miservidor.com");

player.sendScoreboard(sb);
player.removeScoreboard();
```

::: warning
Modificar un scoreboard después de enviarlo **no actualiza** lo que ve el jugador automáticamente. Hay que volver a llamar `player.sendScoreboard(sb)` o usar el [Live Scoreboard](#live-scoreboard).
:::

### Métodos del scoreboard

| Método | Retorna | Descripción |
|---|---|---|
| `scoreboard.create(titulo)` | `Scoreboard` | Crea un nuevo scoreboard con el título dado |
| `sb.setLine(index, texto)` | `boolean` | Setea el texto en la línea `index` (0–14). Rellena con líneas vacías si el índice supera las actuales |
| `sb.addLine(texto)` | `boolean` | Agrega una línea al final. Retorna `false` si se superan las 15 líneas |
| `sb.removeLine(index)` | `boolean` | Elimina la línea en el índice dado |
| `sb.getLines()` | `string[]` | Retorna todas las líneas actuales como array |
| `sb.getLineCount()` | `number` | Cantidad de líneas actuales |
| `sb.getTitle()` | `string` | Título del scoreboard |
| `sb.setDescending()` | — | Invierte el orden visual de las líneas |
| `sb.isDescending()` | `boolean` | Si está en orden descendente |
| `sb.removePadding()` | — | Elimina el padding automático de espacios en cada línea |
| `sb.sendTo(player)` | — | Atajo: envía este scoreboard al jugador dado |

### Métodos en el jugador

| Método | Descripción |
|---|---|
| `player.sendScoreboard(sb)` | Envía el scoreboard al jugador |
| `player.removeScoreboard()` | Quita el scoreboard de la pantalla del jugador |

---

## ScoreboardManager

El `ScoreboardManager` gestiona qué scoreboard tiene asignado cada jugador. Es ideal cuando tenés múltiples jugadores con distintos scoreboards o querés saber en cualquier momento qué está viendo cada uno.

Cada plugin tiene su **propio manager aislado** — `scoreboard.getManager()` siempre retorna el mismo manager para ese plugin.

```js
var mgr = scoreboard.getManager();

events.on("PlayerJoin", function(event) {
    var player = event.getPlayer();
    var sb = scoreboard.create("§6Mi Servidor");
    sb.setLine(0, "§7Bienvenido, §f" + player.getName());
    sb.setLine(1, "§7Jugadores: §f" + server.getPlayerCount());
    mgr.set(player, sb);
});

events.on("PlayerQuit", function(event) {
    mgr.remove(event.getPlayer());
});
```

### Métodos del ScoreboardManager

| Método | Retorna | Descripción |
|---|---|---|
| `scoreboard.getManager()` | `ScoreboardManager` | Retorna el manager del plugin |
| `mgr.set(player, sb)` | `boolean` | Asigna el scoreboard al jugador y lo envía inmediatamente |
| `mgr.remove(player)` | — | Quita el scoreboard asignado al jugador y lo elimina de su pantalla |
| `mgr.get(player)` | `Scoreboard \| null` | Retorna el scoreboardWrapper activo del jugador, o `null` si no tiene |
| `mgr.hasScoreboard(player)` | `boolean` | `true` si el jugador tiene un scoreboard asignado |
| `mgr.getAssignedCount()` | `number` | Cantidad de jugadores con scoreboard asignado |
| `mgr.clearAll()` | — | Quita el registro de todos los jugadores del manager |

---

## Live Scoreboard

El Live Scoreboard se **auto-actualiza automáticamente** en un intervalo definido. En cada tick, el callback recibe un scoreboard limpio para que el plugin llene las líneas, y luego lo reenvía a todos los jugadores asignados.

```js
var live = scoreboard.createLive("§aMi Servidor", function(sb) {
    sb.setLine(0, "§fJugadores: §a" + server.getPlayerCount());
    sb.setLine(1, "§fModo: §eLobby");
    sb.setLine(2, "");
    sb.setLine(3, "§7play.miservidor.com");
}, 2000); // actualiza cada 2 segundos

events.on("PlayerJoin", function(event) {
    live.addPlayer(event.getPlayer());
});

events.on("PlayerQuit", function(event) {
    live.removePlayer(event.getPlayer());
});

function onDisable() {
    live.stop(); // detener el ticker al cerrar el plugin
}
```

::: tip
El intervalo mínimo es **50ms**. Para scoreboards de lobby o HUD se recomienda usar **1000–5000ms** para no saturar el servidor con paquetes innecesarios.
:::

### Métodos del Live Scoreboard

| Método | Retorna | Descripción |
|---|---|---|
| `scoreboard.createLive(titulo, fn, intervalMs)` | `LiveScoreboard` | Crea y arranca un live scoreboard |
| `live.addPlayer(player)` | `boolean` | Agrega el jugador — empezará a recibir updates automáticos |
| `live.removePlayer(player)` | `boolean` | Saca al jugador del live y le quita el scoreboard |
| `live.hasPlayer(player)` | `boolean` | `true` si el jugador está recibiendo updates |
| `live.getPlayerCount()` | `number` | Cantidad de jugadores en el live |
| `live.clearPlayers()` | — | Saca a todos los jugadores sin detener el live |
| `live.stop()` | — | Detiene el auto-update y quita el scoreboard a todos los jugadores |

---

## Ejemplo completo

Plugin de lobby con scoreboard en tiempo real para todos los jugadores:

```js
config.setDefaults({
    "titulo": "§6§lMi Servidor",
    "intervalo": 2000
});

var live;

function onEnable() {
    var titulo = config.getString("titulo", "§6§lMi Servidor");
    var intervalo = config.getInt("intervalo", 2000);

    live = scoreboard.createLive(titulo, function(sb) {
        var hora = new Date().toLocaleTimeString();
        sb.setLine(0, "§7Jugadores: §f" + server.getPlayerCount() + "§7/§f" + server.getMaxPlayers());
        sb.setLine(1, "§7Hora: §f" + hora);
        sb.setLine(2, "");
        sb.setLine(3, "§7play.miservidor.com");
    }, intervalo);

    // Agregar jugadores que ya estén conectados
    var players = server.getPlayers();
    for (var i = 0; i < players.length; i++) {
        live.addPlayer(players[i]);
    }

    events.on("PlayerJoin", function(event) {
        live.addPlayer(event.getPlayer());
    });

    events.on("PlayerQuit", function(event) {
        live.removePlayer(event.getPlayer());
    });
}

function onDisable() {
    if (live) live.stop();
    config.save();
}

module = { onEnable: onEnable, onDisable: onDisable };
```
