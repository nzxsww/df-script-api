---
layout: home

hero:
  name: "Dragonfly Script API"
  text: "Plugins JavaScript para Minecraft Bedrock"
  tagline: Escribe plugins en JavaScript para tu servidor Bedrock. Sistema de eventos estilo Bukkit/Spigot, comandos, configuración YAML y mucho más.
  actions:
    - theme: brand
      text: Primeros pasos →
      link: /guide/getting-started
    - theme: alt
      text: Ver en GitHub
      link: https://github.com/nzxsww/dragonfly-script-api

features:
  - icon: 🎯
    title: Eventos estilo Bukkit
    details: Sistema de eventos con prioridades y cancelación. Los mismos patrones que conocés de Minecraft Java, pero para Bedrock.
  - icon: ⚡
    title: JavaScript puro
    details: Escribe tus plugins en JS sin compilar nada. Coloca tu plugin en la carpeta plugins/ y el servidor lo carga automáticamente.
  - icon: 🛠️
    title: API completa del jugador
    details: Más de 50 métodos para interactuar con los jugadores — posición, vida, inventario, modo de juego, sonidos y mucho más.
  - icon: 💾
    title: Config YAML automática
    details: Cada plugin tiene su propia configuración en YAML. Define defaults, lee y escribe valores, y guarda al cerrar el servidor.
  - icon: 📦
    title: 20 eventos disponibles
    details: PlayerJoin, PlayerChat, BlockBreak, PlayerDeath, PlayerHurt y 15 eventos más para cubrir casi cualquier caso de uso.
  - icon: 🔧
    title: Comandos con autocompletado
    details: Registra comandos desde JS con un solo método. Dragonfly los envía al cliente automáticamente para autocompletado.
---
