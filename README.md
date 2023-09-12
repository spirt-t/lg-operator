Load-Generator Operator
-


A service designed to create and manage load generators in k8s.  
_Swagger_ documentation is available along the way `pkg/lg-operator/lg-operator.swagger.json`.

[[_TOC_]]

## Minimum requirements
To use this tool, you must have at least:
  - Kubernetes cluster,
  - the ability to set up service account *lg-operator* (at least access to k8s with [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)).

## Basic commands

- `make install` : install go dependencies,
- `make test` : run tests,
- `make cover` : display test coverage results,
- `make run` : run the service,
- `make go-generate` : regenerate mocks,
- `make generate` : update the service according to API.

## Configuration setting  

You can configure the startup and operation of the service in the file `config.yaml`.   

<details>
<summary>Config structure :point_down: </summary>

```yaml
service:
   ports:
      http: 7000
      grpc: 7002
   log:
      level: INFO

kubernetes:
   service:
      host: '172.16.128.1'
      port: 443
   namespace: default
   timeouts:
      create: '2m'
      delete: '30s'
   generator:
      port: 8888
      label: load-generator

default_resources:
   cpu:
      request: 1
      limit: 2
   memory:
      request: 1Gi
      limit: 2Gi

cleaning:
   outdated:
      ttl: '24h'
      enabled: true
   completed:
      interval: '5m'
      enabled: true
```
</details>

Field Details:
- *service* section defines parameters for starting the service:
  - *service.ports.http* : port on which the http-server will run
  - *service.ports.grpc* : port on which the grpc-server will run
- *kubernetes* section determines the parameters for connecting to the kubernetes api and generator start parameters:
  - *kubernetes.service* section determines the parameters for connecting to the kubernetes api (depend on your cluster settings):
    - *kubernetes.service.host* - k8s API host
    - *kubernetes.service.port* - k8s API port
  - *kubernetes.namespace* - namespace where your load generators will run
  - *kubernetes.timeouts* section defines timeouts for operations with generators:
    - *kubernetes.timeouts.create* - load generator creation timeout
    - *kubernetes.timeouts.delete* - load generator deletion timeout
  - *kubernetes.generator* sections defines parameters for load generators deployment:
    - *kubernetes.generator.port* - the port on which the generator will run
    - *kubernetes.generator.label* - label to be added to all generator k8s-entities
- *default_resources* defines default resources for load generator if not specified in the request to create  
- *cleaning* sets autovacuum options:
  - *cleaning.outdated* section sets parameters for deleting old generators:
    - *cleaning.outdated.enabled* - enable removal of old generators
    - *cleaning.outdated.ttl* - the maximum lifetime of the generator, after which it will be deleted regardless of the status
  - *cleaning.completed* section sets parameters for deleting completed generators:
    - *cleaning.completed.enabled* - enable removal of completed generators
    - *cleaning.completed.interval* - frequency of deleting completed generators.


## Key features

This service will help you:  
- launch load generators on a kubernetes cluster;
- remove previously launched generators;
- get a list of running generators;
- configure auto-cleanup of load generators after they complete.

You can also easily add or change the functionality of the service in accordance with your needs.

### Creating load generators

You can launch a list of generators by `POST /v1/generators`.  
The service starts pod, service and ingress for each generator.   

[![](https://mermaid.ink/img/pako:eNptkUtrQyEQhf-KzLbefXERKLSL7gqlOzeDTqzc-KiPQAn575kbb-AmrQtR55wzn3oCkyyBgko_naKhV4-uYNBR8EDTUhFflcrYZyzNG58xNnFwU8pUkBV_i_NzFS8f76Mw5iVFTLvd09aphE1m5oIP6EguGRjqMGxkw7eGKmEKYSORk5WCU4_esNVHV6iu3ptiy5vQTo7iFnlN5PiH8q3Hg4wxpnt8H_fpH9yrcLkxK7IUEQMT1oatV5AQqAT0ll_9tHg1tG8KpEHx0mKZNeh4Zh32lj5_owHVSicJPVtmWn_o_vDNem4Mao-HSucLGUacCw?type=png)](https://mermaid.live/edit#pako:eNptkUtrQyEQhf-KzLbefXERKLSL7gqlOzeDTqzc-KiPQAn575kbb-AmrQtR55wzn3oCkyyBgko_naKhV4-uYNBR8EDTUhFflcrYZyzNG58xNnFwU8pUkBV_i_NzFS8f76Mw5iVFTLvd09aphE1m5oIP6EguGRjqMGxkw7eGKmEKYSORk5WCU4_esNVHV6iu3ptiy5vQTo7iFnlN5PiH8q3Hg4wxpnt8H_fpH9yrcLkxK7IUEQMT1oatV5AQqAT0ll_9tHg1tG8KpEHx0mKZNeh4Zh32lj5_owHVSicJPVtmWn_o_vDNem4Mao-HSucLGUacCw)

<details>
<summary>The request must include the following information :point_down: </summary>

```json
{
   "parameters": [
      {
         "image": "string",
         "resources": {
            "memory": {
               "limit": "string",
               "request": "string"
            },
            "cpu": {
               "limit": "string",
               "request": "string"
            }
         },
         "additional_envs": [
            {
               "name": "string",
               "val": "string"
            }
         ],
         "commands": [
            "string"
         ]
      }
   ]
}
```
Field Details:
- *image* : container image name. More info: https://kubernetes.io/docs/concepts/containers/images .
- *resources* : compute Resources required by this container (in our case only memory and cpu). More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ .
- *additional_envs* : list of environment variables to set in the container.
- *commands* : entrypoint array. Not executed within a shell. More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell .

</details>

**Restrictions**
- The service starts load generators only in its k8s cluster.
- Generators run on a fixed port defined in config.
- Generators will not restart in case of a deployment error.

The method will return you the parameters of the launched generators or error.

<details>
<summary>Generator parameters :point_down: </summary>

```json
{
   "load_generators": [
      {
         "name": "string",
         "cluster_ip": "string",
         "external_ip": "string",
         "port": 0,
         "status": "string"
      }
   ]
}
```
Field Details:
- *name* : unique name of generator pod and service;
- *cluster_ip* : *ClusterIP* of deployed load-generator service; available only within the k8s cluster;
- *external_ip* : *ExternalIP* of deployed load-generator service; accessible from outside the k8s cluster;
- *port* : port of generator pod and service; you can use it in *http*-requests with *cluster_ip* or *external_ip*;
- *status* : k8s status of generator pod.

</details>

You now have access to the generator within and outside the k8s cluster! 

### Getting a list of generators
You can find out the parameters of currently running generators (`GET /v1/generators`).  
This method does not require input parameters.
The return value will contain the parameters of the currently running generators described above.

### Deleting load generators
After the generators finished, it is recommended to remove them from the cluster (`DELETE /v1/generators`).
To do this, you must specify a list of generator names that you want to remove.

[![](https://mermaid.ink/img/pako:eNptkcFOwzAMhl_F8pVGXFEOk5DgwA0JcevFSkyJ2ibBSRHTtHcnXdrRjeUQJfb3_7blA5pgGTUm_prYG35y1AmNrYdyyOQg8J5Y6j-SZGdcJJ9h6FSILFSI_8n-IcHj68sNVSCrOvZbYb3nKqB2u7uts4YzC55GrugGqIqlnAbLA2eGGGwDxe_bGW7A-U44paq1nLKE_c1GFhtQP1fp1fiEXbGlAXXZcujvWWQ1vWj2BM-T_lHY4MgykrNlC4dZ02L-5DIr6vK0JH2LrT8WjqYc3vbeoM4ycYNTtJTXjaH-oCGdo8_WlYpL8PgL63aiTQ?type=png)](https://mermaid.live/edit#pako:eNptkcFOwzAMhl_F8pVGXFEOk5DgwA0JcevFSkyJ2ibBSRHTtHcnXdrRjeUQJfb3_7blA5pgGTUm_prYG35y1AmNrYdyyOQg8J5Y6j-SZGdcJJ9h6FSILFSI_8n-IcHj68sNVSCrOvZbYb3nKqB2u7uts4YzC55GrugGqIqlnAbLA2eGGGwDxe_bGW7A-U44paq1nLKE_c1GFhtQP1fp1fiEXbGlAXXZcujvWWQ1vWj2BM-T_lHY4MgykrNlC4dZ02L-5DIr6vK0JH2LrT8WjqYc3vbeoM4ycYNTtJTXjaH-oCGdo8_WlYpL8PgL63aiTQ)

## How to make changes  

To change the service API, you need to:
1. make the appropriate changes to the file `api/lg-operator/lg-operator.proto`;
2. install [Protocol Compiler](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation) (protoc);
3. install dependencies:  
     ```go
      go install google.golang.org/protobuf/cmd/protoc-gen-go
      go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
      go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
      go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
      go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
      ```
4. run the following command:
    ```
    make generate
    ```

