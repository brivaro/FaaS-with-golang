# FaaS-with-golang 🚀 🚀 🚀

Este proyecto implementa una plataforma de **Functions as a Service (FaaS)** utilizando Go. Permite ejecutar funciones serverless en contenedores Docker de forma escalable, manejando la comunicación a través de NATS como sistema de mensajería. La arquitectura del proyecto es modular y cada componente se ha organizado en carpetas especializadas.

## Instalación 💻📥

Para instalar el servicio FaaS, se requiere seguir los siguientes pasos:

- **Primer paso**: Clonar el repositorio
```bash
git clone https://github.com/alerone/FaaS-with-golang.git
```
- **Segundo paso**: Crear un archivo de variables de entorno (``.env``) en la raíz de la carpeta [apiServer](./apiServer) siguiendo la estructura definida en el archivo [.env.example](./apiServer/.env.example) que se puede observar en ese mismo directorio.
- **Tercer paso**: Para que apisix sepa qué ``secret`` utiliza el apiServer para generar los tokens JWT debemos copiar el valor que se haya puesto en el campo ``SECRET`` del .env que se ha creado en el **Segundo paso** de la instalación en el consumidor que se encuentra abajo del todo del archivo [routes.yaml](./apisix_conf/routes.yaml) para que utilice el mismo secret que el apiServer.

```yaml
consumers:
  - username: faas_jwt_consumer
    plugins:
      jwt-auth:
        key: faas_jwt_consumer
        secret: "your_secret_password_from_dot_env"
        algorithm: HS256
#END
```
## Cómo se usa el Function as a Service? 🎮 🖱️

Para utilizar el servicio, tras realizar la instalación, podemos iniciar el sistema distribuido ejecutando ``docker-compose`` desde la raíz del proyecto
```bash
docker-compose up --build
```
Este comando construirá las imágenes de apiServer y worker y se traerá del Docker hub las imágenes de Nats server y Apache apisix. Tras traerse las imágenes lanzará los contenedores de estas imágenes y lanzará además 7 workers (por defecto).

### Ajustar la escalabilidad de workers
Para lanzar más o menos workers, se puede acceder al archivo [docker-compose](./docker-compose.yaml) y en el servicio worker, se puede ajustar la escalabilidad cambiando la propiedad ``replicas``.
```yaml
  worker:
    image: faas-worker
    build: ./worker
    depends_on:
      - nats
      - apiserver
    environment:
      - nats_url=nats://nats:4222
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    networks:
      - faas-network
    deploy:
      replicas: 7
```

### Rutas del servicio FaaS
Como se utiliza el reverse proxy de Apache Apisix, todas las peticiones se deben lanzar al puerto ``9080`` de la máquina que tenga el servidor lanzado con el docker-compose. El docker-comopse está compuesto de tal forma que no hay ninguna salida al exterior (port forwarding) más que por el servicio apisix.

Las distintas rutas del servicio son las siguientes:

- #### ``POST /register``:
Sirve para registrar un usuario en el sistema FaaS. Se requiere subir un JSON con la estructura siguiente:
```json
{
  "Username": "hola1234",
  "Password": "contraseña1234"
}
```
>[!WARNING]
>El nombre de usuario NO debe contener los siguientes carácteres especiales: '.', '*', '>', '$'

- #### ``POST /login``:
Sirve para iniciar sesión con un usuario previamente registrado. Se requiere un JSON con la estructura siguiente:
```json
{
  "Username": "hola1234",
  "Password": "contraseña1234"
}
```
``Resultado``: Token de autorización que se requerirá posteriormente para autenticarse en las demás rutas del sistema.


> [!WARNING]  
> A partir de aquí las rutas requieren pasar el **Token de autenticación** por medio del ``header Authorization`` (Authorization: Bearer <jwt-token>) o utilizando la ``cookie Authorization`` (Authorization=<jwt-token>)


- #### ``GET /validate``:
Para validar que el token funciona correctamente

``Resultado``: Información del usuario que ha iniciado sesión (la contraseña se guarda con una función hash).


- #### ``POST /registerFunction``:
Para registrar una función se debe primero subir una función al Docker hub que sea pública para que pueda ser extraida sencillamente por el servicio FaaS. A continuación, se puede registrar una función mediante el uso de esta ruta y en el body se debe poner un JSON con la siguiente estructura:
```json
{
    "Name": "cuenta_palabras",
    "Data": "bvalrod/cuenta_palabras:v1"
}
```
Donde ``Name`` es el nombre genérico que le quieres dar a la función a registrar y ``Data`` es la imagen de la función que quieres registrar.

``Resultado``: identificador de la función que se requiere para realizar acciones sobre esa función registrada.

- #### ``GET /getFunctions``:
Para ver todas las funciones que ha registrado el usuario.

``Resultado``: lista de funciones registradas por el usuario (con la estructura: ID, Name, Data, CreatedAt, UserID)

- #### ``DELETE /deleteFunction/:id``:
Para borrar una función por su ID (el id va en la ruta de la petición). 

``Resultado``: Mensaje que informa del estado final de la petición
>[!IMPORTANT]
> Un usuario solo puede borrar una función que haya creado él.

- #### ``POST /execute``:
Esta ruta permite al usuario ejecutar una función que haya sido anteriormente registrada. Requiere de un JSON con la siguiente estructura:
```json
{
    "FuncID": "e507855957272e8c91f2cb845698e04f",
    "Parameter": "Hola, Mundo!"
}
```
Donde ``FuncID`` es el identificador de la función, resultado del registro de una función en el sistema y ``Parameter`` es el parámetro (en string) que se le va a pasar a la función (la estructura del parámetro es controlada por el usuario que sube la función, si quiere más de un parámetro puede pasar un json dentro de este campo y tratar en la función la descodificación de json a objeto o como el usuario prefiera). 

``Resultado``: Resultado de la ejecución de la función y el tiempo que ha usado el servidor para ejecutar dicha función.

>[!IMPORTANT]
> Un usuario sólo puede ejecutar una función que ha registrado él mismo.

## Ejemplos de uso 💡 🔍
En la carpeta [/functions](./functions) se encuentran dos ejemplos de funciones que se pueden probar a registrar en el sistema FaaS. La primera función es ``cuenta_palabras`` en la versión v1 de esta imagen, se pasa un parámetro en formato string y devuelve el número de palabras que contiene el parámetro de la función.

### Primer paso: registrar la función
Tras iniciar sesión y poner el token de autorización en la petición, se registra la función lanzando una petición a ``POST /registerFunction`` con el siguiente JSON:
```json
{
    "Name": "cuenta_palabras",
    "Data": "bvalrod/cuenta_palabras:v1"
}
```

``resultado``:
```json
{
    "functionIdentifier": "e2c044c43d3eec8a6dcbcc0c2740cc29"
}
```

### Segundo paso: ejecutar la función
Ejecutamos la función con la ruta ``POST /execute`` con el siguiente JSON en el body de la consulta (importante poner el token de autorización en la consulta)
```json
{
    "FuncID": "e2c044c43d3eec8a6dcbcc0c2740cc29",
    "Parameter": "Hola, Mundo!"
}
```

``Resultado``:
```json
{
    "executionTime": "903.837001ms",
    "result": "2\n",
    "error": ""
}
```

### Funcion multiplica

Para la segunda función de prueba, multiplica, se requiere pasar por parámetros un JSON. Primero registraremos la función en la ruta ``POST /registerFunction`` con el siguiente JSON:
```json
{
    "Name": "cuenta_palabras",
    "Data": "bvalrod/cuenta_palabras:v2"
}
```
Y ahora conseguiremos el "functionIdentifier" para ejecutar la función a continuación. Lanzaremos una petición a ``POST /execute`` con el siguiente JSON en el body:
```json
{
    "FuncID": "5b200e85f063e31566d5e0e9ae8779f2",
    "Parameter": "{\"a\": 23, \"b\": 10}"
}
```

Esto nos dará como resultado, la multiplicación de a * b (en este caso: 230):
```json
{
    "executionTime": "930.710788ms",
    "result": "230\n",
    "error": ""
}
```


## Estructura del Proyecto 🚀

La estructura de directorios está organizada en 4 carpetas principales: apiServer, apisix_config, functions, worker

### API Server

```plaintext
apiServer
├── .env.example                  # Archivo de ejemplo de variables de entorno
├── .env                          # Variables de entorno
├── .gitignore                    # Archivos ignorados por git
├── Dockerfile                    # Imagen Docker del servicio API
├── go.mod                        # Módulo Go con dependencias y versiones
├── go.sum                        # Registro de dependencias
├── main.go                       # Punto de entrada del servicio API
│
├── controllers                   # Lógica de las rutas del API
│   ├── auth.go                   # Rutas de autenticación
│   ├── createRoutes.go           # Configuración de rutas
│   ├── executor.go               # Controlador de ejecución de funciones
│   ├── functionController.go     # Gestión de funciones
│   ├── middleware.go             # Autenticación con JWT
│   └── userControllers.go        # Gestión de usuarios
│
├── dataSource                    # Lógica de acceso a datos almacenados en NATS
│   ├── functionDataSource.go     # CRUD de funciones
│   └── userDataSource.go         # CRUD de usuarios
│
├── initializers                  # Configuración inicial del API
│   ├── loadEnvVariables.go       # Carga variables de entorno
│   └── nclient                   # Inicialización de NATS
│       ├── connectToNats.go      # Conexión al servidor NATS
│       ├── createFunctionKV.go   # Key-value store para funciones
│       ├── createJetStream.go    # Contexto JetStream
│       ├── createJobStream.go    # Stream de tareas
│       ├── createResponseStream.go # TODO: Stream de respuesta asíncrona
│       ├── createUserKV.go       # Key-value store para usuarios
│       └── subscribeToFunctions.go # Suscripción a funciones
│
├── models                        # Modelos de datos del API
│   ├── functionModel.go          # Modelo de funciones
│   ├── natsClient.go             # Modelo del cliente NATS
│   └── userModel.go              # Modelo de usuarios
│
├── repository                    # Conexión entre datos y lógica del servicio
│   ├── functionRepository.go     # Gestión de funciones
│   └── userRepository.go         # Gestión de usuarios
│
├── services                      # Lógica de servicios del API
│   ├── auth                      # Servicios de autenticación
│   │   └── authService.go        # Registro, login y validación
│   ├── executor                  # Servicios de ejecución
│   │   ├── errors.go             # Tipos de errores
│   │   ├── executorService.go    # Lógica de ejecución
│   │   └── types.go              # Modelos de datos
│   └── functions                 # Servicios de gestión de funciones
│       ├── errors.go             # Tipos de errores
│       ├── functionService.go    # Lógica de funciones
│       └── types.go              # Modelos de datos
│
└── utils                         # Funciones útiles
    └── randomString.go           # Generador de strings aleatorios
```

El servicio API permite la conexión entre el usuario y el sistema FaaS por medio de rutas configuradas para registrar usuarios, iniciar sesión, registrar funciones, 
ver funciones, ejecutar funciones y borrar funciones.

El API se levanta en un puerto 8080 aunque el usuario no utiliza ese puerto para conectarse a éste pues antes de llegar al API, se utiliza el ``reverse proxy`` de Apache Apisix
para gestionar el acceso de los usuarios al API.

### Apisix configuration

```plaintext
apisix_conf                       # Configuración de Apache Apisix
├── config.yaml                   # Configuración general del despliegue
└── routes.yaml                   # Configuración de rutas y plugins
```
Apisix funciona de tal manera que con archivos de configuración como son: [config.yaml](./apisix_conf/config.yaml) y [routes.yaml](./apisix_conf/routes.yaml) se levanta el servicio en modo ``standalone`` de forma sencilla con pocos pasos de configuración. El archivo ``config`` maneja el despliegue de apisix para que reciba peticiones al puerto ``9080`` y también registra los plugins que se van a utilizar en el Reverse proxy. En cambio, el archivo ``routes`` guarda las rutas que van a ser utilizadas por apisix para redirigir las peticiones al servicio API de backend. Las rutas hacen uso del plugin de JWT para autenticar los tokens generados por la ruta ``login`` y, además también utilizan el plugin ``limit-count`` para restringir las peticiones repetidas de los usuarios a un máximo de 2 peticiones cada 10 segundos, para evitar grandes cargas al servidor por usuarios que lanzan muchas peticiones seguidas.

### Functions

```plaintext
functions                         # Funciones para probar el FaaS
├── func1.cuentapalabras          # Función 1: cuenta palabras
│   ├── cuenta_palabras.py
│   └── Dockerfile
└── func2.multiplica              # Función 2: multiplica números
    ├── multiplica.py
    └── Dockerfile
```
En la carpeta de funciones existen algunas funciones para probar la funcionalidad que aporta el servicio FaaS.

### Worker

```plaintext
worker                            # Lógica de los workers del servicio
├── Dockerfile                    # Imagen Docker para los workers
├── go.mod                        
├── go.sum
├── main.go                       # Punto de entrada del worker
└── service                       # Funciones del worker
    └── workerService.go          # Funciones para iniciar y manejar workers
```
Esta carpeta [worker](./worker) es la lógica para ejecutar un worker que recibe peticiones del API del servicio FaaS para lanzar contenedores docker, recibir resultados de la ejecución de las funciones y devolver el resultado al API. Gracias a docker-compose podemos lanzar varias instancias de workers para ``escalar horizontalmente`` el servicio.

### Raíz del proyecto

```plaintext
Raíz del proyecto
├── docker-compose.yaml           # Organiza y escala los servicios con Docker
└── README.md                     # Documentación del proyecto
```
En la raíz del proyecto, además de las carpetas que separan las diferentes partes del mismo, encontramos el archivo [docker-compose](./docker-compose.yaml) que organiza el inicio de los servicios del FaaS en contenedores y permite replicar los workers para escalar horizontalmente. 

[Ahora yo tmb soy contributor hahaha]: #
