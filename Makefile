main-test:
	rm database/database.db && \
	touch database/database.db && \
	go run main.go	
