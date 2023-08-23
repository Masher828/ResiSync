FROM golang:1.20

WORKDIR /app

# windows
# 64 -> RESI_SYNC_GOOS=windows RESI_SYNC_GOARCH=amd64

# 32 -> RESI_SYNC_GOOS=windows RESI_SYNC_GOARCH=386

# mac
# 64 -> RESI_SYNC_GOOS=darwin RESI_SYNC_GOARCH=amd64

# 32 -> RESI_SYNC_GOOS=darwin RESI_SYNC_GOARCH=386

# linux
# 64 -> RESI_SYNC_GOOS=linux RESI_SYNC_GOARCH=amd64

# 32 -> RESI_SYNC_GOOS=linux RESI_SYNC_GOARCH=386

ENV GOOS=darwin
ENV GOARCH=amd64

COPY . /app

# ENTRYPOINT [ "chmod", "+x", "/app/dockerscript.sh", ";", "/app/dockerscript.sh"]
