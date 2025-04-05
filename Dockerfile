FROM golang:1.24-bullseye

# Install bash and other dev tools Gitpod needs
RUN apt-get update && \
    apt-get install -y bash curl git vim && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /workspace