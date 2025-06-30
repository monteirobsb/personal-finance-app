<template>
  <div class="dashboard-container" :style="{ backgroundColor: dynamicBackgroundColor }">
    <div class="dashboard-main-content">
      <div class="saldo-disponivel-label">Saldo Disponível</div>
      <div class="valor-principal">
        R$ {{ formatCurrency(currentBalance) }}
      </div>

      <div v-if="isLoading" class="loading-spinner">
        <!-- Adicionar um spinner/loading visual aqui -->
        Carregando...
      </div>

      <div v-if="errorMessage" class="error-message-dashboard">
        {{ errorMessage }}
      </div>

      <div v-if="!isLoading && !errorMessage && projection" class="projection-area">
        <div class="projection-card">
          <span class="projection-icon">📊</span> <!-- Ícone de gráfico/alerta -->
          <p v-if="projection.yellowAlertDay && financialHealthStatus === 'amarelo'">
            Seu ritmo de gastos atual indica que você ficará no amarelo por volta do dia {{ formatAlertDay(projection.yellowAlertDay) }}.
          </p>
          <p v-else-if="projection.redAlertDay && financialHealthStatus === 'vermelho'">
            Atenção! Seu ritmo de gastos atual indica que você poderá ficar no vermelho por volta do dia {{ formatAlertDay(projection.redAlertDay) }}.
          </p>
           <p v-else-if="projection.endOfMonthBalance !== undefined">
            Projeção no fim do mês: R$ {{ formatCurrency(projection.endOfMonthBalance) }}.
          </p>
        </div>
      </div>
       <!-- Área de Projeção (Vazia/Oculta) - Se não houver projeção, este bloco não renderiza -->


      <!-- Simulação de atualização de dados para testar a cor (REMOVER EM PRODUÇÃO) -->
      <!--
      <div class="controls" style="background: rgba(0,0,0,0.3); padding: 10px; margin-top: 20px; border-radius: 5px; color: white;">
        <p>Controles de Simulação (DEV)</p>
        <label> Simular Renda: <input type="number" v-model.number="simulatedIncome" @input="updateFinancialsFromSimulation" /> </label>
        <label> Simular Despesas Fixas: <input type="number" v-model.number="simulatedFixedExpenses" @input="updateFinancialsFromSimulation" /> </label>
        <label> Simular Despesas Variáveis: <input type="number" v-model.number="simulatedVariableExpenses" @input="updateFinancialsFromSimulation" /> </label>
         <small>Dias no Mês (Proj): {{ daysInMonthForProjection }} / Dia Atual (Proj): {{ dayOfMonthForProjection }}</small>
      </div>
      -->
    </div>

    <div v-if="!isLoading && !errorMessage && (financialHealthStatus === 'amarelo' || financialHealthStatus === 'vermelho')" class="mensagem-apoio">
      <p v-if="financialHealthStatus === 'amarelo'">
        Ainda está tudo sob controle, mas vamos ficar de olho nas próximas despesas.
      </p>
      <p v-if="financialHealthStatus === 'vermelho'">
        O sinal está vermelho. Que tal rever os gastos ou buscar novas fontes de renda?
      </p>
    </div>

    <button class="fab" @click="openAddExpenseModal" title="Adicionar Despesa">+</button>

    <AddExpenseModal
      v-if="showAddExpenseModal"
      @close="closeAddExpenseModal"
      @save-expense="handleSaveExpense"
    />
  </div>
</template>

<script>
import axios from 'axios';
import AddExpenseModal from '@/components/AddExpenseModal.vue';

export default {
  name: 'DashboardView',
  components: {
    AddExpenseModal, // Registra o componente modal
  },
  data() {
    return {
      currentBalance: 0,
      totalIncome: 0,
      projection: null,
      financialHealthStatus: 'verde',
      healthPercentage: 100,

      isLoading: true,
      errorMessage: '',
      showAddExpenseModal: false, // Controla a visibilidade do modal
    };
  },
  computed: {
    // availableBalance, netFlow, healthPercentage são agora diretamente do backend
    // ou calculados com base nos dados do backend.
    // A cor de fundo é determinada pelo healthPercentage vindo do backend.
    dynamicBackgroundColor() {
      // Usa o healthPercentage recebido do backend
      const percentage = this.healthPercentage;
      if (percentage > 60) {
        return '#28a745'; // Verde (Bootstrap success green)
      } else if (percentage >= 25 && percentage <= 60) {
        return '#ffc107'; // Amarelo (Bootstrap warning yellow)
      } else {
        return '#dc3545'; // Vermelho (Bootstrap danger red)
      }
    },
  },
  methods: {
    openAddExpenseModal() {
      this.showAddExpenseModal = true;
    },
    closeAddExpenseModal() {
      this.showAddExpenseModal = false;
    },
    async handleSaveExpense(expenseData) {
      // Lógica para chamar a API POST /expenses
      // Esta é a implementação da próxima etapa do plano (lógica de POST /expenses)
      // Por enquanto, apenas logamos e atualizamos o dashboard para simular
      console.log('Despesa a ser salva:', expenseData);
      this.showAddExpenseModal = false; // Fecha o modal

      const token = localStorage.getItem('authToken');
      if (!token) {
        this.errorMessage = 'Sessão expirada. Faça login novamente.';
        this.$router.push('/');
        return;
      }

      try {
        // Adiciona um pequeno delay para simular a chamada de API e ver o loading
        this.isLoading = true;
        // await new Promise(resolve => setTimeout(resolve, 500)); // Simula delay

        const response = await axios.post('/api/expenses', expenseData, {
          headers: { Authorization: `Bearer ${token}` },
        });

        if (response.status === 201) {
          console.log('Despesa salva com sucesso, atualizando dashboard...');
          this.fetchDashboardData(); // Recarrega os dados do dashboard
        } else {
           this.errorMessage = response.data.message || 'Erro ao salvar despesa.';
           this.isLoading = false;
        }
      } catch (error) {
        console.error('Erro ao salvar despesa:', error);
        if (error.response && error.response.data && error.response.data.error) {
          this.errorMessage = error.response.data.error;
        } else {
          this.errorMessage = 'Não foi possível salvar a despesa.';
        }
        this.isLoading = false; // Garante que o loading pare em caso de erro
      }
      // this.fetchDashboardData(); // Atualiza os dados do dashboard após salvar
    },
    formatCurrency(value) {
      if (typeof value !== 'number') {
        return '0,00';
      }
      return value.toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
    },
    formatAlertDay(dateString) {
      if (!dateString) return '';
      try {
        const date = new Date(dateString + 'T00:00:00');
        return `dia ${date.getDate()}`;
      } catch (e) {
        return dateString;
      }
    },
    async fetchDashboardData() {
      this.isLoading = true; // Movido para cá para garantir que é setado antes da chamada
      this.errorMessage = '';
      const token = localStorage.getItem('authToken');

      if (!token) {
        this.errorMessage = 'Sessão expirada ou inválida. Por favor, faça login novamente.';
        this.isLoading = false;
        this.$router.push('/'); // Redireciona para login
        return;
      }

      try {
        const response = await axios.get('/api/balance', {
          headers: { Authorization: `Bearer ${token}` },
        });

        const data = response.data;
        this.currentBalance = data.currentBalance;
        this.totalIncome = data.totalIncome;
        this.projection = data.projection; // Pode ser null
        this.financialHealthStatus = data.financialHealthStatus;
        this.healthPercentage = data.healthPercentage;

        // Para debug (opcional)
        // this.daysInMonthForProjection = data.daysInMonthForProjection;
        // this.dayOfMonthForProjection = data.dayOfMonthForProjection;

      } catch (error) {
        console.error('Erro ao buscar dados do dashboard:', error);
        if (error.response) {
          if (error.response.status === 401) {
            this.errorMessage = 'Sessão expirada. Por favor, faça login novamente.';
            localStorage.removeItem('authToken'); // Limpa token inválido
            this.$router.push('/');
          } else if (error.response.data && error.response.data.error) {
            this.errorMessage = error.response.data.error;
          } else {
            this.errorMessage = 'Falha ao carregar dados do dashboard. Tente novamente mais tarde.';
          }
        } else {
          this.errorMessage = 'Não foi possível conectar ao servidor.';
        }
      } finally {
        this.isLoading = false;
      }
    },
    // updateFinancialsFromSimulation() { // Função de simulação
    //   this.totalIncome = this.simulatedIncome;
    //   // Recalcular com base nos dados simulados para teste de cor
    //   const netFlowSim = this.simulatedIncome - this.simulatedFixedExpenses - this.simulatedVariableExpenses;
    //   this.currentBalance = netFlowSim; // Simplificado para teste

    //   if (this.simulatedIncome > 0) {
    //     this.healthPercentage = (netFlowSim / this.simulatedIncome) * 100;
    //   } else {
    //     this.healthPercentage = 0;
    //   }

    //   if (this.healthPercentage > 60) this.financialHealthStatus = 'verde';
    //   else if (this.healthPercentage >= 25) this.financialHealthStatus = 'amarelo';
    //   else this.financialHealthStatus = 'vermelho';

    //   // Simular projeção para teste visual
    //   if(this.financialHealthStatus === 'amarelo'){
    //     this.projection = { yellowAlertDay: '2023-12-25' }; // Data de exemplo
    //   } else if (this.financialHealthStatus === 'vermelho') {
    //     this.projection = { redAlertDay: '2023-12-20' }; // Data de exemplo
    //   } else {
    //     this.projection = null;
    //   }
    // }
  },
  created() {
    this.fetchDashboardData();
    // this.updateFinancialsFromSimulation(); // Para teste inicial da UI com simulação
  },
};
</script>

<style scoped>
.dashboard-container {
  min-height: 100vh;
  width: 100%;
  display: flex;
  flex-direction: column; /* Garante que o conteúdo principal e a mensagem de apoio fiquem empilhados */
  justify-content: space-between; /* Empurra a mensagem de apoio para baixo se houver espaço */
  align-items: center;
  text-align: center;
  color: white;
  padding: 20px;
  box-sizing: border-box;
  transition: background-color 0.8s ease-in-out;
  position: relative;
}

.dashboard-main-content {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  justify-content: center; /* Centraliza o conteúdo principal verticalmente */
  align-items: center;
  width: 100%;
}


.saldo-disponivel-label {
  font-size: 1.2em; /* Texto pequeno no topo */
  font-weight: 300;
  margin-bottom: 5px;
  color: rgba(255, 255, 255, 0.8);
}

.valor-principal {
  font-size: 5rem; /* Valor Gigante */
  font-weight: bold;
  line-height: 1;
  margin-bottom: 20px; /* Espaço antes da área de projeção */
  word-break: break-all; /* Quebra o valor se for muito grande para a tela */
}

@media (max-width: 600px) {
  .valor-principal {
    font-size: 3.5rem; /* Reduz um pouco em telas menores */
  }
  .saldo-disponivel-label {
    font-size: 1em;
  }
}


.projection-area {
  margin-top: 10px;
  margin-bottom: 20px;
  width: 100%;
  max-width: 450px; /* Limita a largura da área de projeção */
}

.projection-card {
  background-color: rgba(0, 0, 0, 0.25); /* Fundo semi-transparente */
  padding: 15px 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.15);
  font-size: 0.95em;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.projection-icon {
  font-size: 1.5em;
}

.mensagem-apoio {
  font-size: 0.9em;
  color: rgba(255, 255, 255, 0.7);
  padding: 10px 20px;
  width: 100%;
  text-align: center;
  /* position: absolute; */ /* Se quisesse fixar no rodapé absoluto */
  /* bottom: 20px; */
  margin-top: auto; /* Empurra para baixo se o .dashboard-main-content não preencher tudo */
}

.loading-spinner, .error-message-dashboard {
  margin-top: 20px;
  font-size: 1.1em;
}

.error-message-dashboard {
  color: #f8d7da; /* Cor de erro Bootstrap */
  background-color: #721c24; /* Fundo para erro Bootstrap */
  padding: 10px;
  border-radius: 5px;
}

/* FAB (Botão de Ação Flutuante) */
.fab {
  position: fixed;
  bottom: 30px;
  right: 30px;
  background-color: #007bff; /* Azul primário */
  color: white;
  width: 60px;
  height: 60px;
  border-radius: 50%;
  border: none;
  font-size: 28px; /* Tamanho do ícone '+' */
  line-height: 60px; /* Centraliza o '+' verticalmente */
  text-align: center;
  box-shadow: 0 4px 10px rgba(0,0,0,0.25);
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.2s ease;
  z-index: 999; /* Para ficar sobre outros elementos se houver */
}
.fab:hover {
  background-color: #0056b3; /* Azul mais escuro no hover */
  transform: translateY(-2px); /* Leve elevação no hover */
}
.fab:active {
  transform: translateY(0px); /* Volta ao normal ao clicar */
}

</style>
