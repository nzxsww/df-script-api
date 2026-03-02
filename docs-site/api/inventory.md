# Objeto Inventory

El objeto `inventory` representa el inventario de un bloque contenedor (cofre, barril, tolva, horno, etc.) o del jugador. Se obtiene a través de `world.getInventory()` o `player.getInventory()`.

```js
// Inventario de un bloque contenedor (cofre, barril, etc.)
var inv = world.getInventory(x, y, z);
if (inv !== null) {
    // El bloque tiene inventario
}

// Inventario del jugador
var inv = player.getInventory();
```

::: warning Transacciones del mundo
`world.getInventory()` solo funciona desde dentro de un **evento** o un **comando**, no en `onEnable()`. Ver [Objeto World](/api/world#cuándo-usar-world-importante) para más detalles.
:::

## Métodos

| Método | Retorna | Descripción |
|---|---|---|
| `getType()` | `string` | Tipo de inventario: `"player"`, `"chest"`, `"barrel"`, `"hopper"`, `"furnace"`, `"blast_furnace"`, `"smoker"`, `"brewing_stand"`, `"container"` |
| `getSize()` | `number` | Cantidad total de slots del inventario |
| `getItem(slot)` | `{name, count}\|null` | Item en el slot dado. `null` si está vacío |
| `setItem(slot, nombre, cantidad)` | `boolean` | Colocar un item en el slot dado. Retorna `true` si tuvo éxito |
| `addItem(nombre, cantidad)` | `number` | Agrega al primer slot libre. Retorna la cantidad que **no** pudo agregarse (0 = éxito total) |
| `removeItem(nombre, cantidad)` | `boolean` | Remueve la cantidad dada del item. Retorna `true` si tuvo éxito |
| `clear()` | — | Vacía todos los slots del inventario |
| `contains(nombre)` | `boolean` | `true` si el inventario contiene al menos 1 del item dado |
| `count(nombre)` | `number` | Cantidad total del item dado en todo el inventario |
| `getItems()` | `{slot, name, count}[]` | Array con todos los slots no vacíos |
| `setContents(items)` | `boolean` | Reemplaza todo el contenido. `items`: array de `{slot, name, count}`. Los slots no especificados quedan vacíos |

## Ejemplos

### Inspeccionar un cofre

```js
commands.register("inspeccionarcofre", "Ver contenido del cofre frente a ti", function(player, args) {
    var x = Math.floor(player.getX());
    var y = Math.floor(player.getY());
    var z = Math.floor(player.getZ() + 1); // bloque frente al jugador

    var inv = world.getInventory(x, y, z);
    if (inv === null) {
        player.sendMessage("§cNo hay un contenedor ahí.");
        return;
    }

    var items = inv.getItems();
    player.sendMessage("§eContenido del cofre (§f" + items.length + "§e items):");
    for (var i = 0; i < items.length; i++) {
        player.sendMessage("§7Slot §f" + items[i].slot + "§7: §f" + items[i].name + " x" + items[i].count);
    }
});
```

### Llenar un cofre con items

```js
commands.register("llenarcofre", "Llena el cofre frente a ti con diamantes", function(player, args) {
    var x = Math.floor(player.getX());
    var y = Math.floor(player.getY());
    var z = Math.floor(player.getZ() + 1);

    var inv = world.getInventory(x, y, z);
    if (inv === null) {
        player.sendMessage("§cNo hay un contenedor ahí.");
        return;
    }

    inv.clear();
    var sobrante = inv.addItem("minecraft:diamond", 1000);
    player.sendMessage("§bDiamantes agregados! Sobrante: §f" + sobrante);
});
```

### Inventario del jugador

```js
events.on("PlayerJoin", function(event) {
    var p = event.getPlayer();
    var inv = p.getInventory();

    // Ver qué tiene el jugador
    var items = inv.getItems();
    console.log(p.getName() + " tiene " + items.length + " tipos de items");

    // Dar item de bienvenida si no tiene diamantes
    if (!inv.contains("minecraft:diamond")) {
        inv.addItem("minecraft:diamond", 1);
        p.sendMessage("§b¡Diamante de bienvenida!");
    }
});
```

### Verificar y quitar items

```js
commands.register("pagar", "Paga 5 diamantes para acceder", function(player, args) {
    var inv = player.getInventory();

    if (inv.count("minecraft:diamond") < 5) {
        player.sendMessage("§cNecesitás 5 diamantes para acceder.");
        return;
    }

    inv.removeItem("minecraft:diamond", 5);
    player.sendMessage("§a¡Acceso concedido! Se te cobraron 5 diamantes.");
});
```

## Bloques con inventario

| Bloque | Nombre en JS |
|---|---|
| Cofre | `minecraft:chest` |
| Cofre trampa | `minecraft:trapped_chest` |
| Barril | `minecraft:barrel` |
| Tolva | `minecraft:hopper` |
| Horno | `minecraft:furnace` |
| Horno de fundición | `minecraft:blast_furnace` |
| Ahumador | `minecraft:smoker` |
| Dispensador | `minecraft:dispenser` |
| Tolva | `minecraft:dropper` |
