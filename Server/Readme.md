# Specs

* O server.go deverá consumir a API contendo o câmbio: https://economia.awesomeapi.com.br/json/last/USD-BRL - OK
* Sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms - OK
* O endpoint do server.go será: /cotacao e a porta do servidor será a :8080 - OK
* Em seguida deverá retornar no formato JSON o resultado para o cliente - OK
* Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida - OK
* O timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms - OK
* Os contextos deverão retornar erro nos logs caso o tempo de execução seja insuficiente - OK


