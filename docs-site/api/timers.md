# Timers

Los timers permiten ejecutar código después de un delay o de forma repetida. La firma es la misma que en JavaScript del navegador.

## `setTimeout(fn, delay)`

Ejecuta una función **una sola vez** después del delay especificado en milisegundos.

```js
// Ejecutar después de 5 segundos
setTimeout(function() {
    console.log("Pasaron 5 segundos.");
}, 5000);
```

::: warning Orden de argumentos
Primero la **función**, luego el **delay** en ms. Al revés causa bugs silenciosos.
:::

### Usos comunes

```js
// Mensaje de bienvenida con delay
events.on("PlayerJoin", function(event) {
    var player = event.getPlayer();
    var nombre = player.getName();

    // Dar items 3 segundos después de entrar
    setTimeout(function() {
        player.giveItem("minecraft:bread", 10);
        player.sendMessage("§aTe dimos algo de comida para empezar.");
    }, 3000);
});
```

```js
// Mensaje de advertencia antes de reiniciar
setTimeout(function() {
    console.log("El servidor se reinicia en 1 minuto.");
}, 60000);
```

## `setInterval(fn, delay)`

Ejecuta una función **repetidamente** cada `delay` milisegundos. Retorna una función que cancela el interval.

```js
var cancelar = setInterval(function() {
    console.log("Ejecutado cada 10 segundos.");
}, 10000);

// Cancelar después de 1 minuto
setTimeout(function() {
    cancelar(); // detiene el interval
}, 60000);
```

## `clearInterval(fn)`

Cancela un interval usando la función retornada por `setInterval`.

```js
var id = setInterval(function() {
    console.log("tick");
}, 1000);

// Más adelante...
clearInterval(id);
```

## Ejemplo práctico — Anuncio automático

```js
var mensajes = [
    "§e¡Bienvenidos al servidor!",
    "§aUsa §f/hola §apara recibir un saludo.",
    "§b¡Diviértete y respeta las reglas!"
];
var indice = 0;

function onEnable() {
    // Mostrar un anuncio diferente cada 5 minutos
    setInterval(function() {
        var msg = mensajes[indice % mensajes.length];
        console.log("Anuncio: " + msg);
        indice++;
        // Para enviar a todos los jugadores necesitás iterar sobre ellos
    }, 300000); // 5 minutos
}
```

## Notas importantes

- Los timers se ejecutan en **goroutines separadas**. Si accedés a variables compartidas entre timers y callbacks de eventos, podés tener condiciones de carrera. Usá con cuidado.
- Los timers **no se cancelan automáticamente** al desactivar el plugin. Si usás `setInterval`, cancelalo en `onDisable()`.
- El delay mínimo real depende del sistema operativo — delays muy pequeños (< 1ms) pueden no ser precisos.

```js
var intervalId;

function onEnable() {
    intervalId = setInterval(function() {
        console.log("tick");
    }, 1000);
}

function onDisable() {
    if (intervalId) {
        clearInterval(intervalId);
    }
    config.save();
}
```
