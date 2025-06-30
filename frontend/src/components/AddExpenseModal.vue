<template>
  <div class="modal-overlay" @click.self="closeModal">
    <div class="modal-card">
      <button class="close-button" @click="closeModal">&times;</button>

      <div class="modal-header">
        <h2>Adicionar Despesa Rápida</h2>
      </div>

      <div class="modal-body">
        <div class="value-input-section">
          <span class="currency-symbol-modal">R$</span>
          <input
            type="number"
            v-model.number="expenseValue"
            placeholder="0,00"
            class="value-input-modal"
            ref="valueInput"
            inputmode="numeric"
            pattern="[0-9]*"
            step="0.01"
            min="0.01"
          />
        </div>

        <p class="category-label">Selecione uma categoria:</p>
        <div class="category-icons">
          <button
            v-for="category in categories"
            :key="category.name"
            @click="selectCategoryAndSubmit(category)"
            class="category-icon-button"
            :class="{ 'selected': selectedCategory && selectedCategory.name === category.name }"
            :disabled="!expenseValue || expenseValue <= 0"
            :title="category.name"
          >
            <span class="icon">{{ category.icon }}</span>
            <!-- <span class="category-name-icon">{{ category.name }}</span> -->
          </button>
        </div>

        <div class="optional-details-section">
          <button @click="toggleDetails" class="details-toggle-button">
            {{ showDetails ? 'Menos detalhes' : 'Adicionar detalhes' }}
            <span class="arrow">{{ showDetails ? '▲' : '▼' }}</span>
          </button>
          <transition name="slide-fade">
            <div v-if="showDetails" class="details-fields">
              <input
                type="text"
                v-model="expenseDescription"
                placeholder="Descrição (opcional)"
                class="description-input"
              />
              <input
                type="date"
                v-model="expenseDate"
                class="date-input"
              />
            </div>
          </transition>
        </div>
      </div>

      <!-- O botão de submissão principal é removido, pois o clique no ícone da categoria submete -->
      <!-- Se quiser um botão de submissão explícito, pode adicionar aqui -->
       <div class="modal-footer" v-if="showDetails">
         <button
            @click="submitWithDetails"
            class="submit-details-button"
            :disabled="!expenseValue || expenseValue <= 0 || !selectedCategory"
          >
            Salvar Despesa
          </button>
       </div>

    </div>
  </div>
</template>

<script>
export default {
  name: 'AddExpenseModal',
  props: {
    // Pode receber categorias como prop se forem dinâmicas no futuro
  },
  data() {
    return {
      expenseValue: null,
      selectedCategory: null,
      expenseDescription: '',
      expenseDate: new Date().toISOString().slice(0,10), // Default para hoje
      showDetails: false,
      categories: [ // Categorias sugeridas
        { name: 'Alimentação', icon: '🍔' }, // Ou 🍽️
        { name: 'Transporte', icon: '🚗' }, // Ou 🚌
        { name: 'Moradia', icon: '🏠' },
        { name: 'Lazer', icon: '🎬' },    // Ou 🎮
        { name: 'Saúde', icon: '⚕️' },   // Ou ❤️‍🩹
        // { name: 'Outros', icon: '🏷️'}
      ],
    };
  },
  mounted() {
    // Autofocus no campo de valor quando o modal é montado
    this.$nextTick(() => {
      if (this.$refs.valueInput) {
        this.$refs.valueInput.focus();
      }
    });
  },
  methods: {
    closeModal() {
      this.$emit('close');
    },
    selectCategoryAndSubmit(category) {
      if (!this.expenseValue || this.expenseValue <= 0) {
        // Adicionar feedback visual se o valor for inválido
        if(this.$refs.valueInput) this.$refs.valueInput.focus();
        return;
      }
      this.selectedCategory = category;
      // Se os detalhes não estiverem visíveis, submete imediatamente
      if (!this.showDetails) {
        this.submitExpense();
      }
      // Se os detalhes estiverem visíveis, o usuário usará o botão "Salvar Despesa"
    },
    toggleDetails() {
      this.showDetails = !this.showDetails;
    },
    submitWithDetails() {
        if (!this.expenseValue || this.expenseValue <= 0 || !this.selectedCategory) {
            // Idealmente, fornecer feedback mais específico
            alert("Valor e categoria são obrigatórios.");
            return;
        }
        this.submitExpense();
    },
    submitExpense() {
      const expenseData = {
        value: this.expenseValue,
        category: this.selectedCategory.name, // Envia o nome da categoria
        description: this.showDetails ? this.expenseDescription : '',
        date: this.showDetails ? this.expenseDate : new Date().toISOString().slice(0,10),
      };
      // Emitir evento para o Dashboard.vue com os dados da despesa
      this.$emit('save-expense', expenseData);
      this.closeModal(); // Fecha o modal após submeter
      // Resetar campos?
      // this.expenseValue = null;
      // this.selectedCategory = null;
      // this.expenseDescription = '';
      // this.expenseDate = new Date().toISOString().slice(0,10);
      // this.showDetails = false;
    },
  },
};
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.6); /* Fundo escurecido */
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000; /* Sobrepor outros conteúdos */
}

.modal-card {
  background-color: white;
  padding: 25px 30px;
  border-radius: 12px;
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.2);
  width: 90%;
  max-width: 400px;
  position: relative;
  text-align: center;
}

.close-button {
  position: absolute;
  top: 10px;
  right: 15px;
  background: none;
  border: none;
  font-size: 1.8em;
  color: #aaa;
  cursor: pointer;
}
.close-button:hover {
    color: #333;
}

.modal-header h2 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #333;
  font-size: 1.4em;
}

.modal-body {
  margin-bottom: 20px;
}

.value-input-section {
  display: flex;
  align-items: center;
  justify-content: center; /* Centraliza o input de valor */
  margin-bottom: 20px;
  border: 2px solid #007bff; /* Borda destacada */
  border-radius: 8px;
  padding: 5px 15px;
}

.currency-symbol-modal {
  font-size: 2.2em;
  font-weight: bold;
  color: #007bff;
  margin-right: 8px;
}

.value-input-modal {
  flex-grow: 1;
  border: none;
  font-size: 2.2em; /* Campo de valor grande */
  font-weight: bold;
  color: #333;
  outline: none;
  text-align: left;
  min-width: 100px; /* Para não ficar muito pequeno */
  background: transparent;
}
.value-input-modal::placeholder {
    color: #bdc9d8;
    font-weight: normal;
}

.category-label {
  font-size: 0.9em;
  color: #555;
  margin-bottom: 10px;
}

.category-icons {
  display: flex;
  justify-content: space-around; /* Ou space-between */
  align-items: center;
  margin-bottom: 25px;
  gap: 10px; /* Espaço entre ícones */
}

.category-icon-button {
  background-color: #f0f0f0;
  border: 2px solid transparent;
  border-radius: 50%; /* Ícones redondos */
  width: 55px;
  height: 55px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  transition: all 0.2s ease;
  padding: 0;
}
.category-icon-button:hover:not(:disabled) {
  background-color: #e0e0e0;
  border-color: #007bff;
}
.category-icon-button.selected {
  background-color: #007bff;
  border-color: #0056b3;
  color: white; /* Se o ícone for texto/emoji */
}
.category-icon-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}


.category-icon-button .icon {
  font-size: 1.8em; /* Tamanho do emoji/ícone */
}
/* Se quiser nome abaixo do ícone:
.category-icon-button .category-name-icon {
  font-size: 0.7em;
  margin-top: 4px;
  color: #555;
}
.category-icon-button.selected .category-name-icon {
    color: white;
}
*/


.optional-details-section {
  margin-top: 20px;
}

.details-toggle-button {
  background: none;
  border: none;
  color: #007bff;
  cursor: pointer;
  font-size: 0.9em;
  padding: 5px;
}
.details-toggle-button .arrow {
    display: inline-block;
    margin-left: 5px;
    transition: transform 0.2s ease;
}
.details-toggle-button .arrow { /* Quando showDetails é true, rotaciona a seta */
    transform: rotate(0deg);
}
/* Se quiser que a seta rotacione quando aberto:
.details-fields .arrow {
    transform: rotate(180deg);
}
*/


.details-fields {
  margin-top: 15px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.description-input,
.date-input {
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 0.95em;
}
.description-input:focus,
.date-input:focus {
    border-color: #007bff;
    box-shadow: 0 0 0 0.2rem rgba(0,123,255,.25);
    outline: none;
}

/* Animação para os detalhes */
.slide-fade-enter-active, .slide-fade-leave-active {
  transition: all 0.3s ease-out;
}
.slide-fade-enter-from, .slide-fade-leave-to {
  transform: translateY(-10px);
  opacity: 0;
}

.modal-footer {
    margin-top: 20px;
}
.submit-details-button {
    background-color: #28a745;
    color: white;
    padding: 12px 20px;
    border: none;
    border-radius: 6px;
    font-size: 1em;
    cursor: pointer;
    width: 100%;
    transition: background-color 0.2s;
}
.submit-details-button:hover:not(:disabled) {
    background-color: #1e7e34;
}
.submit-details-button:disabled {
    background-color: #ccc;
    cursor: not-allowed;
}

</style>
