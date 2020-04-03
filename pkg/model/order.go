package model

import "time"

type Order struct {
	Id           int64     `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	PodName      string    `json:"podName,omitempty"`
	PodNamespace string    `json:"podNamespace,omitempty"`
	PodIp        string    `json:"podIp,omitempty"`
	Version      string    `json:"version,omitempty"`
	Service      string    `json:"service,omitempty"`
	Type         string    `json:"type,omitempty"`
	Time         time.Time `json:"time,omitempty"`
}

func (o Order) JavaClassName() string {
	return "com.k8s.Order"
}
