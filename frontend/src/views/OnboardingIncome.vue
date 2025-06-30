<template>
  <div class="onboarding-income-container">
    <div class="onboarding-card">
      <h1>Qual é a sua renda mensal?</h1>
      <p>Saber sua renda nos ajuda a entender sua saúde financeira inicial.</p>
      <div class="input-group">
        <span class="currency-symbol">R$</span>
        <input
          type="number"
          v-model.number="monthlyIncome"
          placeholder="Ex: 3500.00"
          min="0"
          step="0.01"
          @input="validateInput"
        />
      </div>
      <button @click="nextStep" :disabled="isNextButtonDisabled" class="next-button">
        Avançar
      </button>
      <p v_if="errorMessage" class="error-message">{{ errorMessage }}</p>
    </div>
  </div>
</template>

<script>
import axios from 'axios'; // Será usado quando a chamada da API for implementada

export default {
  name: 'OnboardingIncome',
  data() {
    return {
      monthlyIncome: null,
      errorMessage: '',
    };
  },
  computed: {
    isNextButtonDisabled() {
      return !this.monthlyIncome || this.monthlyIncome <= 0;
    },
  },
  methods: {
    validateInput() {
      if (this.monthlyIncome < 0) {
        this.monthlyIncome = 0;
      }
      // Poderia adicionar mais validações aqui se necessário
    },
    async nextStep() {
      if (this.isNextButtonDisabled) {
        this.errorMessage = 'Por favor, insira um valor de renda válido.';
        return;
      }
      this.errorMessage = ''; // Limpa mensagens de erro anteriores

      try {
        // Recuperar o token JWT do localStorage (ou de onde estiver armazenado após o login)
        const token = localStorage.getItem('authToken');
        if (!token) {
          this.errorMessage = 'Erro de autenticação. Por favor, faça login novamente.';
          // Idealmente, redirecionar para a página de login
          this.$router.push('/'); // Ou para a tela de login/auth
          return;
        }

        const response = await axios.post('/api/onboarding/income', // O proxy cuidará do /api
          { rendaMensal: this.monthlyIncome },
          { headers: { Authorization: `Bearer ${token}` } }
        );

        if (response.status === 200 || response.status === 201) {
          console.log('Renda salva com sucesso:', response.data);
          // Navegar para a próxima etapa do onboarding (despesas)
          this.$router.push('/onboarding/expenses');
        } else {
          // Tratar outros status de sucesso que não sejam 200/201, se houver
          this.errorMessage = response.data.message || 'Ocorreu um erro ao salvar a renda.';
        }
      } catch (error) {
        console.error('Erro ao salvar renda:', error);
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
.onboarding-income-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f4f7f6; /* Um cinza claro para o fundo */
}

.onboarding-card {
  background-color: white;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  text-align: center;
  width: 100%;
  max-width: 450px;
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

.input-group {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
  border: 1px solid #ccc;
  border-radius: 4px;
  padding: 0 10px;
}

.input-group:focus-within {
  border-color: #007bff; /* Destaque ao focar */
  box-shadow: 0 0 0 0.2rem rgba(0,123,255,.25);
}

.currency-symbol {
  font-size: 1.2em;
  color: #555;
  margin-right: 8px;
}

input[type="number"] {
  flex-grow: 1;
  padding: 12px 0; /* Ajuste o padding para alinhar com o símbolo */
  border: none;
  font-size: 1.2em;
  outline: none; /* Remove o outline padrão, já que o grupo tem focus */
  width: calc(100% - 30px); /* Ajuste para não estourar o container */
}

/* Remove setas do input number no Chrome, Safari, Edge, Opera */
input::-webkit-outer-spin-button,
input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
/* Remove setas do input number no Firefox */
input[type=number] {
  -moz-appearance: textfield;
}

.next-button {
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

.next-button:hover:not(:disabled) {
  background-color: #0056b3;
}

.next-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.error-message {
  color: red;
  margin-top: 16px;
  font-size: 0.9em;
}
</style>
