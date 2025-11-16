package http

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/protobuf/chat/golang/conversationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Conversation struct {
	ID                  string
	ClientUserID        string
	LastMessageSnippet  string
	LastMessageSenderID string
	UpdatedAt           string
	Status              string
}

func (h Handler) ShowSupportPage(c echo.Context) error {
	ctx := grpcContext(c)
	user, _ := getUserFromContext(c)

	// Make a light request just to get the total count for the stat card
	newListReq := &conversationpb.ConversationListRequest{
		CurrentPage: 1,
		PageSize:    1, // We only need the total count
	}
	ownListReq := &conversationpb.ConversationListRequest{
		CurrentPage: 1,
		PageSize:    1, // We only need the total count
	}

	newListResp, err := h.conversationAd.ConversationNewList(ctx, newListReq)
	totalNew := uint64(0)
	if err != nil {
		h.logger.Error("Failed to get new conversation count for stats", "error", err)
	} else {
		totalNew = newListResp.GetTotalNumber()
	}

	ownListResp, err := h.conversationAd.ConversationOwnList(ctx, ownListReq)
	totalMine := uint64(0)
	if err != nil {
		h.logger.Error("Failed to get own conversation count for stats", "error", err)
	} else {
		totalMine = ownListResp.GetTotalNumber()
	}

	data := map[string]interface{}{
		"TotalNew":     totalNew,
		"TotalMine":    totalMine,
		"TemplateName": "support_page",
		"User":         user,
	}

	if isHTMX(c) {
		return c.Render(http.StatusOK, "support_page", data)
	}

	return c.Render(http.StatusOK, "main_layout", data)
}

func (h Handler) ListNewConversationsPartial(c echo.Context) error {
	ctx := grpcContext(c)

	page, _ := strconv.ParseUint(c.QueryParam("page"), 10, 64)
	if page == 0 {
		page = 1
	}

	listReq := &conversationpb.ConversationListRequest{
		CurrentPage: page,
		PageSize:    defaultPageSize,
	}

	listResp, err := h.conversationAd.ConversationNewList(ctx, listReq)
	if err != nil {
		return h.renderGRPCError(c, "ListNewConversationsPartial", err)
	}

	conversations := make([]Conversation, len(listResp.GetConversations()))
	for i, pbConv := range listResp.GetConversations() {
		conversations[i] = convertConversationPbToConversation(pbConv)
	}

	pagination := h.buildConversationPaginationData(listResp, c.Request().URL.Query(), "/support/new-list")

	data := map[string]interface{}{
		"Conversations": conversations,
		"Pagination":    pagination,
	}

	return c.Render(http.StatusOK, "support_new_list_partial", data)
}

func (h Handler) ListOwnConversationsPartial(c echo.Context) error {
	ctx := grpcContext(c)

	page, _ := strconv.ParseUint(c.QueryParam("page"), 10, 64)
	if page == 0 {
		page = 1
	}

	status := c.QueryParam("status")
	var statuses []conversationpb.Status
	if status == "STATUS_OPEN" {
		statuses = append(statuses, conversationpb.Status_STATUS_OPEN)
	} else if status == "STATUS_CLOSED" {
		statuses = append(statuses, conversationpb.Status_STATUS_CLOSED)
	}

	listReq := &conversationpb.ConversationListRequest{
		CurrentPage: page,
		PageSize:    defaultPageSize,
		Statuses:    statuses,
	}

	listResp, err := h.conversationAd.ConversationOwnList(ctx, listReq)
	if err != nil {
		return h.renderGRPCError(c, "ListOwnConversationsPartial", err)
	}

	conversations := make([]Conversation, len(listResp.GetConversations()))
	for i, pbConv := range listResp.GetConversations() {
		conversations[i] = convertConversationPbToConversation(pbConv)
	}

	pagination := h.buildConversationPaginationData(listResp, c.Request().URL.Query(), "/support/own-list")

	data := map[string]interface{}{
		"Conversations": conversations,
		"Pagination":    pagination,
	}

	return c.Render(http.StatusOK, "support_own_list_partial", data)
}

func convertConversationPbToConversation(pb *conversationpb.Conversation) Conversation {
	return Conversation{
		ID:                  pb.GetId(),
		ClientUserID:        pb.GetClientUserId(),
		LastMessageSnippet:  pb.GetLastMessageSnippet(),
		LastMessageSenderID: pb.GetLastMessageSenderId(),
		UpdatedAt:           formatTimestamp(pb.GetUpdatedAt()),
		Status:              pb.GetStatus().String(),
	}
}

func (h Handler) buildConversationPaginationData(resp *conversationpb.ConversationListResponse, query url.Values, baseURL string) PaginationData {
	p := PaginationData{
		TotalNumber: resp.GetTotalNumber(),
		CurrentPage: resp.GetCurrentPage(),
		TotalPage:   resp.GetTotalPage(),
		HasPrev:     resp.GetCurrentPage() > 1,
		HasNext:     resp.GetCurrentPage() < resp.GetTotalPage(),
		BaseURL:     baseURL,
		Status:      query.Get("status"),
	}
	if p.HasPrev {
		p.PrevPage = p.CurrentPage - 1
	}
	if p.HasNext {
		p.NextPage = p.CurrentPage + 1
	}

	buildURL := func(page int) string {
		v := url.Values{}
		v.Set("page", strconv.Itoa(page))
		if query.Get("status") != "" {
			v.Set("status", query.Get("status"))
		}
		return p.BaseURL + "?" + v.Encode()
	}

	maxPagesToShow := 5
	start := int(math.Max(1, float64(int(p.CurrentPage)-maxPagesToShow/2)))
	end := int(math.Min(float64(p.TotalPage), float64(start+maxPagesToShow-1)))

	if end-start+1 < maxPagesToShow && start > 1 {
		start = int(math.Max(1, float64(end-maxPagesToShow+1)))
	}

	p.Pages = []PaginationPage{}
	for i := start; i <= end; i++ {
		p.Pages = append(p.Pages, PaginationPage{
			PageNumber: i,
			URL:        buildURL(i),
			IsCurrent:  i == int(p.CurrentPage),
		})
	}

	return p
}

func formatTimestamp(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	t := ts.AsTime()
	now := time.Now()
	if now.Year() == t.Year() && now.Month() == t.Month() && now.Day() == t.Day() {
		return t.Format("3:04 PM")
	}
	if now.Year() == t.Year() && now.Month() == t.Month() && now.Day()-1 == t.Day() {
		return "Yesterday"
	}
	return t.Format("Jan 2")
}
