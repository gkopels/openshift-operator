#!/bin/bash

IMAGE_REGISTRY=${N3000_IMAGE_REGISTRY:-registry.connect.redhat.com/intel}
IMAGE_VER=${N3000_IMAGE_VERSION:-v1.0.0}
OPAE_VER=${N3000_OPAE_VERSION:-1.3.8-2}


export N3000_DAEMON_IMAGE=${N3000_DAEMON_IMAGE:-${IMAGE_REGISTRY}/n3000-daemon:${IMAGE_VER}--${OPAE_VER}}
export N3000_LABELER_IMAGE=${N3000_LABELER_IMAGE:-${IMAGE_REGISTRY}/n3000-labeler:${IMAGE_VER}}
export N3000_DRIVER_CONTAINER_IMAGE=${N3000_DRIVER_CONTAINER_IMAGE:-${IMAGE_REGISTRY}/n3000-driver:${IMAGE_VER}--${OPAE_VER}}
export N3000_MONITORING_IMAGE=${N3000_MONITORING_IMAGE:-${IMAGE_REGISTRY}/n3000-monitoring:${IMAGE_VER}--${OPAE_VER}}
