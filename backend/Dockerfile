# Estágio de build
FROM golang:1.21-alpine AS builder

# Define o diretório de trabalho dentro do container
WORKDIR /app

# Copia os arquivos go.mod e go.sum para baixar as dependências primeiro (cache layer)
# Se go.sum não existir ainda, ele será criado pelo go mod tidy
COPY go.mod ./
# COPY go.sum ./
# Devido aos problemas com a sandbox, go.sum não foi gerado.
# O go mod tidy abaixo tentará resolver isso.

# Baixa as dependências.
# Se go.sum for copiado e estiver correto, isso será mais rápido.
# O tidy garante que o go.mod e go.sum estejam consistentes.
RUN go mod tidy
RUN go mod download

# Copia o restante do código da aplicação
COPY . .

# Constrói a aplicação Go.
# CGO_ENABLED=0 para builds estáticos sem dependências C.
# -o /app/server para especificar o nome do output.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/server ./main.go

# Estágio de produção
FROM alpine:latest

# Define o diretório de trabalho
WORKDIR /root/

# Copia o executável construído do estágio builder
COPY --from=builder /app/server .

# Copia o .env.example para que possa ser usado como template se necessário
# Em um ambiente de produção real, as variáveis de ambiente seriam injetadas
# de forma segura, e não através de um arquivo .env no container.
COPY .env.example .

# Expõe a porta que a aplicação vai usar (deve corresponder à porta no main.go e docker-compose.yml)
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["./server"]
