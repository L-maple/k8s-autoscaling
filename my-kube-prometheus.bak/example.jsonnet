local mixin = import 'kube-prometheus/kube-prometheus-config-mixins.libsonnet';
local kp = 
(import 'kube-prometheus/kube-prometheus.libsonnet') +
(import 'kube-prometheus/kube-prometheus-kubeadm.libsonnet') +
// (import 'kube-prometheus/kube-prometheus-anti-affinity.libsonnet') +
  {
    _config+:: {
      namespace: 'monitoring',
      prometheus+:: {
        // 那些ns需要授权给到prometheus。
       namespaces+: ['default',"kube-system","monitoring"],
      },
    },
    // 这里替换成自己的私有仓库地址前缀
  } + mixin.withImageRepository('aliuchangjie');

{ ['00namespace-' + name]: kp.kubePrometheus[name] for name in std.objectFields(kp.kubePrometheus) } +
{ ['0prometheus-operator-' + name]: kp.prometheusOperator[name] for name in std.objectFields(kp.prometheusOperator) } +
{ ['node-exporter-' + name]: kp.nodeExporter[name] for name in std.objectFields(kp.nodeExporter) } +
{ ['kube-state-metrics-' + name]: kp.kubeStateMetrics[name] for name in std.objectFields(kp.kubeStateMetrics) } +
{ ['alertmanager-' + name]: kp.alertmanager[name] for name in std.objectFields(kp.alertmanager) } +
{ ['prometheus-' + name]: kp.prometheus[name] for name in std.objectFields(kp.prometheus) } +
{ ['prometheus-adapter-' + name]: kp.prometheusAdapter[name] for name in std.objectFields(kp.prometheusAdapter) } +
{ ['grafana-' + name]: kp.grafana[name] for name in std.objectFields(kp.grafana) }
