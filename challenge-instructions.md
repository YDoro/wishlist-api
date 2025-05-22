📅 Contexto
O aiqfome está expandindo seus canais de integração e precisa de uma API robusta para gerenciar os "produtos favoritos" de usuários na plataforma.


Essa funcionalidade será usada por apps e interfaces web para armazenar e consultar produtos marcados como favoritos pelos clientes. 


A API terá alto volume de uso e integrará com outros sistemas internos e externos.

🔧 Sua missão:
Desenvolva uma API RESTful que permita:

Clientes
- Criar, visualizar, editar e remover clientes.
- Dados obrigatórios: nome e e-mail.
- Um mesmo e-mail não pode se repetir no cadastro.

Favoritos
- Um cliente deve ter uma lista de produtos favoritos.
- Os produtos devem ser validados via API externa (link fornecido abaixo).
- Um produto não pode ser duplicado na lista de um cliente.
- Produtos favoritos devem exibir: ID, título, imagem, preço e review (se houver).

Requisitos de Integração
 - Sugerimos o uso de uma API genérica para buscar produtos. Porém, para facilitar a execução e deixar tudo mais direto ao ponto, recomendamos o uso da seguinte API pública:

🔗 https://fakestoreapi.com/docs

Você pode utilizar especificamente estes dois endpoints:

Listar todos os produtos:
GET https://fakestoreapi.com/products

Buscar produto por ID:
GET https://fakestoreapi.com/products/{id}

Confira algumas dicas aqui
⚖️ Regras Gerais
- A API deve ser pública, mas conter autenticação e autorização.
- Evite duplicidade de dados.
- Estruture bem o código, seguindo boas práticas REST.
- Pense em performance e escalabilidade.
- Documente sua API (OpenAPI/Swagger é bem-vindo, mas opcional).
- Não use IA ou cópias. Será passível de eliminação.

💡 Requisitos Técnicos
Você pode escolher uma das seguintes linguagens:
- Go
- PHP
- Python
- Node.js (Javascript ou Typescript)
- Ou outra linguagem (sem apego!)

Banco de dados sugerido:
- PostgreSQL (preferencial)
- MySQL
- MongoDB

📊 O que esperamos:
Critério	Peso
Correção e funcionamento da API	🔥🔥🔥🔥
Modelagem de dados (clientes/produtos)	🔥🔥🔥
Validação e controle de dados	🔥🔥🔥
Documentação ou instrução de uso	🔥🔥
Segurança básica (auth, validação)	🔥🔥

🖇️ Entrega
- Este desafio deve ser entregue em até 5 dias corridos a partir do recebimento deste documento.
- Suba em um repositório público (GitHub, GitLab, Bitbucket) OU envie um ZIP para monica@aiqfome.com ou (44) 8831-0262
- Inclua instruções claras de como rodar o projeto.
- “Leia-me” é bem-vindo: explique suas escolhas.

Boa sorte e mostre sua identidade no código ❤