# Commands

Commands are registered inside `onEnable()` using `commands.register()`. When a player types `/name` in chat, the callback receives the player and the arguments.

Dragonfly automatically sends registered commands to the client for **autocomplete**.

## Syntax

```js
// Without aliases
commands.register("name", "description", function(player, args) {
    // do something
});

// With aliases
commands.register("name", "description", ["alias1", "alias2"], function(player, args) {
    // do something
});
```

## The callback

```js
function(player, args) {
    // player — the player who ran the command
    // args   — array of strings with the arguments
}
```

## Examples

### Simple command

```js
commands.register("hello", "Greets you", function(player, args) {
    player.sendMessage("§aHello, " + player.getName() + "!");
});
```

### Command with arguments

```js
commands.register("give", "Give items. Usage: /give <item> [count]", function(player, args) {
    if (args.length < 1) {
        player.sendMessage("§cUsage: /give <item> [count]");
        return;
    }

    var item  = args[0];
    var count = args.length >= 2 ? parseInt(args[1]) : 1;

    if (isNaN(count) || count < 1) count = 1;
    if (count > 64) count = 64;

    var ok = player.giveItem(item, count);
    if (ok) {
        player.sendMessage("§aGiven §f" + count + "x §a" + item);
        player.playSound("pop");
    } else {
        player.sendMessage("§cUnknown item or full inventory.");
        player.playSound("deny");
    }
});
```

### Command with aliases

```js
commands.register("gamemode", "Change game mode", ["gm"], function(player, args) {
    if (args.length < 1) {
        player.sendMessage("§cUsage: /gamemode <survival|creative|adventure|spectator>");
        return;
    }
    player.setGameMode(args[0].toLowerCase());
    player.sendMessage("§aGame mode changed to §f" + args[0]);
});
```

## Notes

- Commands can only be run by **players** — not from the server console.
- Arguments are always **strings**. Use `parseInt()` or `parseFloat()` to convert to numbers.
- If no arguments are passed, `args` is an empty array `[]`.
