package main

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strconv"

	"github.com/faunX/ogame"
)

//https://s109-no.ogame.gameforge.com/game/index.php?page=ingame&component=rewarding&tab=rewards&action=fetchRewards&ajax=1&asJson=1&tier=2&token=20d1dd0a7804a0f4ae24e5dc7e190d0d
//var rewardToken = "79e5c225d17c474d94022008450740d4";

/*
   var tab = "rewards";
   var taskToken = "45ae978a5526d37d3788d662ad2a2f90";
   var rewardToken = "79e5c225d17c474d94022008450740d4";
   var selectedTier = 1;
   var tiers = {"1":"100","2":"200","3":"300","4":"400","5":"500","6":"600","7":"700"};
   var currentTier = 2;
   var urlFetchTasks = "https:\/\/s109-no.ogame.gameforge.com\/game\/index.php?page=ingame&component=rewarding&tab=tasks&action=fetchTasks&ajax=1";
   var urlFetchRewards = "https:\/\/s109-no.ogame.gameforge.com\/game\/index.php?page=ingame&component=rewarding&tab=rewards&action=fetchRewards&ajax=1";
*/

type Tasks struct {
	ComponentInitializationSuccessful bool `json:"componentInitializationSuccessful"`
	Rewards                           []struct {
		ID            int    `json:"id"`
		RewardType    int    `json:"rewardType"`
		SpecialReward int    `json:"specialReward"`
		CSSClass      string `json:"cssClass"`
		Name          string `json:"name"`
		Quantity      int    `json:"quantity"`
		Selected      bool   `json:"selected"`
		TierSelected  bool   `json:"tierSelected"`
		TierAvailable bool   `json:"tierAvailable"`
		Selectable    bool   `json:"selectable"`
		Reason        string `json:"reason"`
		ItemImage     string `json:"itemImage"`
	} `json:"rewards"`
	SelectedSpecialReward string `json:"selectedSpecialReward"`
	SelectedTierValue     string `json:"selectedTierValue"`
	PlayerTritium         int    `json:"playerTritium"`
	SelectedTier          int    `json:"selectedTier"`
	AllowedToBuy          bool   `json:"allowedToBuy"`
	AlreadySelected       bool   `json:"alreadySelected"`
	Success               bool   `json:"success"`
	Token                 string `json:"token"`
	NewAjaxToken          string `json:"newAjaxToken"`
}

type RewardResult struct {
	ComponentInitializationSuccessful bool          `json:"componentInitializationSuccessful"`
	SelectedReward                    int           `json:"selectedReward"`
	RewardSelected                    string        `json:"rewardSelected"`
	Rewarded                          string        `json:"rewarded"`
	SelectedTier                      int           `json:"selectedTier"`
	AllOfficers                       bool          `json:"allOfficers"`
	Status                            string        `json:"status"`
	Message                           string        `json:"message"`
	Token                             string        `json:"token"`
	Components                        []interface{} `json:"components"`
	NewAjaxToken                      string        `json:"newAjaxToken"`
}

func GrabRewards(b *ogame.OGame) {
	// Tier Config
	tiers := make(map[string]string)

	tiers["1"] = "1"
	tiers["2"] = "3"
	tiers["3"] = "1"
	tiers["4"] = "1"
	tiers["5"] = "3"
	tiers["6"] = "1"
	tiers["7"] = "1"

	var params url.Values
	var payload url.Values
	params.Add("page", "ingame")
	params.Add("component", "rewarding")
	pageHTML, err := b.GetPageContent(params)
	if err != nil {
		return
	}

	//	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))

	r := regexp.MustCompile(`var rewardToken = "([^"]+)"`)
	m := r.FindStringSubmatch(string(pageHTML))
	if len(m) != 2 {
		err = errors.New("failed to find reward token")
		return
	}
	rewardToken := m[1]

	r2 := regexp.MustCompile(`var taskToken = "([^"]+)"`)
	n := r2.FindStringSubmatch(string(pageHTML))
	if len(m) != 2 {
		err = errors.New("failed to find reward token")
		return
	}
	taskToken := n[1]

	for k, v := range tiers {
		//token=f6a5fc6294176cb49e363ba6d5bcc1b4&asJson=1
		var paramsTasks url.Values
		paramsTasks.Add("page", "ingame")
		paramsTasks.Add("component", "rewarding")
		paramsTasks.Add("tab", "rewards")
		paramsTasks.Add("action", "fetchRewards")
		paramsTasks.Add("ajax", "1")
		paramsTasks.Add("asJson", "1")
		paramsTasks.Add("tier", k)
		paramsTasks.Add("token", taskToken)

		pageTasksHTML, _ := b.GetPageContent(paramsTasks)
		var resultTasks Tasks
		json.Unmarshal(pageTasksHTML, &resultTasks)
		rewardToken = resultTasks.Token
		selectedTierValue, _ := strconv.ParseInt(resultTasks.SelectedTierValue, 10, 64)
		if !resultTasks.AlreadySelected && resultTasks.AllowedToBuy && int64(resultTasks.PlayerTritium) >= selectedTierValue {
			var result RewardResult
			params.Add("tab", "rewards")
			params.Add("action", "submitReward")
			params.Add("asJson", "1")

			payload.Add("selectedReward", v)
			payload.Add("selectedTier", k)
			payload.Add("token", rewardToken)

			pageHTML, err := b.PostPageContent(params, payload)
			if err != nil {

			}

			json.Unmarshal(pageHTML, &result)
			rewardToken = result.Token
		}
	}

	/*
				1. Obtain Token from
				https://s109-no.ogame.gameforge.com/game/index.php?page=ingame&component=rewarding
			    var taskToken = "981c91cb62ac6e82ad59b2c411b43fbd";
		    	var rewardToken = "07764059e506d6d6b271a3b21c04a0bc";
	*/

	// GET: https://s109-no.ogame.gameforge.com/game/index.php?page=ingame&component=rewarding&tab=rewards&action=fetchRewards&ajax=1&tier=2&token=f6a5fc6294176cb49e363ba6d5bcc1b4&asJson=1
	// Response:
	//{"componentInitializationSuccessful":true,"rewards":[{"id":9,"rewardType":6,"specialReward":0,"cssClass":"","name":"Platinum Metal Booster","quantity":2,"selected":false,"tierSelected":false,"tierAvailable":false,"selectable":true,"reason":"","itemImage":"\/cdn\/img\/item-images\/ff1ad1a6d5879cb0ea720199c9eb6518584f0922-large.png"},{"id":10,"rewardType":6,"specialReward":0,"cssClass":"","name":"Platinum Crystal Booster","quantity":2,"selected":false,"tierSelected":false,"tierAvailable":false,"selectable":true,"reason":"","itemImage":"\/cdn\/img\/item-images\/d4e203516d95ae28081a3d985818e2df5a2475d4-large.png"},{"id":11,"rewardType":6,"specialReward":0,"cssClass":"","name":"Platinum Deuterium Booster","quantity":2,"selected":false,"tierSelected":false,"tierAvailable":false,"selectable":true,"reason":"","itemImage":"\/cdn\/img\/item-images\/8245a9d21fb27088b25d48ae024e9382fcea1448-large.png"}],"selectedSpecialReward":"Receive the following additional rewards if the\u00a0<a href=\"https:\/\/s109-no.ogame.gameforge.com\/game\/index.php?page=premium&openDetail=12\" target=\"_self\">Commanding Staff<\/a> is active:","selectedTierValue":"200","playerTritium":100,"selectedTier":2,"allowedToBuy":false,"alreadySelected":false,"urlSubmitReward":"https:\/\/s109-no.ogame.gameforge.com\/game\/index.php?page=ingame&component=rewarding&tab=rewards&action=submitReward&asJson=1","success":true,"token":"809cafd6a50ebbb928c6e8112c03c94a","components":[],"newAjaxToken":"809cafd6a50ebbb928c6e8112c03c94a"}

	// POST: https://s109-no.ogame.gameforge.com/game/index.php?page=ingame&component=rewarding&tab=rewards&action=submitReward&asJson=1
	/*
		QUERY PARAMS:
			page=ingame
			component=rewarding
			tab=rewards
			action=submitReward
			asJson=1
		Payload:
			selectedReward=1&selectedTier=1&token=767d680741d271230e7ed6642ab02dba

		Result:
			{"componentInitializationSuccessful":true,"selectedReward":1,"rewardSelected":"You have already received this reward.","rewarded":"Selected","selectedTier":1,"allOfficers":false,"status":"success","message":"The reward has been added to your account.","token":"07764059e506d6d6b271a3b21c04a0bc","components":[],"newAjaxToken":"07764059e506d6d6b271a3b21c04a0bc"}
	*/
}
