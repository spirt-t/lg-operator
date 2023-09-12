package model

import (
	"time"

	coreV1 "k8s.io/api/core/v1"
)

// LoadGenerator model.
/*
  - Name - name of load-generator pod deployed in k8s;
  - ClusterIP - ClusterIP of deployed load-generator service; available only within the k8s cluster;
  - ExternalIP - ExternalIP of deployed load-generator service; accessible from outside the k8s cluster;
  - Port - port of load-generator service;
  - Status - k8s status of load-generator pod.
*/
type LoadGenerator struct {
	Name       string
	ClusterIP  string
	ExternalIP string
	Port       int32
	Status     coreV1.PodPhase
	CreatedAt  time.Time
}
