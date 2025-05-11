"""Entry‑point that ties everything together and runs forever inside a Pod."""
import time
import logging
import os
import threading
from http.server import HTTPServer, BaseHTTPRequestHandler

from .predictor import ARPredictor
from .scaler import AutoScaler, ScalingConfig
from .metrics import PromClient
from .config import load
from .k8s_client import patch_replicas

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(message)s")


class ReadinessHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/ready":
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b"OK")
        else:
            self.send_response(404)
            self.end_headers()

    def log_message(self, format, *args):
        # отключаем спам от /ready
        if self.path != "/ready":
            super().log_message(format, *args)


def start_readiness_server():
    server = HTTPServer(("", 8081), ReadinessHandler)
    logging.info("Readiness probe available at http://0.0.0.0:8081/ready")
    server.serve_forever()


def run():
    cfg = load(os.getenv("CONFIG", "config.yaml"))
    prom = PromClient(cfg.prometheus_url)
    predictor = ARPredictor(p=2, forgetting=0.99)
    scaler = AutoScaler(ScalingConfig(**cfg.scaling))

    expr_total = 'rate(mc_requests_total[1m])'
    expr_by_instance = 'sum(rate(mc_requests_total[1m])) by (instance)'

    current_replicas = scaler.prev_replicas

    logging.info("Autoscaler started for deployment: %s (namespace: %s)", cfg.deployment, cfg.namespace)
    logging.info("Prometheus URL: %s", cfg.prometheus_url)
    logging.info("Polling interval: %d seconds", cfg.interval)

    threading.Thread(target=start_readiness_server, daemon=True).start()

    while True:
        rps = prom.query_instant(expr_total)
        rps_by_instance = prom.query_instant(expr_by_instance)

        if isinstance(rps_by_instance, dict):
            logging.info("RPS per instance:")
            for instance, val in rps_by_instance.items():
                logging.info("  → %s: %.1f", instance, val)
                logging.info("Total RPS: %.1f", sum(rps_by_instance.values()))
        elif isinstance(rps_by_instance, float):
            logging.info("RPS (single result): %.1f", rps_by_instance)

        if rps is not None:
            predictor.update(rps)
            pred = predictor.predict() or rps
            desired = scaler.optimise(pred, current_replicas)
            if desired != current_replicas:
                logging.info(
                    "Scaling %s/%s from %d → %d (predicted RPS %.1f)",
                    cfg.namespace, cfg.deployment, current_replicas, desired, pred,
                )
                patch_replicas(cfg.namespace, cfg.deployment, desired)
                current_replicas = desired

        time.sleep(cfg.interval)


if __name__ == "__main__":
    run()
