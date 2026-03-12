# 📋 Proto Code Generation Guide

**Importante**: Este documento explica cómo generar el código gRPC a partir del archivo proto.

---

## 🔧 Requisitos Previos

### 1. Instalar Protocol Buffers Compiler

**macOS:**
```bash
brew install protobuf
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt-get update
sudo apt-get install protobuf-compiler
```

**Windows:**
Descargar desde: https://github.com/protocolbuffers/protobuf/releases

Verificar instalación:
```bash
protoc --version
```

---

## 🐍 Generación de Código Python

### 1. Instalar herramientas gRPC para Python

```bash
cd services/python-rag-service
pip install grpcio-tools==1.60.0
```

### 2. Generar código Python

```bash
cd proto
python -m grpc_tools.protoc \
  -I. \
  --python_out=../services/python-rag-service \
  --grpc_python_out=../services/python-rag-service \
  rag_service.proto
```

### 3. Archivos generados

Se crearán los siguientes archivos:
```
services/python-rag-service/
├── rag_service_pb2.py          # Definiciones de mensajes
└── rag_service_pb2_grpc.py     # Definiciones de servicio
```

### 4. Actualizar imports en Python

En `services/python-rag-service/grpc_server/rag_service.py`:

```python
# Agregar imports
from rag_service_pb2 import (
    IndexDocumentResponse,
    DeleteDocumentResponse,
    QueryDocumentResponse
)
from rag_service_pb2_grpc import RAGServiceServicer, add_RAGServiceServicer_to_server
```

En `services/python-rag-service/grpc_server/server.py`:

```python
# Agregar imports
from rag_service_pb2_grpc import add_RAGServiceServicer_to_server
from grpc_server.rag_service import RAGServicer
```

---

## 🐹 Generación de Código Go

### 1. Instalar herramientas gRPC para Go

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Verificar que están en PATH:
```bash
which protoc-gen-go
which protoc-gen-go-grpc
```

### 2. Generar código Go

```bash
cd proto
protoc \
  --go_out=../services/go-user-service/internal/grpc_client \
  --go-grpc_out=../services/go-user-service/internal/grpc_client \
  rag_service.proto
```

### 3. Archivos generados

Se crearán los siguientes archivos:
```
services/go-user-service/internal/grpc_client/pb/
├── rag_service.pb.go           # Definiciones de mensajes
└── rag_service_grpc.pb.go      # Definiciones de servicio
```

### 4. Actualizar imports en Go

En `services/go-user-service/internal/grpc_client/client.go`:

```go
// Agregar imports
import (
    pb "github.com/Morraban-Grid/docmind/services/go-user-service/internal/grpc_client/pb"
)
```

### 5. Actualizar go.mod

```bash
cd services/go-user-service
go get google.golang.org/grpc@latest
go get google.golang.org/protobuf@latest
go mod tidy
```

---

## 🔄 Flujo Completo de Generación

### Opción 1: Script Automatizado (Recomendado)

Crear archivo `generate-proto.sh`:

```bash
#!/bin/bash

set -e

echo "Generating gRPC code from proto files..."

# Python
echo "Generating Python code..."
cd proto
python -m grpc_tools.protoc \
  -I. \
  --python_out=../services/python-rag-service \
  --grpc_python_out=../services/python-rag-service \
  rag_service.proto

# Go
echo "Generating Go code..."
protoc \
  --go_out=../services/go-user-service/internal/grpc_client \
  --go-grpc_out=../services/go-user-service/internal/grpc_client \
  rag_service.proto

cd ..

echo "Proto code generation completed successfully!"
```

Ejecutar:
```bash
chmod +x generate-proto.sh
./generate-proto.sh
```

### Opción 2: Generación Manual

Ejecutar los comandos de Python y Go por separado (ver arriba).

---

## ✅ Verificación

### Python

```bash
cd services/python-rag-service
python -c "import rag_service_pb2; print('Python proto import OK')"
```

### Go

```bash
cd services/go-user-service
go build ./internal/grpc_client
```

---

## 🔐 SEGURIDAD - IMPORTANTE

### ✅ Verificaciones de Seguridad

1. **Nunca commitear archivos generados con credenciales**
   - Los archivos `.pb.go` y `_pb2.py` son seguros (no contienen credenciales)
   - Verificar que `.env` NO está en git

2. **Credenciales en variables de entorno**
   - Todas las credenciales deben estar en `.env`
   - `.env` debe estar en `.gitignore`
   - Usar `.env.example` con placeholders

3. **Verificar .gitignore**
   ```bash
   git check-ignore .env
   # Debe retornar: .env
   ```

---

## 🐛 Troubleshooting

### Error: "protoc: command not found"
- Instalar protobuf compiler (ver arriba)
- Verificar que está en PATH: `protoc --version`

### Error: "protoc-gen-go: command not found"
- Instalar Go plugins: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
- Verificar PATH: `which protoc-gen-go`

### Error: "No such file or directory: rag_service.proto"
- Ejecutar desde directorio `proto/`
- Verificar que el archivo existe: `ls -la rag_service.proto`

### Error: "grpc_tools not found"
- Instalar: `pip install grpcio-tools==1.60.0`
- Verificar: `python -m grpc_tools.protoc --version`

---

## 📝 Próximos Pasos

1. ✅ Generar código proto (este documento)
2. ⏳ Actualizar imports en Python y Go
3. ⏳ Implementar métodos en RAGServicer
4. ⏳ Implementar métodos en RAGClient
5. ⏳ Integrar en flujos de documento
6. ⏳ Testear comunicación gRPC

---

## 📚 Referencias

- Proto file: `proto/rag_service.proto`
- Python gRPC: https://grpc.io/docs/languages/python/
- Go gRPC: https://grpc.io/docs/languages/go/
- Protocol Buffers: https://developers.google.com/protocol-buffers

---

**Status**: ✅ GUÍA COMPLETA PARA GENERACIÓN DE CÓDIGO PROTO
