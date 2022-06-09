# fate-test use guide

```bash
# build
docker build --build-arg SOURCE_TAG=${version}-release -t  federatedai/fate-test:${version}-release .
# run
docker run -it federatedai/fate-test:${version}-release bash
```

## Edit config and check

```bash
fate_test config edit
```

Change  `path(FATE)` to `/data/projects/fate`.
Modify the partyID and corresponding address.

```yaml
      - {address: ${PartyAIP}:30097, parties: [9999]}
      - {address: ${PartyBIP}:30107, parties: [10000]}
```

```bash
fate_test config check
```

## Init flow and pipeline

fateflow_IP is the Guest Party IP

```bash
flow init --ip ${fateflow_IP} --port ${fateflow_Port}
pipeline init --ip ${fateflow_IP} --port ${fateflow_Port}
```

## Run test

```bash
fate_test suite -i test_suite/kubefate/base_testsuite.json
```

