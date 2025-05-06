
"""Adaptive AR(p) predictor using Recursive Least Squares."""
import numpy as np

class ARPredictor:
    def __init__(self, p: int = 2, forgetting: float = 0.99):
        self.p = p
        self.forgetting = forgetting
        self.phi = np.zeros(p)
        self.P = np.eye(p) * 1_000.0
        self.history = []

    # online update with new observation
    def update(self, new_value: float) -> None:
        if len(self.history) < self.p:
            self.history.append(new_value)
            return

        x = np.array(self.history[-self.p:][::-1])
        y_hat = float(self.phi @ x)
        error = new_value - y_hat

        Px = self.P @ x
        k = Px / (self.forgetting + x @ Px)

        self.phi += k * error
        self.P = (self.P - np.outer(k, x) @ self.P) / self.forgetting
        self.history.append(new_value)

    # one‑step forecast
    def predict(self):
        if len(self.history) < self.p:
            return None
        x = np.array(self.history[-self.p:][::-1])
        return float(self.phi @ x)
