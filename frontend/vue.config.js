const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // URL do seu backend Go
        changeOrigin: true,
        pathRewrite: { '^/api': '' }, // Remove /api do path antes de enviar para o backend
      },
    },
    port: 8081, // Porta em que o servidor de desenvolvimento do Vue vai rodar
  },
})
