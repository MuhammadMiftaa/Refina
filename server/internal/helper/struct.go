package helper

type GitHubPlan struct {
	Collaborators int    `json:"collaborators"`
	Name          string `json:"name"`
	PrivateRepos  int    `json:"private_repos"`
	Space         int64  `json:"space"`
}

type GitHubUser struct {
	AvatarURL               string      `json:"avatar_url"`
	Bio                     string      `json:"bio"`
	Blog                    string      `json:"blog"`
	Collaborators           int         `json:"collaborators"`
	Company                 interface{} `json:"company"` // nullable
	CreatedAt               string      `json:"created_at"`
	DiskUsage               int         `json:"disk_usage"`
	Email                   interface{} `json:"email"` // nullable
	EventsURL               string      `json:"events_url"`
	Followers               int         `json:"followers"`
	FollowersURL            string      `json:"followers_url"`
	Following               int         `json:"following"`
	FollowingURL            string      `json:"following_url"`
	GistsURL                string      `json:"gists_url"`
	GravatarID              string      `json:"gravatar_id"`
	Hireable                interface{} `json:"hireable"` // nullable
	HTMLURL                 string      `json:"html_url"`
	ID                      float64     `json:"id"`
	Location                interface{} `json:"location"` // nullable
	Login                   string      `json:"login"`
	Name                    string      `json:"name"`
	NodeID                  string      `json:"node_id"`
	NotificationEmail       interface{} `json:"notification_email"` // nullable
	OrganizationsURL        string      `json:"organizations_url"`
	OwnedPrivateRepos       int         `json:"owned_private_repos"`
	Plan                    GitHubPlan  `json:"plan"`
	PrivateGists            int         `json:"private_gists"`
	PublicGists             int         `json:"public_gists"`
	PublicRepos             int         `json:"public_repos"`
	ReceivedEventsURL       string      `json:"received_events_url"`
	ReposURL                string      `json:"repos_url"`
	SiteAdmin               bool        `json:"site_admin"`
	StarredURL              string      `json:"starred_url"`
	SubscriptionsURL        string      `json:"subscriptions_url"`
	TotalPrivateRepos       int         `json:"total_private_repos"`
	TwitterUsername         interface{} `json:"twitter_username"` // nullable
	TwoFactorAuthentication bool        `json:"two_factor_authentication"`
	Type                    string      `json:"type"`
	UpdatedAt               string      `json:"updated_at"`
	URL                     string      `json:"url"`
	UserViewType            string      `json:"user_view_type"`
}
