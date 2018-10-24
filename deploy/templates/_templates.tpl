{{/* vim: set filetype=mustache: */}}

{{/*
  Return the correct name of the component expect no argument
*/}}
{{- define "component.name" }}
{{- printf "%s-%s" .Release.Name .Chart.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
  Return the name of another component, expect a single string argument
*/}}
{{- define "compute.name"}}
{{- printf "%s-%s" .Release.Name . | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
  Return the fqdn of another component, expect the following object

  object:
    name: somename
    namespace: namespace
*/}}
{{- define "compute.cluster_fqdn" }}
{{- printf "%s-%s.%s.svc.cluster.local" .Release.Name .name .namespace | trunc 63 | trimSuffix "-" -}}
{{- end }}

{{/*
  Return the base 64 representation of a string
*/}}
{{- define "compute.b64" }}
{{- printf "%s" . | b64enc }}
{{- end }}
