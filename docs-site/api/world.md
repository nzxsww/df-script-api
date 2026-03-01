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

## Entidades

| Método | Retorna | Descripción |
|---|---|---|
| `getEntities()` | `entity[]` | Obtener todas las entidades del mundo |
| `getEntitiesInRadius(x, y, z, radio)` | `entity[]` | Entidades dentro del radio esférico dado (en bloques) |
| `removeEntityByUUID(uuid)` | `boolean` | Remover una entidad por su UUID. Retorna `true` si fue encontrada y removida |

**Objeto `entity` retornado:**

| Método | Disponible en | Descripción |
|---|---|---|
| `getUUID()` | Todas | UUID único de la entidad |
| `getType()` | Todas | Tipo de entidad (ej: `"minecraft:item"`, `"minecraft:tnt"`, `"minecraft:arrow"`) |
| `getX/Y/Z()` | Todas | Posición en el mundo |
| `remove()` | Todas | Remover la entidad del mundo |
| `teleport(x, y, z)` | Todas (si soporta) | Teletransportar la entidad |
| `setVelocity(x, y, z)` | Todas (si soporta) | Cambiar velocidad |
| `getHealth()` | Living (mobs) | Vida actual |
| `getMaxHealth()` | Living (mobs) | Vida máxima |
| `setMaxHealth(n)` | Living (mobs) | Cambiar vida máxima |
| `isDead()` | Living (mobs) | Si está muerta |
| `hurt(damage)` | Living (mobs) | Aplicar daño |
| `heal(health)` | Living (mobs) | Curar |
| `knockBack(x, y, z, fuerza, altura)` | Living (mobs) | Empujar |
| `addEffect(nombre, nivel, segundos)` | Living (mobs) | Aplicar efecto de poción |
| `removeEffect(nombre)` | Living (mobs) | Quitar efecto |
| `clearEffects()` | Living (mobs) | Quitar todos los efectos |
| `getSpeed()` | Living (mobs) | Velocidad de movimiento |
| `getName()` | Jugadores | Nombre del jugador |
| `sendMessage(msg)` | Jugadores | Enviar mensaje |
| `sendTitle(texto, subtitulo)` | Jugadores | Mostrar título |
| `disconnect(msg)` | Jugadores | Desconectar |

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
| `"dragonfly:text"` | Texto flotante (creado con `spawnText`) |

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
