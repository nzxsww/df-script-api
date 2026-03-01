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

### Invocar entidades

| Método | Descripción |
|---|---|
| `spawnEntity(tipo, x, y, z)` | Invocar una entidad en la posición dada |
| `spawnEntity(tipo, x, y, z, opciones)` | Invocar una entidad con opciones adicionales |

**Tipos disponibles y sus opciones:**

| Tipo | Opciones | Descripción |
|---|---|---|
| `"lightning"` | — | Rayo |
| `"tnt"` | `{ fuse: 4 }` | TNT activado. `fuse` = segundos hasta explotar (default: 4) |
| `"text"` | `{ text: "§aHola" }` | Texto flotante |
| `"experience_orb"` | `{ amount: 50 }` | Orbe de experiencia con cantidad de XP (default: 1) |
| `"item"` | `{ item: "minecraft:diamond", count: 1 }` | Item en el suelo |

```js
world.spawnEntity("lightning", x, y, z);
world.spawnEntity("tnt", x, y, z, { fuse: 4 });
world.spawnEntity("text", x, y, z, { text: "§a§lBienvenido!" });
world.spawnEntity("experience_orb", x, y, z, { amount: 100 });
world.spawnEntity("item", x, y, z, { item: "minecraft:diamond", count: 5 });
```

### Obtener y remover entidades

| Método | Retorna | Descripción |
|---|---|---|
| `getEntities()` | `entity[]` | Obtener todas las entidades del mundo |
| `getEntitiesInRadius(x, y, z, radio)` | `entity[]` | Entidades dentro del radio esférico dado (en bloques) |
| `removeEntityByUUID(uuid)` | `boolean` | Remover una entidad por su UUID. Retorna `true` si fue encontrada y removida |

Cada entidad retornada es un objeto con métodos propios. Ver la referencia completa en [Objeto Entity](/api/entity).

```js
// Listar todas las entidades del mundo
var entities = world.getEntities();
for (var i = 0; i < entities.length; i++) {
    var e = entities[i];
    console.log(e.getType() + " en " + e.getX() + "," + e.getY() + "," + e.getZ());
}

// Entidades en un radio de 10 bloques alrededor del jugador
events.on("PlayerJoin", function(event) {
    var p = event.getPlayer();
    var nearby = world.getEntitiesInRadius(p.getX(), p.getY(), p.getZ(), 10);
    p.sendMessage("§eEntidades cercanas: §f" + nearby.length);
    for (var i = 0; i < nearby.length; i++) {
        var e = nearby[i];
        // Si es living, mostrar vida
        if (typeof e.getHealth === "function") {
            p.sendMessage("§7- §f" + e.getType() + " §7vida=" + e.getHealth().toFixed(1));
        } else {
            p.sendMessage("§7- §f" + e.getType());
        }
    }
});

// Remover entidad por UUID
world.removeEntityByUUID("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx");

// Remover todas las entidades de tipo TNT
var entities = world.getEntities();
for (var i = 0; i < entities.length; i++) {
    if (entities[i].getType() === "minecraft:tnt") {
        entities[i].remove();
    }
}
```

**Tipos de entidad comunes:**

| Tipo | Descripción |
|---|---|
| `"minecraft:item"` | Item tirado en el suelo |
| `"minecraft:tnt"` | TNT activado |
| `"minecraft:arrow"` | Flecha |
| `"minecraft:xp_orb"` | Orbe de experiencia |
| `"minecraft:lightning_bolt"` | Rayo |
| `"minecraft:falling_block"` | Bloque cayendo |
| `"minecraft:fireworks_rocket"` | Fuego artificial |
| `"dragonfly:text"` | Texto flotante (creado con `world.spawnEntity("text", ...)`) |

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
