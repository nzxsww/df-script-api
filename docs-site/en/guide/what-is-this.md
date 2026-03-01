# What is Dragonfly Script API?

## The problem

[Dragonfly](https://github.com/df-mc/dragonfly) is a Go library for creating Minecraft Bedrock Edition servers. It's very powerful, but has a low-level API: to react to anything a player does (move, chat, break a block) you need to implement an interface with more than 35 methods in Go.

This means that to do something as simple as sending a message to a player when they join the server, you need to write a lot of Go code and recompile every time you make a change.

## The solution

**Dragonfly Script API** adds a layer on top of Dragonfly that allows:

- Writing plugins in **JavaScript** — no compilation, no strict typing
- Using a **Bukkit/Spigot-style event system** familiar to Minecraft Java devs
- Reacting to server events with simple callbacks: `events.on("PlayerJoin", fn)`
- Registering **commands** that players can run with `/`
- Saving **configuration** in YAML files per plugin
- Accessing **50+ player methods** from JS

## Compatibility

| Component | Version |
|---|---|
| Minecraft Bedrock | 1.21.130, 1.21.131, 1.21.132 (protocol 898) |
| Dragonfly | v0.10.10 |
| Go | 1.21+ |
| Node.js | Not required |

::: warning Compatibility note
JS plugins run **inside the Go server** using the [Goja](https://github.com/dop251/goja) engine. You don't need Node.js installed. The JavaScript you write runs on the server, not the client.
:::

## How it works internally

```
Player does something in Minecraft
         ↓
  Dragonfly detects the event
         ↓
  dragonflyHandler translates it
         ↓
  PluginManager dispatches it
         ↓
  Your JS callback receives it
         ↓
  You can read, modify or cancel it
```

JS plugins are fully integrated into the same event system as any other server component. No magic — it's the same system for everyone.
