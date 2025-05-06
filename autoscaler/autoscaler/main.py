
"""Entry‑point that ties everything together and runs forever inside a Pod."""
import time, logging, os
from .predictor import ARPredictor
from .scaler import AutoScaler, ScalingConfig
from .metrics import PromClient
from .config import load
from .k8s_client import patch_replicas

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(message)s")

def run():
    cfg = load(os.getenv("CONFIG", "config.yaml"))
    prom = PromClient(cfg.prometheus_url)
    predictor = ARPredictor(p=2, forgetting=0.99)
    scaler = AutoScaler(ScalingConfig(**cfg.scaling))

    expr = f'rate(http_requests_total{{deployment="{cfg.deployment}"}}[1m])'

    current_replicas = scaler.prev_replicas

    while True:
        rps = prom.query_instant(expr)
        if rps is not None:
            predictor.update(rps)
            pred = predictor.predict() or rps
            desired = scaler.optimise(pred, current_replicas)
            if desired != current_replicas:
                logging.info("Scaling %s/%s from %d → %d (predicted RPS %.1f)", cfg.namespace, cfg.deployment, current_replicas, desired, pred)
                patch_replicas(cfg.namespace, cfg.deployment, desired)
                current_replicas = desired
        time.sleep(cfg.interval)

if __name__ == "__main__":
    run()
