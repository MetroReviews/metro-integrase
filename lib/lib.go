package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/MetroReviews/metro-integrase/types"
)

/*
 * Package ``lib`` provides the main code for the integrase system.
 */

func authReq(r *http.Request, cfg types.ListConfig) bool {
	if r == nil {
		return false
	}

	if r.Header.Get("Authorization") == "" || r.Header.Get("Authorization") != cfg.SecretKey {
		return false
	}

	return true
}

type ListFunction func(bot *types.Bot) error

func coreHandler(fn ListFunction, cfg types.ListConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
			return
		}

		if !authReq(r, cfg) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request: " + err.Error()))
			return
		}

		var bot types.Bot
		err = json.Unmarshal(body, &bot)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Serialization error occured: " + err.Error()))
			return
		}

		adpErr := fn(&bot)

		if adpErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Request handle error:" + adpErr.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK :)"))
	}
}

// Patches a list
func PatchList(cfg types.ListConfig, data types.ListPatch) (*types.ListPatchResp, error) {
	payload, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 20 * time.Second}

	req, err := http.NewRequest("PATCH", types.APIUrl+"/lists/"+cfg.ListID, bytes.NewBuffer(payload))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", cfg.SecretKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var respData types.ListPatchResp

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &respData)

		if err != nil {
			return nil, err
		}

		return &respData, nil
	} else {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}

		return nil, errors.New(string(body))
	}
}

type Router interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request))
}

// Prepares a web server handling all core integrase functions
func Prepare(adp types.ListAdapter, r Router) {
	cfg := adp.GetConfig()

	if cfg.StartupLogs {
		fmt.Println("Starting integrase server")
	}

	if cfg.ListID == "" {
		panic("List ID not set")
	} else if cfg.SecretKey == "" {
		panic("Secret Key not set")
	}

	r.HandleFunc("/claim", coreHandler(adp.ClaimBot, cfg))
	r.HandleFunc("/unclaim", coreHandler(adp.UnclaimBot, cfg))
	r.HandleFunc("/approve", coreHandler(adp.ApproveBot, cfg))
	r.HandleFunc("/deny", coreHandler(adp.DenyBot, cfg))
	r.HandleFunc("/data-request", func(w http.ResponseWriter, r *http.Request) {
		if !authReq(r, cfg) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		botId := r.URL.Query().Get("bot_id")

		if botId == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bot ID is missing"))
			return
		}

		bot, err := adp.DataRequest(botId)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Request handle error:" + err.Error()))
			return
		}

		botStr, err := json.Marshal(bot)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Serialization error occured: " + err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(botStr))
	})

	r.HandleFunc("/data-delete", func(w http.ResponseWriter, r *http.Request) {
		if !authReq(r, cfg) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		botId := r.URL.Query().Get("bot_id")

		if botId == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bot ID is missing"))
			return
		}

		err := adp.DataDelete(botId)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Request handle error:" + err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("All associated data has been deleted from this list according to the lists adapter"))
	})

	if cfg.DomainName != "" {
		if cfg.StartupLogs {
			fmt.Println("Updating Metro Reviews with new routes")
		}

		patched, err := PatchList(cfg, types.ListPatch{
			ClaimBotAPI:     cfg.DomainName + "/claim",
			UnclaimBotAPI:   cfg.DomainName + "/unclaim",
			ApproveBotAPI:   cfg.DomainName + "/approve",
			DenyBotAPI:      cfg.DomainName + "/deny",
			DataRequestAPI:  cfg.DomainName + "/data-request",
			DataDeletionAPI: cfg.DomainName + "/data-delete",
		})

		if cfg.StartupLogs {
			if err != nil {
				fmt.Println("Metro Reviews update failed: ", err)
			} else {
				fmt.Println("Metro Reviews update successful with ", patched.HasUpdated, " updated")
			}
		}
	}

	if cfg.StartupLogs {
		fmt.Println(`Integrase prepared, now you need to run something like the below to start integrase:
		
// Add any middleware here (ex: logging middleware)
log := handlers.LoggingHandler(os.Stdout, r)

http.ListenAndServe("ADDRESS HERE", log)

Don't like these logs? Disable StartupLogs
		`)
	}
}
