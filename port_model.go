package yuubari_go

type PortAPI struct {
	APIResult    int    `json:"api_result"`
	APIResultMsg string `json:"api_result_msg"`
	APIData      struct {
		APIMaterial []struct {
			APIMemberID int `json:"api_member_id"`
			APIID       int `json:"api_id"`
			APIValue    int `json:"api_value"`
		} `json:"api_material"`
		APIDeckPort []struct {
			APIMemberID int    `json:"api_member_id"`
			APIID       int    `json:"api_id"`
			APIName     string `json:"api_name"`
			APINameID   string `json:"api_name_id"`
			APIMission  []int  `json:"api_mission"`
			APIFlagship string `json:"api_flagship"`
			APIShip     []int  `json:"api_ship"`
		} `json:"api_deck_port"`
		APINdock []struct {
			APIMemberID        int    `json:"api_member_id"`
			APIID              int    `json:"api_id"`
			APIState           int    `json:"api_state"`
			APIShipID          int    `json:"api_ship_id"`
			APICompleteTime    int64  `json:"api_complete_time"`
			APICompleteTimeStr string `json:"api_complete_time_str"`
			APIItem1           int    `json:"api_item1"`
			APIItem2           int    `json:"api_item2"`
			APIItem3           int    `json:"api_item3"`
			APIItem4           int    `json:"api_item4"`
		} `json:"api_ndock"`
		APIShip []struct {
			APIID          int   `json:"api_id"`
			APISortno      int   `json:"api_sortno"`
			APIShipID      int   `json:"api_ship_id"`
			APILv          int   `json:"api_lv"`
			APIExp         []int `json:"api_exp"`
			APINowhp       int   `json:"api_nowhp"`
			APIMaxhp       int   `json:"api_maxhp"`
			APISoku        int   `json:"api_soku"`
			APILeng        int   `json:"api_leng"`
			APISlot        []int `json:"api_slot"`
			APIOnslot      []int `json:"api_onslot"`
			APISlotEx      int   `json:"api_slot_ex"`
			APIKyouka      []int `json:"api_kyouka"`
			APIBacks       int   `json:"api_backs"`
			APIFuel        int   `json:"api_fuel"`
			APIBull        int   `json:"api_bull"`
			APISlotnum     int   `json:"api_slotnum"`
			APINdockTime   int   `json:"api_ndock_time"`
			APINdockItem   []int `json:"api_ndock_item"`
			APISrate       int   `json:"api_srate"`
			APICond        int   `json:"api_cond"`
			APIKaryoku     []int `json:"api_karyoku"`
			APIRaisou      []int `json:"api_raisou"`
			APITaiku       []int `json:"api_taiku"`
			APISoukou      []int `json:"api_soukou"`
			APIKaihi       []int `json:"api_kaihi"`
			APITaisen      []int `json:"api_taisen"`
			APISakuteki    []int `json:"api_sakuteki"`
			APILucky       []int `json:"api_lucky"`
			APILocked      int   `json:"api_locked"`
			APILockedEquip int   `json:"api_locked_equip"`
		} `json:"api_ship"`
		APIBasic struct {
			APIMemberID         string      `json:"api_member_id"`
			APINickname         string      `json:"api_nickname"`
			APINicknameID       string      `json:"api_nickname_id"`
			APIActiveFlag       int         `json:"api_active_flag"`
			APIStarttime        int64       `json:"api_starttime"`
			APILevel            int         `json:"api_level"`
			APIRank             int         `json:"api_rank"`
			APIExperience       int         `json:"api_experience"`
			APIFleetname        interface{} `json:"api_fleetname"`
			APIComment          string      `json:"api_comment"`
			APICommentID        string      `json:"api_comment_id"`
			APIMaxChara         int         `json:"api_max_chara"`
			APIMaxSlotitem      int         `json:"api_max_slotitem"`
			APIMaxKagu          int         `json:"api_max_kagu"`
			APIPlaytime         int         `json:"api_playtime"`
			APITutorial         int         `json:"api_tutorial"`
			APIFurniture        []int       `json:"api_furniture"`
			APICountDeck        int         `json:"api_count_deck"`
			APICountKdock       int         `json:"api_count_kdock"`
			APICountNdock       int         `json:"api_count_ndock"`
			APIFcoin            int         `json:"api_fcoin"`
			APIStWin            int         `json:"api_st_win"`
			APIStLose           int         `json:"api_st_lose"`
			APIMsCount          int         `json:"api_ms_count"`
			APIMsSuccess        int         `json:"api_ms_success"`
			APIPtWin            int         `json:"api_pt_win"`
			APIPtLose           int         `json:"api_pt_lose"`
			APIPtChallenged     int         `json:"api_pt_challenged"`
			APIPtChallengedWin  int         `json:"api_pt_challenged_win"`
			APIFirstflag        int         `json:"api_firstflag"`
			APITutorialProgress int         `json:"api_tutorial_progress"`
			APIPvp              []int       `json:"api_pvp"`
			APIMedals           int         `json:"api_medals"`
			APILargeDock        int         `json:"api_large_dock"`
		} `json:"api_basic"`
		APILog []struct {
			APINo      int    `json:"api_no"`
			APIType    string `json:"api_type"`
			APIState   string `json:"api_state"`
			APIMessage string `json:"api_message"`
		} `json:"api_log"`
		APIPBgmID             int `json:"api_p_bgm_id"`
		APIParallelQuestCount int `json:"api_parallel_quest_count"`
		APIDestShipSlot       int `json:"api_dest_ship_slot"`
		APICFlag              int `json:"api_c_flag"`
		APICFlag2             int `json:"api_c_flag2"`
	} `json:"api_data"`
}
