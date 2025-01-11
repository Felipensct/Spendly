# Diretórios e arquivos
PROTO_DIR := proto
PROTO_FILES := $(PROTO_DIR)/auth.proto
OUT_DIR := internal/grpc

# Alvo padrão
.PHONY: all
all: generate

# Geração de código a partir dos arquivos Protobuf
.PHONY: generate
generate:
	@echo "Gerando arquivos Protobuf..."
	@mkdir -p $(OUT_DIR)
	@protoc --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	        --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	        --proto_path=$(PROTO_DIR) $(PROTO_FILES)

# Limpeza dos arquivos gerados
.PHONY: clean
clean:
	@echo "Limpando arquivos gerados..."
	@rm -rf $(OUT_DIR)/*
