#!/bin/bash

gomplate --file ./cluster.yaml.tmpl --out ./cluster.yaml

kind create cluster --config ./cluster.yaml
