{{/* Helm required labels */}}
{{- define "fate.labels" -}}
name: {{ .Values.partyName | quote  }}
partyId: {{ .Values.partyId | quote  }}
owner: kubefate
cluster: fate-exchange
heritage: {{ .Release.Service }}
release: {{ .Release.Name }}
chart: {{ .Chart.Name }}
{{- end -}}

{{/* matchLabels */}}
{{- define "fate.matchLabels" -}}
name: {{ .Values.partyName | quote  }}
partyId: {{ .Values.partyId | quote  }}
{{- end -}}

{{/*
Create the name of the controller service account to use
*/}}
{{- define "serviceAccountName" -}}
{{- if .Values.podSecurityPolicy.enabled -}}
    {{ default .Values.partyName }}
{{- else -}}
    {{ default "default" }}
{{- end -}}
{{- end -}}