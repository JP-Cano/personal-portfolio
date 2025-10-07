build:
	@echo "✨ Building backend..."
	cd backend && go build -o ../dist/backend
	@echo "✨ Building frontend..."
	cd frontend && npm install && npm run build

dev:
	@echo "🚀 Starting backend (dev)..."
	cd backend && go run main.go &
	@echo "🚀 Starting frontend (dev)..."
	cd frontend && npm run dev

clean:
	@echo "🗑️  Cleaning dist/ and node_modules..."
	rm -rf dist
	cd frontend && rm -rf node_modules package-lock.json .svelte-kit