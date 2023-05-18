package brood

type User struct {
	Id              string `json:"id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	NormalizedEmail string `json:"normalized_email"`
	Verified        bool   `json:"verified"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type UserGeneratedToken struct {
	Id        string `json:"access_token"`
	TokenType string `json:"token_type"`
}

type UserToken struct {
	Id        string `json:"id"`
	TokenType string `json:"token_type"`
	CreatedAt string `json:"created_at"`
	UserId    string `json:"user_id"`
	UpdatedAt string `json:"updated_at"`
	Note      string `json:"note"`
	Active    bool   `json:"active"`
}

type UserTokensList struct {
	UserId   string      `json:"user_id"`
	Username string      `json:"username"`
	Tokens   []UserToken `json:"token"`
}

type Group struct {
	Id   string `json:"id"`
	Name string `json:"group_name"`
}

type UserGroup struct {
	GroupID       string `json:"group_id"`
	GroupName     string `json:"group_name"`
	Autogenerated bool   `json:"autogenerated"`
	UserId        string `json:"user_id"`
	UserType      string `json:"user_type"`
}

type UserGroupsList struct {
	Groups []UserGroup `json:"groups"`
}

// Applications

type Application struct {
	Id          string `json:"id"`
	GroupId     string `json:"group_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type ApplicationsList struct {
	Applications []Application `json:"applications"`
}

// Resources
type Resource struct {
	Id            string      `json:"id"`
	ApplicationId string      `json:"application_id"`
	ResourceData  interface{} `json:"resource_data"`
}

type Resources struct {
	Resources []Resource `json:"resources"`
}

type resourceCreateRequest struct {
	ApplicationId string      `json:"application_id"`
	ResourceData  interface{} `json:"resource_data"`
}

type resourceUpdateRequest struct {
	Update   interface{} `json:"update"`
	DropKeys []string    `json:"drop_keys"`
}
