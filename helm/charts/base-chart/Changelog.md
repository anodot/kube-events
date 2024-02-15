# Change Log

## 1.1.12

**Release date:** 2024-01-12

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* chore: add clusterrole, clusterrolebinding templates

### Default value changes

```diff
diff --git a/helm-charts/base-chart/values.yaml b/helm-charts/base-chart/values.yaml
index c27dc66..41adc55 100644
--- a/helm-charts/base-chart/values.yaml
+++ b/helm-charts/base-chart/values.yaml
@@ -12,6 +12,24 @@ serviceAccount:
   # awsRole: "AwsRoleName"
   # awsAccount: "111222333444"
 
+clusterRole:
+  enabled: false
+  # name: role-name
+  # rules:
+  #   - apiGroups: [""]
+  #     resources: ["pods", "replicasets", "deployments", "daemonsets"]
+  #     verbs: ["get", "list", "patch", "watch"]
+  #   - apiGroups: ["apps"]
+  #     resources: ["replicasets", "deployments", "daemonsets"]
+  #     verbs: ["get", "list", "patch", "watch"]
+
+clusterRoleBinding:
+  enabled: false
+  # name: role-name-binding
+  # roleName: role-name
+  # serviceAccountName: serviceaccount-name
+  # namespace: serviceaccount-namespace
+
 deployment:
   labels:
     version: v1
```

## 1.1.11

**Release date:** 2024-01-05

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* feat: add configmap base-chart
* feat: add configmap base-chart

### Default value changes

```diff
diff --git a/helm-charts/base-chart/values.yaml b/helm-charts/base-chart/values.yaml
index fdcc0db..c27dc66 100644
--- a/helm-charts/base-chart/values.yaml
+++ b/helm-charts/base-chart/values.yaml
@@ -147,4 +147,15 @@ dbmigration:
   enabled: false
   image: 
     repository: 932213950603.dkr.ecr.us-east-1.amazonaws.com/recs-migration
-    pullPolicy: Always
\ No newline at end of file
+    pullPolicy: Always
+
+configmap:
+  enabled: false
+  mountPath: /path/to/config
+  data:
+    my-config.yaml: |
+      exampleKey: exampleValue
+  binaryData:
+    enabled: false
+    files:
+      my-binary-file.bin: <base64-encoded-data>
\ No newline at end of file
```

## 1.1.10

**Release date:** 2023-11-29

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* fix: custom port  names

### Default value changes

```diff
diff --git a/helm-charts/base-chart/values.yaml b/helm-charts/base-chart/values.yaml
index 3af3286..fdcc0db 100644
--- a/helm-charts/base-chart/values.yaml
+++ b/helm-charts/base-chart/values.yaml
@@ -84,7 +84,7 @@ podDisruptionBudget:
 
 service:
   serverPort: 8181
-  serverPortName: http-server
+  serverPortName: http
   managementPort: 8080
   managementPortName: http-management
   labels:
```

## 1.1.9

**Release date:** 2023-11-29

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* fix: add optinal svc targetPort

### Default value changes

```diff
# No changes in this release
```

## 1.1.8

**Release date:** 2023-10-27

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* feat: add base-chart dbmigration job
* chore: resources
* chore: remove cpu limit as default
* fix: fix ingress
* fix: fix ingress base chart
* fix: base-chart temp
* fix: remove test part base-chart

### Default value changes

```diff
diff --git a/helm-charts/base-chart/values.yaml b/helm-charts/base-chart/values.yaml
index 138455a..3af3286 100644
--- a/helm-charts/base-chart/values.yaml
+++ b/helm-charts/base-chart/values.yaml
@@ -95,10 +95,9 @@ service:
 
 ingress:
   enabled: false
-  annotations:
-    kubernetes.io/ingress.class: nginx
-  annotations_mgmt:
-    kubernetes.io/ingress.class: nginx
+  className: internal
+  # annotations:
+  #   kubernetes.io/ingress.class: nginx
   paths: /
   # Array of ingress hosts.
   hosts:
@@ -121,7 +120,7 @@ resources:
   # choice for the user. This also increases chances charts run on environments with little
   # resources, such as Minikube. If you do want to specify resources, uncomment the following
   # lines, adjust them as necessary, and remove the curly braces after 'resources:'.
-  limits:
+  limits: 
     cpu: 2000m
     # memory: 1500M
   requests:
@@ -143,3 +142,9 @@ affinity:
   # should be either
   # `kubernetes.io/hostname` or `failure-domain.beta.kubernetes.io/zone`
   topologyKey: "kubernetes.io/hostname"
+
+dbmigration:
+  enabled: false
+  image: 
+    repository: 932213950603.dkr.ecr.us-east-1.amazonaws.com/recs-migration
+    pullPolicy: Always
\ No newline at end of file
```

## 1.1.7

**Release date:** 2023-07-27

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* fix: java change default port

### Default value changes

```diff
diff --git a/helm-charts/base-chart/values.yaml b/helm-charts/base-chart/values.yaml
index d367d49..138455a 100644
--- a/helm-charts/base-chart/values.yaml
+++ b/helm-charts/base-chart/values.yaml
@@ -50,8 +50,6 @@ securityContext:
   # readOnlyRootFilesystem: false
 
 extraEnv:
-  # CONFIG_SERVER_HOST: config-server
-  # CONFIG_SERVER_PORT: 8888
   # SPRING_PROFILES_ACTIVE: dev
 
 #Configuration for java microservices
@@ -85,7 +83,7 @@ podDisruptionBudget:
   enabled: false
 
 service:
-  serverPort: 8888
+  serverPort: 8181
   serverPortName: http-server
   managementPort: 8080
   managementPortName: http-management
```

## 1.1.6

**Release date:** 2023-07-27

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* fix: base-chart add default probe and management port

### Default value changes

```diff
diff --git a/helm-charts/base-chart/values.yaml b/helm-charts/base-chart/values.yaml
index fae531b..d367d49 100644
--- a/helm-charts/base-chart/values.yaml
+++ b/helm-charts/base-chart/values.yaml
@@ -1,7 +1,7 @@
 replicaCount: 1
 
 image:
-  repository: 932213950603.dkr.ecr.us-east-1.amazonaws.com/config-client-service
+  repository: 932213950603.dkr.ecr.us-east-1.amazonaws.com/config-server
   tag: latest
   pullPolicy: Always
 
@@ -20,8 +20,22 @@ deployment:
   terminationGracePeriodSeconds: 30
   readinessProbe:
     enabled: false
+    path: '/actuator/health'
+    port: 8080
+    initialDelaySeconds: 10
+    periodSeconds: 15
+    timeoutSeconds: 5
+    successThreshold: 1
+    failureThreshold: 3
   livenessProbe:
     enabled: false
+    path: '/actuator/health'
+    port: 8080
+    initialDelaySeconds: 10
+    periodSeconds: 15
+    timeoutSeconds: 5
+    successThreshold: 1
+    failureThreshold: 3
   # lifecycleHooks: {}
   lifecycleHooks:
     preStop:
@@ -44,17 +58,25 @@ extraEnv:
 javaConfig:
   enabled: true
   # should be set with suffix 'M'
-  heapSize: "200M"
+  heapSize: "300M"
   # -- JAVA_OPTS env variable.  Ex. - "-Param1"
-  JAVA_OPTS: [ ]
-  #  - "-Param1"
+  JAVA_OPTS:
   #  - "-Param2"
+  #  - "-XX:+HeapDumpOnOutOfMemoryError"
+  #  - "-XX:HeapDumpPath=/heapdump/heapdump.bin"
 
   # -- JMX_OPTS env variable.  Ex. - "-DParam1"
   JMX_OPTS: [ ]
   #  - "-DParam1"
   #  - "-DParam2"
 
+sidecars: [ ]
+  # - name: sidecar
+  #   image: busybox
+  #   command: ["sh", "-c", "tail -f /dev/null"]
+  #   volumeMounts:
+  #   - name: heap-dump
+  #     mountPath: /heapdump
 
 externalSecret:
   enabled: false
@@ -63,10 +85,10 @@ podDisruptionBudget:
   enabled: false
 
 service:
-  serverPort: 8181
+  serverPort: 8888
   serverPortName: http-server
-  # managementPort: 8071
-  # managementPortName: http-management
+  managementPort: 8080
+  managementPortName: http-management
   labels:
     version: v1
   # annotations:
```

## 1.1.5

**Release date:** 2023-07-25

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* fix: base-chart split java part
* fix: add tolerations to base-chart

### Default value changes

```diff
diff --git a/helm-charts/base-chart/values.yaml b/helm-charts/base-chart/values.yaml
index 1e3836a..fae531b 100644
--- a/helm-charts/base-chart/values.yaml
+++ b/helm-charts/base-chart/values.yaml
@@ -36,14 +36,15 @@ securityContext:
   # readOnlyRootFilesystem: false
 
 extraEnv:
-  CONFIG_SERVER_HOST: config-server
-  CONFIG_SERVER_PORT: 8888
+  # CONFIG_SERVER_HOST: config-server
+  # CONFIG_SERVER_PORT: 8888
   # SPRING_PROFILES_ACTIVE: dev
 
 #Configuration for java microservices
-java:
+javaConfig:
+  enabled: true
   # should be set with suffix 'M'
-  heapSize: "10M"
+  heapSize: "200M"
   # -- JAVA_OPTS env variable.  Ex. - "-Param1"
   JAVA_OPTS: [ ]
   #  - "-Param1"
@@ -90,8 +91,8 @@ autoscaler:
   cpu:
     targetAverageUtilization: "80"
   memory:
-    targetAverageValue: "1000M"
-    # targetAverageUtilization: "90"
+    # targetAverageValue: "1000M"
+    targetAverageUtilization: "80"
 
 
 
```

## 1.1.4

**Release date:** 2023-07-18

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* fix: remove var SPRING_PROFILES_ACTIVE

### Default value changes

```diff
diff --git a/helm-charts/base-chart/values.yaml b/helm-charts/base-chart/values.yaml
index 12eccf1..1e3836a 100644
--- a/helm-charts/base-chart/values.yaml
+++ b/helm-charts/base-chart/values.yaml
@@ -38,7 +38,7 @@ securityContext:
 extraEnv:
   CONFIG_SERVER_HOST: config-server
   CONFIG_SERVER_PORT: 8888
-  SPRING_PROFILES_ACTIVE: dev
+  # SPRING_PROFILES_ACTIVE: dev
 
 #Configuration for java microservices
 java:
```

## 1.1.3

**Release date:** 2023-07-12

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* Bump version of base-chart to 1.1.3
* fix: base-chart deployment tweaks

### Default value changes

```diff
diff --git a/helm-charts/base-chart/values.yaml b/helm-charts/base-chart/values.yaml
index b7cb05f..12eccf1 100644
--- a/helm-charts/base-chart/values.yaml
+++ b/helm-charts/base-chart/values.yaml
@@ -9,6 +9,8 @@ nameOverride: ""
 
 serviceAccount:
   create: true
+  # awsRole: "AwsRoleName"
+  # awsAccount: "111222333444"
 
 deployment:
   labels:
```

## 1.1.2

**Release date:** 2023-07-06

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* Bump version of base-chart to 1.1.2
* add missing volume

### Default value changes

```diff
# No changes in this release
```

## 1.1.1

**Release date:** 2023-06-27

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* Bump version of base-chart to 1.1.1

### Default value changes

```diff
# No changes in this release
```

## 1.1.0

**Release date:** 2023-06-27

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* Bump version of base-chart to 1.1.0
* fix: add getJavaOptsLine trim
* fix test#2
* add tests
* resolve conflicts

### Default value changes

```diff
diff --git a/helm-charts/base-chart/values.yaml b/helm-charts/base-chart/values.yaml
index 3d50900..b7cb05f 100644
--- a/helm-charts/base-chart/values.yaml
+++ b/helm-charts/base-chart/values.yaml
@@ -42,8 +42,15 @@ extraEnv:
 java:
   # should be set with suffix 'M'
   heapSize: "10M"
-  # TODO: implement if needed. Just put here as example.
-  additionalJavaOpts: ""
+  # -- JAVA_OPTS env variable.  Ex. - "-Param1"
+  JAVA_OPTS: [ ]
+  #  - "-Param1"
+  #  - "-Param2"
+
+  # -- JMX_OPTS env variable.  Ex. - "-DParam1"
+  JMX_OPTS: [ ]
+  #  - "-DParam1"
+  #  - "-DParam2"
 
 
 externalSecret:
@@ -84,15 +91,7 @@ autoscaler:
     targetAverageValue: "1000M"
     # targetAverageUtilization: "90"
 
-# -- JAVA_OPTS env variable.  Ex. - "-Param1"
-JAVA_OPTS: []
-#  - "-Param1"
-#  - "-Param2"
 
-# -- JMX_OPTS env variable.  Ex. - "-DParam1"
-JMX_OPTS: []
-#  - "-DParam1"
-#  - "-DParam2"
 
 resources:
   # We usually recommend not to specify default resources and to leave this as a conscious
```

## 1.0.2

**Release date:** 2023-06-27

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* Merge remote-tracking branch 'origin/main'

### Default value changes

```diff
diff --git a/helm-charts/base-chart/values.yaml b/helm-charts/base-chart/values.yaml
index 296c4f1..3d50900 100644
--- a/helm-charts/base-chart/values.yaml
+++ b/helm-charts/base-chart/values.yaml
@@ -84,6 +84,16 @@ autoscaler:
     targetAverageValue: "1000M"
     # targetAverageUtilization: "90"
 
+# -- JAVA_OPTS env variable.  Ex. - "-Param1"
+JAVA_OPTS: []
+#  - "-Param1"
+#  - "-Param2"
+
+# -- JMX_OPTS env variable.  Ex. - "-DParam1"
+JMX_OPTS: []
+#  - "-DParam1"
+#  - "-DParam2"
+
 resources:
   # We usually recommend not to specify default resources and to leave this as a conscious
   # choice for the user. This also increases chances charts run on environments with little
@@ -91,10 +101,10 @@ resources:
   # lines, adjust them as necessary, and remove the curly braces after 'resources:'.
   limits:
     cpu: 2000m
-    memory: 1500M
+    # memory: 1500M
   requests:
     cpu: 50m
-    memory: 50M
+    # memory: 50M
 
 deploymentStrategy:
   rollingUpdate:
```

## 1.0.1

**Release date:** 2023-06-27

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* add logic for rolling upgrade strategy. Fix heap size calculation

### Default value changes

```diff
diff --git a/helm-charts/base-chart/values.yaml b/helm-charts/base-chart/values.yaml
index 3a781d4..296c4f1 100644
--- a/helm-charts/base-chart/values.yaml
+++ b/helm-charts/base-chart/values.yaml
@@ -24,13 +24,13 @@ deployment:
   lifecycleHooks:
     preStop:
       exec:
-        command: ['sleep', '10']
+        command: [ 'sleep', '10' ]
 
 securityContext:
   enabled: false
   # userId: 1000
   # groupId: 1337
-## disabled during java debug
+  ## disabled during java debug
   # readOnlyRootFilesystem: false
 
 extraEnv:
@@ -38,9 +38,13 @@ extraEnv:
   CONFIG_SERVER_PORT: 8888
   SPRING_PROFILES_ACTIVE: dev
 
-# Get heapdump
-javaDebug:
-  enabled: false
+#Configuration for java microservices
+java:
+  # should be set with suffix 'M'
+  heapSize: "10M"
+  # TODO: implement if needed. Just put here as example.
+  additionalJavaOpts: ""
+
 
 externalSecret:
   enabled: false
@@ -94,12 +98,16 @@ resources:
 
 deploymentStrategy:
   rollingUpdate:
-    maxSurge: 2
+    maxSurge: 0
     maxUnavailable: 1
   type: RollingUpdate
 
-nodeSelector: {}
+nodeSelector: { }
 
-tolerations: []
+tolerations: [ ]
 
-affinity: {}
+affinity:
+  enabled: true
+  # should be either
+  # `kubernetes.io/hostname` or `failure-domain.beta.kubernetes.io/zone`
+  topologyKey: "kubernetes.io/hostname"
```

## 1.0.2

**Release date:** 2023-06-27

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* Bump version of base-chart to 1.0.2
* fix: base chart add heap memory

### Default value changes

```diff
diff --git a/helm-charts/base-chart/values.yaml b/helm-charts/base-chart/values.yaml
index 3a781d4..a646d81 100644
--- a/helm-charts/base-chart/values.yaml
+++ b/helm-charts/base-chart/values.yaml
@@ -80,6 +80,16 @@ autoscaler:
     targetAverageValue: "1000M"
     # targetAverageUtilization: "90"
 
+# -- JAVA_OPTS env variable.  Ex. - "-Param1"
+JAVA_OPTS: []
+#  - "-Param1"
+#  - "-Param2"
+
+# -- JMX_OPTS env variable.  Ex. - "-DParam1"
+JMX_OPTS: []
+#  - "-DParam1"
+#  - "-DParam2"
+
 resources:
   # We usually recommend not to specify default resources and to leave this as a conscious
   # choice for the user. This also increases chances charts run on environments with little
@@ -87,10 +97,10 @@ resources:
   # lines, adjust them as necessary, and remove the curly braces after 'resources:'.
   limits:
     cpu: 2000m
-    memory: 1500M
+    # memory: 1500M
   requests:
     cpu: 50m
-    memory: 50M
+    # memory: 50M
 
 deploymentStrategy:
   rollingUpdate:
```

## 1.0.1

**Release date:** 2023-06-27

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* Bump version of base-chart to 1.0.1

### Default value changes

```diff
# No changes in this release
```

## 1.0.0

**Release date:** 2023-06-15

![AppVersion: 1.0.0](https://img.shields.io/static/v1?label=AppVersion&message=1.0.0&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* chore: add draft base-chart

### Default value changes

```diff
replicaCount: 1

image:
  repository: 932213950603.dkr.ecr.us-east-1.amazonaws.com/config-client-service
  tag: latest
  pullPolicy: Always

nameOverride: ""

serviceAccount:
  create: true

deployment:
  labels:
    version: v1
  # annotations:
  #   key: value
  terminationGracePeriodSeconds: 30
  readinessProbe:
    enabled: false
  livenessProbe:
    enabled: false
  # lifecycleHooks: {}
  lifecycleHooks:
    preStop:
      exec:
        command: ['sleep', '10']

securityContext:
  enabled: false
  # userId: 1000
  # groupId: 1337
## disabled during java debug
  # readOnlyRootFilesystem: false

extraEnv:
  CONFIG_SERVER_HOST: config-server
  CONFIG_SERVER_PORT: 8888
  SPRING_PROFILES_ACTIVE: dev

# Get heapdump
javaDebug:
  enabled: false

externalSecret:
  enabled: false

podDisruptionBudget:
  enabled: false

service:
  serverPort: 8181
  serverPortName: http-server
  # managementPort: 8071
  # managementPortName: http-management
  labels:
    version: v1
  # annotations:
  #   key: value
  type: ClusterIP

ingress:
  enabled: false
  annotations:
    kubernetes.io/ingress.class: nginx
  annotations_mgmt:
    kubernetes.io/ingress.class: nginx
  paths: /
  # Array of ingress hosts.
  hosts:
    - base-chart.example.com

autoscaler:
  enabled: false
  autoscaleMin: 1
  autoscaleMax: 5
  cpu:
    targetAverageUtilization: "80"
  memory:
    targetAverageValue: "1000M"
    # targetAverageUtilization: "90"

resources:
  # We usually recommend not to specify default resources and to leave this as a conscious
  # choice for the user. This also increases chances charts run on environments with little
  # resources, such as Minikube. If you do want to specify resources, uncomment the following
  # lines, adjust them as necessary, and remove the curly braces after 'resources:'.
  limits:
    cpu: 2000m
    memory: 1500M
  requests:
    cpu: 50m
    memory: 50M

deploymentStrategy:
  rollingUpdate:
    maxSurge: 2
    maxUnavailable: 1
  type: RollingUpdate

nodeSelector: {}

tolerations: []

affinity: {}
```

---
Autogenerated from Helm Chart and git history using [helm-changelog](https://github.com/mogensen/helm-changelog)
