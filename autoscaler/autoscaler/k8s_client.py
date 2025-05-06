
from kubernetes import client, config

def patch_replicas(namespace: str, deployment: str, replicas: int):
    config.load_incluster_config()
    apps = client.AppsV1Api()
    body = {"spec": {"replicas": replicas}}
    apps.patch_namespaced_deployment_scale(name=deployment, namespace=namespace, body=body)
