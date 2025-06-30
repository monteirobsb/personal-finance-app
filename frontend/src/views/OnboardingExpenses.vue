<template>
  <div class="onboarding-expenses-container">
    <div class="onboarding-card">
      <header class="onboarding-header">
        <h1>Quais são suas despesas fixas?</h1>
      </header>
      <div class="onboarding-body">
        <p class="support-text">Inclua contas que não mudam de valor, como aluguel, internet e assinaturas.</p>

        <div v-for="(expense, index) in expenses" :key="index" class="expense-item">
          <input
            type="text"
            v-model="expense.name"
            placeholder="Nome da despesa (ex: Aluguel)"
            class="expense-name-input"
          />
          <div class="input-group value-input-group">
            <span class="currency-symbol">R$</span>
            <input
              type="number"
              v-model.number="expense.value"
              placeholder="0,00"
              min="0.01"
              step="0.01"
              class="expense-value-input"
            />
          </div>
          <button @click="removeExpense(index)" class="remove-expense-button" v-if="expenses.length > 0 && index !== 0">
            &times; <!-- Só mostra o botão de remover se não for o primeiro item ou se houver mais de um -->
          </button>
           <div v-else class="remove-placeholder"></div> <!-- Placeholder para alinhar -->
        </div>

        <button @click="addExpense" class="add-expense-button">
          + Adicionar outra despesa
        </button>
      </div>
      <footer class="onboarding-footer">
        <button @click="finishOnboarding" :disabled="isFinishButtonDisabled" class="finish-button">
          Concluir
        </button>
      </footer>
      <p v-if="errorMessage" class="error-message-footer">{{ errorMessage }}</p>
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
      // Botão Concluir habilitado após a inserção de pelo menos uma despesa válida.
      // Uma despesa é válida se tiver nome e valor > 0.
      // Se a lista estiver vazia OU se todos os itens forem inválidos, o botão fica desabilitado.
      if (this.expenses.length === 0) return true; // Se não há despesas, desabilita.

      // Verifica se PELO MENOS UMA despesa é válida. Se sim, o botão habilita.
      // O wireframe diz "habilitado após a inserção de PELO MENOS UMA despesa"
      // Isso implica que se houver uma válida, mesmo que outras estejam em branco, ele habilita.
      // No entanto, para submissão, provavelmente queremos que TODAS as despesas preenchidas sejam válidas.
      // A lógica atual no `finishOnboarding` já verifica isso antes de submeter.
      // Para o estado do botão, vamos seguir o "pelo menos uma válida".
      const atLeastOneValid = this.expenses.some(
        (exp) => exp.name && exp.value && exp.value > 0
      );
      return !atLeastOneValid; // Habilita se atLeastOneValid for true.
    },
  },
  methods: {
    addExpense() {
      this.expenses.push({ name: '', value: null });
    },
    removeExpense(index) {
      // Permite remover qualquer item, mas o `isFinishButtonDisabled` vai reavaliar.
      // Se o último item válido for removido, o botão Concluir será desabilitado.
      if (this.expenses.length > 0) { // Só remove se houver itens
        this.expenses.splice(index, 1);
        if (this.expenses.length === 0) { // Se todas foram removidas, adiciona uma em branco
            this.addExpense();
        }
      }
    },
    async finishOnboarding() {
      // Antes de submeter, garante que TODAS as despesas que têm algum campo preenchido
      // sejam completamente válidas (nome e valor positivo).
      // Despesas totalmente em branco podem ser ignoradas ou filtradas.
      const validExpensesToSubmit = this.expenses.filter(
        exp => exp.name && exp.value && exp.value > 0
      );

      if (validExpensesToSubmit.length === 0 && this.expenses.some(exp => exp.name || exp.value)) {
        this.errorMessage = 'Por favor, preencha pelo menos uma despesa completamente (nome e valor positivo) ou remova as despesas incompletas.';
        return;
      }
      // Se todas as despesas estiverem em branco (ex: usuário apagou tudo e deixou um item em branco)
      // e o botão estiver habilitado por alguma lógica anterior, mas não há nada para enviar.
      // Embora a lógica do `isFinishButtonDisabled` deva cobrir isso.
      if (validExpensesToSubmit.length === 0 && this.expenses.length > 0 && !this.expenses.some(exp => exp.name || exp.value)) {
         // Nenhuma despesa válida e nenhuma parcialmente preenchida, pode ser um array de {name:'', value:null}
         // O usuário pode querer avançar sem despesas fixas. O backend deve aceitar lista vazia.
         // Se for obrigatório ter despesas, a lógica do botão e aqui precisa mudar.
         // Assumindo que é opcional enviar despesas:
         console.log("Nenhuma despesa fixa para enviar.");
      }


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

      // Usa as despesas filtradas para o payload
      const payload = {
        despesasFixas: validExpensesToSubmit.map(exp => ({ nome: exp.name, valor: exp.value })),
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
  padding: 30px 40px; /* Ajuste no padding */
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  text-align: center;
  width: 100%;
  max-width: 600px;
  display: flex;
  flex-direction: column;
}

.onboarding-header h1 {
  color: #333;
  margin-bottom: 16px;
  font-size: 1.8em;
}

.onboarding-body {
  flex-grow: 1;
  width: 100%;
}

.support-text {
  color: #666;
  margin-bottom: 24px;
  font-size: 0.95em;
}

.expense-item {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  gap: 10px;
}

.expense-name-input {
  flex-grow: 1;
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
}

.value-input-group:focus-within {
  border-color: #007bff;
  box-shadow: 0 0 0 0.2rem rgba(0,123,255,.25);
}

.currency-symbol {
  font-size: 1em;
  color: #555;
  margin-right: 5px;
}

.expense-value-input {
  padding: 10px 0;
  border: none;
  font-size: 1em;
  outline: none;
  width: 90px; /* Ajuste de largura para o valor */
  text-align: right;
}
.expense-value-input::placeholder {
  text-align: left; /* Placeholder à esquerda */
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
  width: 28px; /* Levemente menor */
  height: 28px;
  font-size: 1.1em;
  line-height: 28px;
  text-align: center;
  cursor: pointer;
  transition: background-color 0.3s ease;
}
.remove-placeholder {
  width: 28px; /* Mesmo tamanho do botão para manter alinhamento */
  height: 28px;
}


.remove-expense-button:hover {
  background-color: #d9363e;
}

.add-expense-button {
  background-color: transparent;
  color: #007bff;
  padding: 10px 15px;
  border: 1px dashed #007bff;
  border-radius: 4px;
  font-size: 0.9em;
  cursor: pointer;
  transition: background-color 0.3s ease, color 0.3s ease;
  margin-top: 10px;
  margin-bottom: 20px;
  display: inline-block; /* Para não ocupar largura total */
}

.add-expense-button:hover {
  background-color: #e7f3ff; /* Um azul bem claro */
}

.onboarding-footer {
  padding-top: 20px;
  width: 100%;
}

.finish-button {
  background-color: #007bff;
  color: white;
  padding: 14px 24px; /* Botão um pouco maior */
  border: none;
  border-radius: 4px;
  font-size: 1.1em; /* Texto do botão um pouco maior */
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

.error-message-footer { /* Renomeado para evitar conflito se houver .error-message no body */
  color: red;
  margin-top: 10px; /* Menor margem, pois está abaixo do footer */
  font-size: 0.9em;
  min-height: 1.2em; /* Para evitar pulo de layout */
}
</style>
