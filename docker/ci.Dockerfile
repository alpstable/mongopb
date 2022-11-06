FROM golang:1.19

# Create the working directory.
WORKDIR /app

COPY . .

# Run the tests.
CMD ["go", "test", "-count", "5", "-v", "./..."]
