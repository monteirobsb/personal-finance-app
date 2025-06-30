import { createApp } from 'vue';
import App from './App.vue';
import router from './router'; // Importa a configuração do router
// import store from './store'; // Se você usar Vuex, descomente

const app = createApp(App);

app.use(router); // Diz ao Vue para usar o router
// app.use(store); // Se você usar Vuex, descomente

app.mount('#app');
