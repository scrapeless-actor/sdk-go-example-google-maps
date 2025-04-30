package main

import (
	"fmt"
	"strconv"
)

const (
	GoogleMaps                   = "maps"
	GoogleMapsAutocomplete       = "autocomplete"
	GoogleMapsContributorReviews = "contributorreviews"
	GoogleMapsDirections         = "directions"
	GoogleMapsPhotoMeta          = "photometa"
	GoogleMapsPhotos             = "photos"
	GoogleMapsReviews            = "reviews"
)

type RequestParam struct {
	Q             string `json:"q"`
	Ll            string `json:"ll,omitempty"`
	GoogleDomain  string `json:"google_Domain,omitempty"`
	Hl            string `json:"hl,omitempty"`
	Gl            string `json:"gl,omitempty"`
	Data          string `json:"data,omitempty"`
	PlaceId       string `json:"place_id,omitempty"`
	Type          string `json:"type,omitempty"`
	Start         string `json:"start,omitempty"`
	Engine        string `json:"engine,omitempty"`
	Cp            string `json:"cp,omitempty"`
	ContributorId string `json:"contributor_id,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	Num           string `json:"num,omitempty"`
	SortBy        string `json:"sort_by,omitempty"`
	TopicId       string `json:"topic_id,omitempty"`
	DataId        string `json:"data_id,omitempty"`

	// google_maps_directions
	StartAddr    string `json:"start_addr,omitempty"`
	EndAddr      string `json:"end_addr,omitempty"`
	TravelMode   string `json:"travel_mode,omitempty"`
	DistanceUnit string `json:"distance_unit,omitempty"`
	Avoid        string `json:"avoid,omitempty"`
	Prefer       string `json:"prefer,omitempty"`
	Route        string `json:"route,omitempty"`
	Time         string `json:"time,omitempty"`

	// google_maps_photos
	CategoryId string `json:"category_id,omitempty"`

	// google_maps_photos_meat
	PhotoDataId string `json:"photo_data_id,omitempty"`
}

func (params *RequestParam) FieldValidation() error {
	if params.Engine == "" {
		return fmt.Errorf("engine is required")
	}
	if params.Gl == "" {
		params.Gl = "us"
	}
	if params.Hl == "" {
		params.Hl = "en"
	}
	if params.GoogleDomain == "" {
		params.GoogleDomain = "google.com"
	}
	switch params.Engine {
	case GoogleMaps:
		if params.PlaceId == "" {
			switch params.Type {
			case "":
				return fmt.Errorf("place id is required")
			case "search":
				if params.Q == "" {
					return fmt.Errorf("q is required")
				}
			case "place":
				if params.Data == "" {
					return fmt.Errorf("data is required")
				}
			default:
				return fmt.Errorf("type is not defined")
			}
		}

		if params.Start != "" {
			if params.Ll == "" {
				return fmt.Errorf("missing query `ll` parameter. Required when using pagination")
			}
			twenty := isMultipleOfTwenty(params.Start)
			if !twenty {
				return fmt.Errorf("`start` is not a multiple of twenty")
			}
		}
	case GoogleMapsAutocomplete:
		if params.Q == "" {
			return fmt.Errorf("q is required")
		}
		if params.Ll == "" {
			return fmt.Errorf("ll is required")
		}
	case GoogleMapsContributorReviews:
		if params.ContributorId == "" {
			return fmt.Errorf("contributor_id is required")
		}
	case GoogleMapsReviews:
		if params.DataId == "" && params.PlaceId == "" {
			return fmt.Errorf("data_id or place_id is required")
		}
		if params.TopicId == "" && params.NextPageToken == "" && params.Num != "" {
			return fmt.Errorf("`num` parameter should not be used on the initial page when neither `next_page_token` nor `topic_id` is set. It always returns 8 results")
		}
	case GoogleMapsDirections:
		if params.TravelMode != "3" {
			if params.Prefer != "" {
				params.Prefer = ""
			}
			if params.Route != "" {
				params.Route = ""
			}
			if params.Time == "last_available" {
				params.Time = ""
			}
		}
	case GoogleMapsPhotos:
		if params.DataId == "" {
			return fmt.Errorf("data_id is required")
		}
	case GoogleMapsPhotoMeta:
		if params.PhotoDataId == "" {
			return fmt.Errorf("photo_data_id is required")
		}
	default:
		return fmt.Errorf("unknow engine type")
	}
	return nil
}

func isMultipleOfTwenty(s string) bool {
	num, err := strconv.Atoi(s)
	if err != nil {
		// 如果转换失败，返回 false
		return false
	}
	return num%20 == 0
}

type Response struct {
	LocalResults      []LocalResults     `json:"local_results,omitempty"`
	PlaceResults      *PlaceResults      `json:"place_results,omitempty"`
	SearchInformation *SearchInformation `json:"search_information,omitempty"`
	Suggestions       []Suggestions      `json:"suggestions,omitempty"`
	Contributor       *Contributor       `json:"contributor,omitempty"`
	Reviews           []Reviews          `json:"reviews,omitempty"`
	PlaceInfo         *PlaceInfo         `json:"place_info,omitempty"`
	Topics            []Topics           `json:"topics,omitempty"`
	Directions        []Directions       `json:"directions,omitempty"`
	Durations         []Durations        `json:"durations,omitempty"`
	PlaceInfos        []PlaceInfo        `json:"place_infos,omitempty"`
	User              *User              `json:"user,omitempty"`
	Location          *Location          `json:"location,omitempty"`
	Date              string             `json:"date,omitempty"`
	Categories        []Categories       `json:"categories,omitempty"`
	Photos            []Photos           `json:"photos,omitempty"`
	NextPageToken     string             `json:"next_page_token,omitempty"`
}

type PlaceResults struct {
	Title                 string                   `json:"title,omitempty"`
	PlaceId               string                   `json:"place_id,omitempty"`
	DataId                string                   `json:"data_id,omitempty"`
	DataCid               string                   `json:"data_cid,omitempty"`
	ReviewsLink           string                   `json:"reviews_link,omitempty"`
	PhotosLink            string                   `json:"photos_link,omitempty"`
	GpsCoordinates        *GpsCoordinates          `json:"gps_coordinates,omitempty"`
	PlaceIdSearch         string                   `json:"place_id_search,omitempty"`
	ProviderId            string                   `json:"provider_id,omitempty"`
	Thumbnail             string                   `json:"thumbnail,omitempty"`
	RatingSummary         []RatingSummary          `json:"rating_summary,omitempty"`
	Rating                float64                  `json:"rating,omitempty"`
	Reviews               int                      `json:"reviews,omitempty"`
	Price                 string                   `json:"price,omitempty"`
	Type                  []string                 `json:"type,omitempty"`
	TypeIds               []string                 `json:"type_ids,omitempty"`
	Description           string                   `json:"description,omitempty"`
	Menu                  *Menu                    `json:"menu,omitempty"`
	UnclaimedListing      bool                     `json:"unclaimed_listing,omitempty"`
	OrderOnlineLink       string                   `json:"order_online_link,omitempty"`
	ServiceOptions        *ServiceOptions          `json:"service_options,omitempty"`
	Address               string                   `json:"address,omitempty"`
	Website               string                   `json:"website,omitempty"`
	Phone                 string                   `json:"phone,omitempty"`
	OpenState             string                   `json:"open_state,omitempty"`
	PlusCode              string                   `json:"plus_code,omitempty"`
	WebResultsLink        string                   `json:"web_results_link,omitempty"`
	Hours                 []map[string]interface{} `json:"hours,omitempty"`
	Images                []Image                  `json:"images,omitempty"`
	UserReviews           *UserReview              `json:"user_reviews,omitempty"`
	PopularTimes          *PopularTimes            `json:"popular_times,omitempty"`
	PeopleAlsoSearchFor   []PeopleAlsoSearchFor    `json:"people_also_search_for,omitempty"`
	Extensions            []map[string]interface{} `json:"extensions,omitempty"`
	UnsupportedExtensions []map[string]interface{} `json:"unsupported_extensions,omitempty"`
	QuestionsAndAnswers   []QuestionsAndAnswers    `json:"questions_and_answers,omitempty"`
}

func (p *PlaceResults) IsEmpty() bool {
	return p.Title == "" &&
		p.PlaceId == "" &&
		p.DataId == "" &&
		p.DataCid == "" &&
		p.ReviewsLink == "" &&
		p.PhotosLink == "" &&
		p.GpsCoordinates == nil &&
		p.PlaceIdSearch == "" &&
		p.ProviderId == "" &&
		p.Thumbnail == "" &&
		p.Rating == 0.0 &&
		p.Reviews == 0 &&
		p.Price == "" &&
		len(p.Type) == 0 &&
		len(p.TypeIds) == 0 &&
		p.Description == "" &&
		p.Menu == nil &&
		p.OrderOnlineLink == "" &&
		p.ServiceOptions == nil &&
		p.Address == "" &&
		p.Website == "" &&
		p.Phone == "" &&
		p.OpenState == "" &&
		p.PlusCode == "" &&
		p.Hours == nil &&
		len(p.Images) == 0 &&
		p.UserReviews == nil &&
		p.PopularTimes == nil &&
		len(p.PeopleAlsoSearchFor) == 0
}

type QuestionsAndAnswers struct {
	Question     *Question `json:"question,omitempty"`
	Answer       []Answer  `json:"answer,omitempty"`
	TotalAnswers int       `json:"total_answers,omitempty"`
}

func (p *QuestionsAndAnswers) IsEmpty() bool {
	return p.Question == nil && p.Answer == nil && p.TotalAnswers == 0
}

type Question struct {
	Text     string `json:"text,omitempty"`
	Data     string `json:"data,omitempty"`
	Language string `json:"language,omitempty"`
}

func (q *Question) IsEmpty() bool {
	return q.Text == "" && q.Data == "" && q.Language == ""
}

type Answer struct {
	Text     string `json:"text,omitempty"`
	Data     string `json:"data,omitempty"`
	Language string `json:"language,omitempty"`
}

func (q *Answer) IsEmpty() bool {
	return q.Text == "" && q.Data == "" && q.Language == ""
}

type RatingSummary struct {
	Starts int `json:"starts,omitempty"`
	Amount int `json:"amount,omitempty"`
}

func (r *RatingSummary) IsEmpty() bool {
	return r.Amount == 0
}

type PeopleAlsoSearchFor struct {
	SearchTerm   string        `json:"search_term,omitempty"`
	LocalResults []LocalResult `json:"local_results,omitempty"`
}

func (p *PeopleAlsoSearchFor) IsEmpty() bool {
	return p.SearchTerm == "" && p.LocalResults == nil
}

type LocalResult struct {
	Position       int             `json:"position,omitempty"`
	Title          string          `json:"title,omitempty"`
	DataID         string          `json:"data_id,omitempty"`
	DataCID        string          `json:"data_cid,omitempty"`
	ReviewsLink    string          `json:"reviews_link,omitempty"`
	PhotosLink     string          `json:"photos_link,omitempty"`
	GPSCoordinates *GpsCoordinates `json:"gps_coordinates,omitempty"`
	PlaceIDSearch  string          `json:"place_id_search,omitempty"`
	Rating         float64         `json:"rating,omitempty"`
	Reviews        int             `json:"reviews,omitempty"`
	Thumbnail      string          `json:"thumbnail,omitempty"`
	Type           []string        `json:"type,omitempty"`
}

func (lr LocalResult) IsEmpty() bool {
	if lr.Title != "" ||
		lr.DataID != "" ||
		lr.DataCID != "" ||
		lr.ReviewsLink != "" ||
		lr.PhotosLink != "" ||
		lr.GPSCoordinates != nil ||
		lr.PlaceIDSearch != "" ||
		lr.Rating != 0.0 ||
		lr.Reviews != 0 ||
		lr.Thumbnail != "" ||
		len(lr.Type) > 0 {
		return false
	}
	return true
}

type PopularTimes struct {
	GraphResults map[string]interface{}
	LiveHash     *LiveHash `json:"live_hash,omitempty"`
}

func (p *PopularTimes) IsEmpty() bool {
	return p.GraphResults == nil && p.LiveHash == nil
}

type GraphResultsInfo struct {
	Time          string `json:"time,omitempty"`
	Info          string `json:"info,omitempty"`
	BusynessScore int    `json:"busyness_score"`
}

func (g *GraphResultsInfo) IsEmpty() bool {
	return g.Info == "" && g.BusynessScore == 0 && g.Time == ""
}

type LiveHash struct {
	Info      string `json:"info,omitempty"`
	TimeSpent string `json:"time_spent,omitempty"`
}

func (li *LiveHash) IsEmpty() bool {
	return li.Info == "" && li.TimeSpent == ""
}

type UserReview struct {
	MostRelevant []MostRelevant `json:"most_relevant,omitempty"`
	Summary      []Summary      `json:"summary,omitempty"`
}

func (u *UserReview) IsEmpty() bool {
	return len(u.MostRelevant) == 0 && len(u.Summary) == 0
}

type Summary struct {
	Snippet string `json:"snippet,omitempty"`
}

func (s *Summary) IsEmpty() bool {
	return s.Snippet == ""
}

type MostRelevant struct {
	Username      string  `json:"username,omitempty"`
	Rating        int     `json:"rating,omitempty"`
	ContributorId string  `json:"contributor_id,omitempty"`
	Description   string  `json:"description,omitempty"`
	Date          string  `json:"date,omitempty"`
	Link          string  `json:"link,omitempty"`
	Images        []Image `json:"images,omitempty"`
}

func (m *MostRelevant) IsEmpty() bool {
	return m.Username == "" && m.Rating == 0 && m.ContributorId == "" && m.Description == "" && m.Date == "" && len(m.Images) == 0
}

type Image struct {
	Title     string `json:"title,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

func (i *Image) IsEmpty() bool {
	return i.Title == "" && i.Thumbnail == ""
}

type Menu struct {
	Link   string `json:"link,omitempty"`
	Source string `json:"source,omitempty"`
}

func (m *Menu) IsEmpty() bool {
	return m.Link == "" && m.Source == ""
}

type ServiceOptions struct {
	DineIn   bool `json:"dine_in"`
	Takeout  bool `json:"takeout"`
	Delivery bool `json:"delivery"`
}

func (s *ServiceOptions) IsEmpty() bool {
	return s.DineIn == false && s.Takeout == false && s.Delivery == false
}

type Photos struct {
	Thumbnail   string `json:"thumbnail,omitempty"`
	Image       string `json:"image,omitempty"`
	Video       string `json:"video,omitempty"`
	PhotoDataId string `json:"photo_data_id,omitempty"`
	DataId      string `json:"data_id,omitempty"`
	User        *User  `json:"user,omitempty"`
}

func (p *Photos) IsEmpty() bool {
	return p.Thumbnail == "" && p.Image == "" && p.Video == "" && p.PhotoDataId == "" && p.User == nil
}

type Categories struct {
	Title string `json:"title,omitempty"`
	Id    string `json:"id,omitempty"`
}

func (c *Categories) IsEmpty() bool {
	return c.Title == "" && c.Id == ""
}

type Location struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Name      string  `json:"name,omitempty"`
	Type      string  `json:"type,omitempty"`
}

func (l *Location) IsEmpty() bool {
	return l.Latitude == 0 && l.Longitude == 0 && l.Name == "" && l.Type == ""
}

type Topics struct {
	Keyword  string `json:"keyword,omitempty"`
	Mentions int    `json:"mentions,omitempty"`
	Id       string `json:"id,omitempty"`
}

func (t *Topics) IsEmpty() bool {
	return t.Keyword == "" && t.Mentions == 0 && t.Id == ""
}

type SearchInformation struct {
	LocalResultsState string `json:"local_Results_State,omitempty"`
	QueryDisplayed    string `json:"query_Displayed,omitempty"`
}

type LocalResults struct {
	Position       int                    `json:"position,omitempty"`
	Title          string                 `json:"title,omitempty"`
	PlaceID        string                 `json:"place_id,omitempty"`
	DataID         string                 `json:"data_id,omitempty"`
	DataCid        string                 `json:"data_cid,omitempty"`
	ReviewsLink    string                 `json:"reviews_link,omitempty"`
	PhotosLink     string                 `json:"photos_link,omitempty"`
	GpsCoordinates *GpsCoordinates        `json:"gps_coordinates,omitempty"`
	PlaceIDSearch  string                 `json:"place_id_search,omitempty"`
	ProviderID     string                 `json:"provider_id,omitempty"`
	Rating         float64                `json:"rating,omitempty"`
	Reviews        int                    `json:"reviews,omitempty"`
	Price          string                 `json:"price,omitempty"`
	Type           string                 `json:"type,omitempty"`
	Types          []string               `json:"types,omitempty"`
	TypeID         string                 `json:"type_id,omitempty"`
	TypeIds        []string               `json:"type_ids,omitempty"`
	Address        string                 `json:"address,omitempty"`
	OpenState      string                 `json:"open_state,omitempty"`
	Hours          string                 `json:"hours,omitempty"`
	OperatingHours map[string]interface{} `json:"operating_hours,omitempty"`
	Phone          string                 `json:"phone,omitempty"`
	Website        string                 `json:"website,omitempty"`
	Description    string                 `json:"description,omitempty"`
	ServiceOptions []string               `json:"service_options,omitempty"`
	OrderOnline    string                 `json:"order_online,omitempty"`
	Thumbnail      string                 `json:"thumbnail,omitempty"`
}

func (lr *LocalResults) IsEmpty() bool {
	return lr.Title == "" &&
		lr.PlaceID == "" &&
		lr.DataID == "" &&
		lr.DataCid == "" &&
		lr.ReviewsLink == "" &&
		lr.PhotosLink == "" &&
		lr.GpsCoordinates == nil &&
		lr.PlaceIDSearch == "" &&
		lr.ProviderID == "" &&
		lr.Rating == 0 &&
		lr.Reviews == 0 &&
		lr.Price == "" &&
		lr.Type == "" &&
		len(lr.Types) == 0 &&
		lr.TypeID == "" &&
		len(lr.TypeIds) == 0 &&
		lr.Address == "" &&
		lr.OpenState == "" &&
		lr.Hours == "" &&
		lr.OperatingHours == nil || len(lr.OperatingHours) == 0 &&
		lr.Phone == "" &&
		lr.Website == "" &&
		lr.Description == "" &&
		len(lr.ServiceOptions) == 0 &&
		lr.OrderOnline == "" &&
		lr.Thumbnail == ""
}

type GpsCoordinates struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

func (g *GpsCoordinates) IsEmpty() bool {
	return g.Latitude == 0.0 && g.Longitude == 0.0
}

type Suggestions struct {
	Value       string  `json:"value,omitempty"`
	MapsLink    string  `json:"maps_link,omitempty"`
	Subtext     string  `json:"subtext,omitempty"`
	Type        string  `json:"type,omitempty"`
	Latitude    float64 `json:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty"`
	DataID      string  `json:"data_id,omitempty"`
	ReviewsLink string  `json:"reviews_link,omitempty"`
	PhotosLink  string  `json:"photos_link,omitempty"`
	ProviderId  string  `json:"provider_Id,omitempty"`
}

func (s *Suggestions) IsEmpty() bool {
	if s == nil {
		return true
	}
	return s.Value == "" &&
		s.MapsLink == "" &&
		s.Subtext == "" &&
		s.Type == "" &&
		s.Latitude == 0 &&
		s.Longitude == 0 &&
		s.DataID == "" &&
		s.ReviewsLink == "" &&
		s.PhotosLink == "" &&
		s.ProviderId == ""
}

type Contributor struct {
	Name          string                 `json:"name,omitempty"`
	Thumbnail     string                 `json:"thumbnail,omitempty"`
	Points        int                    `json:"points,omitempty"`
	Level         int                    `json:"level,omitempty"`
	Contributions map[string]interface{} `json:"contributions,omitempty"`
}

func (c *Contributor) IsEmpty() bool {
	return c.Name == "" &&
		c.Thumbnail == "" &&
		c.Points == 0 &&
		c.Level == 0 &&
		// !c.LocalGuide &&
		len(c.Contributions) == 0
}

type Reviews struct {
	PlaceInfo         *PlaceInfo             `json:"place_info,omitempty"`
	Date              string                 `json:"date,omitempty"`
	Snippet           string                 `json:"snippet,omitempty"`
	TranslatedSnippet string                 `json:"translated_snippet,omitempty"`
	ReviewID          string                 `json:"review_id,omitempty"`
	Rating            float64                `json:"rating"`
	Link              string                 `json:"link,omitempty"`
	Likes             int                    `json:"likes"`
	Details           map[string]interface{} `json:"details,omitempty"`
	TranslatedDetails map[string]interface{} `json:"translated_details,omitempty"`
	Images            []ReviewImage          `json:"images,omitempty"`
	Response          *ReviewResponse        `json:"response,omitempty"`
	User              *User                  `json:"user,omitempty"`
	IsoDate           string                 `json:"iso_date,omitempty"`
	IsoDateOfLastEdit string                 `json:"iso_date_of_last_edit,omitempty"`
	Source            string                 `json:"source,omitempty"`
}

func (r *Reviews) IsEmpty() bool {
	return r.Date == "" &&
		r.Snippet == "" &&
		r.TranslatedSnippet == "" &&
		r.ReviewID == "" &&
		r.Rating == 0 &&
		r.Link == "" &&
		r.Likes == 0 &&
		len(r.Details) == 0 &&
		len(r.TranslatedDetails) == 0 &&
		len(r.Images) == 0 &&
		r.IsoDate == "" &&
		r.IsoDateOfLastEdit == "" &&
		r.Source == ""
}

type User struct {
	Name          string `json:"name,omitempty"`
	Link          string `json:"link,omitempty"`
	ContributorID string `json:"contributor_id,omitempty"`
	Thumbnail     string `json:"thumbnail,omitempty"`
	Reviews       int    `json:"reviews,omitempty"`
	Photos        int    `json:"photos,omitempty"`
	Image         string `json:"image,omitempty"`
	UserId        string `json:"user_id,omitempty"`
}

func (u *User) IsEmpty() bool {
	return u.Name == "" &&
		u.Link == "" &&
		u.ContributorID == "" &&
		u.Thumbnail == "" &&
		// !u.LocalGuide &&
		u.Reviews == 0 &&
		u.Photos == 0 &&
		u.Image == "" &&
		u.UserId == ""
}

type PlaceInfo struct {
	Title     string          `json:"title,omitempty"`
	Address   string          `json:"address,omitempty"`
	GpsCoords *GpsCoordinates `json:"gps_coordinates,omitempty"`
	Thumbnail string          `json:"thumbnail,omitempty"`
	DataID    string          `json:"data_id,omitempty"`
	Type      string          `json:"type,omitempty"`
	Ratings   float64         `json:"ratings,omitempty"`
	Reviews   int             `json:"reviews,omitempty"`
}

func (p *PlaceInfo) IsEmpty() bool {
	return p.Title == "" &&
		p.Address == "" &&
		p.Thumbnail == "" &&
		p.DataID == "" &&
		p.Type == "" &&
		p.Ratings == 0.0 &&
		p.Reviews == 0
}

type ReviewImage struct {
	Title     string `json:"title,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Date      string `json:"date,omitempty"`
	Snippet   string `json:"snippet,omitempty"`
	Video     string `json:"video,omitempty"`
}

func (r *ReviewImage) IsEmpty() bool {
	return r.Title == "" &&
		r.Thumbnail == "" &&
		r.Date == "" &&
		r.Snippet == "" &&
		r.Video == ""
}

type ReviewResponse struct {
	Date              string `json:"date,omitempty"`
	Snippet           string `json:"snippet,omitempty"`
	TranslatedSnippet string `json:"translated_snippet,omitempty"`
}

func (r *ReviewResponse) IsEmpty() bool {
	return r.Date == "" &&
		r.Snippet == "" &&
		r.TranslatedSnippet == ""
}

type Directions struct {
	TravelMode            string            `json:"travel_mode,omitempty"`
	Via                   string            `json:"via,omitempty"`
	StartTime             string            `json:"start_time,omitempty"`
	EndTime               string            `json:"end_time,omitempty"`
	Distance              int               `json:"distance,omitempty"`
	Duration              int               `json:"duration,omitempty"`
	ArriveAround          int               `json:"arrive_around,omitempty"`
	LeaveAround           int               `json:"leave_around,omitempty"`
	TypicalDurationRange  string            `json:"typical_duration_range,omitempty"`
	FormattedDistance     string            `json:"formatted_distance,omitempty"`
	FormattedDuration     string            `json:"formatted_duration,omitempty"`
	FormattedArriveAround string            `json:"formatted_arrive_around,omitempty"`
	FormattedLeaveAround  string            `json:"formatted_leave_around,omitempty"`
	Cost                  int               `json:"cost,omitempty"`
	Currency              string            `json:"currency,omitempty"`
	Extensions            []string          `json:"extensions,omitempty"`
	ElevationProfile      *ElevationProfile `json:"elevation_profile,omitempty"`
	Flight                *FlightDetails    `json:"flight,omitempty"`
	Trips                 []Trip            `json:"trips,omitempty"`
}

func (d *Directions) IsEmpty() bool {
	return d.TravelMode == "" &&
		d.Via == "" &&
		d.StartTime == "" &&
		d.EndTime == "" &&
		d.Distance == 0 &&
		d.Duration == 0 &&
		d.ArriveAround == 0 &&
		d.LeaveAround == 0 &&
		d.TypicalDurationRange == "" &&
		d.FormattedDistance == "" &&
		d.FormattedDuration == "" &&
		d.FormattedArriveAround == "" &&
		d.FormattedLeaveAround == "" &&
		d.Cost == 0 &&
		d.Currency == "" &&
		len(d.Extensions) == 0 &&
		len(d.Trips) == 0
}

type ElevationProfile struct {
	Ascent               int    `json:"ascent,omitempty"`
	Descent              int    `json:"descent,omitempty"`
	MaxAltitude          int    `json:"max_altitude,omitempty"`
	MinAltitude          int    `json:"min_altitude,omitempty"`
	FormattedAscent      string `json:"formatted_ascent,omitempty"`
	FormattedDescent     string `json:"formatted_descent,omitempty"`
	FormattedMaxAltitude string `json:"formatted_max_altitude,omitempty"`
	FormattedMinAltitude string `json:"formatted_min_altitude,omitempty"`
}

func (e *ElevationProfile) IsEmpty() bool {
	return e.Ascent == 0 &&
		e.Descent == 0 &&
		e.MaxAltitude == 0 &&
		e.MinAltitude == 0 &&
		e.FormattedAscent == "" &&
		e.FormattedDescent == "" &&
		e.FormattedMaxAltitude == "" &&
		e.FormattedMinAltitude == ""
}

type FlightDetails struct {
	Departure                   string   `json:"departure,omitempty"`
	Arrival                     string   `json:"arrival,omitempty"`
	Date                        string   `json:"date,omitempty"`
	RoundTripPrice              int      `json:"round_trip_price,omitempty"`
	Currency                    string   `json:"currency,omitempty"`
	Airlines                    []string `json:"airlines,omitempty"`
	NonstopDuration             string   `json:"nonstop_duration,omitempty"`
	FormattedNonstopDuration    string   `json:"formatted_nonstop_duration,omitempty"`
	ConnectingDuration          string   `json:"connecting_duration,omitempty"`
	FormattedConnectingDuration string   `json:"formatted_connecting_duration,omitempty"`
	GoogleFlightsLink           string   `json:"google_flights_link,omitempty"`
}

func (f *FlightDetails) IsEmpty() bool {
	return f.Departure == "" &&
		f.Arrival == "" &&
		f.Date == "" &&
		f.RoundTripPrice == 0 &&
		f.Currency == "" &&
		len(f.Airlines) == 0 &&
		f.NonstopDuration == "" &&
		f.FormattedNonstopDuration == "" &&
		f.ConnectingDuration == "" &&
		f.FormattedConnectingDuration == "" &&
		f.GoogleFlightsLink == ""
}

type DirectionDetail struct {
	Title             string          `json:"title,omitempty"`
	Action            string          `json:"action,omitempty"`
	Distance          int             `json:"distance,omitempty"`
	Duration          int             `json:"duration,omitempty"`
	FormattedDistance string          `json:"formatted_distance,omitempty"`
	FormattedDuration string          `json:"formatted_duration,omitempty"`
	GeoPhoto          string          `json:"geo_photo,omitempty"`
	GPSCoordinates    *GpsCoordinates `json:"gps_coordinates,omitempty"`
	Extensions        []string        `json:"extensions,omitempty"`
}

func (d *DirectionDetail) IsEmpty() bool {
	return d.Title == "" &&
		d.Action == "" &&
		d.Distance == 0 &&
		d.Duration == 0 &&
		d.FormattedDistance == "" &&
		d.FormattedDuration == "" &&
		d.GeoPhoto == "" &&
		len(d.Extensions) == 0
}

type ServiceRunBy struct {
	Name             string `json:"name,omitempty"`
	Link             string `json:"link,omitempty"`
	RouteInformation string `json:"route_information,omitempty"`
}

func (s *ServiceRunBy) IsEmpty() bool {
	return s.Name == "" &&
		s.Link == "" &&
		s.RouteInformation == ""
}

type Trip struct {
	TravelMode        string            `json:"travel_mode,omitempty"`
	Title             string            `json:"title,omitempty"`
	Distance          int               `json:"distance,omitempty"`
	Duration          int               `json:"duration,omitempty"`
	FormattedDistance string            `json:"formatted_distance,omitempty"`
	FormattedDuration string            `json:"formatted_duration,omitempty"`
	StartStop         *Stops            `json:"start_stop,omitempty"`
	EndStop           *Stops            `json:"end_stop,omitempty"`
	Stops             []Stops           `json:"stops,omitempty"`
	ServiceRunBy      *ServiceRunBy     `json:"service_run_by,omitempty"`
	Details           []DirectionDetail `json:"details,omitempty"`
}

func (t *Trip) IsEmpty() bool {
	return t.TravelMode == "" &&
		t.Title == "" &&
		t.Distance == 0 &&
		t.Duration == 0 &&
		t.FormattedDistance == "" &&
		t.FormattedDuration == "" &&
		len(t.Stops) == 0 &&
		len(t.Details) == 0
}

type Stops struct {
	Name   string `json:"name,omitempty"`
	StopId string `json:"stop_id,omitempty"`
	Time   string `json:"time,omitempty"`
	DataId string `json:"data_id,omitempty"`
}

func (s *Stops) IsEmpty() bool {
	return s.Name == "" && s.StopId == "" && s.Time == "" && s.DataId == ""
}

type Durations struct {
	TravelMode        string `json:"travel_mode,omitempty"`
	Duration          int    `json:"duration,omitempty"`
	FormattedDuration string `json:"formatted_duration,omitempty"`
}

func (d *Durations) IsEmpty() bool {
	return d.TravelMode == "" &&
		d.Duration == 0 &&
		d.FormattedDuration == ""
}
