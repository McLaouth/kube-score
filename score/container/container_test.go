package container

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/zegl/kube-score/scorecard"
)

func TestOkAllTheSameContainerResourceRequestsEqualLimits(t *testing.T) {
	t.Parallel()
	f := containerResourceRequestsEqualLimits(true)
	s := f(
		corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: "foo",
						Resources: corev1.ResourceRequirements{
							Requests: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("256Mi"),
							},
							Limits: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("256Mi"),
							},
						},
					},
				},
			},
		},
		metav1.TypeMeta{})

	assert.Equal(t, scorecard.GradeAllOK, s.Grade)
	assert.Len(t, s.Comments, 0)
}

func TestOkMultipleContainersContainerResourceRequestsEqualLimits(t *testing.T) {
	t.Parallel()
	f := containerResourceRequestsEqualLimits(true)
	s := f(
		corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				InitContainers: []corev1.Container{
					{
						Name: "foo",
						Resources: corev1.ResourceRequirements{
							Requests: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("256Mi"),
							},
							Limits: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("256Mi"),
							},
						},
					},
				},
				Containers: []corev1.Container{
					{
						Name: "foo",
						Resources: corev1.ResourceRequirements{
							Requests: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("256Mi"),
							},
							Limits: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("256Mi"),
							},
						},
					},
					{
						Name: "foo2",
						Resources: corev1.ResourceRequirements{
							Requests: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("256Mi"),
							},
							Limits: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("256Mi"),
							},
						},
					},
				},
			},
		},
		metav1.TypeMeta{})

	assert.Equal(t, scorecard.GradeAllOK, s.Grade)
	assert.Len(t, s.Comments, 0)
}

func TestOkSameQuantityContainerResourceRequestsEqualLimits(t *testing.T) {
	t.Parallel()
	f := containerResourceRequestsEqualLimits(true)
	s := f(
		corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: "foo",
						Resources: corev1.ResourceRequirements{
							Requests: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("256Mi"),
							},
							Limits: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1000m"),
								"memory": resource.MustParse("0.25Gi"),
							},
						},
					},
				},
			},
		},
		metav1.TypeMeta{})

	assert.Equal(t, scorecard.GradeAllOK, s.Grade)
	assert.Len(t, s.Comments, 0)
}

func TestFailBothContainerResourceRequestsEqualLimits(t *testing.T) {
	t.Parallel()
	f := containerResourceRequestsEqualLimits(true)
	s := f(
		corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: "foo",
						Resources: corev1.ResourceRequirements{
							Requests: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("256Mi"),
							},
							Limits: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("2"),
								"memory": resource.MustParse("512Mi"),
							},
						},
					},
				},
			},
		},
		metav1.TypeMeta{})

	assert.Equal(t, scorecard.GradeCritical, s.Grade)
	assert.Len(t, s.Comments, 2)
	assert.Equal(t, "foo", s.Comments[0].Path)
	assert.Equal(t, "CPU requests does not match limits", s.Comments[0].Summary)
	assert.Equal(t, "Having equal requests and limits is recommended to avoid resource DDOS of the node during spikes. Set resources.requests.cpu == resources.limits.cpu", s.Comments[0].Description)
	assert.Equal(t, "foo", s.Comments[1].Path)
	assert.Equal(t, "Memory requests does not match limits", s.Comments[1].Summary)
	assert.Equal(t, "Having equal requests and limits is recommended to avoid resource DDOS of the node during spikes. Set resources.requests.memory == resources.limits.memory", s.Comments[1].Description)
}

func TestFailCpuInitContainerResourceRequestsEqualLimits(t *testing.T) {
	t.Parallel()
	f := containerResourceRequestsEqualLimits(true)
	s := f(
		corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				InitContainers: []corev1.Container{
					{
						Name: "init",
						Resources: corev1.ResourceRequirements{
							Requests: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("256Mi"),
							},
							Limits: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("2"),
								"memory": resource.MustParse("256Mi"),
							},
						},
					},
				},
				Containers: []corev1.Container{
					{
						Name: "foo",
						Resources: corev1.ResourceRequirements{
							Requests: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("256Mi"),
							},
							Limits: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("256Mi"),
							},
						},
					},
				},
			},
		},
		metav1.TypeMeta{})

	assert.Equal(t, scorecard.GradeCritical, s.Grade)
	assert.Len(t, s.Comments, 1)
	assert.Equal(t, "init", s.Comments[0].Path)
	assert.Equal(t, "CPU requests does not match limits", s.Comments[0].Summary)
	assert.Equal(t, "Having equal requests and limits is recommended to avoid resource DDOS of the node during spikes. Set resources.requests.cpu == resources.limits.cpu", s.Comments[0].Description)
}

func TestOkIgnoreCpuContainerResourceRequestsEqualLimits(t *testing.T) {
	t.Parallel()
	f := containerResourceRequestsEqualLimits(false)
	s := f(
		corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				InitContainers: []corev1.Container{
					{
						Name: "init",
						Resources: corev1.ResourceRequirements{
							Requests: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("256Mi"),
							},
							Limits: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("2"),
								"memory": resource.MustParse("256Mi"),
							},
						},
					},
				},
				Containers: []corev1.Container{
					{
						Name: "foo",
						Resources: corev1.ResourceRequirements{
							Requests: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("256Mi"),
							},
							Limits: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("256Mi"),
							},
						},
					},
				},
			},
		},
		metav1.TypeMeta{})

	assert.Equal(t, scorecard.GradeAllOK, s.Grade)
	assert.Len(t, s.Comments, 0)
}

func TestFailIgnoreCpuContainerResourceRequestsEqualLimits(t *testing.T) {
	t.Parallel()
	f := containerResourceRequestsEqualLimits(false)
	s := f(
		corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: "foo",
						Resources: corev1.ResourceRequirements{
							Requests: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("256Mi"),
							},
							Limits: map[corev1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("2"),
								"memory": resource.MustParse("512Mi"),
							},
						},
					},
				},
			},
		},
		metav1.TypeMeta{})

	assert.Equal(t, scorecard.GradeCritical, s.Grade)
	assert.Len(t, s.Comments, 1)
	assert.Equal(t, "foo", s.Comments[0].Path)
	assert.Equal(t, "Memory requests does not match limits", s.Comments[0].Summary)
	assert.Equal(t, "Having equal requests and limits is recommended to avoid resource DDOS of the node during spikes. Set resources.requests.memory == resources.limits.memory", s.Comments[0].Description)
}
