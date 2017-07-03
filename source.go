package bitbucket_pull_request_resource

type Source struct {
	Endpoint   string `json:"endpoint"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Repository string `json:"repository"`
	Owner      string `json:"owner"`
}
