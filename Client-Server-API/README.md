# Desafio Go - Client & Server com Context, SQLite e HTTP

Este projeto implementa dois programas em Go (`server.go` e `client.go`) para atender ao desafio proposto.  
Ele utiliza **Context** para controle de timeout, **SQLite** para persistência, e **HTTP** para comunicação entre cliente e servidor.

---

## 📌 Descrição do Desafio

- **client.go**  
  - Realiza uma requisição HTTP para `server.go` solicitando a cotação do dólar.  
  - Timeout de **300ms** para receber resposta.  
  - Salva o valor atual do câmbio (campo `bid`) no arquivo `cotacao.txt` no formato:
    ```
    Dólar: {valor}
    ```
  - Loga erro caso o tempo de execução seja insuficiente.

- **server.go**  
  - Consome a API pública:  
    ```
    https://economia.awesomeapi.com.br/json/last/USD-BRL
    ```
  - Timeout de **200ms** para buscar cotação da API.  
  - Persiste no SQLite o valor recebido (campo `bid`) com timeout de **10ms**.  
  - Endpoint `/cotacao` na porta **8080**.  
  - Retorna JSON com a cotação para o cliente.

---

## 📂 Estrutura do Projeto

```
.
├── client.go
├── server.go
├── database.go
├── go.mod
├── go.sum
└── README.md
```

---

## 🛠 Pré-requisitos

- **Go 1.18+** instalado
- **SQLite3** instalado
- Acesso à internet para consumir a API de câmbio

---

## 🚀 Instalação

```bash
# Clonar o repositório
git clone https://github.com/jpfigueredo/full-cycle-challenges
cd Client-Server-API

# Inicializar módulo Go
go mod init <projeto>
go mod tidy
```

---

## 📦 Dependências

```bash
go get github.com/gin-gonic/gin
go get gorm.io/driver/sqlite
go get gorm.io/gorm
```

---

## ▶️ Executando o Servidor

```bash
go run server.go
```

O servidor estará disponível em:
```
http://localhost:8080/cotacao
```

---

## ▶️ Executando o Cliente

Em outro terminal:

```bash
go run client.go
```

O cliente irá:
1. Solicitar a cotação ao servidor.
2. Salvar no arquivo `cotacao.txt` o valor do dólar no formato:
   ```
   Dólar: 5.1234
   ```

---

## ⚠️ Tratamento de Timeouts

- **Client → Server**: 300ms
- **Server → API**: 200ms
- **Server → DB**: 10ms

Se algum tempo for excedido:
- A função logará o erro.
- Nenhum dado será salvo ou retornado (dependendo do ponto de falha).

---

## 📄 Exemplo de Resposta do Servidor

```json
{
  "bid": "5.1234"
}
```

---

## 📜 Licença
Este projeto é de uso livre para estudo e prática.
