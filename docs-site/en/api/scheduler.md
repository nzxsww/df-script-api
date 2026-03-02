# Scheduler API

The Scheduler API lets you execute JS tasks **with an active transaction (Tx)** so you can safely use `world` and `server`, similar to Bukkit.

All callbacks run in a safe context **without deadlocks**.

> ⚠️ Do not use `world` or `server.getPlayers()` inside `setInterval` / `setTimeout`. Use `scheduler` instead.

---

## API

```js
scheduler.run(fn)
scheduler.runLater(ticks, fn)
scheduler.runRepeating(ticks, fn) // returns task.cancel()
```

- **ticks:** one tick equals **50ms** (20 ticks = 1s)
- **fn:** receives `(world, server)` as arguments

---

## Real-world examples

### 1) Run a task immediately
```js
scheduler.run(function(world, server) {
    world.setBlock(0, 64, 0, "minecraft:diamond_block");
});
```

### 2) Run once after 2 seconds
```js
scheduler.runLater(40, function(world, server) {
    server.broadcast("§aEvent started!");
});
```

### 3) Repeat every second
```js
var task = scheduler.runRepeating(20, function(world, server) {
    server.broadcast("§7Players: §f" + server.getPlayerCount());
});

// stop later
// task.cancel();
```

### 4) Use world + player with scoreboard
```js
scheduler.runRepeating(5, function(world, server) {
    var players = server.getPlayers();
    for (var i = 0; i < players.length; i++) {
        var p = players[i];
        var sb = scoreboard.create("§aServer");
        sb.setLines([
            "§7Coords: " + Math.floor(p.getX()) + " " + Math.floor(p.getY()) + " " + Math.floor(p.getZ()),
            "§7Block below: " + world.getBlock(Math.floor(p.getX()), Math.floor(p.getY()) - 1, Math.floor(p.getZ()))
        ]);
        p.sendScoreboard(sb);
    }
});
```

### 5) Access chest inventories
```js
scheduler.run(function(world) {
    var inv = world.getInventory(0, 64, 0);
    if (inv) {
        inv.addItem("minecraft:diamond", 1);
    }
});
```

---

## Why use scheduler?

- `world` and `server` require an active **Tx**.
- `setTimeout` and `setInterval` have no Tx → can deadlock.
- Scheduler always runs inside `World.Exec()` with safe Tx.

---

## Cancellation

```js
var task = scheduler.runRepeating(20, function(world, server) {
    server.broadcast("tick");
});

// cancel anytime
// task.cancel();
```
