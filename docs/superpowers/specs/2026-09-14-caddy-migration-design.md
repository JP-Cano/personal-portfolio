# Migrar de nginx + certbot a Caddy

## Contexto

`net::ERR_CERT_DATE_INVALID` recurrente en juanpablocano.com fue causado por un gap
estructural: `certbot` renueva el certificado en disco, pero `nginx` solo lo lee una vez al
arrancar, así que nunca se enteraba de la renovación sin un reload manual. Ese bug puntual ya
se arregló (ver PR #3 — loop de reload cada 6h). Esta migración va más allá: reemplaza
`nginx` + `certbot` por **Caddy**, que maneja obtención, renovación y aplicación del
certificado dentro de un solo proceso, eliminando de raíz esta *clase* de bug (no puede haber
un gap entre "el certificado se renovó" y "el servidor lo está usando" si es el mismo proceso
el que hace ambas cosas).

De paso, esta migración resuelve dos problemas adicionales encontrados durante la
investigación:

1. **Drift entre config de dev y de producción.** Hoy `nginx/conf.d/default.conf` (dev, HTTP
   puro) y la config de producción (generada por `scripts/ssl-setup.sh` desde
   `nginx/templates/default.prod.conf.template`) son archivos distintos que pueden pisarse
   entre sí por error (`scp -r nginx ...`). Con Caddy, un solo `Caddyfile` sirve ambos casos
   usando variables de entorno — no hay dos archivos que puedan desincronizarse.
2. **Puerto incorrecto del frontend en la config de dev.** `nginx/conf.d/default.conf` apunta
   a `frontend:4321`, pero dentro de Docker el frontend siempre escucha en el puerto que fije
   la env var `PORT` (`docker-compose.yml` la fija en `3000` — confirmado en
   `frontend/astro.config.mjs`, adaptador `@astrojs/node` en modo `standalone`, que lee
   `process.env.PORT` directamente). `4321` es el puerto por defecto de `astro dev` corriendo
   *fuera* de Docker, no el del contenedor — la config de dev estaba, y sigue estando,
   apuntando al puerto equivocado. El Caddyfile nuevo usa `3000` en ambos casos.

## Alcance

Incluye:
- Reemplazar los servicios `nginx` y `certbot` por un servicio `caddy` en `docker-compose.yml`.
- Un `Caddyfile` único en la raíz del repo, unificado para desarrollo local y producción.
- Actualizar `.env.example` (nueva variable `SITE_ADDRESS`, se mantiene `SSL_EMAIL`).
- Eliminar `scripts/ssl-setup.sh` (ya no hace falta — Caddy provisiona el certificado solo).
- Limpiar referencias a `nginx`/`certbot` en `scripts/setup-droplet.sh`.
- Actualizar `docs/DEPLOYMENT.md` (quitar el paso manual de SSL, actualizar troubleshooting).
- Agregar un paso de sync (`docker-compose.yml`, `Caddyfile`) a `.github/workflows/deploy.yml`
  para que el droplet deje de depender de `scp` manual.
- Eliminar `nginx/` y `certbot/` del repo y (documentado, manual) del droplet.
- Corregir el healthcheck del servicio `frontend` en `docker-compose.yml`, que hoy prueba el
  puerto `4321` (equivocado, ver arriba) — se encuentra en el mismo archivo que ya estamos
  reescribiendo y es necesario para que el healthcheck no reporte "unhealthy" perpetuamente.

Fuera de alcance:
- Cambiar el proveedor de hosting o la arquitectura de contenedores más allá de la capa de
  reverse proxy/TLS.
- Tocar `backend/` o `frontend/` (excepto el healthcheck de `frontend` en `docker-compose.yml`
  mencionado arriba).

## Diseño

### 1. `docker-compose.yml`

Se eliminan los servicios `nginx` y `certbot` completos. Se agrega:

```yaml
  caddy:
    image: caddy:2.11.4-alpine
    container_name: portfolio-caddy
    restart: unless-stopped
    environment:
      - SITE_ADDRESS=${SITE_ADDRESS:-:80}
      - SSL_EMAIL=${SSL_EMAIL:-}
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:ro
      - ./data/caddy/data:/data
      - ./data/caddy/config:/config
    depends_on:
      - frontend
      - backend
    networks:
      - portfolio-network
    healthcheck:
      test: [ "CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost/health" ]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

`./data/caddy/data` guarda certificados y cuenta ACME (equivalente a lo que hoy es
`certbot/conf`); `./data/caddy/config` guarda el autosave de configuración de Caddy. Ambos
bind-mounts, mismo patrón que ya usa `./data/backend`.

El servicio `frontend` corrige su healthcheck (bug preexistente, mismo archivo):

```yaml
    healthcheck:
      test: [ "CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:3000" ]
```
(antes probaba `:4321`, el puerto equivocado — ver Contexto)

### 2. `Caddyfile` (nuevo, raíz del repo)

```caddyfile
{
	email {$SSL_EMAIL}
}

{$SITE_ADDRESS::80} {
	# Doble ":" a propósito: {$VAR:default} de Caddy separa en el PRIMER ":",
	# así que {$SITE_ADDRESS:80} daría default = "80" (hostname suelto), no
	# ":80" (binding de puerto). {$SITE_ADDRESS::80} sí da default = ":80".
	encode gzip zstd

	handle /health {
		respond "healthy" 200
	}

	handle /api/admin/* {
		reverse_proxy frontend:3000
	}

	handle /api/* {
		reverse_proxy backend:8080
	}

	handle /swagger/* {
		reverse_proxy backend:8080
	}

	handle /certifications/* {
		reverse_proxy backend:8080
	}

	handle {
		reverse_proxy frontend:3000
	}
}
```

Comportamiento:
- **Local** (`SITE_ADDRESS` sin definir en `.env`): Caddy sirve HTTP plano en `:80`, sin TLS,
  igual que hoy. No hace ninguna llamada saliente a Let's Encrypt.
- **Producción** (`SITE_ADDRESS=juanpablocano.com www.juanpablocano.com` en el `.env` del
  droplet): Caddy detecta que es un dominio público, obtiene el certificado automáticamente
  al cargar la configuración (no espera al primer request), sirve HTTP→HTTPS redirect por defecto, y renueva sin
  intervención — todo dentro de este mismo proceso.
- Caddy agrega automáticamente `X-Forwarded-For`, `X-Forwarded-Proto`, `Host`, y reenvía
  cookies y upgrades de WebSocket en `reverse_proxy` sin configuración adicional (a diferencia
  de nginx, que requería declarar cada header a mano).
- El bloque global `{ email {$SSL_EMAIL} }` registra la cuenta ACME con ese correo — mismo
  propósito que `SSL_EMAIL` tenía para `certbot --email`.

### 3. `.env.example`

Reemplazar:
```env
# Domain (update these!)
DOMAIN=yourdomain.com
SSL_EMAIL=your-email@example.com
```
por:
```env
# Site address for Caddy. Leave unset/empty for local development (Caddy serves
# plain HTTP on :80, no TLS). In production, set to your real domain(s),
# space-separated (NOT comma-separated — Caddyfile site addresses are
# whitespace-delimited), e.g.: SITE_ADDRESS="yourdomain.com www.yourdomain.com"
# Caddy then automatically obtains and renews a Let's Encrypt certificate.
SITE_ADDRESS=

# Email for Let's Encrypt certificate registration (used by Caddy)
SSL_EMAIL=your-email@example.com
```

### 4. Archivos eliminados

- `scripts/ssl-setup.sh` — su función entera (generar config HTTPS, correr certbot, reiniciar
  nginx) ahora la hace Caddy automáticamente al arrancar.
- `nginx/` completo (`nginx.conf`, `conf.d/default.conf`, `templates/default.prod.conf.template`).
- En `scripts/setup-droplet.sh`: quitar la creación de `certbot/www`, `certbot/conf`,
  `nginx/conf.d` del `mkdir -p` (línea ~102); agregar `data/caddy/data`, `data/caddy/config`.

### 5. `docs/DEPLOYMENT.md`

- Eliminar por completo "Step 6: Set Up SSL Certificate" (ya no existe ese paso manual).
- Renumerar los pasos siguientes.
- Step 4 (copiar archivos): reemplazar `nginx scripts` por `Caddyfile scripts` en el `scp`
  de ejemplo, y quitar la advertencia sobre no pisar `nginx/conf.d/default.conf` (ya no aplica
  — no hay archivo generado en el droplet que un `scp` pueda pisar; el `Caddyfile` es el mismo
  en cualquier lado).
- Sección de troubleshooting SSL: reemplazar los comandos de `certbot`/`nginx -s reload` por
  los equivalentes de Caddy:
  ```bash
  # Ver logs de Caddy (incluye actividad de emisión/renovación ACME)
  docker compose logs caddy --tail 100

  # Forzar un reload de la config (poco común, Caddy detecta cambios de Caddyfile solo)
  docker exec portfolio-caddy caddy reload --config /etc/caddy/Caddyfile

  # Confirmar el certificado servido
  echo | openssl s_client -connect yourdomain.com:443 -servername yourdomain.com 2>/dev/null | openssl x509 -noout -dates
  ```

### 6. `.github/workflows/deploy.yml`

Agregar un paso de sync antes de `docker compose pull`, usando `appleboy/scp-action` con una
lista explícita de archivos (nunca un directorio completo sin filtrar):

```yaml
      - name: Sync compose & Caddy config to droplet
        uses: appleboy/scp-action@v1.0.0
        with:
          host: ${{ secrets.DROPLET_HOST }}
          username: ${{ secrets.DROPLET_USERNAME }}
          key: ${{ secrets.DROPLET_SSH_KEY }}
          source: "docker-compose.yml,Caddyfile"
          target: "/opt/portfolio"
          strip_components: 0
```

Como el `Caddyfile` no necesita sustitución de placeholders (usa variables de entorno nativas
de Caddy vía `{$VAR}`), la copia es directa — no hace falta ningún paso de generación. El paso
existente `docker compose pull && docker compose up -d --remove-orphans` ya recrea cualquier
servicio cuyo config cambió, así que no hace falta un reload explícito adicional.

### 7. Corte en el droplet (manual, una sola vez)

No tengo acceso SSH al droplet, así que este paso lo ejecuta el usuario. Documentar en el PR:

1. En el `.env` del droplet, agregar `SITE_ADDRESS="juanpablocano.com www.juanpablocano.com"`
   (separado por espacio, no coma — así es como Caddyfile espera múltiples direcciones de
   sitio). **`SSL_EMAIL` debe tener un valor real y no puede estar vacío** — si está vacío,
   Caddy no arranca en absoluto (falla al parsear su config), y esto pasaría justo después de
   que el paso 3 ya tiró nginx/certbot. Confirmar con:
   `grep '^SSL_EMAIL=' /opt/portfolio/.env` antes de continuar.
2. Copiar el `docker-compose.yml` y `Caddyfile` nuevos al droplet (manual esta primera vez;
   los siguientes deploys ya lo hacen solos vía el paso de CI agregado).
3. `docker compose up -d --remove-orphans` — esto detiene y remueve `nginx`/`certbot` como
   huérfanos (ya no están en el compose file) y levanta `caddy`. Downtime esperado: unos
   minutos mientras Caddy obtiene el certificado nuevo de Let's Encrypt.
4. Verificar: `docker compose logs caddy` sin errores de ACME, `curl -I https://juanpablocano.com`
   responde `200`, certificado con expiry ~90 días adelante. (Nota: `curl -I` es HEAD —
   contra `/api/v1/health` específicamente da 404 por una particularidad preexistente del
   router del backend con HEAD, no relacionada a esta migración; usar GET ahí si se necesita
   probar ese endpoint: `curl -s -o /dev/null -w "%{http_code}\n" https://juanpablocano.com/api/v1/health`.)
5. Limpieza (opcional, no bloqueante): `rm -rf nginx/ certbot/` en `/opt/portfolio` del droplet.

## Testing / Verificación

**Local:**
```bash
docker compose up -d
curl -I http://localhost/health          # 200 "healthy"
curl -s -o /dev/null -w "%{http_code}\n" http://localhost/api/v1/health   # proxied al backend (GET — HEAD 404s on this route, pre-existing Gin behavior)
curl -I http://localhost                 # proxied al frontend (puerto 3000 correcto)
```
Confirmar que `docker compose ps` reporta `caddy` y `frontend` como `healthy` (el fix del
puerto del healthcheck del frontend se verifica aquí).

**Producción (droplet, después del corte manual):**
```bash
docker compose logs caddy --tail 50 | grep -i "certificate obtained\|error"
echo | openssl s_client -connect juanpablocano.com:443 -servername juanpablocano.com 2>/dev/null | openssl x509 -noout -dates
curl -s -o /dev/null -w "%{http_code}\n" https://juanpablocano.com/api/v1/health
```
Abrir `https://juanpablocano.com` en el navegador y confirmar el candado sin advertencias.

**CI:** confirmar que un push a `main` sincroniza `docker-compose.yml` y `Caddyfile` al
droplet antes de que corra `docker compose up -d` (revisar logs del job de GitHub Actions).

## Riesgos / notas

- Caddy pedirá certificados **nuevos** a Let's Encrypt (no reutiliza los de `certbot/`) — no
  hay problema de rate limit porque es un dominio con historial normal de emisiones, muy por
  debajo del límite de 50 certificados/semana por dominio registrado.
- Mientras el corte no se ejecute en el droplet (paso manual, sección 7), el sitio sigue
  sirviendo con nginx/certbot tal cual quedó después del PR #3 — no hay nada roto por este
  cambio hasta que se aplique manualmente.
- `caddy reload` desde dentro del contenedor requiere que el `Caddyfile` montado ya tenga los
  cambios (el volumen es `:ro`); para cambios de Caddyfile, un `docker compose up -d --force-recreate caddy`
  es más simple que `caddy reload` y es lo que ya hace el pipeline de CI.
