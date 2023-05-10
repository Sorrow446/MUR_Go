package main

type Meta struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   struct {
		Offset  int `json:"offset"`
		Limit   int `json:"limit"`
		Total   int `json:"total"`
		Count   int `json:"count"`
		Results []struct {
			ID        int `json:"id"`
			IssueMeta struct {
				ID                          int    `json:"id"`
				CatalogID                   int    `json:"catalog_id"`
				Title                       string `json:"title"`
				SeriesID                    int    `json:"series_id"`
				SeriesTitle                 string `json:"series_title"`
				ReleaseDate                 string `json:"release_date"`
				ReleaseDateFormatted        string `json:"release_date_formatted"`
				ReleaseDateDigital          string `json:"release_date_digital"`
				ReleaseDateDigitalFormatted string `json:"release_date_digital_formatted"`
				Description                 string `json:"description"`
				Format                      string `json:"format"`
				Thumbnail                   struct {
					Path      string `json:"path"`
					Extension string `json:"extension"`
				} `json:"thumbnail"`
				Creators struct {
					StringList        string `json:"string_list"`
					CreatorsAnalytics string `json:"creators_analytics"`
					ExtendedList      []struct {
						ID       int    `json:"id"`
						FullName string `json:"full_name"`
						Role     string `json:"role"`
					} `json:"extended_list"`
					ExtendedListRoles []string `json:"extended_list_roles"`
				} `json:"creators"`
			} `json:"issue_meta"`
			Pages []struct {
				ID        int    `json:"id"`
				BookID    int    `json:"book_id"`
				ColorFile string `json:"color_file"`
				Type      string `json:"type"`
				Height    int    `json:"height"`
				Width     int    `json:"width"`
				IsLegacy  bool   `json:"is_legacy"`
				IsAd      bool   `json:"is_ad"`
				Sequence  int    `json:"sequence"`
				Panels    []struct {
					ID                 int    `json:"id"`
					Type               string `json:"type"`
					Sequence           int    `json:"sequence"`
					X                  int    `json:"x"`
					Y                  int    `json:"y"`
					Height             int    `json:"height"`
					Width              int    `json:"width"`
					MaskColor          string `json:"mask_color"`
					TransitionsForward []struct {
						Type      string  `json:"type"`
						Duration  float64 `json:"duration"`
						Direction string  `json:"direction"`
					} `json:"transitions_forward"`
				} `json:"panels"`
			} `json:"pages"`
			PrevNextIssue struct {
				PrevIssueMeta struct {
					ID                          int    `json:"id"`
					CatalogID                   int    `json:"catalog_id"`
					Title                       string `json:"title"`
					SeriesID                    int    `json:"series_id"`
					SeriesTitle                 string `json:"series_title"`
					ReleaseDate                 string `json:"release_date"`
					ReleaseDateFormatted        string `json:"release_date_formatted"`
					ReleaseDateDigital          string `json:"release_date_digital"`
					ReleaseDateDigitalFormatted string `json:"release_date_digital_formatted"`
					Description                 string `json:"description"`
					Format                      string `json:"format"`
					Thumbnail                   struct {
						Path      string `json:"path"`
						Extension string `json:"extension"`
					} `json:"thumbnail"`
					Creators struct {
						StringList        string `json:"string_list"`
						CreatorsAnalytics string `json:"creators_analytics"`
						ExtendedList      []struct {
							ID       int    `json:"id"`
							FullName string `json:"full_name"`
							Role     string `json:"role"`
						} `json:"extended_list"`
						ExtendedListRoles []string `json:"extended_list_roles"`
					} `json:"creators"`
				} `json:"prev_issue_meta"`
				UpcomingIssueMeta struct {
					ID                   int    `json:"id"`
					CatalogID            int    `json:"catalog_id"`
					Title                string `json:"title"`
					SeriesID             int    `json:"series_id"`
					SeriesTitle          string `json:"series_title"`
					ReleaseDate          string `json:"release_date"`
					ReleaseDateFormatted string `json:"release_date_formatted"`
					Description          string `json:"description"`
					Format               string `json:"format"`
					Thumbnail            struct {
						Path      string `json:"path"`
						Extension string `json:"extension"`
					} `json:"thumbnail"`
					Creators struct {
						StringList        string `json:"string_list"`
						CreatorsAnalytics string `json:"creators_analytics"`
						ExtendedList      []struct {
							ID       int    `json:"id"`
							FullName string `json:"full_name"`
							Role     string `json:"role"`
						} `json:"extended_list"`
						ExtendedListRoles []string `json:"extended_list_roles"`
					} `json:"creators"`
				} `json:"upcoming_issue_meta"`
			} `json:"prev_next_issue"`
			ExternalLinks struct {
				Login      string `json:"login"`
				Checkout   string `json:"checkout"`
				Gateway    string `json:"gateway"`
				MarvelSite string `json:"marvel_site"`
			} `json:"external_links"`
			IsFree bool `json:"is_free"`
		} `json:"results"`
	} `json:"data"`
}

type AssetMeta struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   struct {
		Offset  int `json:"offset"`
		Limit   int `json:"limit"`
		Total   int `json:"total"`
		Count   int `json:"count"`
		Results []struct {
			ID        int `json:"id"`
			AuthState struct {
				LoggedIn   bool `json:"logged_in"`
				Subscriber bool `json:"subscriber"`
				UserID     int  `json:"user_id"`
			} `json:"auth_state"`
			Pages []struct {
				ID     int `json:"id"`
				Assets struct {
					Source    string `json:"source"`
					Thumbnail string `json:"thumbnail"`
					Uncanny   string `json:"uncanny"`
				} `json:"assets"`
			} `json:"pages"`
		} `json:"results"`
	} `json:"data"`
}
