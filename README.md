[![Build and Test](https://github.com/nettojulio/ufape-crawler-golang/actions/workflows/release.yml/badge.svg)](https://github.com/nettojulio/ufape-crawler-golang/actions/workflows/release.yml)
[![codecov](https://codecov.io/gh/nettojulio/ufape-crawler-golang/graph/badge.svg)](https://codecov.io/gh/nettojulio/ufape-crawler-golang)
[![Go Report Card](https://goreportcard.com/badge/github.com/nettojulio/ufape-crawler-golang)](https://goreportcard.com/report/github.com/nettojulio/ufape-crawler-golang)
![GitHub release(including pre-releases)](https://img.shields.io/github/v/release/nettojulio/ufape-crawler-golang?include_prereleases&cache_bust=1)
[![GitHub license](https://img.shields.io/github/license/nettojulio/ufape-crawler-golang)](https://github.com/nettojulio/ufape-crawler-golang/blob/main/LICENSE.md)
[![Go Version](https://img.shields.io/github/go-mod/go-version/nettojulio/ufape-crawler-golang)](https://go.dev/)

# 🚀 UFAPE Crawler Golang

![Banner do Projeto](https://cdn.dribbble.com/userupload/42462891/file/original-2f612076f7073b798d9b17f647e8d0f2.gif)

Um projeto Go para coletar dados de uma URL específica! Este README fornecerá uma visão geral,
instruções de instalação e uso, e muito mais.

---

## 💡 Sobre o Projeto

Este projeto foi desenvolvido visando analisar o site institucional
da [Universidade Federal do Agreste de Pernambuco](https://ufape.edu.br/), aplicando conhecimento adquirido sobre Grafos
na Disciplina Algoritmos e Estrutura de Dados II. É a API para coleta de dados oferecendo alta performance, facilitar a extração, definido tarefas e escopos.

### ✨ Principais Features

* **Processamento de Dados Rápido**: Otimizado para alta performance.
* **API RESTful**: Interface simples para interação.
* **Configuração Flexível**: Permite fácil adaptação a diferentes ambientes.
* **Testes Abrangentes**: Garantindo a estabilidade e confiabilidade.

---

## 🛠️ Tecnologias Utilizadas

* [Go Lang](https://golang.org/) - Linguagem de programação
* [Echo](https://echo.labstack.com/) - Framework para requisições HTTP
* [Swagger](https://swagger.io/) - Gerador de documentação
* [GoDotEnv](https://github.com/joho/godotenv) - Carregamento de variáveis de ambiente

---

## 🛫 Como Começar

Siga estas instruções para colocar o projeto em funcionamento em sua máquina local para fins de desenvolvimento e teste.

### Pré-requisitos

Certifique-se de ter o seguinte instalado:

* Go (versão 1.24 ou superior)

### Instalação

1. **Clone o repositório:**

   HTTPS

   ```bash
   git clone https://github.com/nettojulio/ufape-crawler-golang.git
   ```

   SSH

   ```bash
   git clone git@github.com:nettojulio/ufape-crawler-golang.git
   ```

2. **Instale as dependências:**
   ```bash
   go mod tidy
   ```

3. **Configure as variáveis de ambiente:**
   Crie um arquivo `.env` na raiz do projeto baseado no `.env.example`:
   ```
   # Exemplo de .env
   APP_PORT=8080
   APP_HOST=localhost:8080
   ```

4. **Execute o Projeto:**
   ```bash
   go run cmd/main.go
   ```
   O projeto estará rodando em `http://localhost:8080` (ou na porta configurada) e espera por requisições HTTP.
   A documentação do Swagger com rotas e detalhes estará disponível em `http://localhost:8080/swagger/index.html`.

---

## 🧪 Rodando os Testes

Para rodar os testes unitários e de integração, execute o seguinte comando:

```bash
go test ./...
```

---

## ⚙️ Configuracoes iniciais para commits padronizados

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
