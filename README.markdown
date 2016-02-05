## Kubernetes Zookeeper as a Contained Service

To run this locally without staging data:

1. `make`
1. `docker run -e MYID=1 -e <container_id>`

To run this locally while creating new folder paths:

1. `make`
1. `docker run -d -e MYID=1 -e LAB=1 -e BUCKET_1=Folder1 -e BUCKET_2=Folder2 <container_id>`

To run this within a Kubernetes cluster:

1. `make`, `docker tag <container_id>`, and `docker push` the Dockerfile.
1. Service `.metadata.name` field must always be `zookeeper-XX`, with `XX` being a two digit number.
1. Replication Controller `.spec.template.spec.containers.env` field must always contain `MYID=X`, with X being the server ID integer.

* If you are trying to stage data, make sure the `LAB` and `BUCKET_X` variables are added to the Replication Controllers.

This container exposes ports 2181, 2888, and 3888.

:exclamation: Note: for best performance use with an emtpyDir volume mount that uses Memory or tmpfs mount from the host.

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

The `zkcfg.go` program is what actually configures Zookeeper. If changes need to be made in regards to Zookeeper configurations, create a pull request and implement them here.