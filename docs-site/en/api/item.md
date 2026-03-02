# Item Object

The `item` object represents an `ItemStack` (Bukkit-like). It is obtained from inventories (`getItem`, `getItems`) or created with `item.create()`.

```js
var it = item.create("minecraft:diamond_sword", 1);
```

## Create an item

| Method | Description |
|---|---|
| `item.create(name, count)` | Create an item with the given count |

```js
var it = item.create("minecraft:diamond", 64);
```

## Methods

| Method | Returns | Description |
|---|---|---|
| `getName()` | `string` | Item name (e.g. `minecraft:diamond`) |
| `getCount()` | `number` | Count |
| `setCount(n)` | `item` | Change count and return new item |
| `getDisplayName()` | `string` | Custom display name |
| `setDisplayName(name)` | `item` | Set custom display name |
| `getLore()` | `string[]` | Lore lines |
| `setLore(lines)` | `item` | Set lore |
| `getDurability()` | `number` | Current durability |
| `setDurability(n)` | `item` | Set durability |
| `getMaxDurability()` | `number` | Max durability |
| `isEnchanted()` | `boolean` | Has enchantments |
| `getEnchantments()` | `{name, level}[]` | Enchantments list |
| `addEnchantment(name, level)` | `item` | Add enchantment |
| `removeEnchantment(name)` | `item` | Remove enchantment |
| `clone()` | `item` | Clone item |

## Supported enchantments

| Name |
|---|
| `sharpness` |
| `efficiency` |
| `unbreaking` |
| `silk_touch` |
| `power` |
| `punch` |
| `flame` |
| `infinity` |
| `protection` |
| `fire_protection` |
| `blast_protection` |
| `projectile_protection` |
| `feather_falling` |
| `thorns` |
| `respiration` |
| `aqua_affinity` |
| `depth_strider` |
| `mending` |
| `vanishing` |
| `fire_aspect` |
| `knockback` |

## Examples

### Create and customize an item

```js
var sword = item.create("minecraft:diamond_sword", 1)
    .setDisplayName("§bEpic Sword")
    .setLore(["§7Forged in dragons", "§7Level 10"])
    .addEnchantment("sharpness", 5)
    .addEnchantment("unbreaking", 3);
```

### Use an item in an inventory

```js
var inv = player.getInventory();
inv.addItemStack(sword);
```
