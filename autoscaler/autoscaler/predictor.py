"""Adaptive AR(p) predictor using Recursive Least Squares with safety checks."""
import numpy as np

class ARPredictor:
    def __init__(self, p: int = 2, forgetting: float = 0.99):
        self.p = p
        self.forgetting = forgetting
        self.phi = np.zeros(p)
        self.P = np.eye(p) * 1000.0
        self.history = []

    def update(self, new_value: float) -> None:
        if len(self.history) < self.p:
            self.history.append(new_value)
            return

        x = np.array(self.history[-self.p:][::-1])
        y_hat = float(self.phi @ x)
        error = new_value - y_hat

        Px = self.P @ x
        denominator = self.forgetting + x @ Px
        if denominator == 0:  # edge case
            return

        k = Px / denominator
        self.phi += k * error

        self.P = (self.P - np.outer(k, x) @ self.P) / self.forgetting
        self.history.append(new_value)

        # стабилизация коэффициентов (опционально)
        self.phi = np.clip(self.phi, -10, 10)

    def predict(self):
        if len(self.history) < self.p:
            return None
        x = np.array(self.history[-self.p:][::-1])
        y_pred = float(self.phi @ x)
        return max(y_pred, 0.0)  # защита от отрицательного предсказания