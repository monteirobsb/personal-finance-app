# Estágio 1: Build da aplicação Vue.js
FROM node:18-alpine AS builder

# Define o diretório de trabalho no container
WORKDIR /app

# Copia o package.json e package-lock.json (se existir)
# Isso aproveita o cache do Docker se esses arquivos não mudarem
COPY package.json ./
# COPY package-lock.json ./
# (package-lock.json não foi gerado na sandbox devido à falha do npm install)

# Instala as dependências
# O npm ci seria preferível se package-lock.json estivesse presente e atualizado
RUN npm install

# Copia todos os outros arquivos do projeto para o diretório de trabalho
COPY . .

# Constrói a aplicação para produção
# Isso criará uma pasta 'dist' com os arquivos estáticos
RUN npm run build

# Estágio 2: Servir a aplicação com Nginx
FROM nginx:1.25-alpine

# Remove a configuração padrão do Nginx
RUN rm /etc/nginx/conf.d/default.conf

# Copia a configuração customizada do Nginx (será criada a seguir)
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copia os arquivos construídos do estágio 'builder' para o diretório www do Nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# Expõe a porta 80 (porta padrão do Nginx)
EXPOSE 80

# Comando para iniciar o Nginx
CMD ["nginx", "-g", "daemon off;"]
