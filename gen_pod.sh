#!/bin/bash

cat <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: yahagi-go
spec:
  containers:
  - image: ${YUUBARIGO_IMAGE:-yuubarigo}:$(git rev-parse --short HEAD)
    args:
    - -debug
    - -interval
    - "2"
    - -retry
    - "10"
    - -kcp
    - http://127.0.0.1:8081
    - -expedNotify
    - -tgToken
    - "${TGTOKEN}"
    - -tgUser
    - "${TGUSER}"
    name: yuubarigo
    ports:
    - name: yuubarigo
      containerPort: 8099
      protocol: TCP
  - image: ${KCP_IMAGE:-kccacheproxy}:v2.6.3
    name: kccacheproxy
    ports:
    - name: kcp
      containerPort: 8081
      protocol: TCP
    volumeMounts:
    - mountPath: /app/cache
      name: kcp
  hostNetwork: true
  volumes:
  - name: kcp
    hostPath:
      path: /regi/kcp
      type: Directory
EOF
