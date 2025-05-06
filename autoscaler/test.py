import requests
import pandas as pd

url = "http://prometheus.m22.oliwio-pp.rnd.mtt/api/v1/query_range"
params = {
    "query": 'sum(rate(process_cpu_seconds_total{app="number-api"}[1m]))',
    "start": "2025-04-06T08:00:00Z",
    "end": "2025-04-06T09:00:00Z",
    "step": "15"  # шаг в секундах
}

resp = requests.get(url, params=params)
print(resp)
results = resp.json()["data"]["result"]

# Преобразуем в DataFrame
if results:
    values = results[0]["values"]
    df = pd.DataFrame(values, columns=["timestamp", "value"])
    df["timestamp"] = pd.to_datetime(df["timestamp"], unit="s")
    df["value"] = df["value"].astype(float)
    df.to_csv("metrics.csv", index=False)
    print("✅ Сохранено в metrics.csv")
else:
    print("Нет данных по метрике")
