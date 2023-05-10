
# KubeFATE deploys GPU-enabled FATE

## Prerequisites

- KubeFATE v1.11.1+
- Kubenertes support GPU (<https://kubernetes.io/docs/tasks/manage-gpus/scheduling-gpus/>)

## How to configure

The algorithm and device in cluster.yaml must be changed to NN and GPU. (Currently FATE only supports the use of GPU for NN algorithms)

```yaml
algorithm: NN
device: GPU
```

Then the resource of the python pod is allocated at least 1 GPU resource. (GPU computing is only in the pod of fateflow)

```bash
python:
  resources:
    requests:
      nvidia.com/gpu: 1
    limits:
      nvidia.com/gpu: 1
```

Here is an example [cluster-gpu.yaml](../k8s-deploy/examples/party-9999/cluster-gpu.yaml).

Then deploy the cluster defined by cluster.yaml, and you can use FATE to run GPU tasks.
