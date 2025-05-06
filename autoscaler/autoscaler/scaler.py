"""Decision engine that converts predicted load to replica count."""
from dataclasses import dataclass


@dataclass
class ScalingConfig:
    target_rps_per_pod: float = 100.0
    hysteresis: int = 1
    alpha: float = 0.6  # latency weight
    beta: float = 0.3  # cost weight
    gamma: float = 0.1  # churn weight
    max_replicas: int = 20
    min_replicas: int = 1


class AutoScaler:
    """Simple cost‑based autoscaling policy."""

    def __init__(self, cfg: ScalingConfig):
        self.cfg = cfg
        self.prev_replicas = cfg.min_replicas

    def optimise(self, predicted_rps: float, current_replicas: int):
        best_cost = float("inf")
        best_rep = current_replicas

        for r in range(self.cfg.min_replicas, self.cfg.max_replicas + 1):
            rps_per = predicted_rps / r if r else predicted_rps
            latency = max(0.0, (rps_per - self.cfg.target_rps_per_pod) / self.cfg.target_rps_per_pod)
            cost = r
            churn = abs(r - self.prev_replicas)
            total = self.cfg.alpha * latency + self.cfg.beta * cost + self.cfg.gamma * churn
            if total < best_cost:
                best_cost, best_rep = total, r

        # apply hysteresis
        if abs(best_rep - current_replicas) < self.cfg.hysteresis:
            best_rep = current_replicas

        self.prev_replicas = best_rep
        return best_rep
