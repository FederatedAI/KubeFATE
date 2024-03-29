
# Copyright 2019-2022 VMware, Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

{{/* Images Suffix */}}

{{- define "images.spark-worker.suffix" -}}
{{- if eq .Values.algorithm "NN" -}}
-nn
{{- end -}}
{{- if eq .Values.algorithm "ALL" -}}
-all
{{- end -}}
{{- if eq .Values.device "IPCL" -}}
-ipcl
{{- end -}}
{{- if eq .Values.device "GPU" -}}
-gpu
{{- end -}}
{{- end -}}
