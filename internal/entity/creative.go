package entity

type Creative struct {
	ID          int     `json:"id"`
	CampaignID  int     `json:"campaign_id"`
	Price       float64 `json:"price"`
	Duration    int     `json:"duration"`
	PlaylistHLS string  `json:"playlist_hls"`
}
