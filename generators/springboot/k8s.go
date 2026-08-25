package springboot

type K8sConfig struct {
	pg projectGenerator
}

func NewK8sConfig(pg projectGenerator) *K8sConfig {
	return &K8sConfig{pg: pg}
}

func (k K8sConfig) generate(pc ProjectConfig) error {
	if !pc.K8sSupport {
		return nil
	}

	templateMap := map[string]string{
		"k8s/config.yaml.tmpl":  "k8s/config.yaml",
		"k8s/db.yaml.tmpl":      "k8s/db.yaml",
		"k8s/app.yaml.tmpl":     "k8s/app.yaml",
		"k8s/ingress.yaml.tmpl": "k8s/ingress.yaml",
	}
	if pc.EmailSupport {
		templateMap["k8s/mailpit.yaml.tmpl"] = "k8s/mailpit.yaml"
	}
	if pc.RedisCachingSupport {
		templateMap["k8s/redis.yaml.tmpl"] = "k8s/redis.yaml"
	}
	if pc.RabbitMQSupport {
		templateMap["k8s/rabbitmq.yaml.tmpl"] = "k8s/rabbitmq.yaml"
	}
	if pc.OpenTelemetrySupport {
		templateMap["k8s/otel.yaml.tmpl"] = "k8s/otel.yaml"
	}

	for tmpl, filePath := range templateMap {
		if err := k.pg.executeTemplate(pc, tmpl, filePath); err != nil {
			return err
		}
	}
	return nil
}
