# --- Builder Stage ---
    FROM golang:1.23.0-alpine as builder

    WORKDIR /app
    
    # Copy go module files and download dependencies
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy all source files
    COPY . .
    
    # Build the application
    RUN go build -o /armur-server ./cmd/server/main.go
    
# --- Final Stage ---
    FROM armur-tools as final
    
    WORKDIR /app
    
    # Copy binary from builder stage
    COPY --from=builder /armur-server /armur-server
    COPY --from=builder /app/pkg/common/cwd.json /app/pkg/common/cwd.json
    # Copy necessary configs
    COPY rule_config /app/rule_config
    
    
    # Expose the port
    EXPOSE 4500
    
    # Run the application
    CMD ["/armur-server"]