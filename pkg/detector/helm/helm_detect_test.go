package helm

import (
	"reflect"
	"testing"

	"github.com/Checkmarx/kics/pkg/model"
	"github.com/rs/zerolog"
)

func TestEngine_detectHelmLine(t *testing.T) { //nolint
	type args struct {
		file          *model.FileMetadata
		searchKey     string
		logWithFields *zerolog.Logger
		outputLines   int
	}

	tests := []struct {
		name string
		args args
		want model.VulnerabilityLines
	}{
		{
			name: "test_detect_helm_line",
			args: args{
				file: &model.FileMetadata{
					ID:       "1",
					ScanID:   "console",
					Document: model.Document{},
					Kind:     model.KindHELM,
					FileName: "test-connection.yaml",
					HelmID:   "# KICS_HELM_ID_0",
					OriginalData: `# KICS_HELM_ID_0:
apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "test_helm.fullname" . }}-test-connection"
  labels:
    {{- include "test_helm.labels" . | nindent 4 }}
  annotations:
	"helm.sh/hook": test
spec:
  containers:
    - name: wget
      image: busybox
	  command: ['wget']
	  args: ['{{ include "test_helm.fullname" . }}:{{ .Values.service.port }}']
    restartPolicy: Never
`,
					Content: ``,
				},
				searchKey:     "KICS_HELM_ID_0.metadata.name={{RELEASE-NAME-test_helm-test-connection}}.spec.containers",
				logWithFields: &zerolog.Logger{},
				outputLines:   1,
			},
			want: model.VulnerabilityLines{
				Line: 10,
				VulnLine: []model.VulnLines{
					{
						Position: 10,
						Line:     "  containers:",
					},
				},
				LineWithVulnerabilty: "  containers:",
			},
		},
		{
			name: "test_dup_values",
			args: args{
				file: &model.FileMetadata{
					ID:       "1",
					ScanID:   "console",
					Document: model.Document{},
					Kind:     model.KindHELM,
					FileName: "test-dup_values.yaml",
					IDInfo: map[int]interface{}{0: map[int]int{0: 0, 1: 1, 2: 2, 3: 3, 4: 4,
						5: 5, 6: 6, 7: 7, 8: 8, 9: 9, 10: 10, 11: 11, 12: 12, 13: 13, 14: 14, 15: 15, 16: 16, 17: 17,
						18: 18, 19: 19, 21: 21, 22: 22}},
					HelmID: "# KICS_HELM_ID_0",
					OriginalData: `# KICS_HELM_ID_0:
apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "test_helm.fullname" . }}-test-connection"
  labels:
    {{- include "test_helm.labels" . | nindent 4 }}
  annotations:
	"helm.sh/hook": test
spec:
  containers:
    - name: wget
      image: busybox
	  command: ['wget']
	  args: ['{{ include "test_helm.fullname" . }}:{{ .Values.service.port }}']
    restartPolicy: Never
  containers:
    - name: wget2
      image: busybox
	  command: ['wget']
	  args: ['{{ include "test_helm.fullname" . }}:{{ .Values.service.port }}']
    restartPolicy: Never
`,
					Content: ``,
				},
				searchKey:     "KICS_HELM_ID_0.metadata.name={{RELEASE-NAME-test_helm-test-connection}}.spec.containers",
				logWithFields: &zerolog.Logger{},
				outputLines:   1,
			},
			want: model.VulnerabilityLines{
				Line: 9,
				VulnLine: []model.VulnLines{
					{
						Position: 9,
						Line:     "spec:",
					},
				},
				LineWithVulnerabilty: "spec:",
			},
		},
		{
			name: "test_detect_helm_with_dups",
			args: args{
				file: &model.FileMetadata{
					ID:       "1",
					ScanID:   "console",
					Document: model.Document{},
					Kind:     model.KindHELM,
					FileName: "test-dups.yaml",
					HelmID:   "# KICS_HELM_ID_1",
					OriginalData: `# KICS_HELM_ID_0:
apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "test_helm.fullname" . }}-test-connection"
  labels:
    {{- include "test_helm.labels" . | nindent 4 }}
  annotations:
	"helm.sh/hook": test
spec:
  containers:
    - name: wget
      image: busybox
	  command: ['wget']
	  args: ['{{ include "test_helm.fullname" . }}:{{ .Values.service.port }}']
    restartPolicy: Never
---
# KICS_HELM_ID_1:
apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "test_helm.fullname" . }}-test-dups"
  labels:
    {{- include "test_helm.labels" . | nindent 4 }}
  annotations:
	"helm.sh/hook": test
spec:
  containers:
    - name: wget
      image: busybox
	  command: ['wget']
	  args: ['{{ include "test_helm.fullname" . }}:{{ .Values.service.port }}']
    restartPolicy: Never
`,
					Content: ``,
				},
				searchKey:     "KICS_HELM_ID_1.metadata.name={{RELEASE-NAME-test_helm-test-connection}}.spec.containers",
				logWithFields: &zerolog.Logger{},
				outputLines:   1,
			},
			want: model.VulnerabilityLines{
				Line: 26,
				VulnLine: []model.VulnLines{
					{
						Position: 26,
						Line:     "  containers:",
					},
				},
				LineWithVulnerabilty: "  containers:",
			},
		},
	}

	for _, tt := range tests {
		detector := DetectKindLine{}
		t.Run(tt.name, func(t *testing.T) {
			got := detector.DetectLine(tt.args.file, tt.args.searchKey, tt.args.logWithFields, tt.args.outputLines)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("detectHelmLine() = %v, want = %v", got, tt.want)
			}
		})
	}
}
