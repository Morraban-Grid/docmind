# 📋 ITERACIÓN 6: gRPC Integration Between Services

**Fecha**: 2026-03-12
**Status**: ✅ EN PROGRESO

---

## 🎯 Objetivo

Implementar comunicación gRPC entre el servicio Go y el servicio Python para:
- Indexing de documentos (Go → Python)
- Eliminación de embeddings (Go → Python)
- Consultas RAG (Go → Python)

---

## 📁 ARCHIVOS Y CARPETAS CREADOS

### ✅ NUEVOS ARCHIVOS CREADOS

#### Proto Definition
```
proto/
└── rag_service.proto                  ✅ NUEVO
    ├── IndexDocumentRequest
    ├── IndexDocumentResponse
    ├── DeleteDocumentRequest
    ├── DeleteDocumentResponse
    ├── QueryDocumentRequest
    ├── QueryDocumentResponse
    └── RAGService (3 métodos)
```

#### Python gRPC Server
```
services/python-rag-service/grpc_server/
├── __init__.py                        ✅ NUEVO
├── rag_service.py                     ✅ NUEVO (RAGServicer class)
└── server.py                          ✅ NUEVO (RAGServiceGrpcServer class)
```

#### Go gRPC Client
```
services/go-user-service/internal/grpc_client/
├── README.md                          ✅ NUEVO (Instrucciones)
└── client.go                          ✅ NUEVO (RAGClient class)
```

#### Configuration Files
```
services/python-rag-service/
├── .env.example                       ✅ MODIFICADO (agregadas variables gRPC)
└── requirements.txt                   ✅ MODIFICADO (agregadas dependencias gRPC)

services/go-user-service/
└── .env.example                       ✅ NUEVO (variables de configuración)

services/python-rag-service/
└── config.py                          ✅ MODIFICADO (agregadas variables gRPC)
```

---

## 🔧 COMPONENTES IMPLEMENTADOS

### 1. **Proto Definition** (`proto/rag_service.proto`)

**Mensajes:**
- `IndexDocumentRequest` - Solicitud de indexing
- `IndexDocumentResponse` - Respuesta de indexing
- `DeleteDocumentRequest` - Solicitud de eliminación
- `DeleteDocumentResponse` - Respuesta de eliminación
- `QueryDocumentRequest` - Solicitud de consulta
- `QueryDocumentResponse` - Respuesta de consulta

**Servicio RAGService:**
- `IndexDocument(IndexDocumentRequest) → IndexDocumentResponse`
- `DeleteDocument(DeleteDocumentRequest) → DeleteDocumentResponse`
- `QueryDocument(QueryDocumentRequest) → QueryDocumentResponse`

### 2. **Python gRPC Server** (`grpc_server/`)

**RAGServicer:**
- `index_document()` - Indexa documento
- `delete_document()` - Elimina embeddings
- `query_document()` - Consulta documentos (placeholder para Iteración 7)

**RAGServiceGrpcServer:**
- `start()` - Inicia servidor gRPC
- `stop()` - Detiene servidor gRPC

### 3. **Go gRPC Client** (`internal/grpc_client/`)

**RAGClient:**
- `NewRAGClient()` - Crea cliente gRPC
- `IndexDocument()` - Llama IndexDocument
- `DeleteDocument()` - Llama DeleteDocument
- `QueryDocument()` - Llama QueryDocument
- `HealthCheck()` - Verifica disponibilidad
- `Close()` - Cierra conexión

---

## 🔐 SEGURIDAD

### ✅ Medidas de Seguridad Implementadas

1. **Credenciales Seguras**
   - ✅ `.env.example` con placeholders (NO credenciales reales)
   - ✅ Variables de entorno para todas las credenciales
   - ✅ `.gitignore` previene commit de `.env`

2. **Comunicación gRPC**
   - ✅ Estructura lista para TLS (actualmente insecure para desarrollo)
   - ✅ Timeouts configurados para prevenir bloqueos
   - ✅ Manejo de errores de conexión

3. **Logging**
   - ✅ Sin credenciales en logs
   - ✅ Logging de eventos de conexión
   - ✅ Logging de errores con contexto

4. **Validación**
   - ✅ Validación de parámetros en servicer
   - ✅ Manejo de excepciones
   - ✅ Mensajes de error descriptivos

---

## 📊 CONFIGURACIÓN

### Variables de Entorno - Python

```env
# gRPC Server
GRPC_HOST=0.0.0.0
GRPC_PORT=50051
```

### Variables de Entorno - Go

```env
# Python gRPC Service
PYTHON_GRPC_HOST=localhost
PYTHON_GRPC_PORT=50051
```

---

## ⏱️ TIMEOUTS

| Operación | Timeout |
|-----------|---------|
| IndexDocument | 30 segundos |
| DeleteDocument | 10 segundos |
| QueryDocument | 45 segundos |
| HealthCheck | 5 segundos |

---

## 📝 PRÓXIMOS PASOS

### Generación de Código Proto

**Para Go:**
```bash
cd proto
protoc --go_out=../services/go-user-service/internal/grpc_client \
       --go-grpc_out=../services/go-user-service/internal/grpc_client \
       rag_service.proto
```

**Para Python:**
```bash
cd proto
python -m grpc_tools.protoc -I. --python_out=../services/python-rag-service \
       --grpc_python_out=../services/python-rag-service \
       rag_service.proto
```

### Integración en Flujos

1. **Document Upload Flow**
   - Después de upload a MinIO
   - Llamar gRPC IndexDocument
   - Actualizar estado a "indexed"

2. **Document Deletion Flow**
   - Antes de eliminar de BD
   - Llamar gRPC DeleteDocument
   - Eliminar embeddings de Qdrant

3. **Query Flow**
   - Recibir query en Go
   - Llamar gRPC QueryDocument
   - Retornar respuesta al usuario

---

## ✅ VALIDACIONES

- ✅ Proto file válido
- ✅ Python gRPC server structure correcta
- ✅ Go gRPC client structure correcta
- ✅ Configuración de variables de entorno
- ✅ Dependencias actualizadas
- ✅ Seguridad validada

---

## 🎯 ESTADO

**Estructura**: ✅ COMPLETADA
**Generación de código**: ⏳ PENDIENTE (usuario debe ejecutar protoc)
**Integración en flujos**: ⏳ PRÓXIMA FASE

---

## 📚 REFERENCIAS

- Proto file: `proto/rag_service.proto`
- Python server: `services/python-rag-service/grpc_server/`
- Go client: `services/go-user-service/internal/grpc_client/`
- Documentación: `services/go-user-service/internal/grpc_client/README.md`

**Status**: ✅ ESTRUCTURA COMPLETADA, LISTA PARA GENERACIÓN DE CÓDIGO
