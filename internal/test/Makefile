# Определяем цель для удаления и создания базы данных, а также запуска тестов
reset-db-test:
	rm -f database/database.db && \
	touch database/database.db && \
	go test -v document_folder_test.go

