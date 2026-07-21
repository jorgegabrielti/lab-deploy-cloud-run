# Clima por CEP (Google Cloud Run)

Sistema em Go que recebe um CEP de 8 dígitos, busca a localização correspondente via ViaCEP e retorna a temperatura atual em Celsius, Fahrenheit e Kelvin.

## Cloud Run URL

Serviço implantado no Google Cloud Run:
`https://desafio-clima-cep-897720935574.us-central1.run.app/`

### Exemplos de Uso

**Sucesso (200 OK):**
```bash
curl -s "https://desafio-clima-cep-897720935574.us-central1.run.app/?zipcode=01001000"
```

Resposta:
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

**CEP com formato inválido (422 Unprocessable Entity):**
```bash
curl -i "https://desafio-clima-cep-897720935574.us-central1.run.app/?zipcode=123"
```
Retorno: `HTTP 422` com o corpo `invalid zipcode`.

**CEP não encontrado (404 Not Found):**
```bash
curl -i "https://desafio-clima-cep-897720935574.us-central1.run.app/?zipcode=99999999"
```
Retorno: `HTTP 404` com o corpo `can not find zipcode`.

## Execução Local com Docker

Build da imagem:
```bash
docker build -t desafio-clima-cep .
```

Execução do container:
```bash
docker run --rm -p 8080:8080 desafio-clima-cep
```

Testar localmente:
```bash
curl -s "http://localhost:8080/?zipcode=01001000"
```

## Testes Automatizados

Para rodar a suíte de testes:

```bash
go test -v ./...
```

Os testes cobrem:
- Validação e sanitização do formato de CEP (`pkg/viacep`)
- Conversões de temperatura Celsius para Fahrenheit e Kelvin (`pkg/weather`)
- Handlers HTTP cobrindo os status 200, 404 e 422 (`pkg/handler`)

## Deploy no Cloud Run

```bash
gcloud run deploy desafio-clima-cep --source . --region us-central1 --allow-unauthenticated
```
