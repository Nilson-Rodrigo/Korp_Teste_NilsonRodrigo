// Package database contém a configuração e conexão com o banco de dados.
//
// Este arquivo é um placeholder para futura integração com PostgreSQL.
// Atualmente, a aplicação utiliza repositórios em memória.
// Para migrar para PostgreSQL, implemente as interfaces de repositório
// definidas em ports/repository.go utilizando um driver como pgx ou database/sql.
package database

// PostgresConfig contém as configurações de conexão com o PostgreSQL.
type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ConnectionString retorna a string de conexão formatada para o PostgreSQL.
func (c PostgresConfig) ConnectionString() string {
	return "host=" + c.Host +
		" port=" + itoa(c.Port) +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode
}

// itoa converte um inteiro para string sem dependências externas.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	negative := n < 0
	if negative {
		n = -n
	}
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	if negative {
		result = "-" + result
	}
	return result
}
