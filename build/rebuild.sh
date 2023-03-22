#!/bin/bash
docker build -f build/Dockerfile_cloud . -t main
kubectl delete -f build/main.yaml 
kubectl create -f build/main.yaml 