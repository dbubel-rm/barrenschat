start:
	docker-compose -f docker-compose.yaml up --build -d
stop:
	docker-compose -f docker-compose.yaml down
db:
	# rm cmd/bishopfox-api/bfscans.db
	sqlite3 cmd/bishopfox-api/bfscans.db < schema.sql 
test:
	cd cmd/bishopfox-api/handlers
	go test -v -failfast ./...

.PHONY: start stop db test
