TARGET = pidstop
SRC_DIR = cmd
OUT_DIR = bin
GOFLAGS = -v
GOCMD = go
BUILD_CMD = CGO_ENABLED=0 $(GOCMD) build $(GOFLAGS) -o $(OUT_DIR)/$(TARGET) ./$(SRC_DIR)/.

.PHONY: all build release clean

# Default target
all: build

# Build the executable
build:
	@mkdir -p $(OUT_DIR)
	$(BUILD_CMD)

# Build a release version with static linking
release:
	@mkdir -p $(OUT_DIR)
	CGO_ENABLED=0 $(BUILD_CMD)

# Run the program
run: build
	@./$(OUT_DIR)/$(TARGET)

# Clean generated files
clean:
	@rm -rf $(OUT_DIR)
