# FATE Algorithm and Computational Acceleration Selection

As a federated learning framework, FATE supports many algorithms. It is also the needs of many enterprises to select the corresponding algorithm and accelerator card according to the business. Currently, KubeFATE supports the selection of algorithms and accelerator cards.

Whether in docker-compose or k8s deployment, the following two parameters can be selected:

- `algorithm` Algorithm choice
- `device` Computing Device Selection

## Algorithm

The choice of algorithm consists of two options:

- `Basic`
    Basic is the default option, which includes dependencies related to the removal of nn (including homo_nn and hetero_nn) algorithms.
- `NN`
    NN contains all the dependencies required for nn to include (homo_nn and hetero_nn). ***NN can only be used when computing is Eggroll***
- `LLM`
    LLM expresses support for [FATE-LLM](https://github.com/FederatedAI/FATE-LLM).
- `ALL`
    ALL represents all algorithms, including basic NN and LLM.

## Device

Device selection consists of an option:

- `CPU`
    The CPU is a computing device that uses the CPU as a FATE computing device.
- `IPCL`
    The CPU is a computing device that uses the CPU as a FATE computing device.
- `GPU`
    The GPU is a computing device that uses the GPU as a FATE computing device.

## Support matrix

Various combinations currently supported by KubeFATE.
| Device \ Algorithm | Basic | NN | LLM | ALL |
|---|---|---|---|---|
| CPU | EGGROLL&Spark | EGGROLL | - | - |
| IPCL| EGGROLL&SPark | - | - | - |
| GPU | - | EGGROLL | - | EGGROLL |
