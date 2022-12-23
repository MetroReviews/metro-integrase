package types

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

/* Package `types` provides the interfaces that all bot lists wishing to use the official integrase  (integration code)
 * should follow
 */

const APIUrl = "https://catnip.metrobots.xyz"

type Bot struct {
	BotID    string `json:"bot_id"`
	CanAdd   bool   `json:"can_add"`
	Reviewer string `json:"reviewer"`
	Reason   string `json:"reason"`
}

func (b Bot) Resolve() (*FullBot, error) {
	req, err := http.NewRequest("GET", "https://catnip.metrobots.xyz"+b.BotID, nil)

	if err != nil {
		return nil, err
	}

	client := http.Client{Timeout: 5 * time.Second}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, errors.New("not found")
	}

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var bot FullBot

	err = json.Unmarshal(bytes, &bot)

	if err != nil {
		return nil, err
	}

	return &bot, nil
}

type FullBot struct {
	BotID           string `json:"bot_id"`
	Reviewer        string `json:"reviewer"`
	Username        string `json:"username"`
	Description     string `json:"description"`
	LongDescription string `json:"long_description"`
	Owner           string `json:"owner"`

	// The following fields are optionally set
	Banner      string   `json:"banner,omitempty"`
	ExtraOwners []string `json:"extra_owners,omitempty"` // Usually set
	ListSource  string   `json:"list_source,omitempty"`  // Usually set
	CrossAdd    bool     `json:"cross_add,omitempty"`    // In rare cases, this may not be set
	Website     string   `json:"website,omitempty"`      // Usually set
	Support     string   `json:"support,omitempty"`
	Donate      string   `json:"donate,omitempty"`
	Library     string   `json:"library,omitempty"`
	Prefix      string   `json:"prefix,omitempty"`
	Invite      string   `json:"invite,omitempty"`
	NSFW        bool     `json:"nsfw,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	ReviewNote  string   `json:"review_note,omitempty"`
	Limited     bool     `json:"limited"`
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
	// List ID (required)
	ListID string
	// Secret Key (required)
	SecretKey string
	// Domain name (optional, if not specified, auto-registration will be disabled)
	DomainName string
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
