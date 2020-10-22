package score

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
	"testing"

	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
)

func TestPodSecurityContext(test *testing.T) {
	test.Parallel()

	b := func(b bool) *bool { return &b }
	i := func(i int64) *int64 { return &i }

	tests := []struct {
		ctx             *corev1.SecurityContext
		podCtx          *corev1.PodSecurityContext
		expectedGrade   scorecard.Grade
		expectedComment *scorecard.TestScoreComment
	}{
		// No security context set
		{
			ctx:           nil,
			expectedGrade: scorecard.GradeCritical,
			expectedComment: &scorecard.TestScoreComment{
				Path:        "foobar",
				Summary:     "Container has no configured security context",
				Description: "Set securityContext to run the container in a more secure context.",
			},
		},
		// All required variables set correctly
		{
			ctx: &corev1.SecurityContext{
				ReadOnlyRootFilesystem: b(true),
				RunAsGroup:             i(23000),
				RunAsUser:              i(33000),
				RunAsNonRoot:           b(true),
				Privileged:             b(false),
			},
			expectedGrade: scorecard.GradeAllOK,
		},
		// Read only file system is explicitly false
		{
			ctx: &corev1.SecurityContext{
				ReadOnlyRootFilesystem: b(false),
			},
			expectedGrade: scorecard.GradeCritical,
			expectedComment: &scorecard.TestScoreComment{
				Path:        "foobar",
				Summary:     "The pod has a container with a writable root filesystem",
				Description: "Set securityContext.readOnlyRootFilesystem to true",
			},
		},
		{
			ctx: &corev1.SecurityContext{
				ReadOnlyRootFilesystem: b(false),
			},
			expectedGrade: scorecard.GradeCritical,
			expectedComment: &scorecard.TestScoreComment{
				Path:        "foobar",
				Summary:     "The pod has a container with a writable root filesystem",
				Description: "Set securityContext.readOnlyRootFilesystem to true",
			},
		},

		// Context is non-null, but has all null values
		{
			ctx:           &corev1.SecurityContext{},
			expectedGrade: scorecard.GradeCritical,
			expectedComment: &scorecard.TestScoreComment{
				Path:        "foobar",
				Summary:     "The pod has a container with a writable root filesystem",
				Description: "Set securityContext.readOnlyRootFilesystem to true",
			},
		},
		// Context is non nul, but has all null values
		{
			ctx:           &corev1.SecurityContext{},
			expectedGrade: scorecard.GradeCritical,
			expectedComment: &scorecard.TestScoreComment{
				Path:        "foobar",
				Summary:     "The container is running with a low user ID",
				Description: "A userid above 10 000 is recommended to avoid conflicts with the host. Set securityContext.runAsUser to a value > 10000",
			},
		},
		// Context is non nul, but has all null values
		{
			ctx:           &corev1.SecurityContext{},
			expectedGrade: scorecard.GradeCritical,
			expectedComment: &scorecard.TestScoreComment{
				Path:        "foobar",
				Summary:     "The container running with a low group ID",
				Description: "A groupid above 10 000 is recommended to avoid conflicts with the host. Set securityContext.runAsGroup to a value > 10000",
			},
		},
		// PodSecurityContext is set, assert that the values are inherited
		{
			ctx: &corev1.SecurityContext{
				ReadOnlyRootFilesystem: b(true),
				RunAsNonRoot:           b(true),
				Privileged:             b(false),
			},
			podCtx: &corev1.PodSecurityContext{
				RunAsUser:  i(20000),
				RunAsGroup: i(20000),
			},
			expectedGrade: scorecard.GradeAllOK,
		},
		// PodSecurityContext is set, assert that the values are inherited
		// The container ctx has invalid values
		{
			ctx: &corev1.SecurityContext{
				ReadOnlyRootFilesystem: b(true),
				RunAsNonRoot:           b(true),
				Privileged:             b(false),
				RunAsUser:              i(4),
				RunAsGroup:             i(5),
			},
			podCtx: &corev1.PodSecurityContext{
				RunAsUser:  i(20000),
				RunAsGroup: i(20000),
			},
			expectedGrade: scorecard.GradeCritical,
			expectedComment: &scorecard.TestScoreComment{
				Path:        "foobar",
				Summary:     "The container running with a low group ID",
				Description: "A groupid above 10 000 is recommended to avoid conflicts with the host. Set securityContext.runAsGroup to a value > 10000",
			},
		},

		// Privileged defaults to "false"
		{
			ctx: &corev1.SecurityContext{
				ReadOnlyRootFilesystem: b(true),
				RunAsNonRoot:           b(true),
			},
			podCtx: &corev1.PodSecurityContext{
				RunAsUser:  i(20000),
				RunAsGroup: i(20000),
			},
			expectedGrade: scorecard.GradeAllOK,
		},

		// Privileged explicitly set to "false"
		{
			ctx: &corev1.SecurityContext{
				ReadOnlyRootFilesystem: b(true),
				RunAsNonRoot:           b(true),
				Privileged:             b(false),
			},
			podCtx: &corev1.PodSecurityContext{
				RunAsUser:  i(20000),
				RunAsGroup: i(20000),
			},
			expectedGrade: scorecard.GradeAllOK,
		},

		// Privileged explicitly set to "true"
		{
			ctx: &corev1.SecurityContext{
				ReadOnlyRootFilesystem: b(true),
				RunAsNonRoot:           b(true),
				Privileged:             b(true),
			},
			podCtx: &corev1.PodSecurityContext{
				RunAsUser:  i(20000),
				RunAsGroup: i(20000),
			},
			expectedGrade: scorecard.GradeCritical,
			expectedComment: &scorecard.TestScoreComment{
				Path:        "foobar",
				Summary:     "The container is privileged",
				Description: "Set securityContext.privileged to false. Privileged containers can access all devices on the host, and grants almost the same access as non-containerized processes on the host.",
			},
		},
	}

	for caseID, tc := range tests {
		test.Logf("Running caseID=%d", caseID)

		s := appsv1.StatefulSet{
			TypeMeta: metav1.TypeMeta{
				Kind:       "StatefulSet",
				APIVersion: "apps/v1",
			},
			Spec: appsv1.StatefulSetSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						SecurityContext: tc.podCtx,
						Containers: []corev1.Container{
							{
								Name:            "foobar",
								SecurityContext: tc.ctx,
							},
						},
					},
				},
			},
		}

		output, err := yaml.Marshal(s)
		assert.Nil(test, err, "caseID=%d", caseID)

		comments := testExpectedScoreReader(test, bytes.NewReader(output), "Container Security Context", tc.expectedGrade)

		if tc.expectedComment != nil {
			assert.Contains(test, comments, *tc.expectedComment, "caseID=%d", caseID)
		}
	}
}

func TestContainerSecurityContextPrivilegied(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-security-context-privilegied.yaml", "Container Security Context", scorecard.GradeCritical)
}

func TestContainerSecurityContextLowUser(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-security-context-low-user-id.yaml", "Container Security Context", scorecard.GradeCritical)
}

func TestContainerSecurityContextLowGroup(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-security-context-low-group-id.yaml", "Container Security Context", scorecard.GradeCritical)
}

func TestPodSecurityContextInherited(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "security-inherit-pod-security-context.yaml", "Container Security Context", scorecard.GradeAllOK)
}

func TestContainerSecurityContextAllGood(t *testing.T) {
	t.Parallel()
	c := testExpectedScore(t, "pod-security-context-all-good.yaml", "Container Security Context", scorecard.GradeAllOK)
	assert.Empty(t, c)
}

func TestContainerSeccompMissing(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-seccomp-profile"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-seccomp-no-annotation.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container Seccomp Profile", scorecard.GradeWarning)
}

func TestContainerSeccompAllGood(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-seccomp-profile"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-seccomp-annotated.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container Seccomp Profile", scorecard.GradeAllOK)
}

func TestServiceIgnoreNamespace(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["site"] = struct{}{}

	s, err := testScore(config.Configuration{
		VerboseOutput:     0,
		AllFiles:          []ks.NamedReader{testFile("service-externalname.yaml")},
		IgnoredNamespaces: structMap,
	})
	assert.Nil(t, err)
	assert.Len(t, s, 1)

	tested := false

	for _, o := range s {
		for _, c := range o.Checks {
			if c.Check.ID == "service-targets-pod" {
				assert.True(t, c.Skipped)
				assert.Equal(t, scorecard.GradeAllOK, c.Grade)
				tested = true
			}
		}
	}
	assert.True(t, tested)
}
