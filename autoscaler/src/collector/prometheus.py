import requests
import yaml

class PrometheusClient:
    def __init__(self, url="http://prometheus:9090"):
        self.url = url
        with open("config/settings.yaml") as f:
            self.config = yaml.safe_load(f)['metrics']

    def get_metrics(self):
        metrics = {}
        for metric in self.config:
            try:
                response = requests.get(
                    f"{self.url}/api/v1/query",
                    params={'query': metric['query']}
                )
                data = response.json()['data']['result']
                metrics[metric['name']] = self._parse(data)
            except Exception as e:
                print(f"Error fetching {metric['name']}: {str(e)}")
        return metrics

    def _parse(self, data):
        return {item['metric']['instance']: float(item['value'][1]) for item in data}