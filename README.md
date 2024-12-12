# Proyecto Go 

Este proyecto es una aplicación desarrollada en Go que utiliza gRPC para gestionar y comunicar datos sobre los campeones de League of Legends.

## Funcionalidad

El servicio GRPC permite realizar las siguientes acciones:
- **Consultar un campeón**: Obtener información específica de un campeón por su nombre.
- **Listar campeones**: Obtener una lista completa de campeones.
- **Agregar campeones**: Enviar una lista de nuevos campeones para ser añadidos.
- **Consultar campeones por tipo**: Obtener campeones que coincidan con un tipo específico.

Los datos de cada campeón tienen el siguiente formato:
```json
{
  "name": "Ziggs",
  "rol": "MidLine",
  "type": "Mage"
}
```

---

## Cómo ejecutar el proyecto localmente

### Prerrequisitos
1. Tener instalado [Go](https://go.dev/dl/).
2. Tener configurado el entorno de GRPC en tu máquina local (incluyendo `protoc`).

### Pasos para ejecutar:
1. Clona este repositorio:
   ```bash
   git clone <https://github.com/BEDECKTT/ProyectGo.git>
   ```
2. Navega al directorio del proyecto:
   ```bash
   cd <directorio del proyecto>
   ```
3. Genera los archivos GRPC a partir del archivo `.proto`:
   ```bash
   protoc --go_out=. --go-grpc_out=. proto/champions.proto
   ```
4. Ejecuta el servidor local:
   ```bash
   go run server.go
   ```
5. Accede al servicio utilizando un cliente GRPC o herramientas como `grpcurl`.

---

## Cómo interactuar con el servicio

El servicio GRPC permite interactuar utilizando un cliente GRPC. Se puede probar con herramientas como `grpcurl` o implementando un cliente en Go.

### Ejemplos de uso:

1. **Consultar información de un campeón:**
   ```bash
   grpcurl -d '{"name":"Ziggs"}' localhost:50051 champs.ChampService/GetChampInfo
   ```

2. **Listar todos los campeones:**
   ```bash
   grpcurl -d '{}' localhost:50051 champs.ChampService/GetChampList
   ```

3. **Agregar campeones:**
   ```bash
   grpcurl -d '{"name":"Ziggs", "rol":"MidLine", "type":"Mage"}' localhost:50051 champs.ChampService/AddChamps
   ```

4. **Consultar campeones por tipo:**
   ```bash
   grpcurl -d '{"type":"Mage"}' localhost:50051 champs.ChampService/GetChampsbyType
   ```

---

## Entrada en el portafolio

- **Nombre del Proyecto**: League of Legends Champions GRPC Service
- **Descripción**: Gestión de información sobre campeones del juego League of Legends utilizando un servicio GRPC.
- **Enlace al Repositorio**: [GitHub Repository](<https://github.com/BEDECKTT/ProyectGo.git>)
- **Archivo `.proto`**: [champions.proto](<https://drive.google.com/drive/folders/1x40RpVD1lX-1F3OucHP62gDxyOoZR6JQ?usp=sharing>)

---

## Archivo `.proto`
El archivo `.proto` para la definición del servicio se encuentra en la raíz del proyecto y tiene el siguiente contenido:

```proto
syntax = "proto3";

package champs;

option go_package =  "champslol-grcp/proto;champs";

service ChampService {

    rpc GetChampInfo(ChampRequest) returns (ChampResponse);
    rpc GetChampList( Empty ) returns ( stream ChampResponse );
    rpc AddChamps( stream NewChampRequest) returns ( AddChampResponse );
    rpc GetChampsbyType( stream ChampTypeRequest ) returns ( stream ChampResponse );
}

message ChampRequest{
    string name = 1;
}

message ChampResponse {
  string name = 1;
  string rol = 2;
  string type = 3;
}

message NewChampRequest {
  string name = 1;
  string rol = 2;
  string type = 3;
}

message AddChampResponse {
  int32 count = 1;
}

message Empty{}

message ChampTypeRequest {
    string type = 1;
}
```


