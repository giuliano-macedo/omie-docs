# [Gerador não oficial de documentações da API da Omie](https://giuliano-macedo.github.io/omie-docs/)
![Coverage](https://img.shields.io/badge/Coverage-92.2%25-brightgreen)

Gere documentações OpenAPI, coleções do Postman e coleções do [Bruno](https://github.com/usebruno/bruno) da API da omie através da [documentação Omie](https://developer.omie.com.br/service-list/).

## Uso

Use `go-task run` para rodar rodar o crawler, conversor e bundler do openai, por padrão o programa irá salvar todos os arquivos necessários na pasta `./bundle`

Use `go-task build` para gerar o binário `bin/generate_openapi`
