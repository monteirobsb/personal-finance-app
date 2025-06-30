# Personal Finance App

Este é um aplicativo de finanças pessoais full-stack, projetado para ajudar os usuários a gerenciar suas receitas, despesas fixas e variáveis, e obter projeções sobre sua saúde financeira.

O backend é construído com Golang (Gin, GORM) e o frontend com Vue.js. Os dados são armazenados em um banco de dados PostgreSQL. A aplicação é totalmente containerizada usando Docker e Docker Compose.

## Pré-requisitos

Antes de começar, garanta que você tem os seguintes softwares instalados:

*   [Docker](https://www.docker.com/get-started)
*   [Docker Compose](https://docs.docker.com/compose/install/) (geralmente vem com a instalação do Docker Desktop)

## Configuração

1.  **Clone o repositório:**
    ```bash
    git clone <URL_DO_REPOSITORIO>
    cd personal-finance-app
    ```

2.  **Configure as variáveis de ambiente do Backend:**
    Navegue até o diretório `backend/` e copie o arquivo de exemplo `.env.example` para um novo arquivo chamado `.env`.
    ```bash
    cd backend
    cp .env.example .env
    ```
    Abra o arquivo `backend/.env` e revise as variáveis. Especialmente, defina um `JWT_SECRET` forte e único para produção. Os valores padrão são adequados para desenvolvimento local.

    *   `DB_HOST=db`
    *   `DB_PORT=5432`
    *   `DB_USER=user`
    *   `DB_PASSWORD=password`
    *   `DB_NAME=personalfinancedb`
    *   `JWT_SECRET="seu_segredo_jwt_super_secreto_aqui"`
    *   `PORT=8080`

    Volte para a raiz do projeto:
    ```bash
    cd ..
    ```

## Como Executar a Aplicação

Com o Docker e Docker Compose instalados e as variáveis de ambiente configuradas, você pode iniciar todos os serviços com um único comando a partir da raiz do projeto:

```bash
docker-compose up -d --build
```

*   `--build`: Constrói as imagens Docker para o backend e frontend antes de iniciar os containers. Use isso na primeira vez ou quando houver mudanças nos Dockerfiles ou no código que exijam um novo build.
*   `-d`: Executa os containers em modo detached (em segundo plano).

Para parar os serviços:

```bash
docker-compose down
```

Para ver os logs dos serviços (útil para debugging):

```bash
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f db
```

## Serviços e Portas

A aplicação é composta por três serviços principais gerenciados pelo Docker Compose:

*   **Backend (API Golang):**
    *   Disponível em: `http://localhost:8080`
    *   Responsável pela lógica de negócios, autenticação, e interação com o banco de dados.
*   **Frontend (Vue.js App):**
    *   Disponível em: `http://localhost:8081`
    *   Interface do usuário construída com Vue.js e servida pelo Nginx.
*   **Banco de Dados (PostgreSQL):**
    *   Porta interna: `5432` (acessível pelos outros containers na rede Docker)
    *   Porta externa (mapeada para o host): `5432`
    *   Os dados são persistidos em um volume Docker chamado `postgres_data`.

## Estrutura do Projeto

```
.
├── backend/            # Código fonte da API em Golang
│   ├── Dockerfile      # Dockerfile para construir a imagem do backend
│   ├── go.mod          # Módulo Go e dependências
│   ├── main.go         # Ponto de entrada da API
│   ├── handlers/       # Handlers HTTP (Gin)
│   ├── models/         # Modelos de dados (GORM structs)
│   ├── database/       # Lógica de conexão com o banco de dados
│   ├── middleware/     # Middlewares (ex: autenticação JWT)
│   ├── .env.example    # Exemplo de variáveis de ambiente para o backend
│   └── .env            # Arquivo de variáveis de ambiente (não versionado se contiver segredos)
├── frontend/           # Código fonte da aplicação Vue.js
│   ├── Dockerfile      # Dockerfile para construir e servir o frontend com Nginx
│   ├── nginx.conf      # Configuração do Nginx para servir a SPA Vue
│   ├── package.json    # Dependências e scripts do Node.js
│   ├── vue.config.js   # Configuração do Vue CLI (ex: proxy de desenvolvimento)
│   ├── public/         # Arquivos estáticos públicos (ex: index.html)
│   └── src/            # Código fonte Vue (componentes, views, router, etc.)
│       ├── main.js     # Ponto de entrada da aplicação Vue
│       ├── App.vue     # Componente raiz Vue
│       ├── router/     # Configuração do Vue Router
│       ├── views/      # Componentes de página (rotas)
│       └── components/ # Componentes reutilizáveis
├── docker-compose.yml  # Arquivo de orquestração do Docker Compose
└── README.md           # Este arquivo
```

## Próximos Passos (Exemplos)

*   Implementar testes unitários e de integração.
*   Adicionar mais funcionalidades (ex: edição de despesas, relatórios, categorias personalizadas).
*   Refinar a UI/UX.
*   Configurar HTTPS para produção.
```
