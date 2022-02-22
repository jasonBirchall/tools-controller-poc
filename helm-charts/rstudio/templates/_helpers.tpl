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

e.g. username-rstudio.example.com
*/}}
{{- define "host" -}}
{{- (printf "%s-rstudio.%s" .Values.username .Values.toolsDomain) | lower -}}
{{- end -}}


{{/*
unidle-key: used by unidler to find RStudio resources

will be the host if <= 63 characters (no change in old `alpha`
cluster) or the part before the first dot (".") otherwise
(to avoid problems with new, long, domain)

examples:
host="alice-rstudio.example.com" => returns "alice-rstudio.example.com"
host="alice-rstudio.verylongdomain[...]example.com" => returns "alice-rstudio"
*/}}
{{- define "unidle_key" -}}
    {{- $label := include "host" . -}}
    {{- if gt (len $label) 63 -}}
        {{- $parts := split "." $label -}}
        {{- $label = $parts._0 -}}
    {{- end -}}
    {{- $label | quote -}}
{{- end -}}
