# Objeto World

El objeto `world` está disponible globalmente en cada plugin. Permite interactuar con el mundo de Minecraft: colocar bloques, invocar entidades, generar partículas y trabajar con los jugadores conectados.

## Bloques

| Método | Retorna | Descripción |
|---|---|---|
| `setBlock(x, y, z, nombre)` | — | Colocar un bloque en la posición dada |
| `getBlock(x, y, z)` | `string` | Obtener el nombre del bloque en la posición dada |
| `getHighestBlock(x, z)` | `number` | Obtener la Y del bloque más alto en X,Z |

```js
world.setBlock(100, 64, 100, "minecraft:stone");
world.setBlock(100, 65, 100, "minecraft:grass_block");

var bloque = world.getBlock(100, 64, 100); // "minecraft:stone"
var y = world.getHighestBlock(100, 100);   // ej: 65
```

## Partículas

| Método | Descripción |
|---|---|
| `spawnParticle(x, y, z, nombre)` | Generar una partícula en la posición dada |

**Partículas disponibles:**

| Nombre | Descripción |
|---|---|
| `"flame"` | Llama de antorcha/fuego |
| `"lava"` | Burbuja de lava |
| `"water_drip"` | Goteo de agua |
| `"lava_drip"` | Goteo de lava |
| `"explosion"` | Explosión grande |
| `"bone_meal"` | Uso de harina de hueso |
| `"evaporate"` | Evaporación de agua |
| `"snowball"` | Poof de bola de nieve |
| `"egg_smash"` | Huevo rompiéndose |
| `"entity_flame"` | Entidad en llamas |

```js
world.spawnParticle(100, 65, 100, "flame");
world.spawnParticle(100, 65, 100, "explosion");
```

## Entidades

| Método | Descripción |
|---|---|
| `spawnLightning(x, y, z)` | Invocar un rayo en la posición dada |
| `spawnTNT(x, y, z, fuse)` | Invocar un TNT activado. `fuse` es segundos hasta explotar |
| `spawnText(x, y, z, texto)` | Crear un texto flotante en la posición dada |
| `spawnExperienceOrb(x, y, z, cantidad)` | Generar un orbe de experiencia |

```js
// Invocar un rayo
world.spawnLightning(100, 64, 100);

// TNT con mecha de 4 segundos
world.spawnTNT(100, 65, 100, 4);

// Texto flotante (soporta códigos de color)
world.spawnText(0, 70, 0, "§a§l¡Bienvenido al servidor!");

// Orbe de experiencia
world.spawnExperienceOrb(100, 64, 100, 50);
```

## Jugadores

| Método | Retorna | Descripción |
|---|---|---|
| `getPlayers()` | `player[]` | Obtener todos los jugadores conectados como objetos player |
| `getPlayerCount()` | `number` | Cantidad de jugadores conectados |
| `broadcast(msg)` | — | Enviar un mensaje a todos los jugadores conectados |

```js
// Broadcast a todos
world.broadcast("§aEl servidor reinicia en 5 minutos!");

// Iterar sobre todos los jugadores
var jugadores = world.getPlayers();
for (var i = 0; i < jugadores.length; i++) {
    jugadores[i].sendTitle("§aEvento!", "§7Empezando ahora");
    jugadores[i].addEffect("speed", 2, 60);
}

// Cantidad de jugadores
var cantidad = world.getPlayerCount();
world.broadcast("§7Hay §f" + cantidad + " §7jugadores conectados.");
```
