# ⚠️ Suposiciones y Limitaciones

Este documento describe las **suposiciones de diseño** y las **limitaciones actuales** del servicio `admira-service`.

## ✅ Suposiciones de Diseño

* **Matching por UTM**: se asume que los parámetros UTM son consistentes entre Ads y CRM.
* **Lead = Click**: cada click de Ads se considera un lead potencial para los cálculos.
* **Currency única**: todos los montos monetarios están en la misma moneda (USD).
* **Timezones**: todas las fechas se procesan en UTC para consistencia.

## ⚠️ Limitaciones Actuales

* **Persistencia**: los datos se almacenan en memoria (se pierden al reiniciar el servicio).
* **Escalabilidad**: diseñado para cargas moderadas (\~100 RPS).
* **Cache**: no implementado para endpoints de métricas.
* **Autenticación**: no requiere autenticación (solo para desarrollo).

## 📊 Limitaciones de Datos

* **UTMs incompletos**: algunos registros pueden tener UTMs parciales o faltantes.
* **Data Latency**: los datos de CRM pueden llegar con retraso frente a los de Ads.
* **Attribution Window**: ventana de atribución fija (no configurable).
* **Currency Conversion**: no soporta conversión entre monedas.


## 🔗 Dependencias Externas

* **Mocky.io**: los endpoints deben estar disponibles y responder en <30s.
* **Rate Limiting**: no hay protección contra límites de APIs externas.
* **SSL**: se requieren certificados válidos para conexiones HTTPS.


## 🚫 Consideraciones de Producción

El servicio **NO DEBE usarse en producción sin**:

* Base de datos persistente (ej. PostgreSQL).
* Sistema de autenticación (JWT u OAuth).
* Rate limiting y protección DDoS.
* Monitoring y alerting (Prometheus/Grafana).
* Backup y recovery procedures.
