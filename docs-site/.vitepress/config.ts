import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Dragonfly Script API',
  description: 'Motor de scripting para plugins JavaScript en servidores Minecraft Bedrock',
  base: '/df-script-api/',

  head: [
    ['link', { rel: 'icon', href: '/favicon.ico' }],
  ],

  locales: {
    root: {
      lang: 'es',
      label: 'Español',
      title: 'Dragonfly Script API',
      description: 'Motor de scripting para plugins JavaScript en servidores Minecraft Bedrock',
      themeConfig: {
        logo: '🐉',

        nav: [
          { text: 'Inicio', link: '/' },
          { text: 'Guía', link: '/guide/getting-started' },
          { text: 'API', link: '/api/events' },
          { text: 'Ejemplos', link: '/examples/example-plugin' }
        ],

        sidebar: [
          {
            text: 'Guía',
            items: [
              { text: '¿Qué es esto?', link: '/guide/what-is-this' },
              { text: 'Primeros pasos', link: '/guide/getting-started' },
              { text: 'Estructura de un plugin', link: '/guide/plugin-structure' },
            ]
          },
          {
            text: 'Referencia de API',
            items: [
              { text: 'Eventos', link: '/api/events' },
              { text: 'Objeto Player', link: '/api/player' },
              { text: 'Comandos', link: '/api/commands' },
              { text: 'Configuración YAML', link: '/api/config' },
              { text: 'Timers', link: '/api/timers' },
              { text: 'Console', link: '/api/console' },
            ]
          },
          {
            text: 'Ejemplos',
            items: [
              { text: 'Plugin completo', link: '/examples/example-plugin' },
            ]
          }
        ],

        socialLinks: [
          { icon: 'github', link: 'https://github.com/nzxsww/dragonfly-script-api' }
        ],

        search: { provider: 'local' },

        footer: {
          message: 'Dragonfly Script API',
          copyright: 'Compatible con Minecraft Bedrock 1.21.130 - 1.21.132 (protocolo 898)'
        },

      }
    },

    en: {
      lang: 'en',
      label: 'English',
      title: 'Dragonfly Script API',
      description: 'JavaScript plugin scripting engine for Minecraft Bedrock servers',
      themeConfig: {
        logo: '🐉',

        nav: [
          { text: 'Home', link: '/en/' },
          { text: 'Guide', link: '/en/guide/getting-started' },
          { text: 'API', link: '/en/api/events' },
          { text: 'Examples', link: '/en/examples/example-plugin' }
        ],

        sidebar: [
          {
            text: 'Guide',
            items: [
              { text: 'What is this?', link: '/en/guide/what-is-this' },
              { text: 'Getting Started', link: '/en/guide/getting-started' },
              { text: 'Plugin Structure', link: '/en/guide/plugin-structure' },
            ]
          },
          {
            text: 'API Reference',
            items: [
              { text: 'Events', link: '/en/api/events' },
              { text: 'Player Object', link: '/en/api/player' },
              { text: 'Commands', link: '/en/api/commands' },
              { text: 'YAML Config', link: '/en/api/config' },
              { text: 'Timers', link: '/en/api/timers' },
              { text: 'Console', link: '/en/api/console' },
            ]
          },
          {
            text: 'Examples',
            items: [
              { text: 'Full Plugin Example', link: '/en/examples/example-plugin' },
            ]
          }
        ],

        socialLinks: [
          { icon: 'github', link: 'https://github.com/nzxsww/dragonfly-script-api' }
        ],

        search: { provider: 'local' },

        footer: {
          message: 'Dragonfly Script API',
          copyright: 'Compatible with Minecraft Bedrock 1.21.130 - 1.21.132 (protocol 898)'
        },

      }
    }
  }
})
