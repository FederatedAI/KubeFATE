# How to use KubeFATE to configure FATE polling

FATE's polling is divided into server side and client side.

**what's polling?**

In Federated Learning scenario, a party (or site) will send and receive data from / to other parties (or sites).

Normally, *Two-Way mode is recommended*. In this mode, a party (site) needs to listen to a public port and provide it to exchange or other direct connected parties (sites).

In some specific cases where a party (or site) is not allowed or not willing to provide a public port, Polling mode can be used. In Polling mode, the pattern of sending data is the same as Two-Way mode. But when receiving data, this specific party (or site) will be acting as a polling client actively polls (i.e. fetches) data from the nearest RollSite.

## server

If you want to be the server side of polling, you can configure it like this:

```bash
# cluster.yaml

...
rollsite: 
  ...
  polling:
    enabled: true
    type: server
    clientList:                 # polling client list
    - partID: 9999    
    - partID: 10000
    concurrency: 60
```



## client

If you want to be the client side of polling, you can configure it like this:

```bash
# cluster.yaml

...
rollsite: 
  ...
  polling:
    enabled: true
    type: server
    server:
      ip: 192.168.100.1           # polling server ip
      port: 9370                  # polling server port
```