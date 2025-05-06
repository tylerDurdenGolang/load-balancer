
import requests, time, logging

class PromClient:
    """Minimal Prometheus HTTP API client."""
    def __init__(self, base_url: str):
        self.base = base_url.rstrip("/")

    def query(self, expr: str):
        r = requests.get(f"{self.base}/api/v1/query", params={"query": expr})
        r.raise_for_status()
        return r.json()['data']['result']

    def query_instant(self, expr: str):
        res = self.query(expr)
        if not res:
            return None
        return float(res[0]['value'][1])
