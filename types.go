package zammad

import (
	"fmt"
	"net/http"
	"time"
)

type (
	// Client is used to query Zammad. It is safe to use concurrently. If you (inadvertly) added
	// multiple authencation options that will be applied in the order, basic auth, token based, and
	// then oauth. Where the last one set, wins.
	Client struct {
		Client   Doer
		Username string
		Password string
		Token    string
		OAuth    string
		Url      string
		// FromFunc is used to set the From HTTP header, if you want to act on behalf of another user.
		// See https://docs.zammad.org/en/latest/api/intro.html#actions-on-behalf-of-other-users. If not nil
		// *and* returning a non empty string, this value will be used in the request.
		FromFunc func() string
	}

	// ErrorResponse is the response returned by Zammad when an error occured.
	ErrorResponse struct {
		Description      string `json:"error"`
		DescriptionHuman string `json:"error_human"`
	}

	// Doer is an interface that allows mimicking a *http.Client.
	Doer interface {
		Do(*http.Request) (*http.Response, error)
	}
)

type SearchTicket struct {
	ID                      int        `json:"id"`
	GroupID                 int        `json:"group_id"`
	PriorityID              int        `json:"priority_id"`
	StateID                 int        `json:"state_id"`
	Number                  string     `json:"number"`
	Title                   string     `json:"title"`
	OwnerID                 int        `json:"owner_id"`
	CustomerID              int        `json:"customer_id"`
	Note                    *string    `json:"note"`
	FirstResponseAt         *time.Time `json:"first_response_at"`
	CloseAt                 *time.Time `json:"close_at"`
	CloseInMin              *int       `json:"close_in_min"`
	CloseDiffInMin          *int       `json:"close_diff_in_min"`
	UpdatedAt               time.Time  `json:"updated_at"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedByID             int        `json:"updated_by_id"`
	CreatedByID             int        `json:"created_by_id"`
	LastOwnerUpdateAt       *time.Time `json:"last_owner_update_at"`
	LastContactAt           *time.Time `json:"last_contact_at"`
	LastContactAgentAt      *time.Time `json:"last_contact_agent_at"`
	LastContactCustomerAt   *time.Time `json:"last_contact_customer_at"`
	PendingTime             *time.Time `json:"pending_time"`
	EscalationAt            *time.Time `json:"escalation_at"`
	Type                    *string    `json:"type"`
	ArticleCount            int        `json:"article_count"`
	CreateArticleTypeID     int        `json:"create_article_type_id"`
	CreateArticleSenderID   int        `json:"create_article_sender_id"`
	ArticleIDs              []int      `json:"article_ids"`
	TicketTimeAccountingIDs []int      `json:"ticket_time_accounting_ids"`
	// Custom fields
	Measure                   *string     `json:"massnahme"`
	ProcessType               *string     `json:"vorgangsart"`
	CustomerName              *string     `json:"name"`
	SalesPipeline             *string     `json:"verkaufspipeline"`
	DamageCauser              *string     `json:"schadensverursacher"`
	GoogleProfileMessageType  *string     `json:"google_business_profile_message_type"`
	GoogleProfileStarRating   *string     `json:"google_business_profile_star_rating"`
	GoogleProfileReviewRating *string     `json:"google_business_profile_review_rating"`
	MoeveID                   *string     `json:"_kv_nummer"`
	Preferences               Preferences `json:"preferences"`
}

type Preferences struct {
	ChannelID int `json:"channel_id"`
}

type Asset struct {
	Tickets []SearchTicket `json:"Tickets"`
}

type TicketSearchResponse struct {
	Assets Asset `json:"assets"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprint(r.Description)
}
