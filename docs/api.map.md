# API Map - Book Me Server

Este documento mapeia as rotas, payloads e respostas esperadas da API do `book-me-server`.

## Base URL
`/` (ou configurado no servidor)

---

## Rotas Públicas

## Tested -> [v]
### 1. Status do Servidor 
- **Método:** `GET`
- **URL:** `/public/status`
- **Descrição:** Verifica a saúde do servidor.
- **Payload:** Nenhum
- **Resposta de Sucesso (200):**
  ```json
  {
    "Code": 200,
    "Message": "Server health ok!"
  }
  ```

## Tested -> [v]
### 2. Registro de Usuário
- **Método:** `POST`
- **URL:** `/public/register`
- **Descrição:** Registra um novo usuário (cliente ou prestador) e seu endereço inicial.
- **Payload:**
  ```json
  {
    "user": {
      "nome": "string",
      "email": "string",
      "senha": "string",
      "confirmaSenha": "string",
      "cpf": "string",
      "cnpj": "string",
      "telefone": "string",
      "userType": "customer | provider",
      "cep": "string",
      "estado": "string",
      "cidade": "string",
      "rua": "string",
      "logradouro": "string"
    },
    "address": {
      "street": "string",
      "city": "string",
      "state": "string",
      "postal_code": "string",
      "country": "string"
    }
  }
  ```
- **Resposta de Sucesso (201):**
  ```json
  {
    "user": {
      "id": "uuid",
      "nome": "string",
      "email": "string"
    },
    "access_token": "string",
    "refresh_token": "string",
    "code": 201,
    "message": "User registered successfully"
  }
  ```
- **Possíveis Erros:**
  - `400 Bad Request`: Payload inválido ou campos obrigatórios ausentes.
  - `500 Internal Server Error`: Falha ao salvar no banco.

## Tested -> [v]
### 3. Login
- **Método:** `POST`
- **URL:** `/public/login`
- **Descrição:** Autentica um usuário e retorna tokens JWT.
- **Payload:**
  ```json
  {
    "email": "string",
    "senha": "string"
  }
  ```
- **Resposta de Sucesso (200):**
  ```json
  {
    "status": "success",
    "code": 200,
    "message": "Login realizado com sucesso",
    "data": {
      "access_token": "string",
      "refresh_token": "string",
      "user": {
        "id": "uuid",
        "nome": "string",
        "email": "string",
        "role": "customer | provider"
      }
    }
  }
  ```
- **Possíveis Erros:**
  - `401 Unauthorized`: E-mail ou senha incorretos.
  - `400 Bad Request`: Dados inválidos.

## Tested -> [v]
### 4. Listar Categorias
- **Método:** `GET`
- **URL:** `/public/categories`
- **Descrição:** Lista todas as categorias de serviços disponíveis.
- **Payload:** Nenhum
- **Resposta de Sucesso (200):**
  ```json
  [
    {
      "id": "uuid",
      "name": "string",
      "description": "string"
    }
  ]
  ```

## Tested -> []
### 5. Atualizar Token (Refresh)
- **Método:** `POST`
- **URL:** `/refresh`
- **Descrição:** Gera um novo Access Token a partir de um Refresh Token válido.
- **Payload:**
  ```json
  {
    "refresh_token": "string"
  }
  ```
- **Resposta de Sucesso (200):**
  ```json
  {
    "access_token": "string"
  }
  ```
- **Possíveis Erros:**
  - `401 Unauthorized`: Refresh token inválido, expirado ou sessão não encontrada.

---

## Rotas Protegidas (Requer Header `Authorization: Bearer <token>`)

## Tested -> [v]
### 1. Obter Meus Dados
- **Método:** `GET`
- **URL:** `/me`
- **Descrição:** Retorna os dados do usuário autenticado.
- **Payload:** Nenhum
- **Resposta de Sucesso (200):**
  ```json
  {
    "id": "uuid",
    "nome": "string",
    "email": "string",
    "type": "customer | provider"
  }
  ```
## Tested -> [v]
### 2. Listar Meus Endereços
- **Método:** `GET`
- **URL:** `/addresses/me`
- **Descrição:** Retorna a lista de endereços vinculados ao usuário autenticado.
- **Payload:** Nenhum
- **Resposta de Sucesso (200):**
  ```json
  [
    {
      "id": "uuid",
      "street": "string",
      "city": "string",
      "state": "string",
      "postal_code": "string",
      "country": "string"
    }
  ]
  ```

## Tested -> [v]
### 3. Criar Serviço (Apenas Prestadores)
- **Método:** `POST`
- **URL:** `/provider/services`
- **Descrição:** Cria um novo serviço oferecido pelo prestador.
- **Payload:**
  ```json
  {
    "title": "string",
    "description": "string",
    "price_base": 0.0,
    "price_type": "string",
    "duration_minutes": 0,
    "category_id": "uuid",
    "is_active": true,
    "is_in_place": true,
    "address_id": "uuid | null"
  }
  ```
- **Resposta de Sucesso (201):**
  ```json
  {
    "id": "uuid",
    "title": "string",
    "description": "string",
    "price_base": 0.0,
    "price_type": "string",
    "duration_minutes": 0,
    "is_active": true,
    "is_in_place": true,
    "provider": { "id": "uuid", "email": "string" },
    "category": { "id": "uuid", "name": "string" }
  }
  ```

## Tested -> [v]
### 4. Listar Serviços do Prestador
- **Método:** `GET`
- **URL:** `/provider/services`
- **Descrição:** Lista os serviços gerenciados pelo prestador autenticado.
- **Query Params:** `limit`, `page`, `category_id`
- **Resposta de Sucesso (200):** Lista de objetos de serviço (mesmo formato da criação).

## Tested -> [v]
### 5. Atualizar Serviço
- **Método:** `PUT`
- **URL:** `/provider/services/{id}`
- **Descrição:** Atualiza os dados de um serviço existente.
- **Payload:** Mesmo formato do POST (Criação).
- **Resposta de Sucesso (200):** Objeto de serviço atualizado.

## Tested -> [v]
### 6. Deletar Serviço
- **Método:** `DELETE`
- **URL:** `/provider/services/{id}`
- **Descrição:** Remove um serviço do prestador.
- **Payload:** Nenhum
- **Resposta de Sucesso (204):** No Content.

## Tested -> [v]
### 7. Listar Serviços (Visão Cliente)
- **Método:** `GET`
- **URL:** `/customer/services`
- **Descrição:** Lista serviços disponíveis para contratação.
- **Query Params:** `limit`, `page`, `category_id`
- **Resposta de Sucesso (200):** Lista de objetos de serviço.
