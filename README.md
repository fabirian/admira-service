# 🎯 Admira Service

Servicio de ETL en Go para integrar datos de Ads y CRM, generando métricas de marketing y revenue para Admira.

## 🚀 Quick Start

### Prerrequisitos  
- Go 1.21+  
- Docker (opcional)  
- Docker Compose (opcional)

### Ejecutar con Docker (Recomendado)  
```  
# Clonar repositorio  
git clone https://github.com/fabirian/admira-service  
cd admira-service

# Configurar variables de entorno  
cp .env.example .env  
# Editar .env con tus URLs de Mocky

# Ejecutar con Docker Compose  
docker-compose up --build  
```

### Ejecutar localmente  
```  
# Configurar entorno  
cp .env.example .env  
# Editar .env con tus URLs

# Instalar dependencias  
go mod download

# Ejecutar  
go run ./cmd/server/main.go  
```

### Ejecutar tests  
```  
go test ./test/... -v  
```

## 📊 Endpoints API

### Ingestar datos  
```  
POST /api/v1/ingest/run?since=2025-08-01  
```  
Params:  
`since\`: Fecha opcional para ingestar datos desde (YYYY-MM-DD)

### Obtener métricas por canal  
```  
GET /api/v1/metrics/channel?channel=google\_ads&from=2025-08-01&to=2025-08-31&limit=10&offset=0  
```  
Params:  
`channel\`: Canal de ads (google\_ads, meta\_ads, etc.)  
`from\`: Fecha inicio (YYYY-MM-DD)  
`to\`: Fecha fin (YYYY-MM-DD)  
`limit\`: Límite de resultados (opcional)  
`offset\`: Offset para paginación (opcional)

### Obtener métricas por campaña  
```  
GET /api/v1/metrics/funnel?utm\_campaign=back\_to\_school&from=2025-08-01&to=2025-08-31  
```  
Params:  
`utm\_campaign\`: Nombre de la campaña UTM  
`from\`, \`to\`, \`limit\`, \`offset\`: Igual que above

### Exportar datos (Opcional)  
```  
POST /api/v1/export/run?date=2025-08-01  
```

### Health Checks  
```  
GET /healthz  # Health check básico  
GET /readyz   # Check de dependencias  
```

## 🎯 Ejemplos de Uso

### 1. Ingestar datos desde una fecha específica  
```  
curl -X POST "http://localhost:8080/api/v1/ingest/run?since=2025-08-01"  
```
### 2. Obtener métricas por canal  
```  
curl "http://localhost:8080/api/v1/metrics/channel?channel=google\_ads&from=2025-08-01&to=2025-08-31"  
```

### 3. Obtener métricas por campaña UTM  
```  
curl "http://localhost:8080/api/v1/metrics/funnel?utm\_campaign=back\_to\_school&from=2025-08-01"  
```

### 4. Exportar datos con firma HMAC  
```  
curl -X POST "http://localhost:8080/api/v1/export/run?date=2025-08-01"  
```

### 5. Health Checks  
```  
curl "http://localhost:8080/healthz"  
curl "http://localhost:8080/readyz"  
```

### Ejemplo de Respuesta Exitosa:  
```  
 {  
   "date": "2025-08-01",  
   "channel": "google\_ads",  
   "campaign\_id": "C-1001",  
   "clicks": 1200,  
   "impressions": 45000,  
   "cost": 350.75,  
   "leads": 25,  
   "opportunities": 8,  
   "closed\_won": 3,  
   "revenue": 15000,  
   "cpc": 0.292,  
   "cpa": 14.03,  
   "cvr\_lead\_to\_opp": 0.32,  
   "cvr\_opp\_to\_won": 0.375,  
   "roas": 42.75  
 }  
```

## 🔧 Configuración

### Variables de Entorno (.env)  
```env  
# URLs de APIs externas  
ADS\_API\_URL=https://api.mocki.io/v2/blsr85p2  
CRM\_API\_URL=https://api.mocki.io/v2/jme1ox5o

# Exportación (Opcional)  
SINK\_URL=https://api.ejemplo.com/sink  
SINK\_SECRET=admira\_secret\_example

# Configuración del servidor  
PORT=8080  
TIMEOUT=30s  
MAX\_RETRIES=3  
RETRY\_DELAY=1s

# Logging  
LOG\_LEVEL=info  
LOG\_FORMAT=json

# Entorno  
ENVIRONMENT=development  
```

## 🏗️ Arquitectura

```  
admira-service/  
├── 📁 cmd/                 # Punto de entrada  
├── 📁 internal/            # Código interno  
│   ├── api/               # Handlers y routers  
│   ├── config/            # Configuración  
│   ├── etl/              # Procesamiento ETL  
│   ├── models/           # Structs de datos  
│   └── service/          # Lógica de negocio  
├── 📁 pkg/                # Librerías exportables  
│   ├── ads/              # Client de Ads API  
│   ├── crm/             # Client de CRM API  
│   └── metrics/         # Cálculo de métricas  
├── 📁 test/              # Tests unitarios  
└── 📁 deploy/            # Configuración Docker  
```

## 🧪 Testing

```  
# Ejecutar todos los tests  
go test ./... -v

# Ejecutar tests específicos  
go test ./test/... -v  
go test ./pkg/metrics/... -v  
```

## 📦 Deployment

### Docker  
```  
docker build -t admira-service .  
docker run -p 8080:8080 --env-file .env admira-service  
```

### Kubernetes (Opcional)  
Ver directorio \`deploy/kubernetes/\`

## ⚠️ Suposiciones y Limitaciones

### Suposiciones de Diseño  
1. Matching por UTM: Se asume que los parámetros UTM son consistentes entre Ads y CRM  
2. Un leads = Un click: Cada click de Ads se considera un lead potencial para cálculos  
3. Currency: Todos los montos monetarios están en la misma currency (USD)  
4. Timezones: Todas las fechas se procesan en UTC para consistencia

### Limitaciones Actuales  
1. Almacenamiento Volátil: Datos en memoria (se pierden al reiniciar el servicio)  
2. Escalabilidad: Diseñado para cargas moderadas (~100 RPS)  
3. Persistence: No hay base de datos persistente (solo memoria)  
4. Cache: No implementado para endpoints de métricas  
5. Autenticación: No requiere autenticación en endpoints (para desarrollo)

### Limitaciones de Datos  
1. UTMs Incompletos: Algunos registros pueden tener UTMs parciales o missing  
2. Data Latency: Los datos de CRM pueden tener delay vs datos de Ads  
3. Attribution Window: Ventana de atribución fija (no configurable)  
4. Currency Conversion: No soporta conversión entre monedas

### Dependencias Externas  
1. Mocky.io: Los endpoints deben estar disponibles y responder en \<30s  
2. Rate Limiting: Sin protección contra rate limiting de APIs externas  
3. SSL: Certificados SSL válidos requeridos para conexiones HTTPS

### Consideraciones de Producción  
⛔ **NO USAR EN PRODUCCIÓN SIN:**

- 🔴 Base de datos persistente (PostgreSQL)
- 🔴 Sistema de autenticación (JWT/OAuth)  
- 🔴 Rate limiting y protección DDoS
- 🔴 Monitoring y alerting (Prometheus/Grafana)
- 🔴 Backup y recovery procedures

## 🛠️ Desarrollo

### Estructura de commits  
`feat:\` Nueva funcionalidad  
`fix:\` Corrección de bugs  
`docs:\` Documentación  
`test:\` Tests  
`refactor:\` Refactorización

### Logging  
Logs estructurados con niveles (debug, info, warn, error) y request IDs.

## 🧑‍💻 Autor

Arlex Fabian Galindez Rivera
📧 Contacto: fabirir99@gmail.com
