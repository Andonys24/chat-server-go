# chat-server-go

Servidor de chat en tiempo real implementado en Go con arquitectura concurrente y protocolo personalizado.

## 🚀 Características

- **Arquitectura Full-Duplex**: Cliente y servidor se comunican de forma asíncrona
- **Concurrencia Segura**: Uso de canales de Go para evitar race conditions
- **Mensajería Privada**: Envío de mensajes directos entre usuarios
- **Broadcast Global**: Difusión de mensajes a todos los usuarios conectados
- **Gestión de Usuarios**: Listado en tiempo real de usuarios conectados
- **Protocolo Personalizado**: Sistema eficiente `HEADER|CONTENT` con delimitadores
- **UX Mejorada**: Cliente con colores ANSI y limpieza de pantalla fluida
- **Logs con Timestamp**: Trazabilidad completa de eventos del servidor
- **Apagado Ordenado**: Limpieza automática de conexiones y recursos

## 📋 Requisitos

- Go 1.16 o superior
- Sistema operativo compatible con conexiones TCP (Linux, macOS, Windows)

## 🛠️ Instalación

1. Clona el repositorio:
```sh
git clone <tu-repositorio>
cd chat-server-go
```

2. Crea un archivo `.env` en la raíz del proyecto (opcional):
```env
PORT=8080
HOST=localhost
MAX_CONNECTIONS=50
```

3. Instala las dependencias:
```sh
go mod download
```

## 🎮 Uso

### Iniciar el Servidor

```sh
go run cmd/server/main.go
```

El servidor iniciará en `localhost:8080` (por defecto) y mostrará:
```
[SISTEMA] Servidor iniciado en localhost:8080
[SISTEMA] Máximo de conexiones: 50
ADMIN >
```

Para detener el servidor, escribe `exit` en el prompt administrativo.

### Conectar un Cliente

```sh
go run cmd/client/main.go
```

Al conectar, se te pedirá ingresar un nickname (3-12 caracteres, debe empezar con letra).

## 📡 Comandos del Cliente

| Comando | Descripción | Formato |
|---------|-------------|---------|
| `/all <mensaje>` | Envía mensaje a todos los usuarios | `/all Hola a todos` |
| `/msg <usuario> <mensaje>` | Envía mensaje privado | `/msg juan Hola Juan` |
| `/users` | Lista usuarios conectados | `/users` |
| `/clear` | Limpia la consola | `/clear` |
| `/exit` | Desconecta del servidor | `/exit` |

## 🏗️ Arquitectura

```
chat-server-go/
├── cmd/
│   ├── client/
│   │   └── main.go          # Punto de entrada del cliente
│   └── server/
│       └── main.go          # Punto de entrada del servidor
├── internal/
│   ├── chat/
│   │   ├── handler.go       # Gestión del ciclo de vida de conexiones
│   │   ├── hub.go           # Gestor concurrente de usuarios y mensajes
│   │   ├── protocol.go      # Constantes del protocolo de comunicación
│   │   ├── transport.go     # Serialización/deserialización de mensajes
│   │   ├── user.go          # Modelo de usuario
│   │   └── validators.go    # Validación de nicknames
│   └── config/
│       └── config.go        # Carga de configuración y variables de entorno
└── go.mod
```

### Componentes Principales

#### Hub ([internal/chat/hub.go](internal/chat/hub.go))
- Gestiona conexiones concurrentes mediante canales
- Mantiene `map[string]*User` para búsquedas O(1)
- Implementa patrón Petición/Respuesta para operaciones seguras
- Canales principales:
  - `Register`: Registro de nuevos usuarios
  - `Unregister`: Eliminación de usuarios desconectados
  - `Broadcast`: Mensajes globales
  - `PrivateMsg`: Mensajes privados
  - `UserRequest`: Solicitudes de lista de usuarios

#### Handler ([internal/chat/handler.go](internal/chat/handler.go))
- Gestiona el ciclo de vida completo de cada conexión
- Valida nicknames con `IsValidNickname`
- Enruta comandos al Hub correspondiente
- Implementa limpieza automática con `defer`
- Niveles de log: `[SISTEMA]`, `[CONEXIÓN]`, `[DESCONEXIÓN]`, `[ADVERTENCIA]`, `[ERROR]`

#### Transport ([internal/chat/transport.go](internal/chat/transport.go))
- Protocolo personalizado `HEADER|CONTENT` con `bufio`
- Serialización eficiente con delimitadores de nueva línea
- Métodos `Send()` y `Receive()` para comunicación bidireccional

#### Protocol ([internal/chat/protocol.go](internal/chat/protocol.go))
Constantes del protocolo:
- Comandos: `CMD_ENTER`, `CMD_EXIT`, `CMD_ALL`, `CMD_MESSAGE`, `CMD_USERS`, `CMD_CLEAN_CONSOLE`
- Respuestas: `RESP_OK_ENTER`, `RESP_ERROR_ENTER`, `RESP_MSG_FROM`, `RESP_USERS_LIST`, `RESP_INFO`
- Tipos de información: `INFO_TYPE_ENTER`, `INFO_TYPE_EXIT`

## 🔒 Concurrencia y Seguridad

- **Sin Mutex Globales**: Uso exclusivo de canales para sincronización
- **Patrón Petición/Respuesta**: Previene race conditions en acceso a datos compartidos
- **Limpieza Automática**: Liberación de recursos en desconexiones normales y abruptas
- **Inicialización Correcta**: Todos los canales inicializados en `NewHub`
- **Deserialización Robusta**: `strings.SplitN` para soportar `|` en contenido de mensajes

## 🎨 Experiencia de Usuario (UX)

### Cliente
- **Colores ANSI**: Mensajes temáticos diferenciados por color
- **Limpieza Fluida**: Secuencias `\033[H\033[2J` para limpieza real de pantalla
- **Prompt Dinámico**: Uso de `\r` y `\033[K` para evitar colisiones visuales
- **Arquitectura Full-Duplex**: Goroutine para escucha + hilo principal para entrada
- **Traductor de Protocolo**: Convierte respuestas del servidor en mensajes amigables

### Servidor
- **Timestamps Automáticos**: Integración con `log.Printf` para trazabilidad
- **Prompt Administrativo**: `ADMIN >` para gestión del servidor
- **Arranque Priorizado**: Servicios activos antes de habilitar consola administrativa

## ⚙️ Configuración

Variables de entorno (archivo `.env` o valores predeterminados):

| Variable | Descripción | Valor por Defecto |
|----------|-------------|-------------------|
| `PORT` | Puerto del servidor | `8080` |
| `HOST` | Host del servidor | `localhost` |
| `MAX_CONNECTIONS` | Conexiones simultáneas máximas | `50` |

La configuración se carga desde [internal/config/config.go](internal/config/config.go) con recuperación segura si el archivo `.env` falta.

## 🧪 Validaciones

### Nickname ([internal/chat/validators.go](internal/chat/validators.go))
- Longitud: 3-12 caracteres
- Formato: Debe empezar con letra
- Caracteres permitidos: Alfanuméricos
- Validación vía expresión regular

## 🔄 Flujo de Conexión

1. Cliente se conecta al servidor
2. Servidor solicita nickname
3. Cliente envía nickname → Validación
4. Si es válido: `RESP_OK_ENTER` → Registro en Hub → Broadcast `INFO_TYPE_ENTER`
5. Si es inválido: `RESP_ERROR_ENTER` → Desconexión
6. Cliente entra en bucle de comandos
7. Al desconectar: Limpieza automática → Broadcast `INFO_TYPE_EXIT`

## 📄 Licencia

Este proyecto está bajo la Licencia MIT.

## 👤 Autor

Desarrollado como proyecto educativo de sistemas concurrentes en Go.

---

**Nota**: Este servidor está diseñado con fines educativos. Para uso en producción, considera añadir autenticación, cifrado TLS y validaciones adicionales de seguridad.
