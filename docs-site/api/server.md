# Objeto Server

El objeto `server` está disponible globalmente en cada plugin. Permite interactuar con el servidor de forma global: obtener jugadores de todos los mundos, enviar mensajes a todos, y consultar información del servidor. Es el equivalente al objeto `Server` de la API de Bukkit.

> **Diferencia con `world`:** El objeto `world` interactúa con un mundo específico (bloques, entidades, partículas). El objeto `server` opera a nivel global del servidor (todos los jugadores, sin importar en qué mundo estén).

## Jugadores

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayers()` | `player[]` | Todos los jugadores conectados en el servidor |
| `getPlayerCount()` | `number` | Cantidad de jugadores conectados |
| `getMaxPlayers()` | `number` | Límite máximo de jugadores configurado |
| `getPlayer(nombre)` | `player\|null` | Buscar jugador por nombre. Retorna `null` si no está conectado |
| `getPlayerByXUID(xuid)` | `player\|null` | Buscar jugador por XUID de Xbox. Retorna `null` si no está conectado |

```js
// Cantidad de jugadores
var count = server.getPlayerCount();
var max = server.getMaxPlayers();
console.log("Jugadores: " + count + "/" + max);

// Buscar jugador específico
var jugador = server.getPlayer("Notch");
if (jugador !== null) {
    jugador.sendMessage("§aHola Notch!");
}

// Iterar sobre todos los jugadores
var jugadores = server.getPlayers();
for (var i = 0; i < jugadores.length; i++) {
    jugadores[i].sendMessage("§eMensaje para todos!");
}
```

## Mensajes

| Método | Descripción |
|---|---|
| `broadcast(msg)` | Enviar un mensaje de chat a todos los jugadores conectados |
| `broadcastTitle(texto, subtitulo)` | Enviar un título grande a todos los jugadores conectados |

```js
server.broadcast("§a[Anuncio] §fEl evento comienza en 1 minuto!");
server.broadcastTitle("§c§l¡EVENTO!", "§7Preparate para jugar");
```

## Información del servidor

| Método | Retorna | Descripción |
|---|---|---|
| `getName()` | `string` | Nombre del plugin que hace la llamada |
| `shutdown()` | — | Cerrar el servidor de forma segura |

```js
console.log("Plugin: " + server.getName());

// Cerrar el servidor (usar con cuidado)
server.shutdown();
```
