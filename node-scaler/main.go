package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    // appsv1 "k8s.io/api/apps/v1"
    "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    // Создаём клиент для Kubernetes API
    clientset, err := getKubeClient()
    if err != nil {
        log.Fatalf("Ошибка при создании клиента Kubernetes: %v", err)
    }

    // Имя Namespace и Deployment, который хотим скейлить
    namespace := "default"
    deploymentName := "my-crud-deployment"

    // Порог условной метрики (для примера)
    var highThreshold float64 = 0.8  // 80% загрузки
    var lowThreshold float64 = 0.2   // 20% загрузки

    // Период проверки
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    log.Println("Запускаем сервис динамического масштабирования...")

    for {
        select {
        case <-ticker.C:
            // 1. Получаем текущие метрики (упрощённо - заглушка).
            //    Реально можно подключиться к Prometheus, k8s Metrics API и т.п.
            currentLoad := getCurrentLoad()

            // 2. Получаем текущий Deployment, чтобы узнать число реплик
            deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.Background(), deploymentName, v1.GetOptions{})
            if err != nil {
                log.Printf("Не удалось получить Deployment %s: %v", deploymentName, err)
                continue
            }

            currentReplicas := *deployment.Spec.Replicas
            log.Printf("Текущая загрузка: %.2f, текущее число реплик: %d", currentLoad, currentReplicas)

            // 3. Логика скейлинга
            if currentLoad > highThreshold {
                // Увеличиваем число реплик
                newReplicas := currentReplicas + 1
                log.Printf("Слишком высокая нагрузка (%.2f). Масштабируем до %d реплик.", currentLoad, newReplicas)
                err = scaleDeployment(clientset, namespace, deploymentName, newReplicas)
                if err != nil {
                    log.Printf("Ошибка масштабирования Deployment: %v", err)
                }
            } else if currentLoad < lowThreshold && currentReplicas > 1 {
                // Уменьшаем число реплик (но не даём опуститься до 0)
                newReplicas := currentReplicas - 1
                if newReplicas < 1 {
                    newReplicas = 1
                }
                log.Printf("Низкая нагрузка (%.2f). Уменьшаем до %d реплик.", currentLoad, newReplicas)
                err = scaleDeployment(clientset, namespace, deploymentName, newReplicas)
                if err != nil {
                    log.Printf("Ошибка масштабирования Deployment: %v", err)
                }
            } else {
                log.Println("Уровень нагрузки в норме. Ничего не делаем.")
            }
        }
    }
}

// scaleDeployment изменяет число реплик у Deployment
func scaleDeployment(clientset *kubernetes.Clientset, namespace, name string, replicas int32) error {
    // Получаем текущий Deployment
    deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, v1.GetOptions{})
    if err != nil {
        return fmt.Errorf("не удалось получить Deployment: %w", err)
    }

    // Меняем поле spec.replicas
    deployment.Spec.Replicas = &replicas

    // Применяем изменения
    _, err = clientset.AppsV1().Deployments(namespace).Update(context.Background(), deployment, v1.UpdateOptions{})
    if err != nil {
        return fmt.Errorf("не удалось обновить Deployment: %w", err)
    }

    log.Printf("Deployment %s успешно масштабирован до %d реплик.", name, replicas)
    return nil
}

// getKubeClient возвращает клиент Kubernetes из окружения или из kubeconfig
func getKubeClient() (*kubernetes.Clientset, error) {
    // Попытаемся сначала загрузить in-cluster config (если сервис крутится внутри k8s)
    config, err := rest.InClusterConfig()
    if err != nil {
        // Если не удалось, пробуем из локального kubeconfig
        kubeconfig := os.Getenv("KUBECONFIG")
        if kubeconfig == "" {
            // Можно задать путь к kubeconfig по умолчанию, например, ~/.kube/config
            kubeconfig = clientcmd.RecommendedHomeFile
        }
        configFromFile, err2 := clientcmd.BuildConfigFromFlags("", kubeconfig)
        if err2 != nil {
            return nil, fmt.Errorf("не удалось загрузить ни in-cluster конфигурацию, ни kubeconfig файл: %v, %v", err, err2)
        }
        config = configFromFile
    }

    // Создаём клиент
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, fmt.Errorf("ошибка при создании clientset: %w", err)
    }

    return clientset, nil
}

// getCurrentLoad – заглушка, возвращающая случайную "загрузку" от 0 до 1
// В реальном проекте тут будет логика опроса Prometheus / k8s Metrics API / Redis и т.д.
func getCurrentLoad() float64 {
    // Допустим, случайное число имитирует загрузку
    return float64(time.Now().UnixNano()%100) / 100
}
