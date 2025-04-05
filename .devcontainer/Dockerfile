# Force pulling the ARM64 variant of the official Go image
FROM --platform=linux/arm64 golang:1.24-bullseye

RUN apt-get update -y \
  && apt-get install -y --no-install-recommends \
     git curl ca-certificates apt-utils bash sudo \
  && rm -rf /var/lib/apt/lists/*