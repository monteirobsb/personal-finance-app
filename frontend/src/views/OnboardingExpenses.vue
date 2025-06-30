<template>
  <div class="onboarding-expenses-container">
    <div class="onboarding-card">
      <h1>Quais são suas despesas fixas mensais?</h1>
      <p>Liste suas despesas recorrentes, como aluguel, assinaturas, etc.</p>

      <div v-for="(expense, index) in expenses" :key="index" class="expense-item">
        <input
          type="text"
          v-model="expense.name"
          placeholder="Nome da Despesa (ex: Aluguel)"
          class="expense-name-input"
        />
        <div class="input-group value-input-group">
          <span class="currency-symbol">R$</span>
          <input
            type="number"
            v-model.number="expense.value"
            placeholder="Valor"
            min="0.01"
            step="0.01"
            class="expense-value-input"
          />
        </div>
        <button @click="removeExpense(index)" class="remove-expense-button" v_if="expenses.length > 1">
          &times;
        </button>
      </div>

      <button @click="addExpense" class="add-expense-button">
        + Adicionar Despesa
      </button>

      <button @click="finishOnboarding" :disabled="isFinishButtonDisabled" class="finish-button">
        Concluir
      </button>
      <p v-if="errorMessage" class="error-message">{{ errorMessage }}</p>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'OnboardingExpenses',
  data() {
    return {
      expenses: [
        { name: '', value: null }, // Começa com um item de despesa
      ],
      errorMessage: '',
    };
  },
  computed: {
    isFinishButtonDisabled() {
      // Desabilita se alguma despesa não tiver nome ou valor, ou se o valor não for positivo
      return this.expenses.some(
        (exp) => !exp.name || !exp.value || exp.value <= 0
      );
    },
  },
  methods: {
    addExpense() {
      this.expenses.push({ name: '', value: null });
    },
    removeExpense(index) {
      if (this.expenses.length > 1) { // Garante que sempre haja pelo menos um item
        this.expenses.splice(index, 1);
      }
    },
    async finishOnboarding() {
      if (this.isFinishButtonDisabled) {
        this.errorMessage = 'Por favor, preencha todas as despesas corretamente (nome e valor positivo).';
        return;
      }
      this.errorMessage = '';

      const token = localStorage.getItem('authToken');
      if (!token) {
        this.errorMessage = 'Erro de autenticação. Por favor, faça login novamente.';
        this.$router.push('/'); // Ou para a tela de login/auth
        return;
      }

      const payload = {
        // O backend espera um array de objetos com "nome" e "valor"
        despesasFixas: this.expenses.map(exp => ({ nome: exp.name, valor: exp.value })),
      };

      try {
        const response = await axios.post('/api/onboarding/fixed-expenses', payload, {
          headers: { Authorization: `Bearer ${token}` },
        });

        if (response.status === 201 || response.status === 200) { // 201 é mais comum para criação
          console.log('Despesas salvas com sucesso:', response.data);
          // Navegar para o dashboard principal
          this.$router.push('/dashboard');
        } else {
          this.errorMessage = response.data.message || 'Ocorreu um erro ao salvar as despesas.';
        }
      } catch (error) {
        console.error('Erro ao salvar despesas:', error);
        if (error.response && error.response.data && error.response.data.error) {
          this.errorMessage = error.response.data.error;
        } else {
          this.errorMessage = 'Não foi possível conectar ao servidor. Tente novamente mais tarde.';
        }
      }
    },
  },
};
</script>

<style scoped>
.onboarding-expenses-container {
  display: flex;
  justify-content: center;
  align-items: flex-start; /* Alinha no topo para listas longas */
  min-height: 100vh;
  background-color: #f4f7f6;
  padding-top: 40px; /* Espaço no topo */
  padding-bottom: 40px; /* Espaço em baixo para rolagem */
}

.onboarding-card {
  background-color: white;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  text-align: center;
  width: 100%;
  max-width: 600px; /* Um pouco mais largo para acomodar os campos lado a lado */
}

h1 {
  color: #333;
  margin-bottom: 16px;
}

p {
  color: #666;
  margin-bottom: 24px;
  font-size: 0.95em;
}

.expense-item {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  gap: 10px; /* Espaço entre os elementos do item */
}

.expense-name-input {
  flex-grow: 1; /* Ocupa mais espaço */
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 1em;
}

.expense-name-input:focus {
  border-color: #007bff;
  box-shadow: 0 0 0 0.2rem rgba(0,123,255,.25);
  outline: none;
}

.value-input-group {
  display: flex;
  align-items: center;
  border: 1px solid #ccc;
  border-radius: 4px;
  padding: 0 10px;
  /* flex-basis: 150px; */ /* Define uma base de largura para o valor */
}

.value-input-group:focus-within {
  border-color: #007bff;
  box-shadow: 0 0 0 0.2rem rgba(0,123,255,.25);
}

.currency-symbol {
  font-size: 1em; /* Ajustado para ser similar ao input */
  color: #555;
  margin-right: 5px;
}

.expense-value-input {
  padding: 10px 0;
  border: none;
  font-size: 1em;
  outline: none;
  width: 80px; /* Largura fixa para o valor */
}

input[type="number"]::-webkit-outer-spin-button,
input[type="number"]::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
input[type=number] {
  -moz-appearance: textfield;
}

.remove-expense-button {
  background-color: #ff4d4f;
  color: white;
  border: none;
  border-radius: 50%;
  width: 30px;
  height: 30px;
  font-size: 1.2em;
  line-height: 30px; /* Centraliza o X */
  text-align: center;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.remove-expense-button:hover {
  background-color: #d9363e;
}

.add-expense-button {
  background-color: #28a745; /* Verde para adicionar */
  color: white;
  padding: 10px 15px;
  border: none;
  border-radius: 4px;
  font-size: 0.9em;
  cursor: pointer;
  transition: background-color 0.3s ease;
  margin-top: 10px; /* Espaço acima */
  margin-bottom: 20px; /* Espaço abaixo antes do botão de concluir */
  display: block; /* Faz ocupar a largura disponível se não houver outros botões na linha */
  margin-left: auto;
  margin-right: auto;
}

.add-expense-button:hover {
  background-color: #1e7e34;
}

.finish-button {
  background-color: #007bff;
  color: white;
  padding: 12px 24px;
  border: none;
  border-radius: 4px;
  font-size: 1em;
  cursor: pointer;
  transition: background-color 0.3s ease;
  width: 100%;
}

.finish-button:hover:not(:disabled) {
  background-color: #0056b3;
}

.finish-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.error-message {
  color: red;
  margin-top: 16px;
  font-size: 0.9em;
}
</style>
