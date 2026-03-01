# Primeros pasos

## Requisitos

- **Go 1.21+** instalado
- El servidor compilado (`go build -o server.exe .`)
- Un editor de texto para escribir JS

## Crear tu primer plugin

### 1. Crear la carpeta del plugin

Dentro de la carpeta `plugins/` del servidor, crea una carpeta con el nombre de tu plugin:

```
plugins/
└── mi-primer-plugin/    ← creá esta carpeta
```

### 2. Crear `plugin.yml`

Dentro de esa carpeta, crea el archivo `plugin.yml` con la información de tu plugin:

```yaml
name: MiPrimerPlugin
version: 1.0.0
author: TuNombre
description: Mi primer plugin para Dragonfly Script API
main: index.js
api-version: 1.0.0
```

| Campo | Descripción |
|---|---|
| `name` | Nombre único del plugin. Aparece en los logs del servidor. |
| `version` | Versión en formato semver (1.0.0) |
| `author` | Tu nombre o nick |
| `description` | Descripción corta del plugin |
| `main` | Archivo JS principal (por defecto: `index.js`) |
| `api-version` | Versión de la API (usar `1.0.0`) |

### 3. Crear `index.js`

Crea el archivo `index.js` en la misma carpeta:

```js
// Este código se ejecuta al cargar el plugin (antes de onEnable)
console.log("Cargando: " + plugin.name + " v" + plugin.version);

function onEnable() {
    // Aquí registrás tus eventos y comandos
    console.log("¡Plugin habilitado!");

    events.on("PlayerJoin", function(event) {
        var player = event.getPlayer();
        player.sendMessage("§aHola " + player.getName() + ", bienvenido!");
    });
}

function onDisable() {
    // Limpieza al cerrar el servidor
    console.log("Plugin deshabilitado.");
}

// Exportar el ciclo de vida (obligatorio)
module = {
    onEnable: onEnable,
    onDisable: onDisable
};
```

### 4. Arrancar el servidor

```bash
./server.exe
```

En la consola deberías ver:

```
[Loader] Cargando plugin: MiPrimerPlugin v1.0.0 (autor: TuNombre)
[MiPrimerPlugin] Cargando: MiPrimerPlugin v1.0.0
[MiPrimerPlugin] ¡Plugin habilitado!
```

### 5. Conectarse con Minecraft

Conectate con **Minecraft Bedrock 1.21.130 - 1.21.132** a `localhost:19132`. Al entrar, recibirás el mensaje de bienvenida.

::: tip
Para conectarte desde la misma PC, usá la dirección `127.0.0.1` puerto `19132`.
:::

## Estructura de archivos resultante

```
plugins/
└── mi-primer-plugin/
    ├── plugin.yml     ← metadatos del plugin
    ├── index.js       ← lógica del plugin
    └── config.yml     ← creado automáticamente si usás config.save()
```

## Próximos pasos

- [Estructura detallada de un plugin →](./plugin-structure)
- [Ver todos los eventos disponibles →](/api/events)
- [Ver todos los métodos del jugador →](/api/player)
- [Registrar comandos →](/api/commands)
