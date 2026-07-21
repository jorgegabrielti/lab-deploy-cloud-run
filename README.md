# ☀️ Desafio: Clima por CEP no Google Cloud Run

Sistema em Go desenvolvido para o desafio do curso **Go Expert (Full Cycle)**. O sistema recebe um CEP de 8 dígitos, identifica a cidade via **ViaCEP**, busca a temperatura atual via **WeatherAPI / Open-Meteo** e retorna as temperaturas em Celsius, Fahrenheit e Kelvin.

---

## 🌐 URL de Acesso no Google Cloud Run

> **URL Ativa no Cloud Run:**
> `https://desafio-clima-cep-897720935574.us-central1.run.app/`

### Exemplos de Requisição na Nuvem:

- **Sucesso (200 OK):**
  ```bash
  curl -s "https://desafio-clima-cep-897720935574.us-central1.run.app/?zipcode=01001000"
  ```
  **Resposta:**
  ```json
  {
    "temp_C": 28.5,
    "temp_F": 83.3,
    "temp_K": 301.65
  }
  ```

- **CEP com formato inválido (422 Unprocessable Entity):**
  ```bash
  curl -i "https://desafio-clima-cep-897720935574.us-central1.run.app/?zipcode=123"
  # Retorna HTTP 422 com o corpo: "invalid zipcode"
  ```

- **CEP não encontrado (404 Not Found):**
  ```bash
  curl -i "https://desafio-clima-cep-897720935574.us-central1.run.app/?zipcode=99999999"
  # Retorna HTTP 404 com o corpo: "can not find zipcode"
  ```

---

## 📋 Especificações da API (Contrato HTTP)

| Cenário | Condição | Status Code | Resposta |
|---|---|---|---|
| **Sucesso** | CEP válido (8 dígitos) e localizado | `200 OK` | `{"temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.65}` |
| **Formato inválido** | CEP com tamanho $\neq 8$ ou letras | `422 Unprocessable Entity` | `invalid zipcode` |
| **CEP não encontrado** | CEP de 8 dígitos inexistente no ViaCEP | `404 Not Found` | `can not find zipcode` |

---

## 📐 Fórmulas de Conversão

- **Fahrenheit:** $F = C \times 1.8 + 32$
- **Kelvin:** $K = C + 273.15$

---

## 🐳 Execução Local via Docker

### 1. Build da imagem Docker
```bash
docker build -t desafio-clima-cep .
```

### 2. Rodar o container
```bash
docker run --rm -p 8080:8080 desafio-clima-cep
```

### 3. Testar a aplicação localmente
```bash
# Teste de Sucesso
curl -s "http://localhost:8080/?zipcode=01001000"

# Teste Formato Inválido (422)
curl -i "http://localhost:8080/?zipcode=123"

# Teste CEP Inexistente (404)
curl -i "http://localhost:8080/?zipcode=99999999"
```

---

## 🧪 Como Rodar os Testes Automatizados

```bash
go test -v ./...
```

### Testes Implementados:
- **`pkg/weather`**: Testes unitários das conversões Celsius $\rightarrow$ Fahrenheit e Kelvin.
- **`pkg/viacep`**: Testes da validação de CEP (formato de 8 dígitos, sanitização) e integração com mock de servidor.
- **`pkg/handler`**: Testes de integração HTTP cobrindo os cenários `200 OK`, `422 invalid zipcode` e `404 can not find zipcode`.

---

## 🚀 Deploy no Google Cloud Run

```bash
# 1. Autenticação no GCP
gcloud auth login

# 2. Deploy direto do código-fonte para o Cloud Run
gcloud run deploy desafio-clima-cep \
  --source . \
  --region us-central1 \
  --allow-unauthenticated
```
