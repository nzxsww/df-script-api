# Inventory Object

The `inventory` object represents the inventory of a container block (chest, barrel, hopper, furnace, etc.) or a player. It is obtained through `world.getInventory()` or `player.getInventory()`.

```js
// Inventory of a container block (chest, barrel, etc.)
var inv = world.getInventory(x, y, z);
if (inv !== null) {
    // The block has an inventory
}

// Player inventory
var inv = player.getInventory();
```

::: warning World transactions
`world.getInventory()` only works from inside an **event** or a **command**, not in `onEnable()`. See [World Object](/en/api/world#when-to-use-world-important) for details.
:::

## Methods

| Method | Returns | Description |
|---|---|---|
| `getType()` | `string` | Inventory type: `"player"`, `"chest"`, `"barrel"`, `"hopper"`, `"furnace"`, `"blast_furnace"`, `"smoker"`, `"brewing_stand"`, `"container"` |
| `getSize()` | `number` | Total number of slots in the inventory |
| `getItem(slot)` | `{name, count}\|null` | Item in the given slot. `null` if empty |
| `setItem(slot, name, count)` | `boolean` | Place an item in the given slot. Returns `true` on success |
| `addItem(name, count)` | `number` | Add to the first available slot. Returns the amount that **could not** be added (0 = full success) |
| `removeItem(name, count)` | `boolean` | Remove the given amount of the item. Returns `true` on success |
| `clear()` | — | Empty all slots in the inventory |
| `contains(name)` | `boolean` | `true` if the inventory contains at least 1 of the given item |
| `count(name)` | `number` | Total count of the given item across all slots |
| `getItems()` | `{slot, name, count}[]` | Array of all non-empty slots |
| `setContents(items)` | `boolean` | Replace all contents. `items`: array of `{slot, name, count}`. Unspecified slots are cleared |

## Examples

### Inspecting a chest

```js
commands.register("inspectchest", "View chest contents in front of you", function(player, args) {
    var x = Math.floor(player.getX());
    var y = Math.floor(player.getY());
    var z = Math.floor(player.getZ() + 1); // block in front of player

    var inv = world.getInventory(x, y, z);
    if (inv === null) {
        player.sendMessage("§cNo container there.");
        return;
    }

    var items = inv.getItems();
    player.sendMessage("§eChest contents (§f" + items.length + "§e items):");
    for (var i = 0; i < items.length; i++) {
        player.sendMessage("§7Slot §f" + items[i].slot + "§7: §f" + items[i].name + " x" + items[i].count);
    }
});
```

### Filling a chest

```js
commands.register("fillchest", "Fill the chest in front of you with diamonds", function(player, args) {
    var x = Math.floor(player.getX());
    var y = Math.floor(player.getY());
    var z = Math.floor(player.getZ() + 1);

    var inv = world.getInventory(x, y, z);
    if (inv === null) {
        player.sendMessage("§cNo container there.");
        return;
    }

    inv.clear();
    var leftover = inv.addItem("minecraft:diamond", 1000);
    player.sendMessage("§bDiamonds added! Leftover: §f" + leftover);
});
```

### Player inventory

```js
events.on("PlayerJoin", function(event) {
    var p = event.getPlayer();
    var inv = p.getInventory();

    // Check what the player has
    var items = inv.getItems();
    console.log(p.getName() + " has " + items.length + " item types");

    // Give welcome item if player has no diamonds
    if (!inv.contains("minecraft:diamond")) {
        inv.addItem("minecraft:diamond", 1);
        p.sendMessage("§b¡Welcome diamond!");
    }
});
```

### Check and remove items

```js
commands.register("pay", "Pay 5 diamonds to gain access", function(player, args) {
    var inv = player.getInventory();

    if (inv.count("minecraft:diamond") < 5) {
        player.sendMessage("§cYou need 5 diamonds to access.");
        return;
    }

    inv.removeItem("minecraft:diamond", 5);
    player.sendMessage("§aAccess granted! 5 diamonds charged.");
});
```

## Blocks with inventory

| Block | Name in JS |
|---|---|
| Chest | `minecraft:chest` |
| Trapped chest | `minecraft:trapped_chest` |
| Barrel | `minecraft:barrel` |
| Hopper | `minecraft:hopper` |
| Furnace | `minecraft:furnace` |
| Blast furnace | `minecraft:blast_furnace` |
| Smoker | `minecraft:smoker` |
| Dispenser | `minecraft:dispenser` |
| Dropper | `minecraft:dropper` |
