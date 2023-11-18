# Cards RESTFULL API

### Alcance
Desarrollar una RESTful API relacionada con las cartas de Pokémon que permite ser
consumida desde un Front-end o una App con las siguientes acciones:
+ Crear una carta.
+ Actualizar una carta.
+ Devolver una carta.
+ Devolver todas las cartas.
+ Borrar una carta.

### La carta debe contener la siguiente información:
+ Nombre
+ HP [^1]
+ ¿Es primera edición?
+ Expansion (Base Set, Jungle, Fossil, Base Set 2)  
    > [!WARNING] No hay validación que rechaze un valor fuera del listado.
+ Tipo (Agua, Fuego, Hierba, Eléctrico, etc.)  
    > [!WARNING] No hay validación que rechaze un valor fuera del listado.
+ Rareza (Común, No Común, Rara)  
    > [!WARNING] No hay validación que rechaze un valor fuera del listado.
+ Precio
+ Imagen de la carta * sólo se carga una ruta.
+ Fecha de creación de la carta  

[^1]: 1 Representa la salud de los Pokemons y siempre es un múltiplo de 10

### Notas
+ La solución se puede escribir en cualquier lenguaje, preferentemente Golang, PHP.
    > [!NOTE] Elegí el lenguage: Go
+ El motor de base de datos debe ser MySQL o DynamoDB
    > [!NOTE] Elegí el motor: MySQL
+ Sentirse libre de hacer cualquier suposición que se necesite acerca de los campos
requeridos, la forma de estructurar la data y las validaciones necesarias.
    > [!NOTE] Todos los campos son requeridos
+ La data debe ser persistente.
+ La solución deberá estar en un repositorio accesible para el colaborador de Agree que
haya solicitado el challenge.
    > [!NOTE] Repositorio: 
+ Una vez finalizado se debe enviar al colaborador de Agree un documento con las
suposiciones que se hicieron, qué es lo que se hizo y cómo ejecutar la solución.
+ Adjuntar documentación de los endpoints de la API (ej: Swagger).  
    > [!WARNING] Para ver la documentacion de los endpoints:
    1) Ir a "https://editor.swagger.io"
    2) File
    3) Impor file
    4) seleccionar el archivo "swagger.json"
    
    > [!CAUTION] La documentacion tiene 3 errores y no estan documentados, ni el paginado, ni los filtros :(

### Bonus
+ Autenticación [❌]

+ Filtros [✅] agregar el query string en la URL --> ?filter=key.value --> ejemplo: ?filter=Name.pokemon_2
+ Paginado [✅]
> [!WARNING]Ejemplo de estructura del query string usando paginado y filtros:  
```
localhost:3000/api/v1/cards?page_number=1&page_size=10&filter=Name.pokemon_2
```

+ Tests Unitarios [✅]
> [!NOTE] Sólo hay un test de un médto del servicio. Elegí el servicio porque donde estimo que debería estar la lógica de la app.


+ Deploy de la API en un servicio Cloud (AWS, Azure, Google Cloud, etc.) [✅]
+ Utilización de los servicios de AWS (API Gateway, EC2, ECS, EKS, S3, etc.) [✅]
+ Serverless [✅]
> [!WARNING]El deploy funciona pero tengo un error de configuracion en la lambda que no logre solucionar a tiempo.  
Al hacer un request al endpoint para chechear la salud tira error:  
"GET - https://w0v0g2nwea.execute-api.us-east-1.amazonaws.com/prod/app/v1/ping"  
En cloudwatch veo el log que dice: "fork/exec /var/task/main: no such file or directory: PathError null"


## Aclaraciones
+ Hay un archivo Makefile para automatizar algunos comandos.  
    Al escribir en la terminal:
    ```
    $make help
    ```
    Verán que lista y explica los comandos.
+ Por una cuestión de simplicidad los ids se solicitan mediante query string.
+ Dentro de la carpeta "scripts" hay archivo sql para crear la base de datos y una tabla.  
    ```
    cards_mysql.sql
    ```