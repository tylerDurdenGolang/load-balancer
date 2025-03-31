class ETS:
    def __init__(self, alpha=0.3, beta=0.1):
        self.alpha = alpha
        self.beta = beta
        self.trend = 0

    def fit(self, data):
        self.level = data[0]
        self.trend = (data[1] - data[0]) / 1
        for i in range(2, len(data)):
            prev_level = self.level
            self.level = self.alpha * data[i] + (1 - self.alpha) * (self.level + self.trend)
            self.trend = self.beta * (self.level - prev_level) + (1 - self.beta) * self.trend

    def predict(self, steps):
        return [self.level + (i+1)*self.trend for i in range(steps)]

class AutoRegressor:
    def __init__(self, lag=3):
        self.lag = lag
        self.weights = [1.0/lag] * lag

    def fit(self, data):
        for _ in range(100):
            for i in range(self.lag, len(data)):
                pred = sum(w * data[i-j-1] for j, w in enumerate(self.weights))
                error = data[i] - pred
                self.weights = [w + 0.01*error*data[i-j-1] for j, w in enumerate(self.weights)]

    def predict(self, last_values, steps):
        return [sum(w * v for w, v in zip(self.weights, last_values[-self.lag:]))
                for _ in range(steps)]