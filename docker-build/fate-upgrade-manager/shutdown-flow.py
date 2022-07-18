import sys
import time
from kubernetes import client, config

APP_NAME = "python"


def shutdown_flow(namespace, api, app):
    if type(app) == client.V1Deployment:
        update_deployment(namespace, api, app)
    else:
        update_sts(namespace, api, app)
    for _ in range(60):
        if type(app) == client.V1Deployment:
            app = api.read_namespaced_deployment_status(APP_NAME, namespace)
        else:
            app = api.read_namespaced_stateful_set_status(APP_NAME, namespace)
        print("the ready replicas number is %s" % app.status.ready_replicas)
        if not app.status.ready_replicas:
            # The default grace time period is 30 seconds
            print("sleep for another 30 seconds, make sure the flow's pod is down")
            time.sleep(30)
            return 0
        else:
            print("wait for 10 seconds and will recheck")
            time.sleep(10)
    print("cannot shutdown the flow's pod")
    return -1


def get_flow_app(namespace, api):
    deployments = api.list_namespaced_deployment(namespace)
    stss = api.list_namespaced_stateful_set(namespace)
    for deployment in deployments.items:
        if deployment.metadata.name == APP_NAME:
            return deployment
    for sts in stss.items:
        if sts.metadata.name == APP_NAME:
            return sts
    return None


def update_deployment(namespace, api, app):
    # Update container image
    app.spec.replicas = 0
    api.patch_namespaced_deployment(
        name=APP_NAME, namespace=namespace, body=app
    )


def update_sts(namespace, api, app):
    app.spec.replicas = 0
    api.patch_namespaced_stateful_set(
        name=APP_NAME, namespace=namespace, body=app
    )


if __name__ == '__main__':
    _, namespace = sys.argv
    config.load_incluster_config()
    apps_v1 = client.AppsV1Api()
    app = get_flow_app(namespace, apps_v1)
    if not app:
        print("cannot find any deployment/sts whose name is 'python'")
        exit()
    return_code = shutdown_flow(namespace, apps_v1, app)
    exit(return_code)
