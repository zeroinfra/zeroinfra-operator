package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RedisSpec 定义了 Redis 的期望状态
type RedisSpec struct {
	// +kubebuilder:validation:Required
	// Redis 镜像
	Image string `json:"image"`

	// Redis 端口，默认为 6379
	Port int32 `json:"port,omitempty"`

	// 密码字符串（直接存储密码）
	Password string `json:"password,omitempty"`

	// 资源限制
	Resources Resources `json:"resources,omitempty"`

	// 环境变量注入（key-value 格式）
	EnvVars map[string]string `json:"envVars,omitempty"`

	// 持久化存储配置（可选）
	PersistentVolume PersistentVolumeSpec `json:"persistentVolume,omitempty"`
}

// Resources 定义了 Redis 的资源限制
type Resources struct {
	Limits   ResourceLimit `json:"limits,omitempty"`
	Requests ResourceLimit `json:"requests,omitempty"`
}

// ResourceLimit 定义了 CPU 和内存的限制
type ResourceLimit struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// PersistentVolumeSpec 定义了持久化存储的配置
type PersistentVolumeSpec struct {
	// 存储大小
	Size string `json:"size,omitempty"`
	// 存储类
	StorageClass string `json:"storageClass,omitempty"`
}

// RedisConditionType 定义了 Redis 的状态类型（枚举）
type RedisConditionType string

const (
	// RedisReady 表示 Redis 实例是否就绪
	RedisReady RedisConditionType = "RedisReady"

	// RedisConfigured 表示 Redis 配置是否完成
	RedisConfigured RedisConditionType = "RedisConfigured"

	// RedisPersistentVolumeBound 表示持久化卷是否绑定成功
	RedisPersistentVolumeBound RedisConditionType = "RedisPersistentVolumeBound"
)

// RedisCondition 定义了 Redis 的条件状态
type RedisCondition struct {
	Type               RedisConditionType     `json:"type"`
	Status             corev1.ConditionStatus `json:"status"`
	LastTransitionTime metav1.Time            `json:"lastTransitionTime,omitempty"`
	Reason             string                 `json:"reason,omitempty"`
	Message            string                 `json:"message,omitempty"`
}

// RedisStatus 定义了 Redis 的观测状态
type RedisStatus struct {
	// Pod 名称
	PodName string `json:"podName,omitempty"`

	// 条件列表
	Conditions []RedisCondition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="RedisReady",type="string",JSONPath=".status.conditions[?(@.type==\"RedisReady\")].status"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Redis 是自定义资源对象
type Redis struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RedisSpec   `json:"spec,omitempty"`
	Status RedisStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RedisList 是 Redis 资源的列表
type RedisList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Redis `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Redis{}, &RedisList{})
}
