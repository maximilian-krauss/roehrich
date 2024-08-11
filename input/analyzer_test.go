package input

import "testing"

func TestGetMRInfoWithValidInput(t *testing.T) {
	result, err := GetMRInfo("https://git.gitlabinstance.local/sub_group/long.project-name/-/merge_requests/1234/diffs?commit_id=8843d7f92416211de9ebb963ff4ce28125932878")
	if err != nil {
		t.Fatal(err)
	}
	var (
		projectId   = "1234"
		projectName = "sub_group/long.project-name"
	)
	if result.Id != projectId {
		t.Fatalf("wrong id: got %s, wanted %s", result.Id, projectId)
	}
	if result.ProjectName != projectName {
		t.Fatalf("wrong project name: got %s, wanted %s", result.ProjectName, projectName)
	}
}

func TestGetMRInfoWithInvalidInput(t *testing.T) {
	var testCases = []struct {
		name  string
		input string
		want  string
	}{
		{"completely wrong input", "foobar", "parse \"foobar\": invalid URI for request"},
		{"missing mr id", "https://git.gitlabinstance.local/subgroup/project-name/-/merge_requests/", "https://git.gitlabinstance.local/subgroup/project-name/-/merge_requests/ is not a valid merge request url"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := GetMRInfo(testCase.input)
			if result != nil {
				t.Error("expected nil result")
			}
			if err.Error() != testCase.want {
				t.Errorf("wrong error: got %s, wanted %s", err.Error(), testCase.want)
			}
		})
	}

}
