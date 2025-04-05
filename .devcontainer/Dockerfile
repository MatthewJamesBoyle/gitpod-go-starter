# Base it on official golang + Debian Bullseye
FROM golang:1.24-bullseye

# Make sure we have a shell, user management tools, etc.
# (especially needed if Dev Containers tries to create a user or uses sudo)
RUN apt-get update -y \
    && apt-get install -y --no-install-recommends \
       git \
       curl \
       ca-certificates \
       apt-utils \
       bash \
       sudo \
    && rm -rf /var/lib/apt/lists/*