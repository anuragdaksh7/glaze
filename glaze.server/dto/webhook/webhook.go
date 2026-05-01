package webhookDto

type PushPayload struct {
	After      string `json:"after"`
	Ref        string `json:"ref"`
	Repository struct {
		ID int64 `json:"id"`
	} `json:"repository"`
	HeadCommit struct {
		Message string `json:"message"`
		Author  struct {
			Name string `json:"name"`
		} `json:"author"`
	} `json:"head_commit"`
}
