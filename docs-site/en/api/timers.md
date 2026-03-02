# Timers

Timers let you run code after a delay or repeatedly. The signature is the same as in browser JavaScript.

## `setTimeout(fn, delay)`

Runs a function **once** after the specified delay in milliseconds.

```js
setTimeout(function() {
    console.log("5 seconds passed.");
}, 5000);
```

::: warning Argument order
First the **function**, then the **delay** in ms. Reversed order causes silent bugs.
:::

## `setInterval(fn, delay)`

Runs a function **repeatedly** every `delay` milliseconds. Returns a function that cancels the interval.

```js
var cancel = setInterval(function() {
    console.log("Running every 10 seconds.");
}, 10000);

setTimeout(function() {
    cancel(); // stop after 1 minute
}, 60000);
```

## `clearInterval(fn)`

Cancels an interval using the function returned by `setInterval`.

```js
var id = setInterval(function() {
    console.log("tick");
}, 1000);

clearInterval(id);
```

## Notes

- Always cancel intervals in `onDisable()` to avoid goroutine leaks.
- Very small delays (< 1ms) may not be precise depending on the OS.

::: warning Do not use world-reading methods from timers
Timers run in separate goroutines **without an active world transaction**. Calling methods like `world.getEntities()`, `world.getEntitiesInRadius()` or `world.getBlock()` from a timer can cause a **deadlock** (the server freezes).

```js
// ❌ Incorrect — may cause deadlock
setInterval(function() {
    var entities = world.getEntities(); // DO NOT do this
}, 5000);

// ✅ Correct — only use inside events or commands
events.on("PlayerJoin", function(event) {
    var entities = world.getEntities(); // OK
});
```

From timers you can safely use: `console.log`, `server.broadcast`, `player.sendMessage` (if you have a player reference), `config`, and **write** methods like `world.spawnEntity` or `world.setBlock`.
:::

```js
var intervalId;

function onEnable() {
    intervalId = setInterval(function() {
        console.log("tick");
    }, 1000);
}

function onDisable() {
    if (intervalId) clearInterval(intervalId);
    config.save();
}
```
