# About Helm Chart and KubeFATE
KubeFATE is designed to deploy and manage clusters based on Helm Chart. In KubeFATE,
1. receives the customized setting through `cluster.yaml` or `cluster-serving.yaml`;
2. render `values.yaml` in Chart Template with `Go Template`. 

* Note: The content of `values.yaml` will persistent into MySQL in the KubeFATE service. If you going to do some hacks on the above process, make sure you understand this. 

In another word, the Helm Chart of KubeFATE is a twofold chart, which is specified for FATE and FATE-Serving.

# Prerequisite
For developing KubeFATE charts, we need to understand:
1. Go template: https://golang.org/pkg/text/template/
2. Helm Chart Developer's Guide: https://helm.sh/docs/chart_template_guide/

They are both complicated and large but take it easy. There are only small feature sets of them using in KubeFATE's chart. I suggest you go through the Quick Tutorials of them, then get your hand dirty. Look up the usages of something you met problems.

# KubeFATE's Chart Structure
All charts using in KubeFATE are located at https://github.com/FederatedAI/KubeFATE/tree/gh-pages/package. And you can find the developed one from https://github.com/FederatedAI/KubeFATE/tree/master/helm-charts. We suggested you develop your version based on one chart in https://github.com/FederatedAI/KubeFATE/tree/gh-pages/package, which are fine verified. 

Unzip one KubeFATE's Chart, you can find a `templates` folder and 4 files"
1. `Chart.yaml`: a YAML file containing information about the chart;
2. `value.yaml`: the default values for chart according to Helm Chart standard;
3. `values-template-example.yaml`: the example files of what values-template will look like;
4. `value-template.yaml`: core file for developed a customized KubeFATE's chart. It bridges the `cluster.yaml`/`cluser-serving.yaml` and templates. The values set in `cluster.yaml`/`cluser-serving.yaml` will passed to here, and set all the variables using in the templates. It follows the `Go Template` standard. 

* Note: `value.yaml` and the rendered `value-template.yaml` will be merged as the "VALUEs" to the chart templates.

## `templates` folder
In `templates` folder, the template yaml file combined with values will generate valid Kubernetes manifest files for each `FATE` or `FATE-Serving` component.

e.g. For `FATE` v1.4.2, there are following templates locating in `template` folder:
1. clustermanager-module.yaml: cluster-manager module
2. eggroll-config.yaml: eggroll module
3. mysql-module.yaml: MySQL module including the creating SQL
4. nodemanager.yaml: node-manager module
5. rollsite-module.yaml: roll-site module
6. python-module.yaml: python, fate-flow, fate-board, notebook/client module and related python libs.
we supported both ingress and istio gateway for HTTP traffic,
7. ingress.yaml: ingress for fate-board and notebook/client
8. istio.yaml: istio compatible service gateway for fate-board and notebook/client. But this feature is **experimental**.
and in some KubeFATE's charts, there is an optional NOTE.txt file locating in `template` folder to describe short usage notes.

All the config of the `FATE` and `FATE-Serving` are setting as ConfigMap in each template yaml. If you are going to change the default config of the components, also find the corresponding template yaml.
**Note: the template files just how we construct the resources, not how the pod or service looks like in Kubernetes.**

# Build KubeFATE's chart
We provides a `Makefile` in the repo, but it is very straghtforward to call the helm command:
```
release: lint
	helm package ./FATE
	helm package ./FATE-Serving
lint:
	helm lint ./FATE
	helm lint ./FATE-Serving
```
You can contribute your own chart and call `helm lint` (https://helm.sh/docs/helm/helm_lint/) and `helm package` (https://helm.sh/docs/helm/helm_package/) as well.
