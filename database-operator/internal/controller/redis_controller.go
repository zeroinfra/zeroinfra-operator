package controller

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	databasev1 "github.com/zeroinfra/database-operator/api/v1"
)

const redisFinalizer = "database.example.com/finalizer"

// RedisReconciler 是 Redis 自定义资源的控制器
type RedisReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile 方法处理 Redis 自定义资源的创建、更新和删除逻辑
func (r *RedisReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	redis := &databasev1.Redis{}
	if err := r.Get(ctx, req.NamespacedName, redis); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Redis resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Redis resource")
		return ctrl.Result{}, err
	}

	// 处理删除逻辑
	if !redis.ObjectMeta.DeletionTimestamp.IsZero() {
		if r.contains(redis.GetFinalizers(), redisFinalizer) {
			if err := r.cleanupResources(ctx, redis); err != nil {
				logger.Error(err, "Failed to clean up resources during deletion")
				return ctrl.Result{}, err
			}
			redis.SetFinalizers(r.remove(redis.GetFinalizers(), redisFinalizer))
			if err := r.Update(ctx, redis); err != nil {
				logger.Error(err, "Failed to remove finalizer from Redis resource")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// 添加 Finalizer
	if !r.contains(redis.GetFinalizers(), redisFinalizer) {
		redis.SetFinalizers(append(redis.GetFinalizers(), redisFinalizer))
		if err := r.Update(ctx, redis); err != nil {
			logger.Error(err, "Failed to add finalizer to Redis resource")
			return ctrl.Result{}, err
		}
	}

	// 使用 CreateOrUpdate 模式处理 Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      redis.Name,
			Namespace: redis.Namespace,
		},
	}
	op, err := controllerutil.CreateOrUpdate(ctx, r.Client, deployment, func() error {
		// 设置控制器引用
		controllerutil.SetControllerReference(redis, deployment, r.Scheme)

		// 生成环境变量
		envVars := []corev1.EnvVar{}
		reservedKeys := map[string]struct{}{
			"REDIS_PASSWORD": {},
		}
		for key, value := range redis.Spec.EnvVars {
			if _, exists := reservedKeys[key]; exists {
				logger.Info("Skipping reserved environment variable key", "key", key)
				continue
			}
			envVars = append(envVars, corev1.EnvVar{Name: key, Value: value})
		}

		// 配置 Resources 资源配置建议
		var resources corev1.ResourceRequirements
		defaultLimits := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("500m"),  // 默认 CPU 上限
			corev1.ResourceMemory: resource.MustParse("512Mi"), // 默认内存上限
		}
		defaultRequests := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("250m"),  // 默认 CPU 请求
			corev1.ResourceMemory: resource.MustParse("256Mi"), // 默认内存请求
		}
		if redis.Spec.Resources.Limits.CPU != "" && redis.Spec.Resources.Limits.Memory != "" {
			resources.Limits = corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(redis.Spec.Resources.Limits.CPU),
				corev1.ResourceMemory: resource.MustParse(redis.Spec.Resources.Limits.Memory),
			}
		} else {
			resources.Limits = defaultLimits
		}
		if redis.Spec.Resources.Requests.CPU != "" && redis.Spec.Resources.Requests.Memory != "" {
			resources.Requests = corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(redis.Spec.Resources.Requests.CPU),
				corev1.ResourceMemory: resource.MustParse(redis.Spec.Resources.Requests.Memory),
			}
		} else {
			resources.Requests = defaultRequests
		}

		// 配置 PodSpec
		containers := []corev1.Container{
			{
				Name:  "redis",
				Image: redis.Spec.Image,
				Ports: []corev1.ContainerPort{{
					ContainerPort: redis.Spec.Port,
				}},
				Env: append([]corev1.EnvVar{
					{
						Name:  "REDIS_PASSWORD",
						Value: redis.Spec.Password,
					},
				}, envVars...),
				Resources: resources,
			},
		}

		// 处理持久化卷
		if redis.Spec.PersistentVolume.Size != "" {
			volumeName := "data"
			containers[0].VolumeMounts = []corev1.VolumeMount{
				{
					Name:      volumeName,
					MountPath: "/data",
				},
			}
			deployment.Spec.Template.Spec.Volumes = []corev1.Volume{
				{
					Name: volumeName,
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: fmt.Sprintf("%s-pvc", redis.Name),
						},
					},
				},
			}
		} else {
			containers[0].VolumeMounts = nil
			deployment.Spec.Template.Spec.Volumes = nil
		}

		deployment.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: map[string]string{"app": redis.Name},
		}
		deployment.Spec.Template.ObjectMeta.Labels = map[string]string{"app": redis.Name}

		return nil
	})
	if err != nil {
		logger.Error(err, "Failed to reconcile Deployment")
		return ctrl.Result{}, err
	}
	if op != controllerutil.OperationResultNone {
		logger.Info("Deployment successfully reconciled", "operation", op)
	}

	// 使用 CreateOrUpdate 模式处理 Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      redis.Name,
			Namespace: redis.Namespace,
		},
	}
	op, err = controllerutil.CreateOrUpdate(ctx, r.Client, service, func() error {
		// 设置控制器引用
		controllerutil.SetControllerReference(redis, service, r.Scheme)

		service.Spec.Selector = map[string]string{"app": redis.Name}
		service.Spec.Ports = []corev1.ServicePort{
			{
				Port: redis.Spec.Port,
				Name: "redis",
			},
		}
		return nil
	})
	if err != nil {
		logger.Error(err, "Failed to reconcile Service")
		return ctrl.Result{}, err
	}
	if op != controllerutil.OperationResultNone {
		logger.Info("Service successfully reconciled", "operation", op)
	}

	// 处理 PVC
	if redis.Spec.PersistentVolume.Size != "" {
		pvcName := fmt.Sprintf("%s-pvc", redis.Name)
		pvc := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pvcName,
				Namespace: redis.Namespace,
			},
		}
		op, err = controllerutil.CreateOrUpdate(ctx, r.Client, pvc, func() error {
			// 设置控制器引用
			controllerutil.SetControllerReference(redis, pvc, r.Scheme)

			pvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
			pvc.Spec.Resources.Requests[corev1.ResourceStorage] = resource.MustParse(redis.Spec.PersistentVolume.Size)
			return nil
		})
		if err != nil {
			logger.Error(err, "Failed to reconcile PVC")
			return ctrl.Result{}, err
		}
		if op != controllerutil.OperationResultNone {
			logger.Info("PVC successfully reconciled", "operation", op)
		}
	} else {
		// 如果禁用了持久化存储，删除 PVC
		pvc := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-pvc", redis.Name),
				Namespace: redis.Namespace,
			},
		}
		if err := r.Get(ctx, types.NamespacedName{Name: pvc.Name, Namespace: redis.Namespace}, pvc); err == nil {
			if err := r.Delete(ctx, pvc); err != nil {
				logger.Error(err, "Failed to delete PVC when PersistentVolume is disabled")
				return ctrl.Result{}, err
			}
			logger.Info("Deleted PVC as PersistentVolume is disabled", "PVC.Name", pvc.Name)
		} else if !errors.IsNotFound(err) {
			logger.Error(err, "Failed to check PVC existence when PersistentVolume is disabled")
			return ctrl.Result{}, err
		}
	}

	// 更新状态
	if err := r.updateStatus(ctx, redis); err != nil {
		logger.Error(err, "Failed to update Redis status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// cleanupResources 清理 Redis 相关资源
func (r *RedisReconciler) cleanupResources(ctx context.Context, redis *databasev1.Redis) error {
	logger := log.FromContext(ctx)
	// 删除 Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      redis.Name,
			Namespace: redis.Namespace,
		},
	}
	if err := r.Delete(ctx, deployment); err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete Deployment: %w", err)
	}

	// 删除 Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      redis.Name,
			Namespace: redis.Namespace,
		},
	}
	if err := r.Delete(ctx, service); err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete Service: %w", err)
	}

	// 删除 PVC
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-pvc", redis.Name),
			Namespace: redis.Namespace,
		},
	}
	if err := r.Get(ctx, types.NamespacedName{Name: pvc.Name, Namespace: redis.Namespace}, pvc); err == nil {
		if err := r.Delete(ctx, pvc); err != nil {
			return fmt.Errorf("failed to delete PVC: %w", err)
		}
		logger.Info("Deleted PVC during cleanup", "PVC.Name", pvc.Name)
	} else if !errors.IsNotFound(err) {
		return fmt.Errorf("failed to get PVC during cleanup: %w", err)
	}

	return nil
}

// updateDeployment 更新 Deployment 配置，此处没有使用
func (r *RedisReconciler) updateDeployment(ctx context.Context, deployment *appsv1.Deployment, redis *databasev1.Redis) error {
	// 示例更新逻辑：检查镜像是否需要更新
	if deployment.Spec.Template.Spec.Containers[0].Image != redis.Spec.Image {
		deployment.Spec.Template.Spec.Containers[0].Image = redis.Spec.Image
		return r.Update(ctx, deployment)
	}
	return nil
}

// updateService 更新 Service 配置，此处没有使用
func (r *RedisReconciler) updateService(ctx context.Context, service *corev1.Service, redis *databasev1.Redis) error {
	// 示例更新逻辑：检查端口是否需要更新
	if service.Spec.Ports[0].Port != redis.Spec.Port {
		service.Spec.Ports[0].Port = redis.Spec.Port
		return r.Update(ctx, service)
	}
	return nil
}

// updateStatus 更新 Redis 的状态
func (r *RedisReconciler) updateStatus(ctx context.Context, redis *databasev1.Redis) error {
	status := databasev1.RedisStatus{
		Conditions: []databasev1.RedisCondition{
			{
				Type:               databasev1.RedisReady,
				Status:             corev1.ConditionTrue,
				LastTransitionTime: metav1.Now(),
				Reason:             "PodReady",
				Message:            "Redis Pod is ready",
			},
			{
				Type:               databasev1.RedisConfigured,
				Status:             corev1.ConditionTrue,
				LastTransitionTime: metav1.Now(),
				Reason:             "ConfigurationComplete",
				Message:            "Redis configuration completed",
			},
		},
	}

	if redis.Spec.PersistentVolume.Size != "" {
		status.Conditions = append(status.Conditions, databasev1.RedisCondition{
			Type:               databasev1.RedisPersistentVolumeBound,
			Status:             corev1.ConditionTrue,
			LastTransitionTime: metav1.Now(),
			Reason:             "RedisPersistentVolumeBound",
			Message:            "Persistent volume is bound successfully",
		})
	}

	if !reflect.DeepEqual(redis.Status, status) {
		redis.Status = status
		return r.Status().Update(ctx, redis)
	}
	return nil
}

// getPersistentVolumes 获取持久化卷配置
func (r *RedisReconciler) getPersistentVolumes(redis *databasev1.Redis) []corev1.Volume {
	if redis.Spec.PersistentVolume.Size != "" {
		return []corev1.Volume{{
			Name: "data",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: fmt.Sprintf("%s-pvc", redis.Name),
				},
			},
		}}
	}
	return nil
}

func (r *RedisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasev1.Redis{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.PersistentVolumeClaim{}). // 添加对 PVC 的拥有关系
		Complete(r)
}

// contains 检查切片中是否包含指定值
func (r *RedisReconciler) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// remove 从切片中移除指定值
func (r *RedisReconciler) remove(slice []string, item string) []string {
	result := slice[:0]
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}
