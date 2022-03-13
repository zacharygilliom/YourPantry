all:
	make -j 2 backend frontend
.PHONY: all

backend:
	go run cmd/main.go
.PHONY: backend

frontend:
	cd frontend/ && python3 -m http.server
.PHONY: frontend
