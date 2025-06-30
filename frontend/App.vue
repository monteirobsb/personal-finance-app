<template>
  <div id="app">
    <router-view /> <!-- Onde os componentes de rota serão renderizados -->
  </div>
</template>

<script>
export default {
  name: 'App',
};
</script>

<style>
/* Estilos globais podem ir aqui ou em um arquivo CSS separado importado no main.js */
body {
  margin: 0; /* Remove margens padrão do body */
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50; /* Cor de texto padrão, pode ser sobrescrita por componentes */
}

#app {
  /* #app ocupa toda a altura da viewport se necessário,
     mas o conteúdo interno (como o dashboard) controlará seu próprio min-height */
}

/* Você pode querer adicionar estilos mais globais aqui, como resets ou um tema base. */
</style>
