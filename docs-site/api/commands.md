# Comandos

Los comandos se registran dentro de `onEnable()` usando `commands.register()`. Cuando el jugador escribe `/nombre` en el chat, el callback recibe el jugador y los argumentos.

Dragonfly envía los comandos registrados al cliente automáticamente, por lo que aparecen en el **autocompletado** del juego.

## Sintaxis

```js
// Sin aliases
commands.register("nombre", "descripción", function(player, args) {
    // hacer algo
});

// Con aliases
commands.register("nombre", "descripción", ["alias1", "alias2"], function(player, args) {
    // hacer algo
});
```

| Parámetro | Tipo | Descripción |
|---|---|---|
| `nombre` | `string` | Nombre del comando (sin `/`) |
| `descripción` | `string` | Descripción que aparece en el autocompletado |
| `aliases` | `array` (opcional) | Nombres alternativos para el mismo comando |
| `callback` | `function` | Función que se ejecuta cuando el jugador usa el comando |

## El callback

El callback recibe dos parámetros:

```js
function(player, args) {
    // player — el jugador que ejecutó el comando
    // args   — array de strings con los argumentos
}
```

```js
commands.register("tp", "Teleportar a coordenadas. Uso: /tp <x> <y> <z>", function(player, args) {
    if (args.length < 3) {
        player.sendMessage("§cUso: /tp <x> <y> <z>");
        return;
    }

    var x = parseFloat(args[0]);
    var y = parseFloat(args[1]);
    var z = parseFloat(args[2]);

    if (isNaN(x) || isNaN(y) || isNaN(z)) {
        player.sendMessage("§cLas coordenadas deben ser números.");
        return;
    }

    player.teleport(x, y, z);
    player.sendMessage("§aTeleportado a §f" + x + ", " + y + ", " + z);
    player.playSound("pop");
});
```

## Ejemplos

### Comando simple

```js
commands.register("hola", "Te saluda", function(player, args) {
    player.sendMessage("§aHola, " + player.getName() + "!");
});
```

### Comando con subcomandos

```js
commands.register("gamemode", "Cambiar modo de juego", ["gm"], function(player, args) {
    if (args.length < 1) {
        player.sendMessage("§cUso: /gamemode <survival|creative|adventure|spectator>");
        return;
    }

    var modo = args[0].toLowerCase();
    var modos = ["survival", "creative", "adventure", "spectator"];

    if (modos.indexOf(modo) === -1) {
        player.sendMessage("§cModo inválido. Opciones: " + modos.join(", "));
        return;
    }

    player.setGameMode(modo);
    player.sendMessage("§aModo cambiado a §f" + modo);
});
```

### Comando con número variable de args

```js
commands.register("give", "Dar items. Uso: /give <item> [cantidad]", function(player, args) {
    if (args.length < 1) {
        player.sendMessage("§cUso: /give <item> [cantidad]");
        return;
    }

    var item  = args[0];
    var count = args.length >= 2 ? parseInt(args[1]) : 1;

    if (isNaN(count) || count < 1) count = 1;
    if (count > 64) count = 64;

    var ok = player.giveItem(item, count);
    if (ok) {
        player.sendMessage("§aSe te dieron §f" + count + "x §a" + item);
        player.playSound("pop");
    } else {
        player.sendMessage("§cItem desconocido o inventario lleno.");
        player.playSound("deny");
    }
});
```

## Notas

- Los comandos solo pueden ser ejecutados por **jugadores** — no por la consola del servidor.
- Los argumentos siempre son **strings**. Convertí a número con `parseInt()` o `parseFloat()` si necesitás.
- Si el jugador no pasa ningún argumento, `args` es un array vacío `[]`.
- Los aliases funcionan exactamente igual que el comando principal.
