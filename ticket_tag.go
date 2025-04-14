package zammad

import (
	"fmt"
	"net/http"
)

type object string

const (
	// ObjectTicket is the object type for tickets.
	ObjectTicket object = "Ticket"
)

type TicketTag struct {
	TicketID int    `json:"o_id"`
	Name     string `json:"item"`
	Object   object `json:"object"`
}
 
func (c *Client) AddTagToTicket(ticketID int, tag string) error {

	t := TicketTag{
		TicketID: ticketID,
		Name:     tag,
		Object:   ObjectTicket,
	}
	req, err := c.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", c.Url, "/api/v1/tags/add"), t)
	if err != nil {
		return err
	}

	if err = c.sendWithAuth(req, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) TicketTagByTicket(ticketID int) ([]Tag, error) {
	var tags struct {
		Tags []string
	}

	req, err := c.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", c.Url, fmt.Sprintf("/api/v1/tags?object=Ticket&o_id=%d", ticketID)), nil)
	if err != nil {
		return nil, err
	}

	if err = c.sendWithAuth(req, &tags); err != nil {
		return nil, err
	}

	tags1 := make([]Tag, len(tags.Tags))
	for i := range tags.Tags {
		tags1[i] = Tag{Name: tags.Tags[i]}
	}

	return tags1, nil
}
