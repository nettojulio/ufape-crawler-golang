# 🚀 Golang Template

![Banner do Projeto](https://cdn.dribbble.com/userupload/42462891/file/original-2f612076f7073b798d9b17f647e8d0f2.gif)

Um projeto Go incrível para resolver problemas XYZ e tornar sua vida mais fácil! Este README fornecerá uma visão geral,
instruções de instalação e uso, e muito mais.

---

## 💡 Sobre o Projeto

Este projeto foi desenvolvido com o objetivo
de [descrever brevemente o propósito principal do projeto, o que ele faz e quais problemas ele resolve]. Ele utiliza as
melhores práticas de Go e foi projetado para
ser [mencione qualidades como: escalável, performático, fácil de usar, etc.].

### ✨ Principais Features

* **Processamento de Dados Rápido**: Otimizado para alta performance.
* **API RESTful**: Interface simples para interação.
* **Configuração Flexível**: Permite fácil adaptação a diferentes ambientes.
* **Testes Abrangentes**: Garantindo a estabilidade e confiabilidade.

---

## 🛠️ Tecnologias Utilizadas

* [Go Lang](https://golang.org/) - Linguagem de programação
* [Gorilla Mux](https://github.com/gorilla/mux) - Roteador de requisições HTTP
* [GORM](https://gorm.io/) - ORM para Go
* [PostgreSQL](https://www.postgresql.org/) - Banco de Dados
* [Docker](https://www.docker.com/) - Containerização

---

## 🚀 Como Começar

Siga estas instruções para colocar o projeto em funcionamento em sua máquina local para fins de desenvolvimento e teste.

### Pré-requisitos

Certifique-se de ter o seguinte instalado:

* Go (versão 1.20 ou superior)
* Docker (opcional, para rodar o banco de dados)

### Instalação

1. **Clone o repositório:**
   ```bash
   git clone [https://github.com/seu-usuario/meu-super-projeto-go.git](https://github.com/seu-usuario/meu-super-projeto-go.git)
   cd meu-super-projeto-go
   ```

2. **Instale as dependências:**
   ```bash
   go mod tidy
   ```

3. **Configuração do Banco de Dados (com Docker):**
   ```bash
   docker-compose up -d postgres
   ```
   Ou configure seu banco de dados PostgreSQL manualmente.

4. **Configure as variáveis de ambiente:**
   Crie um arquivo `.env` na raiz do projeto baseado no `config.example.env`:
   ```
   # Exemplo de .env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=seu_usuario
   DB_PASSWORD=sua_senha
   DB_NAME=seu_banco_de_dados
   API_PORT=8080
   ```

5. **Rode as Migrações do Banco de Dados:**
   ```bash
   go run main.go migrate
   ```

6. **Execute o Projeto:**
   ```bash
   go run main.go
   ```
   O projeto estará rodando em `http://localhost:8080` (ou na porta configurada).

---

## 🧪 Rodando os Testes

Para rodar os testes unitários e de integração, execute o seguinte comando:

```bash
go test ./...
```

---

## ⚙️ Configuracoes iniciais

1. **Baixe e instale o git-chglog** [Link GitHub](https://github.com/git-chglog/git-chglog)
2. **Execute o comando**

   ```bash
    git-chglog --init
   ```

3. **Instale as dependencias npm**

   ```bash
    npm i
   ```

4. **Faca commits utilizando o script npm**
