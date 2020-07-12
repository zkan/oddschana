init:
	sqlite3 thaichana.db < _scripts/init.sql
run:
	PORT=8000 DB_CONN=thaichana.db go run main.go
