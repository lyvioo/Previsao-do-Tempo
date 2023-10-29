# Aplicativo de Previsão do Tempo em Go e MongoDB

Este projeto é uma aplicação web escrita em Go. Ele utiliza a API OpenWeatherMap para buscar previsões do tempo para cidades específicas, armazena os resultados em um banco de dados MongoDB e exibe de forma estilizada em uma página web.

## Funcionalidades

- Consulta previsões do tempo para cidades pré-determinadas via API do OpenWeatherMap.
- Armazena as previsões do tempo em um banco de dados MongoDB.
- Exibe as previsões em uma página web formatada.

## Pré-requisitos

- Go (recomendado Go 1.17 ou superior)
- MongoDB (recomendado v5.0 ou superior)
- Chave de API válida do OpenWeatherMap

## Como executar

1. Clone este repositório:
   ```bash
   git clone https://github.com/lyvioo/Previsao-do-Tempo
Navegue até a pasta do projeto:

bash
Copy code
Configure sua chave da API:
Abra o arquivo main.go e atualize a variável apiKey com sua chave de API do OpenWeatherMap.

Execute o programa:

bash
Copy code
go run main.go
Acesse a previsão do tempo em seu navegador:
Abra seu navegador e vá para http://localhost:8080/previsao.

## Estrutura do Projeto
main.go: O arquivo principal contendo as rotas, lógica de conexão com o MongoDB, fetch da API do OpenWeatherMap e renderização do template.
template.html: Template HTML usado para exibir as previsões do tempo.
Considerações
O projeto foi desenvolvido principalmente para fins didáticos. Certifique-se de que o MongoDB esteja em execução e acessível na URI mongodb://localhost:27017 antes de iniciar a aplicação.

## Contribuições
Contribuições são sempre bem-vindas! Sinta-se à vontade para abrir um issue ou enviar um pull request.

## Licença
Este projeto está licenciado sob a licença MIT.