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
