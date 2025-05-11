from locust import FastHttpUser, task, constant_throughput

PAYLOAD = {
    "expression": "sin(x)",
    "lower": 0, "upper": 6.28, "samples": 1000,
}

class Integrator(FastHttpUser):
    host = "localhost:8080"
    wait_time = constant_throughput(10)   # 10 задач/с на VU :contentReference[oaicite:4]{index=4}

    @task
    def integrate(self):
        self.client.post("/integrate", json=PAYLOAD)
