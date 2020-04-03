#!/bin/bash

../proxy/bazel-bin/src/envoy/envoy -c ./manifests/envoy-org-config.yaml --drain-time-s 1 -l debug --concurrency 1 --base-id 40 --parent-shutdown-time-s 1 --restart-epoch 0