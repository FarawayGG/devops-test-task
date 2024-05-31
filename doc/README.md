# Answer

1. [Dockerfile](../docker/Dockerfile)
2. [docker-compose](../docker/docker-compose.yml)
3. VictoriaMetrics or Prometheus (+ Thanos), Alertmanager, Grafana, метрики с Redis с помощью [redis-exporter](https://github.com/oliver006/redis_exporter), доступность эндпоинта приложения с помощью [blackbox-exporter](https://github.com/prometheus/blackbox_exporter)
4. [Helm chart](../helm/app/)
```
helm install app ./helm/app/
```
5. Не совсем понял суть вопроса, но предположу, что должно быть всё, что нужно для минимальной работоспособности и мониторинга кластера: для метрик и алертов VictoriaMetrics or Prometheus (+ Thanos), Alertmanager, Grafana; для логов ELK/EFK/EVK при желании Splunk; на счет трасировок для такого сервиса не уверен, ну пускай будет Jaeger; Ingress Controller; Cert-Manager. Но на мой взгляд это всё может быть опционально, да и решений полно.
