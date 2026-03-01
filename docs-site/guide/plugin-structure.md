# Estructura de un plugin

## Ciclo de vida

Todo plugin tiene tres fases:

```
Servidor arranca
      ↓
 Código raíz del script se ejecuta    ← setDefaults(), variables globales
      ↓
   onEnable() se llama                ← registrar eventos y comandos
      ↓
   [El servidor está corriendo]
      ↓
   onDisable() se llama              ← guardar config, cleanup
      ↓
 Servidor se cierra
```

::: warning Importante
`events.on()` y `commands.register()` **deben llamarse dentro de `onEnable()`**, no en el nivel raíz del script. El nivel raíz se ejecuta antes de que el sistema esté listo para recibir registros.
:::

## Variables globales disponibles

Estas variables están disponibles en cualquier parte del script:

| Variable | Tipo | Descripción |
|---|---|---|
| `plugin` | objeto | Metadatos del plugin |
| `console` | objeto | Logger del servidor |
| `events` | objeto | Para registrar listeners |
| `commands` | objeto | Para registrar comandos |
| `config` | objeto | Configuración YAML del plugin |
| `setTimeout` | función | Ejecutar código tras un delay |
| `setInterval` | función | Ejecutar código repetidamente |
| `clearInterval` | función | Cancelar un interval |

### `plugin`

```js
plugin.name        // "MiPlugin"
plugin.version     // "1.0.0"
plugin.author      // "TuNombre"
plugin.dataFolder  // "plugins/mi-plugin"
```

## Template completo

```js
// =========================================
// Mi Plugin — Dragonfly Script API
// =========================================

console.log("Cargando: " + plugin.name + " v" + plugin.version);

// Configuración por defecto
// El archivo config.yml se crea en plugins/mi-plugin/config.yml
config.setDefaults({
    "prefix":          "§a[MiPlugin]§r",
    "welcome-message": "¡Bienvenido al servidor!",
    "max-players":     20
});

// =========================================
// Ciclo de vida
// =========================================

function onEnable() {
    // Leer configuración
    var prefix  = config.getString("prefix", "§f[?]");
    var welcome = config.getString("welcome-message", "Bienvenido");

    // Registrar eventos
    events.on("PlayerJoin", function(event) {
        var player = event.getPlayer();
        player.sendMessage(prefix + " " + welcome);
    });

    // Registrar comandos
    commands.register("hola", "Dice hola", function(player, args) {
        player.sendMessage(prefix + " ¡Hola, " + player.getName() + "!");
    });

    console.log("Plugin habilitado con prefix: " + prefix);
}

function onDisable() {
    config.save(); // guardar config al cerrar
    console.log("Plugin deshabilitado.");
}

// Exportar (obligatorio)
module = {
    onEnable:  onEnable,
    onDisable: onDisable
};
```

## Colores en mensajes

Minecraft Bedrock usa el símbolo `§` seguido de un código para aplicar colores:

| Código | Color/Formato |
|---|---|
| `§a` | Verde |
| `§c` | Rojo |
| `§e` | Amarillo |
| `§f` | Blanco |
| `§7` | Gris |
| `§b` | Aqua |
| `§d` | Magenta |
| `§r` | Reset (volver al default) |
| `§l` | Negrita |
| `§o` | Cursiva |

```js
player.sendMessage("§aVerde §cRojo §eAmarillo §r Normal");
```
