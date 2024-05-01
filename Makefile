# test rpc call from master to worker
.PHONY: all master worker clean build-all


BINDIR := bin
# Build master executable
master:
	@echo "Building master..."
	@go build -o $(BINDIR)/master ./master/main.go
	@echo "Master built successfully!"

# Build worker executable
worker:
	@echo "Building worker..."
	@go build -o $(BINDIR)/worker ./worker/main.go
	@echo "Worker built successfully!"


build-all: master worker

start-up:
	@echo "Starting zookeeper,master1,master2,worker"
	docker-compose up --build -d
	@echo "Zookeeper,master1,master2,worker started successfully!"

log:
	docker-compose logs -n 20

log-master1:
	docker-compose logs -n 20 master1

log-master2:
	docker-compose logs -n 20 master2

log-worker:
	docker-compose logs -n 20 worker


rpc-call:
	 docker-compose exec worker /bin/bash -c "go test -v ./master"

stop-master1:
	docker-compose stop master1
	@echo "Master1 stopped successfully!"

stop-master2:
	docker-compose stop master2
	@echo "Master2 stopped successfully!"

start-master1:
	docker-compose start master1
	@echo "Master1 started successfully!"

start-master2:
	docker-compose start master2
	@echo "Master2 started successfully!"

stop-and-remove-all:
	docker-compose down --remove-orphans
	@echo "All stopped and removed successfully!"

clean:
	@echo "Cleaning up..."
	@docker-compose down --rmi all --volumes
	@echo "Cleanup complete!"
