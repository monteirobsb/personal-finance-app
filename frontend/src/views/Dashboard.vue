<template>
  <div class="dashboard-container" :style="{ backgroundColor: dynamicBackgroundColor }">
    <div class="dashboard-content">
      <!-- Conteúdo do Dashboard virá aqui, conforme Wireframe 3 -->
      <!-- Por enquanto, apenas uma mensagem e a demonstração da cor de fundo -->
      <h1>Dashboard Principal</h1>
      <p>Saldo Disponível: R$ {{ availableBalance.toFixed(2) }}</p>
      <p>Total de Receitas: R$ {{ totalIncome.toFixed(2) }}</p>
      <p>Total de Despesas Fixas: R$ {{ totalFixedExpenses.toFixed(2) }}</p>
      <p>Saúde Financeira ({{ healthPercentage.toFixed(2) }}%): {{ financialHealthStatus }}</p>

      <!-- Simulação de atualização de dados para testar a cor -->
      <div class="controls">
        <label>
          Simular Renda:
          <input type="number" v-model.number="simulatedIncome" @input="updateFinancials" />
        </label>
        <label>
          Simular Despesas Fixas:
          <input type="number" v-model.number="simulatedFixedExpenses" @input="updateFinancials" />
        </label>
      </div>
    </div>
  </div>
</template>

<script>
// import axios from 'axios'; // Para buscar dados reais no futuro

export default {
  name: 'DashboardView', // Renomeado para evitar conflito com 'Dashboard' se houver componente com esse nome
  data() {
    return {
      totalIncome: 5000, // Valor de exemplo, virá da API
      totalFixedExpenses: 2000, // Valor de exemplo, virá da API
      // Para simulação e teste da cor de fundo:
      simulatedIncome: 5000,
      simulatedFixedExpenses: 2000,
      // financialData: null, // Para armazenar dados buscados da API
      // isLoading: true,
      // errorMessage: '',
    };
  },
  computed: {
    availableBalance() {
      // No futuro, isso pode incluir outras despesas variáveis, etc.
      // Por agora, é a diferença entre renda e despesas fixas.
      return this.totalIncome - this.totalFixedExpenses;
    },
    netFlow() {
      // Receitas - Despesas Fixas
      return this.totalIncome - this.totalFixedExpenses;
    },
    healthPercentage() {
      if (this.totalIncome <= 0) {
        return 0; // Evita divisão por zero e considera 0% se não há renda
      }
      // Percentual do (Receita - Despesas Fixas) em relação à Receita Total
      // Se Despesas > Receitas, o netFlow é negativo, o que resultará em < 0%
      const percentage = (this.netFlow / this.totalIncome) * 100;
      return Math.max(percentage, 0); // Garante que não seja negativo, mínimo 0% para a lógica de cor
                                     // Se quiser mostrar percentuais negativos, ajuste aqui e na lógica de cor.
                                     // Para a lógica de cores pedida (Ex: <25% = vermelho), um netFlow negativo já cairia aí.
    },
    dynamicBackgroundColor() {
      const percentage = this.healthPercentage;
      if (percentage > 60) {
        return '#4CAF50'; // Verde
      } else if (percentage >= 25 && percentage <= 60) {
        return '#FFC107'; // Amarelo/Laranja
      } else {
        return '#F44336'; // Vermelho
      }
    },
    financialHealthStatus() {
      const percentage = this.healthPercentage;
      if (percentage > 60) {
        return 'Saudável';
      } else if (percentage >= 25 && percentage <= 60) {
        return 'Atenção';
      } else {
        return 'Crítico';
      }
    }
  },
  methods: {
    updateFinancials() {
      // Atualiza os valores principais com base na simulação
      this.totalIncome = this.simulatedIncome;
      this.totalFixedExpenses = this.simulatedFixedExpenses;
    },
    // async fetchFinancialData() {
    //   this.isLoading = true;
    //   this.errorMessage = '';
    //   const token = localStorage.getItem('authToken');
    //   if (!token) {
    //     this.errorMessage = 'Autenticação necessária.';
    //     this.isLoading = false;
    //     this.$router.push('/'); // Ou login
    //     return;
    //   }
    //   try {
    //     // Exemplo: buscar renda e despesas. Você precisará de endpoints para isso.
    //     // const incomeResponse = await axios.get('/api/user/income', { headers: { Authorization: `Bearer ${token}` } });
    //     // const expensesResponse = await axios.get('/api/user/fixed-expenses', { headers: { Authorization: `Bearer ${token}` } });
    //     // this.totalIncome = incomeResponse.data.monthlyIncome || 0;
    //     // this.totalFixedExpenses = expensesResponse.data.reduce((sum, exp) => sum + exp.value, 0) || 0;
    //     // this.simulatedIncome = this.totalIncome; // Atualiza simulação
    //     // this.simulatedFixedExpenses = this.totalFixedExpenses; // Atualiza simulação
    //
    //     // Por enquanto, usamos dados de exemplo:
    //     this.simulatedIncome = this.totalIncome;
    //     this.simulatedFixedExpenses = this.totalFixedExpenses;
    //
    //   } catch (error) {
    //     console.error('Erro ao buscar dados financeiros:', error);
    //     this.errorMessage = 'Falha ao carregar dados do dashboard.';
    //     if (error.response && error.response.status === 401) {
    //       this.$router.push('/'); // Token inválido/expirado
    //     }
    //   } finally {
    //     this.isLoading = false;
    //   }
    // }
  },
  created() {
    // this.fetchFinancialData(); // Chamaria para buscar dados reais
    this.updateFinancials(); // Inicializa com os valores de simulação (que são os de exemplo)
  },
};
</script>

<style scoped>
.dashboard-container {
  min-height: 100vh;
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  color: white; /* Cor do texto padrão para bom contraste com fundos coloridos */
  padding: 20px;
  box-sizing: border-box;
  transition: background-color 0.8s ease-in-out; /* Transição suave da cor de fundo */
}

.dashboard-content {
  background-color: rgba(0, 0, 0, 0.2); /* Um overlay semi-transparente para legibilidade */
  padding: 30px;
  border-radius: 10px;
  box-shadow: 0 5px 15px rgba(0,0,0,0.2);
}

h1 {
  font-size: 2.5em;
  margin-bottom: 20px;
}

p {
  font-size: 1.2em;
  margin-bottom: 10px;
}

.controls {
  margin-top: 30px;
  padding: 20px;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.controls label {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 1em;
}

.controls input[type="number"] {
  padding: 8px;
  border-radius: 4px;
  border: 1px solid #ccc;
  margin-left: 10px;
  width: 100px; /* Largura para os inputs de simulação */
  color: #333; /* Cor do texto dentro do input */
}
</style>
