import { createRouter, createWebHistory } from 'vue-router';
import OnboardingIncome from '../views/OnboardingIncome.vue';
import OnboardingExpenses from '../views/OnboardingExpenses.vue';
import DashboardView from '../views/Dashboard.vue'; // Nome do componente é DashboardView

// Simulação de uma tela de Login/Autenticação inicial, caso o usuário não esteja autenticado
// ou para onde redirecionar se o token não existir.
// Você pode criar este componente depois.
const LoginView = { template: '<div><h1>Login Page</h1><p>Imagine um formulário de login aqui ou a tela de /auth/request-code.</p><router-link to="/onboarding/income">Ir para Onboarding Renda (simulado)</router-link></div>' };

const routes = [
  {
    path: '/',
    name: 'Login', // Ou uma tela de boas-vindas que redireciona
    component: LoginView, // Temporário, idealmente seria sua tela de /auth/request-code ou similar
    // meta: { requiresGuest: true } // Para redirecionar se já estiver logado
  },
  {
    path: '/onboarding/income',
    name: 'OnboardingIncome',
    component: OnboardingIncome,
    meta: { requiresAuth: true }, // Exemplo de meta para rotas protegidas
  },
  {
    path: '/onboarding/expenses',
    name: 'OnboardingExpenses',
    component: OnboardingExpenses,
    meta: { requiresAuth: true },
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardView,
    meta: { requiresAuth: true },
  },
  // Adicione um catch-all para rotas não encontradas, se desejar
  // { path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFoundComponent },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL || '/'), // BASE_URL é configurado pelo Vue CLI
  routes,
});

// Exemplo de Navigation Guard (gancho de navegação)
router.beforeEach((to, from, next) => {
  const isAuthenticated = !!localStorage.getItem('authToken'); // Verifica se o token existe

  if (to.matched.some(record => record.meta.requiresAuth) && !isAuthenticated) {
    // Se a rota requer autenticação e o usuário não está autenticado,
    // redireciona para a página de login (ou a primeira etapa da autenticação).
    // Aqui, vamos redirecionar para a rota raiz '/' que configuramos como LoginView.
    next({ name: 'Login' });
  } else if (to.matched.some(record => record.meta.requiresGuest) && isAuthenticated) {
    // Se a rota é para visitantes (ex: login) e o usuário já está autenticado,
    // redireciona para o dashboard.
    next({ name: 'Dashboard' });
  }
  else {
    // Caso contrário, permite a navegação.
    next();
  }
});

export default router;
