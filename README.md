# go-template

## Copiar variables de entorno
```bash
cp .env.example .env
```

## Levantar todo
```bash
make up
```

## Probar Health
```curl
curl -s localhost:8080/healthz
```

## Crear un usuario
```curl
curl -s -X POST localhost:8080/api/v1/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Leo","email":"leo@example.com","phone":"+54..."}'
```

## Listar usuarios
```curl
curl -s 'localhost:8080/api/v1/users?page=1&size=10'
```

## Notificar (RabbitMQ)
```bash
curl -s -X POST localhost:8080/api/v1/notify \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"<uuid-devuelto-en-create>","message":"hola!"}'
```
UI de Rabbit: http://localhost:15672
(guest/guest)
Redpanda (Kafka) expuesto en localhost:19092 si querés usar rpk o un viewer.



configs: 
si pasamos variables de entorno estas tienen prioridad y se sobreescriben de las que les pase
desde el config.yaml o yml.
Entonces estos pasos son:
- Lee config.yaml si existe.
- Aplica defaults.
- Sobrescribe con .env si se define.
- Sobrescribe con variables de entorno activas.

por ejemplo si ejecuto esto en consola:

```bash
export APP_ENV=prod
export HTTP_PORT=:9090
export KAFKA_BROKERS=redpanda:9092,kafka2:9092
make up
```