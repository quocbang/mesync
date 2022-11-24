# MES Sync

- [MES Sync](#mes-sync)
  - [Compilation](#compilation)
  - [Generate TLS Certificate and Key](#generate-tls-certificate-and-key)
  - [References](#references)
  - [Flags](#flags)
  - [Server Settings](#server-settings)
    - [db setting](#db-setting)
    - [cloud settings](#cloud-settings)
    - [uploadRetryParameters](#uploadretryparameters)
    - [factoryContainerInfo](#factorycontainerinfo)
    - [blobToken](#blobtoken)
    - [Server Settings Example](#server-settings-example)

MES provides the API to synchronize the ERP data or migrates the MES data to ERP.

## Compilation

You can find the available server configuration list by typing:

```console
 go run . -h
```

**In Order to run server using TLS handshake, both `tls-cert` and `tls-key` configurations must be assigned.**

**Use `--config-path` flag to specify the [Server Settings](#server-settings) file path.**

## Generate TLS Certificate and Key

In our case, we only need to generate server certificate and key in order to use TLS handshake.

- Generate Server Certificate

  ```console
  openssl genpkey -algorithm RSA -out server.key
  ```

  ```console
  openssl req -new -nodes -key server.key -out server.csr -subj "/C=TW/OU=IT/O=KENDA/CN=localhost" -config ./openssl.cnf -extensions v3_req
  ```

  ```console
  openssl x509 -req -days 3650 -in server.csr -out server.pem -CA ca.pem -CAkey ca.key -CAcreateserial -extfile ./openssl.cnf -extensions v3_req
  ```

## References

reference: [generate TLS Cert Tutorial](https://segmentfault.com/a/1190000038212054)

## Flags

- grpc-host : the IP to listen on
- grpc-port : this is an internal setting, you could use it when the default port is full on your system
- gateway-host : the IP to listen on
- gateway-port : gateway port to bind
- dev-mode : development mode
- auth-key : authentication key
- timeout : request timeout, need to fetch unit time
- config-path : server settings path
- tls-cert : path to TLS certificate (PUBLIC). To enable TLS handshake, you must set this value
- tls-key : path to TLS certificate key (PRIVATE), To enable TLS handshake, you must set this value
- root-ca : path to the root certificate

## Server Settings

| Field | Description |
| --- | --- |
| kenda | Specify Kenda MES settings. See [db setting](#db-setting) for details. |
| cloud | Specify Azure cloud settings. See [cloud settings](#cloud-settings) for details. |

### db setting

| Field    | Description                                                  |
| -------- | ------------------------------------------------------------ |
| schema   | Specify the database schema to be accessed.                  |
| name     | Specify the database name to be accessed.                    |
| address  | Specify the IP on which the server is to listen for connections from client applications. |
| port     | The TCP port the server listens on.                          |
| userName | User name for database connection.                           |
| password | User password for database connection.                       |

### cloud settings

| Field                                           | Description                                                  |
| ----------------------------------------------- | ------------------------------------------------------------ |
| [uploadRetryParameters](#uploadRetryParameters) | settings about failed tasks                                  |
| factoryContainerInfo                            | determine the container and account information related to the factory id |
| cloud                                           | cloud [db setting](#db setting)                              |

### uploadRetryParameters

| Field                | Description                                 |
| -------------------- | ------------------------------------------- |
| maxTimes             | limitation of retry times                   |
| interval             | retry frequency.  (minutes)                 |
| maxConcurrentJobs    | max number of async tasks                   |
| unsuccessStoragePath | a path to save failed tasks waited to retry |

### factoryContainerInfo

| Field         | Description          |
| ------------- | -------------------- |
| factoryID     |                      |
| containerName |                      |
| blobToken     | [details](blobToken) |

### blobToken

| Field       | Description                                                  |
| ----------- | ------------------------------------------------------------ |
| accountName | azure account name. [how to find it](https://stackoverflow.com/questions/6985921/where-can-i-find-my-azure-account-name-and-account-key). |
| sasToken    | shared access signature. [reference](https://docs.microsoft.com/en-us/azure/storage/common/storage-sas-overview) |

### Server Settings Example

```yaml
kenda:
  KY:
    schema: schemaName
    name: DBName
    address: DBAddress
    port: 5432
    userName: user
    password: password
  KU:
    schema: schemaName
    name: DBName
    address: DBAddress
    port: 5432
    userName: user
    password: password

cloud:
  cloud:
    schema: cloudSchema
    name: cloudDBName
    address: CloudDBAddress
    port: 5432
    userName: user
    password: password

  uploadRetryParameters:
    maxTimes: 5
    interval: 30
    maxConcurrentJobs: 5
    unsuccessStoragePath: ./upload_retry

  factoryContainerInfo:
    - factoryID: factory1
      containerName: container1
      blobToken:
        accountName: account1
        sasToken: sv=1999-09-09&ss=xx&srt=o&sp=xxxxxxxxx&se=2022-2-22T14:59:04Z&st=2022-05-31T06:59:04Z&sip=xx.xxx.xx.xx&spr=httpssig=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

    - factoryID: factory2
      containerName: container2
      blobToken:
        accountName: account2
        sasToken: sv=1999-09-09&ss=xx&srt=o&sp=xxxxxxxxx&se=2022-2-22T14:59:04Z&st=2022-05-31T06:59:04Z&sip=xx.xxx.xx.xx&spr=httpssig=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```
