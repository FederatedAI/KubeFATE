import sys
import time
from kubernetes import client, config

APP_NAME = "python"

def shutdown_flow(namespace, api, app):
    if type(app) == client.V1Deployment:
        update_deployment(namespace, api, app)
    else:
        update_sts(namespace, api, app)
    app = get_flow_app(namespace, api)
    if app.spec.replicas == 0:
        print("change the replicas to 0 successfully")
        # The default grace period is 30 seconds.
        time.sleep(30)
        return 0
    else:
        print("failed to change the replicas to 0")
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
