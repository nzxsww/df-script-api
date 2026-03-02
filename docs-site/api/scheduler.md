# Scheduler API

El Scheduler API permite ejecutar tareas JS **con transacción activa** (Tx) para que puedas usar `world` y `server` de forma segura, similar a Bukkit.

Todos los callbacks se ejecutan en un contexto seguro y **sin deadlocks**.

> ⚠️ No usar `world` ni `server.getPlayers()` desde timers (`setInterval`, `setTimeout`). Para eso usá `scheduler`.

---

## API

```js
scheduler.run(fn)
scheduler.runLater(ticks, fn)
scheduler.runRepeating(ticks, fn) // retorna task.cancel()
```

- **ticks:** cada tick equivale a **50ms** (20 ticks = 1s).
- **fn:** recibe `(world, server)` como argumentos.

---

## Ejemplos reales

### 1) Ejecutar una tarea inmediata
```js
scheduler.run(function(world, server) {
    world.setBlock(0, 64, 0, "minecraft:diamond_block");
});
```

### 2) Ejecutar una vez después de 2 segundos
```js
scheduler.runLater(40, function(world, server) {
    server.broadcast("§a¡Comenzó el evento!");
});
```

### 3) Repetir cada segundo
```js
var task = scheduler.runRepeating(20, function(world, server) {
    server.broadcast("§7Jugadores: §f" + server.getPlayerCount());
});

// detener luego
// task.cancel();
```

### 4) Usar world + player con scoreboard
```js
scheduler.runRepeating(5, function(world, server) {
    var players = server.getPlayers();
    for (var i = 0; i < players.length; i++) {
        var p = players[i];
        var sb = scoreboard.create("§aServidor");
        sb.setLines([
            "§7Coords: " + Math.floor(p.getX()) + " " + Math.floor(p.getY()) + " " + Math.floor(p.getZ()),
            "§7Bloque debajo: " + world.getBlock(Math.floor(p.getX()), Math.floor(p.getY()) - 1, Math.floor(p.getZ()))
        ]);
        p.sendScoreboard(sb);
    }
});
```

### 5) Acceder a inventarios de cofres
```js
scheduler.run(function(world) {
    var inv = world.getInventory(0, 64, 0);
    if (inv) {
        inv.addItem("minecraft:diamond", 1);
    }
});
```

---

## ¿Por qué usar scheduler?

- `world` y `server` necesitan **Tx activa**.
- `setTimeout` y `setInterval` no tienen Tx → pueden causar deadlocks.
- Scheduler siempre ejecuta dentro de `World.Exec()` con Tx segura.

---

## Cancelación

```js
var task = scheduler.runRepeating(20, function(world, server) {
    server.broadcast("tick");
});

// cancelar cuando quieras
// task.cancel();
```
