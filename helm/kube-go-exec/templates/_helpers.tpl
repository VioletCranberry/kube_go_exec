{{/*
Define chart name
*/}}
{{- define "kube-go-exec.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Define common labels
*/}}
{{- define "kube-go-exec.labels" -}}
helm.sh/chart: {{ include "kube-go-exec.chart" . }}
{{ include "kube-go-exec.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Define selector labels
*/}}
{{- define "kube-go-exec.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kube-go-exec.chart" . }}
app.kubernetes.io/app: {{ .Release.Name }}
{{- end }}

{{/*
Define service account namespace
*/}}
{{- define "kube-go-exec.serviceAccountNamespace" -}}
{{- if .Values.serviceAccount.create }}
{{- default .Values.serviceAccount.namespace }}
{{- else }}
{{- default "default" .Values.serviceAccount.namespace }}
{{- end }}
{{- end }}

{{/*
Define service account name
*/}}
{{- define "kube-go-exec.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "kube-go-exec.chart" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
