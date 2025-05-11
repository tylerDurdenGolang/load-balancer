from locust import LoadTestShape
class Step(LoadTestShape):
    stages = [
        {"duration":  30, "users": 200,  "spawn_rate": 200},
        {"duration":  35, "users": 900,  "spawn_rate": 1400},
        {"duration": 180, "users": 900,  "spawn_rate": 0},
    ]
    def tick(self):
        run = self.get_run_time()
        for s in self.stages:
            if run < s["duration"]:
                return (s["users"], s["spawn_rate"])
        return None
