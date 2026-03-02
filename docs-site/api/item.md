# Objeto Item

El objeto `item` representa un `ItemStack` al estilo Bukkit. Se obtiene desde inventarios (`getItem`, `getItems`) o se crea con `item.create()`.

```js
var it = item.create("minecraft:diamond_sword", 1);
```

## Crear un item

| Método | Descripción |
|---|---|
| `item.create(nombre, cantidad)` | Crea un item con la cantidad dada |

```js
var it = item.create("minecraft:diamond", 64);
```

## Métodos disponibles

| Método | Retorna | Descripción |
|---|---|---|
| `getName()` | `string` | Nombre del item (ej: `minecraft:diamond`) |
| `getCount()` | `number` | Cantidad |
| `setCount(n)` | `item` | Cambia la cantidad y retorna el nuevo item |
| `getDisplayName()` | `string` | Nombre personalizado (display name) |
| `setDisplayName(name)` | `item` | Cambia el nombre personalizado |
| `getLore()` | `string[]` | Lore actual |
| `setLore(lines)` | `item` | Cambia el lore |
| `getDurability()` | `number` | Durabilidad actual |
| `setDurability(n)` | `item` | Cambia la durabilidad |
| `getMaxDurability()` | `number` | Durabilidad máxima |
| `isEnchanted()` | `boolean` | Si tiene encantamientos |
| `getEnchantments()` | `{name, level}[]` | Lista de encantamientos |
| `addEnchantment(nombre, nivel)` | `item` | Agrega encantamiento |
| `removeEnchantment(nombre)` | `item` | Quita encantamiento |
| `clone()` | `item` | Clona el item |

## Encantamientos soportados

| Nombre |
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

## Ejemplos

### Crear y personalizar un item

```js
var sword = item.create("minecraft:diamond_sword", 1)
    .setDisplayName("§bEspada Épica")
    .setLore(["§7Forjada en dragones", "§7Nivel 10"])
    .addEnchantment("sharpness", 5)
    .addEnchantment("unbreaking", 3);
```

### Usar un item en un inventario

```js
var inv = player.getInventory();
inv.addItemStack(sword);
```
