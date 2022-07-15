import sys
import time
from kubernetes import client, config
from enum import Enum

APP_NAME = "python"


class AppType(Enum):
    DEPLOYMENT = 1
    STS = 2


def shutdown_flow(namespace, api, app_type):
    if app_type == AppType.DEPLOYMENT:
        resp = update_deployment(namespace, api)
    else:
        resp = update_sts(namespace, api)
    if resp.status.replicas == 0:
        print("Shut down flow succeed")
    else:
        print("Shut down flow failed")
        exit(-1)


def check_flow_type(namespace, api):
    deployments = api.list_namespaced_deployment(namespace)
    for deployment in deployments.items:
        if deployment.metadata.name == APP_NAME:
            return AppType.DEPLOYMENT
    return AppType.STS


def update_deployment(namespace, api):
    # Update container image
    body = client.V1Deployment(
        spec=client.V1DeploymentSpec(
            replicas=0
        )
    )
    return api.patch_namespaced_deployment(
        name=APP_NAME, namespace=namespace, body=body
    )


def update_sts(namespace, api):
    body = client.V1StatefulSet(
        spec=client.V1StatefulSetSpec(
            replicas=0
        )
    )
    return api.patch_namespaced_stateful_set(
        name=APP_NAME, namespace=namespace, body=body
    )


def wait_flow_down(namespace, api):
    flow_down = False
    while True:
        if flow_down:
            return
        rss = api.list_namespaced_replica_set(namespace)
        for rs in rss.items:
            if rs.metadata.name.startwith(APP_NAME):
                if rs.status.ready_replicas == 0:
                    flow_down = True
                else:
                    print("flow is still up, will check 10 seconds later")
                    time.sleep(10)
                break
    return


if __name__ == '__main__':
    _, namespace = sys.argv
    config.load_kube_config()
    apps_v1 = client.AppsV1Api()
    app_type = check_flow_type(namespace, apps_v1)
    shutdown_flow(namespace, apps_v1, app_type)
    wait_flow_down(namespace, apps_v1)
