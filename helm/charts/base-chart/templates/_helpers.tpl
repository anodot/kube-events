{{/*
Expand the name of the chart.
*/}}
{{- define "base-chart.name" -}}
{{- default .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}


{{/*
Create a default fully qualified app name.
*/}}
{{- define "base-chart.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
 Setting memory limits for java applicaiton as HEAP_SIZE + 10%
*/}}
{{- define "getRequestMemory" -}}
{{- $mem := . -}}
{{- $mem = trimSuffix "M" $mem -}}
{{- $mem = div (mul (atoi $mem) 110) 100 -}}
{{- print $mem -}}
{{- end -}}

{{/*
 Setting memory limits for java applicaiton as HEAP_SIZE + 30%
*/}}
{{- define "getLimitMemory" -}}

{{- $mem := . -}}
{{- $mem = trimSuffix "M" $mem -}}
{{- $mem = div (mul (atoi $mem) 130) 100 -}}
{{- print $mem -}}
{{- end -}}

{{/*
Ensure values for heapSize are set in MB
*/}}
{{- define "checkHeapSizeFromat" -}}
    {{- if not (regexMatch "^[0-9]+[M]$" .) }}
      {{- fail (printf "heapSize format error: %s is not a valid heapSize option. Use only 'M' for Megabytes, example: heapSize: 100M" .) }}
  {{- end }}
{{- end }}

{{- define "getJavaOptsLine" -}}
{{- $heapSize := printf "-Xmx%s -Xms%s" .Values.javaConfig.heapSize .Values.javaConfig.heapSize  }}
{{- $res := trim (cat $heapSize (join " " .Values.javaConfig.JAVA_OPTS )) -}}
{{- print $res -}}
{{- end -}}


{{/*
 If # relicas is set to 1, we can not allow having downtime while doing rolloing uprade.
*/}}
{{- define "getMaxSurge" -}}
{{- $res := ternary 1 0  (eq (int .Values.replicaCount) 1)  -}}
{{- print $res -}}
{{- end -}}

{{- define "getMaxUnavailable" -}}
{{- $res := ternary 0 1 (eq (int .Values.replicaCount) 1)  -}}
{{- print $res -}}
{{- end -}}