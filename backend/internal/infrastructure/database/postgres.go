// Package database contém a configuração e conexão com o banco de dados.
//
// Este arquivo é um placeholder para futura integração com PostgreSQL.
// Atualmente, a aplicação utiliza repositórios em memória.
// Para migrar para PostgreSQL, implemente as interfaces de repositório
// definidas em ports/repository.go utilizando um driver como pgx ou database/sql.
package database

import "strconv"

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
		" port=" + strconv.Itoa(c.Port) +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode
}
