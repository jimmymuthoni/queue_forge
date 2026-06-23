db_login:
	psql ${DATABASE_URL}

db_create_migration:
	migrate create -ext sql -dir migrations -seq ${name}

db_run_migrations:
	migrate -database ${DATABASE_URL} -path migrations up