package tester

import (
	"fmt"
	"strings"

	"github.com/MetroReviews/metro-integrase/types"
)

var DefaultBot = types.Bot{
	BotID:           "968734728465289248", // Metro Reviews bot
	Reviewer:        "510065483693817867", // Toxic Dev
	Username:        "Metro Reviews",
	Description:     "Metro Reviews is a pog pog pog bot. This bot is purely for testing purposes.",
	LongDescription: strings.Repeat("Metro Reviews is a pog pog pog bot. This bot is purely for testing purposes.\n\n", 100),
	Owner:           "728871946456137770",           // Burgerking
	ExtraOwners:     []string{"564164277251080208"}, // Select
	Invite:          "https://discord.com/api/oauth2/authorize?client_id=968734728465289248&permissions=8&scope=bot%20applications.commands",
	Website:         "https://metroreviews.xyz",
	Github:          "https://github.com/MetroReviews",
	Support:         "https://discord.gg/49DE35a5eJ",
	NSFW:            true,
	ReviewNote:      "This bot is purely for testing purposes. It is not a real bot.",
	Tags:            []string{"Utility", "Moderation"},
	Reason:          "Test reason",
	ListSource:      "3b50d5e8-d0a0-4e63-aff7-f81068e9ad36", // IBL List ID
}

type Tester struct {
	Adapter types.ListAdapter
}

func (t *Tester) Claim() {
	err := t.Adapter.ClaimBot(&DefaultBot)
	fmt.Println("Claimed bot with error:", err)
}

func (t *Tester) Unclaim() {
	err := t.Adapter.UnclaimBot(&DefaultBot)
	fmt.Println("Unclaimed bot with error:", err)
}

func (t *Tester) Approve() {
	err := t.Adapter.ApproveBot(&DefaultBot)
	fmt.Println("Approved bot with error:", err)
}

func (t *Tester) Deny() {
	err := t.Adapter.DenyBot(&DefaultBot)
	fmt.Println("Denied bot with error:", err)
}
