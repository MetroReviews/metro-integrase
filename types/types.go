package types

/* Package `types` provides the interfaces that all bot lists wishing to use the official integrase  (integration code)
 * should follow
 */

const APIUrl = "https://catnip.metrobots.xyz"

type Bot struct {
	BotID           string   `json:"bot_id"`
	Reviewer        string   `json:"reviewer"`
	Username        string   `json:"username"`
	Description     string   `json:"description"`
	LongDescription string   `json:"long_description"`
	NSFW            bool     `json:"nsfw"`
	Tags            []string `json:"tags"` // Auto set to []string{} if CrossAdd is false
	Owner           string   `json:"owner"`
	ExtraOwners     []string `json:"extra_owners"`
	ListSource      string   `json:"list_source"`
	Reason          string   `json:"reason,omitempty"`
	CrossAdd        *bool    `json:"cross_add,omitempty"` // In rare cases, this may not be set

	// The following fields are optionally set
	Website string `json:"website,omitempty"`
	Github  string `json:"github,omitempty"`
	Support string `json:"support,omitempty"`
	Donate  string `json:"donate,omitempty"`
	Library string `json:"library,omitempty"`
	Prefix  string `json:"prefix,omitempty"`
	Invite  string `json:"invite,omitempty"`
}

type ListPatch struct {
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	Domain        string `json:"domain,omitempty"`
	ClaimBotAPI   string `json:"claim_bot_api,omitempty"`
	UnclaimBotAPI string `json:"unclaim_bot_api,omitempty"`
	ApproveBotAPI string `json:"approve_bot_api,omitempty"`
	DenyBotAPI    string `json:"deny_bot_api,omitempty"`
	// Upcoming
	DataRequestAPI string `json:"data_request_api,omitempty"`
	// Upcoming
	DataDeletionAPI string `json:"data_deletion_api,omitempty"`
	ResetSecretKey  bool   `json:"reset_secret_key,omitempty"`
	Icon            string `json:"icon,omitempty"`
}

type ListPatchResp struct {
	HasUpdated []string `json:"has_updated,omitempty"`
	SecretKey  string   `json:"secret_key,omitempty"`
}

// Core structs
type ListConfig struct {
	// Logs on startup
	StartupLogs bool
	// Logs for requests
	RequestLogs bool
	// List ID (required)
	ListID string
	// Secret Key (required)
	SecretKey string
	// Which IP/Port to bind to
	BindAddr string
	// Domain name (optional, if not specified, auto-registration will be disabled)
	DomainName string
}

type ListAdapter interface {
	// Calling get config should return the configuration for the list. It should not produce any side effects
	GetConfig() ListConfig
	// Calling claim bot should claim the bot if it is present but not add the bot if not
	ClaimBot(adp *ListAdapter, bot *Bot) error
	// Calling unclaim bot should unclaim the bot if it is present but not add the bot if not
	UnclaimBot(adp *ListAdapter, bot *Bot) error
	// Calling approve bot should approve the bot if it is not present and should add the bot if it is
	ApproveBot(adp *ListAdapter, bot *Bot) error
	// Calling deny bot should deny the bot if it is present but not add the bot if not
	DenyBot(adp *ListAdapter, bot *Bot) error
	// Upcoming: calling data deletion request should delete the bot and all associated information
	DataDelete(adp *ListAdapter, id string) error
	// Upcoming: calling data request should return the bot and all associated information as a map
	DataRequest(adp *ListAdapter, id string) (map[string]interface{}, error)
}
