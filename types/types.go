package types

/* Package `types` provides the interfaces that all bot lists wishing to use the official integrase  (integration code)
 * should follow
 */

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

type ListConfig struct {
	StartupLogs bool
}

type ListAdapter interface {
	// Calling get config should return the configuration for the list. It should not produce any side effects
	GetConfig() ListConfig
	// Calling claim bot should claim the bot if it is present but not add the bot if not
	ClaimBot(bot *Bot) error
	// Calling unclaim bot should unclaim the bot if it is present but not add the bot if not
	UnclaimBot(bot *Bot) error
	// Calling approve bot should approve the bot if it is not present and should add the bot if it is
	ApproveBot(bot *Bot) error
	// Calling deny bot should deny the bot if it is present but not add the bot if not
	DenyBot(bot *Bot) error
	// Upcoming: calling data deletion request should delete the bot and all associated information
	DataDelete(id string) error
	// Upcoming: calling data request should return the bot and all associated information as a map
	DataRequest(id string) (map[string]interface{}, error)
}
