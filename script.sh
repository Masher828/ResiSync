docker build -t resi_sync_golang:1 .

docker run --name resi_sync_golang_container --rm -d -v ./:/app resi_sync_golang:1 tail -f /dev/null

docker exec resi_sync_golang_container chmod +x /app/dockerscript.sh

docker exec resi_sync_golang_container /app/dockerscript.sh

docker rm -f resi_sync_golang_container

