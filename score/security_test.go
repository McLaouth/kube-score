package score

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
)

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

func TestContainerSecurityContextUserGroupIDAllGood(t *testing.T) {
	t.Parallel()
	structMap := make(map[string]struct{})
	structMap["container-security-context-user-group-id"] = struct{}{}
	c := testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-security-context-all-good.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container Security Context User Group ID", scorecard.GradeAllOK)
	assert.Empty(t, c)
}

func TestContainerSecurityContextUserGroupIDLowGroup(t *testing.T) {
	t.Parallel()
	optionalChecks := make(map[string]struct{})
	optionalChecks["container-security-context-user-group-id"] = struct{}{}
	comments := testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-security-context-low-group-id.yaml")},
		EnabledOptionalTests: optionalChecks,
	}, "Container Security Context User Group ID", scorecard.GradeCritical)
	assert.Contains(t, comments, scorecard.TestScoreComment{
		Path:        "foobar",
		Summary:     "The container running with a low group ID",
		Description: "A groupid above 10 000 is recommended to avoid conflicts with the host. Set securityContext.runAsGroup to a value > 10000",
	})
}

func TestContainerSecurityContextUserGroupIDLowUser(t *testing.T) {
	t.Parallel()
	optionalChecks := make(map[string]struct{})
	optionalChecks["container-security-context-user-group-id"] = struct{}{}
	comments := testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-security-context-low-user-id.yaml")},
		EnabledOptionalTests: optionalChecks,
	}, "Container Security Context User Group ID", scorecard.GradeCritical)
	assert.Contains(t, comments, scorecard.TestScoreComment{
		Path:        "foobar",
		Summary:     "The container is running with a low user ID",
		Description: "A userid above 10 000 is recommended to avoid conflicts with the host. Set securityContext.runAsUser to a value > 10000",
	})
}

func TestContainerSecurityContextUserGroupIDNoSecurityContext(t *testing.T) {
	t.Parallel()
	optionalChecks := make(map[string]struct{})
	optionalChecks["container-security-context-user-group-id"] = struct{}{}
	comments := testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-security-context-nosecuritycontext.yaml")},
		EnabledOptionalTests: optionalChecks,
	}, "Container Security Context User Group ID", scorecard.GradeCritical)
	assert.Contains(t, comments, scorecard.TestScoreComment{
		Path:        "foobar",
		Summary:     "Container has no configured security context",
		Description: "Set securityContext to run the container in a more secure context.",
	})
}

func TestContainerSecurityContextPrivilegedAllGood(t *testing.T) {
	t.Parallel()
	structMap := make(map[string]struct{})
	structMap["container-security-context-privileged"] = struct{}{}
	c := testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-security-context-all-good.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container Security Context Privileged", scorecard.GradeAllOK)
	assert.Empty(t, c)
}

func TestContainerSecurityContextPrivilegedPrivileged(t *testing.T) {
	t.Parallel()
	optionalChecks := make(map[string]struct{})
	optionalChecks["container-security-context-privileged"] = struct{}{}
	comments := testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-security-context-privileged.yaml")},
		EnabledOptionalTests: optionalChecks,
	}, "Container Security Context Privileged", scorecard.GradeCritical)
	assert.Contains(t, comments, scorecard.TestScoreComment{
		Path:        "foobar",
		Summary:     "The container is privileged",
		Description: "Set securityContext.privileged to false. Privileged containers can access all devices on the host, and grants almost the same access as non-containerized processes on the host.",
	})
}

func TestContainerSecurityContextReadOnlyRootFilesystemAllGood(t *testing.T) {
	t.Parallel()
	structMap := make(map[string]struct{})
	structMap["container-security-context-readonlyrootfilesystem"] = struct{}{}
	c := testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-security-context-all-good.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container Security Context ReadOnlyRootFilesystem", scorecard.GradeAllOK)
	assert.Empty(t, c)
}

func TestContainerSecurityContextReadOnlyRootFilesystemWriteable(t *testing.T) {
	t.Parallel()
	optionalChecks := make(map[string]struct{})
	optionalChecks["container-security-context-readonlyrootfilesystem"] = struct{}{}
	comments := testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-security-context-writeablerootfilesystem.yaml")},
		EnabledOptionalTests: optionalChecks,
	}, "Container Security Context ReadOnlyRootFilesystem", scorecard.GradeCritical)
	assert.Contains(t, comments, scorecard.TestScoreComment{
		Path:        "foobar",
		Summary:     "The pod has a container with a writable root filesystem",
		Description: "Set securityContext.readOnlyRootFilesystem to true",
	})
}

func TestContainerSecurityContextReadOnlyRootFilesystemNoSecurityContext(t *testing.T) {
	t.Parallel()
	optionalChecks := make(map[string]struct{})
	optionalChecks["container-security-context-readonlyrootfilesystem"] = struct{}{}
	comments := testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-security-context-nosecuritycontext.yaml")},
		EnabledOptionalTests: optionalChecks,
	}, "Container Security Context ReadOnlyRootFilesystem", scorecard.GradeCritical)
	assert.Contains(t, comments, scorecard.TestScoreComment{
		Path:        "foobar",
		Summary:     "Container has no configured security context",
		Description: "Set securityContext to run the container in a more secure context.",
	})
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
