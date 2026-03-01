# ¿Qué es Dragonfly Script API?

## El problema

[Dragonfly](https://github.com/df-mc/dragonfly) es una librería en Go para crear servidores de Minecraft Bedrock Edition. Es muy poderosa, pero tiene una API de bajo nivel: para reaccionar a cualquier cosa que haga un jugador hay que implementar una interfaz con más de 35 métodos en Go.

Esto significa que para hacer algo tan simple como enviarle un mensaje a un jugador cuando entra al servidor, hay que escribir bastante código en Go y compilarlo cada vez que hacés un cambio.

## La solución

**Dragonfly Script API** agrega una capa encima de Dragonfly que permite:

- Escribir plugins en **JavaScript** — sin compilar, sin tipado estricto
- Usar un sistema de **eventos estilo Bukkit/Spigot** familiar para devs de Minecraft Java
- Reaccionar a eventos del servidor con simples callbacks: `events.on("PlayerJoin", fn)`
- Registrar **comandos** que los jugadores pueden ejecutar con `/`
- Guardar **configuración** en archivos YAML por plugin
- Acceder a **50+ métodos** del jugador desde JS

## Compatibilidad

| Componente | Versión |
|---|---|
| Minecraft Bedrock | 1.21.130, 1.21.131, 1.21.132 (protocolo 898) |
| Dragonfly | v0.10.10 |
| Go | 1.21+ |
| Node.js | No requerido (JS se ejecuta en el servidor via Goja) |

::: warning Nota de compatibilidad
Los plugins JS se ejecutan **dentro del servidor Go** usando el motor [Goja](https://github.com/dop251/goja). No necesitás Node.js instalado. El JavaScript que escribís se ejecuta en el servidor, no en el cliente.
:::

## ¿Cómo funciona por dentro?

```
Jugador hace algo en Minecraft
         ↓
  Dragonfly detecta el evento
         ↓
  dragonflyHandler lo traduce
         ↓
  PluginManager lo despacha
         ↓
  Tu callback JS lo recibe
         ↓
  Podés leer, modificar o cancelar
```

Los plugins JS están completamente integrados al mismo sistema de eventos que cualquier otro componente del servidor. No hay magia — es el mismo sistema para todos.
