class ScalingOptimizer:
    def __init__(self, cost_per_node=10, risk_coeff=0.7):
        self.cost_per_node = cost_per_node
        self.risk_coeff = risk_coeff

    def calculate(self, current_nodes, forecast, thresholds):
        best = current_nodes
        min_loss = float('inf')

        for nodes in range(
                max(1, current_nodes - 2),
                min(thresholds['max_nodes'], current_nodes + 2) + 1
        ):
            capacity = nodes * thresholds['cpu_threshold']
            risk = sum(max(0, f - capacity)**2 for f in forecast)
            loss = nodes * self.cost_per_node + risk * self.risk_coeff
            if loss < min_loss:
                min_loss = loss
                best = nodes
        return best