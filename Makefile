create_postgres_database:
	docker exec -it postgres createdb --username=admin --owner=admin ${POSTGRES_DB}

drop_postgres_database:
	docker exec -it postgres dropdb ${POSTGRES_DB}

create_clickhouse_database:
	docker exec -it clickhouse clickhouse-client --query "CREATE DATABASE IF NOT EXISTS $(CLICKHOUSE_DB);"

drop_clickhouse_database:
	docker exec -it clickhouse clickhouse-client --query "DROP DATABASE IF EXISTS $(CLICKHOUSE_DB);"

migration_up:
	migrate -path migrations/postgres/sql -database ${POSTGRES_URL} -verbose up

migration_down:
	migrate -path migrations/postgres/sql -database ${POSTGRES_URL} -verbose down

