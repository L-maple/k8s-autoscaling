local mixin = import 'kube-prometheus/kube-prometheus-config-mixins.libsonnet';
local kp = (import 'kube-prometheus/kube-prometheus.libsonnet') +
		(import 'kube-prometheus/kube-prometheus-kubeadm.libsonnet') + 
		(import 'kube-prometheus/kube-prometheus-all-namespaces.libsonnet') + 
		(import 'kube-prometheus/kube-prometheus-custom-metrics.libsonnet') + {
		_config+:: {
			namespace: 'monitoring',

			prometheus+:: {
				namespaces: [],
			},
		},
	}  + mixin.withImageRepository('aliuchangjie');

{ ['setup/0namespace-' + name]: kp.kubePrometheus[name] for name in std.objectFields(kp.kubePrometheus) } +
{
  ['setup/prometheus-operator-' + name]: kp.prometheusOperator[name]
  for name in std.filter((function(name) name != 'serviceMonitor'), std.objectFields(kp.prometheusOperator))
} +
// serviceMonitor is separated so that it can be created after the CRDs are ready
{ 'prometheus-operator-serviceMonitor': kp.prometheusOperator.serviceMonitor } +
{ ['node-exporter-' + name]: kp.nodeExporter[name] for name in std.objectFields(kp.nodeExporter) } +
{ ['kube-state-metrics-' + name]: kp.kubeStateMetrics[name] for name in std.objectFields(kp.kubeStateMetrics) } +
{ ['alertmanager-' + name]: kp.alertmanager[name] for name in std.objectFields(kp.alertmanager) } +
{ ['prometheus-' + name]: kp.prometheus[name] for name in std.objectFields(kp.prometheus) } +
{ ['prometheus-adapter-' + name]: kp.prometheusAdapter[name] for name in std.objectFields(kp.prometheusAdapter) } +
{ ['grafana-' + name]: kp.grafana[name] for name in std.objectFields(kp.grafana) }


