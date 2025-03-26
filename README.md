
🔍 # Monitoramento de APIs

Aplicação que monitora APIs registradas no banco de dados, salvando logs de disponibilidade. Exibe apenas os logs do dia para acompanhamento rápido.


## Tecnologias utilizadas

* Go
* PostgreSQL
* pgx (biblioteca para conexão com PostgreSQL)
* emoji (para exibição de ícones no terminal)


## Funcionalidades

- Monitora APIs cadastradas no banco de dados
- Salva logs das respostas da requisição
- Exibe logs do dia


## Pré-requisitos

Antes de rodar o projeto, certifique-se de ter instalado:

* [Go](https://go.dev/dl/) - Linguagem utilizada no projeto
* [PostgreSQL](https://www.postgresql.org/download/) - Banco de dados 

## Como rodar o projeto

1. Clone o repositório

```bash
  git clone https://github.com/paulacolaca/goAPI
```
2. Instale o Go

3. Configure o banco de dados
Crie um banco no PostgreSQL e adicione as tabelas necessárias:

```sql
-- Tabela com as APIs que serão monitoradas
CREATE TABLE listaapis
(
    id SERIAL PRIMARY KEY,
    ds_api text NOT NULL,
    lk_url text NOT NULL   
)

INSERT INTO listaapis (ds_api, lk_url)
VALUES
('GitHub API','https://api.github.com'),
('JSONPlaceholder','https://jsonplaceholder.typicode.com')
('IPify', 'https://api.ipify.org')
('Postman Echo','https://postman-echo.com/get')
('HTTPBin', 'https://httpbin.org/get')
('Dog API', 'https://dog.ceo/api/breeds/image/random')
('The Cat API', 'https://api.thecatapi.com/v1/images/search')
('MyAnimeList', 'https://myanimelist.net/')
('MMO Api', 'https://www.mmobomb.com/api')
('PokéApi', 'https://pokeapi.co/')
('ReqRes400', 'https://reqres.in/api/login')
('PublicAPIs', 'https://api.publicapis.org')

-- Tabela que irá armazenar os logs
CREATE TABLE logapi
(
    id SERIAL PRIMARY KEY,
    lk_url text NOT NULL,
    vl_status integer NOT NULL,
    vl_temporesposta double precision,
    dt_requisicao timestamp,    
)
```
4. Configure o arquivo .env
O arquivo .env (não incluído no repositório) deve conter suas credenciais do banco de dados. Crie o arquivo e adicione:

```env
DB_Token="postgres://USER:PASSWORD@localhost:5432/DBNAME"
```
No lugar de "USER" adicione o usuário do BD, se não tiver criado um use o padrão do Postgre, "postgres", insira então a sua senha em "PASSWORD" e, por fim, em "DBNAME" o nome do banco de dados.

5. Execute o projeto
```bash
go run main.go
```

## Licença

Este projeto está sob a licença [MIT](https://choosealicense.com/licenses/mit/)


