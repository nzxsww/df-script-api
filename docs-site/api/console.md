# Console

El objeto `console` permite imprimir mensajes en la consola del servidor. Todos los mensajes incluyen automáticamente el prefijo `[NombrePlugin]`.

## Métodos

| Método | Cuándo usarlo |
|---|---|
| `console.log(msg)` | Mensajes informativos generales |
| `console.info(msg)` | Información importante |
| `console.warn(msg)` | Advertencias |
| `console.error(msg)` | Errores |

```js
console.log("Plugin cargado correctamente");
// Salida: [MiPlugin] Plugin cargado correctamente

console.info("Servidor iniciado en puerto 19132");
// Salida: [MiPlugin] [INFO] Servidor iniciado en puerto 19132

console.warn("Config no encontrada, usando defaults");
// Salida: [MiPlugin] [WARN] Config no encontrada, usando defaults

console.error("No se pudo guardar la config");
// Salida: [MiPlugin] [ERROR] No se pudo guardar la config
```

## Concatenar valores

```js
var jugadores = 5;
var max = 20;
console.log("Jugadores conectados: " + jugadores + "/" + max);

var player = event.getPlayer();
console.log(player.getName() + " se unió. Latencia: " + player.getLatency() + "ms");
```
