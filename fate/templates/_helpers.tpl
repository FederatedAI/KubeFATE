{{/* Helm required labels */}}
{{- define "fate.labels" -}}
name: {{ .Values.partyName | quote  }}
partyId: {{ .Values.partyId | quote  }}
owner: kubefate
cluster: fate
heritage: {{ .Release.Service }}
release: {{ .Release.Name }}
chart: {{ .Chart.Name }}
{{- end -}}

{{/* matchLabels */}}
{{- define "fate.matchLabels" -}}
name: {{ .Values.partyName | quote  }}
partyId: {{ .Values.partyId | quote  }}
{{- end -}}
