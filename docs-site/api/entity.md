# Objeto Entity

El objeto `entity` representa una entidad individual en el mundo (item, TNT, flecha, mob, jugador, texto flotante, etc.). Se obtiene a través de los métodos del objeto `world` o desde eventos de entidades.

```js
// Obtener entidades via world
var entities = world.getEntities();
var entity = entities[0];

// Obtener entidades cercanas
var nearby = world.getEntitiesInRadius(x, y, z, 10);
```

## Métodos base (todas las entidades)

Estos métodos están disponibles en **cualquier** tipo de entidad.

| Método | Retorna | Descripción |
|---|---|---|
| `getUUID()` | `string` | UUID único de la entidad |
| `getType()` | `string` | Tipo de entidad (ej: `"minecraft:item"`, `"minecraft:tnt"`) |
| `getX()` | `number` | Posición X |
| `getY()` | `number` | Posición Y |
| `getZ()` | `number` | Posición Z |
| `remove()` | — | Remover la entidad del mundo |
| `teleport(x, y, z)` | — | Teletransportar la entidad (si lo soporta) |
| `setVelocity(x, y, z)` | — | Cambiar la velocidad de la entidad (si lo soporta) |

```js
var entities = world.getEntities();
for (var i = 0; i < entities.length; i++) {
    var e = entities[i];
    console.log(e.getType() + " UUID=" + e.getUUID());
    console.log("  pos: " + e.getX().toFixed(1) + "," + e.getY().toFixed(1) + "," + e.getZ().toFixed(1));
}
```

## Métodos de entidades vivas (Living / mobs)

Disponibles cuando `typeof entity.getHealth === "function"`. Aplica a mobs, animales y jugadores.

| Método | Retorna | Descripción |
|---|---|---|
| `getHealth()` | `number` | Vida actual |
| `getMaxHealth()` | `number` | Vida máxima |
| `setMaxHealth(n)` | — | Cambiar vida máxima |
| `isDead()` | `boolean` | Si la entidad está muerta |
| `hurt(damage)` | — | Aplicar daño de ataque |
| `heal(health)` | — | Curar la entidad |
| `knockBack(x, y, z, fuerza, altura)` | — | Empujar la entidad desde una posición fuente |
| `getSpeed()` | `number` | Velocidad de movimiento |
| `addEffect(nombre, nivel, segundos)` | — | Aplicar efecto de poción |
| `removeEffect(nombre)` | — | Quitar un efecto activo |
| `clearEffects()` | — | Quitar todos los efectos activos |

```js
var nearby = world.getEntitiesInRadius(x, y, z, 10);
for (var i = 0; i < nearby.length; i++) {
    var e = nearby[i];
    if (typeof e.getHealth === "function") {
        console.log(e.getType() + " vida=" + e.getHealth().toFixed(1) + "/" + e.getMaxHealth().toFixed(1));
        e.addEffect("slowness", 1, 10);
        e.hurt(5);
    }
}
```

## Métodos exclusivos de jugadores

Disponibles cuando `typeof entity.getName === "function"`. El jugador también tiene todos los métodos de entidad viva.

| Método | Retorna | Descripción |
|---|---|---|
| `getName()` | `string` | Nombre del jugador |
| `sendMessage(msg)` | — | Enviar mensaje de chat |
| `sendTitle(texto, subtitulo)` | — | Mostrar título en pantalla |
| `disconnect(msg)` | — | Desconectar al jugador |

```js
var entities = world.getEntities();
for (var i = 0; i < entities.length; i++) {
    var e = entities[i];
    if (typeof e.getName === "function") {
        // Es un jugador
        e.sendMessage("§aHola " + e.getName() + "!");
        e.sendTitle("§aEvento", "§7Comenzando...");
    }
}
```

## Tipos de entidad comunes

| Tipo | Descripción | Living |
|---|---|---|
| `"minecraft:item"` | Item tirado en el suelo | No |
| `"minecraft:tnt"` | TNT activado | No |
| `"minecraft:arrow"` | Flecha | No |
| `"minecraft:xp_orb"` | Orbe de experiencia | No |
| `"minecraft:lightning_bolt"` | Rayo | No |
| `"minecraft:falling_block"` | Bloque cayendo | No |
| `"minecraft:fireworks_rocket"` | Fuego artificial | No |
| `"dragonfly:text"` | Texto flotante (creado con `world.spawnText`) | No |

## Verificar el tipo antes de usar métodos opcionales

```js
var entities = world.getEntities();
for (var i = 0; i < entities.length; i++) {
    var e = entities[i];

    // Siempre disponible
    console.log(e.getType() + " en " + e.getX().toFixed(0) + "," + e.getY().toFixed(0) + "," + e.getZ().toFixed(0));

    // Solo en mobs/jugadores
    if (typeof e.getHealth === "function") {
        console.log("  vida: " + e.getHealth().toFixed(1));
    }

    // Solo en jugadores
    if (typeof e.getName === "function") {
        console.log("  jugador: " + e.getName());
    }
}
```

## Remover entidades

```js
// Remover por UUID (desde world)
world.removeEntityByUUID("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx");

// Remover desde el objeto entity directamente
var entities = world.getEntities();
for (var i = 0; i < entities.length; i++) {
    if (entities[i].getType() === "minecraft:tnt") {
        entities[i].remove(); // Eliminar todos los TNT
    }
}

// Remover entidades de item en un radio de 5 bloques
var nearby = world.getEntitiesInRadius(x, y, z, 5);
for (var i = 0; i < nearby.length; i++) {
    if (nearby[i].getType() === "minecraft:item") {
        nearby[i].remove();
    }
}
```
