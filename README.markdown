## Kubernetes Zookeeper as a Service

To run this locally:

1. `make`
1. `docker run -e MYID=1 -e ZOOKEEPER_01_SERVICE_HOST=$localhostIP <container_id>`

To run this within a Kubernetes cluster:

1. `make`, `docker tag <container_id>`, and `docker push` the Dockerfile.
1. Service `.metadata.name` field must always be `zookeeper-XX`, with `XX` being a two digit number.
1. Replication Controller `.spec.template.spec.containers.env` field must always contain `MYID=X`, with X being the server ID integer.

This container exposes ports 2181, 2888, and 3888.

## Examples

Here is an example service configuration.

```json
{
  "apiVersion": "v1",
  "kind": "Service",
  "metadata": {
    "name": "zookeeper-01",
    "version": "3.5.0"
  },
  "spec": {
    "ports": [
      {
        "name": "client",
        "port": 2181
      },
      {
        "name": "followers",
        "port": 2888
      },
      {
        "name": "election",
        "port": 3888
      }
    ],
    "selector": {
      "app": "zookeeper",
      "server-id": "1"
    }
  }
}
```

Here is an example replication controller configuration:

```json
{
  "apiVersion": "v1",
  "kind": "ReplicationController",
  "metadata": {
    "name": "zookeeper-01",
    "version": "3.5.0"
  },
  "spec": {
    "replicas": 1,
    "template": {
      "metadata": {
        "labels": {
          "app": "zookeeper",
          "server-id": "1"
        }
      },
      "spec": {
        "containers": [
          {
            "env": [
              {
                "name": "MYID",
                "value": "1"
              }
            ],
            "image": "zookeeper:3.5.0",
            "name": "server",
            "ports": [
              {
                "containerPort": 2181
              },
              {
                "containerPort": 2888
              },
              {
                "containerPort": 3888
              }
            ]
          }
        ]
      }
    }
  }
}
```

### Improving

The `zkcfg.go` file contains the little script that reads the k8s environmental variables, and parses, and then writes them to the appropriate files. If any locations of the Zookeeper configurations are changed, make sure to update the location flags in `run.sh` for `zkcfg` so they are written to the right location.
