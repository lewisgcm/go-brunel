package notify

import (
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared"
)

const (
	gitLabPrivateTokenHeader = "Private-Token"
)

type GitLabNotify struct {
	Secret     string
	URL        string
	Repository store.JobStore
}

func (notify *GitLabNotify) Notify(id shared.JobID) error {
	// client := &http.Client{}
	// stateText := "pending"

	//job, err := notify.Repository.FindAllByJobID(id)
	//if err != nil {
	//	return errors.Wrap(err, "error getting job")
	//}
	//
	//switch job.State {
	//case shared.JobStateProcessing:
	//	stateText = "running"
	//case shared.JobStateSuccess:
	//	stateText = "success"
	//case shared.JobStateFailed:
	//	stateText = "failed"
	//case shared.JobStateCancelled:
	//	stateText = "canceled"
	//}

	//form := url.Values{
	//	"state": {stateText},
	//}
	//req, err := http.NewRequest(
	//	http.MethodPost,
	//	fmt.Sprintf("%s/projects/%s%%2F%s/statuses/%s", notify.URL, job.Repository.Project, job.Repository.Name, job.Commit.Revision),
	//	strings.NewReader(form.Encode()),
	//)
	//if err != nil {
	//	return errors.Wrap(err, "error building request")
	//}
	//
	//req.Header.Set(gitLabPrivateTokenHeader, notify.Secret)
	//
	//// TODO we could add better error handling here, i.e log the gitlab message if there is an error
	//response, err := client.Do(req)
	//if response.StatusCode >= 300 {
	//	err = errors.New("received non 200 status from gitlab API")
	//}

	return nil // errors.Wrap(err, "error posting status")
}
