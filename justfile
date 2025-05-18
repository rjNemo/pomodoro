default:
  @just --list

# Build the binary
build:
  @go build -ldflags="-s -w" -o dist/pom .

# Run the binary
run: build
  @./dist/pom

dev:
  @air
