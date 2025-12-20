build:
	@echo "Building backend..."
	cd backend && make build
	@echo "Building frontend..."
	cd frontend && pnpm install && pnpm run build

dev:
	@echo "Starting backend (dev)..."
	cd backend && make dev &
	@echo "Starting frontend (dev)..."
	cd frontend && pnpm run dev

clean:
	@echo "Cleaning dist/ and node_modules..."
	rm -rf dist
	cd frontend && rm -rf node_modules package-lock.json