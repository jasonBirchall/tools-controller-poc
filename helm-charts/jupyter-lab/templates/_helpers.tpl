{{/*
Create a default fully qualified app name.
We truncate at 24 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 24 | trimSuffix "-" -}}
{{- end -}}


{{/*
Host

e.g. username-jupyter-lab.example.com
*/}}
{{- define "host" -}}
{{- (printf "%s-jupyter-lab.%s" .Values.Username .Values.toolsDomain) | lower -}}
{{- end -}}


{{/*
unidle-key: used by unidler to find RStudio resources

will be the host if <= 63 characters (no change in old `alpha`
cluster) or the part before the first dot (".") otherwise
(to avoid problems with new, long, domain)

examples:
host="alice-jupyter-lab.example.com" => returns "alice-jupyter-lab.example.com"
host="alice-jupyter-lab.verylongdomain[...]example.com" => returns "alice-jupyter-lab"
*/}}
{{- define "unidle_key" -}}
    {{- $label := include "host" . -}}
    {{- if gt (len $label) 63 -}}
        {{- $parts := split "." $label -}}
        {{- $label = $parts._0 -}}
    {{- end -}}
    {{- $label | quote -}}
{{- end -}}
