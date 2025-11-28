package http

import (
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/protobuf/chat/golang/conversationpb"
)

const defaultConvPageSize = 15 // Conversations are lighter, maybe show more

// ShowSupportService renders the main support page shell.
func (h Handler) ShowSupportService(c echo.Context) error {
	ctx := grpcContext(c)
	user, _ := getUserFromContext(c)

	listReq := &conversationpb.ConversationListRequest{
		CurrentPage: 1,
		PageSize:    1,
	}
	listResp, err := h.conversationSvc.ConversationNewList(ctx, listReq)
	totalNew := uint64(0)
	if err != nil {
		h.logger.Error("Failed to get new conversation count for stats", "error", err)
	} else {
		totalNew = listResp.GetTotalNumber()
	}

	accessToken, _ := getAccessTokenFromCookie(c, h.logger)

	chatWsURL := "ws://localhost:2530/chats/supports"

	data := map[string]interface{}{
		"TotalNewConversations": totalNew,
		"TemplateName":          "support_page",
		"User":                  user,
		"AllStatuses":           GetAllConversationStatuses(),
		"WebSocketToken":        accessToken,
		"ChatWsURL":             chatWsURL,
	}

	if isHTMX(c) {
		return c.Render(http.StatusOK, "support_page", data)
	}

	return c.Render(http.StatusOK, "main_layout", data)
}

// ListNewConversationsPartial renders the list of "New" conversations.
func (h Handler) ListNewConversationsPartial(c echo.Context) error {
	ctx := grpcContext(c)
	user, _ := getUserFromContext(c)

	page, _ := strconv.ParseUint(c.QueryParam("page"), 10, 64)
	if page == 0 {
		page = 1
	}

	// New conversations are always sorted DESC (newest first).
	listReq := &conversationpb.ConversationListRequest{
		CurrentPage:   page,
		PageSize:      defaultConvPageSize,
		SortDirection: conversationpb.SortDirection_SORT_DIRECTION_DESC,
	}

	listResp, err := h.conversationSvc.ConversationNewList(ctx, listReq)
	if err != nil {
		return h.renderGRPCError(c, "ListNewConversationsPartial", err)
	}

	conversations := make([]Conversation, len(listResp.GetConversations()))
	for i, pbConv := range listResp.GetConversations() {
		conversations[i] = convertConversationPbToConversation(pbConv)
	}

	// We build pagination data
	pagination := h.buildSupportPaginationData(listResp,
		"/support/list/new", // Base URL for this list.
		int(conversationpb.SortDirection_SORT_DIRECTION_DESC),
		nil, // No status filters for "New" list.
	)

	data := map[string]interface{}{
		"Conversations": conversations,
		"Pagination":    pagination,
		"CurrentUserID": user.ID,
	}

	// This partial is reusable for both lists.
	return c.Render(http.StatusOK, "support_list_partial", data)
}

// ListMyConversationsPartial renders the list of "My" (assigned) conversations.
func (h Handler) ListMyConversationsPartial(c echo.Context) error {
	ctx := grpcContext(c)
	user, _ := getUserFromContext(c)

	// 1. Parse Query Parameters.
	page, _ := strconv.ParseUint(c.QueryParam("page"), 10, 64)
	if page == 0 {
		page = 1
	}

	sortDirInt, _ := strconv.ParseInt(c.QueryParam("sort_direction"), 10, 32)
	if sortDirInt == 0 {
		sortDirInt = 2 // Default to DESC (Newest First).
	}

	// Get status filters.
	statusStrings := c.QueryParams()["statuses"]
	statuses := ParseStatusesFromForm(statusStrings)

	// If no statuses are provided, default to OPEN.
	if len(statuses) == 0 {
		statuses = []conversationpb.Status{conversationpb.Status_STATUS_OPEN}
		// We update statusStrings to reflect this default for pagination links.
		statusStrings = []string{"open"}
	}

	listReq := &conversationpb.ConversationListRequest{
		CurrentPage:   page,
		PageSize:      defaultConvPageSize,
		SortDirection: conversationpb.SortDirection(sortDirInt),
		Statuses:      statuses,
		// AssignedSupportID is handled by the gRPC server (ConversationOwnList).
	}

	listResp, err := h.conversationSvc.ConversationOwnList(ctx, listReq)
	if err != nil {
		return h.renderGRPCError(c, "ListMyConversationsPartial", err)
	}

	conversations := make([]Conversation, len(listResp.GetConversations()))
	for i, pbConv := range listResp.GetConversations() {
		conversations[i] = convertConversationPbToConversation(pbConv)
	}

	pagination := h.buildSupportPaginationData(listResp,
		"/support/list/my", // Base URL for this list.
		int(sortDirInt),
		statusStrings, // Pass filters to pagination builder.
	)

	data := map[string]interface{}{
		"Conversations": conversations,
		"Pagination":    pagination,
		"CurrentUserID": user.ID,
	}

	// Reusing the same partial
	return c.Render(http.StatusOK, "support_list_partial", data)
}

// buildSupportPaginationData is a helper to create the pagination struct for conversation lists.
// (Copied from user.go and adapted).
func (h Handler) buildSupportPaginationData(
	resp *conversationpb.ConversationListResponse,
	baseURL string,
	sortDir int,
	statusStrings []string,
) PaginationData {
	p := PaginationData{
		TotalNumber:   resp.GetTotalNumber(),
		CurrentPage:   resp.GetCurrentPage(),
		TotalPage:     resp.GetTotalPage(),
		HasPrev:       resp.GetCurrentPage() > 1,
		HasNext:       resp.GetCurrentPage() < resp.GetTotalPage(),
		BaseURL:       baseURL, // Use the provided base URL.
		SortDirection: sortDir,
		// Note: We don't use 'Query' (username) here, but we use 'Statuses'
		// We store the raw strings for building the URL
		// This is a bit of a hack; ideally PaginationData would be more generic.
		Roles: statusStrings, // Re-using the 'Roles' field for 'statuses'.
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
		v.Set("sort_direction", strconv.Itoa(p.SortDirection))
		for _, s := range p.Roles { // Re-using 'Roles' as 'statuses'.
			v.Add("statuses", s)
		}
		return p.BaseURL + "?" + v.Encode()
	}

	// Build page number links (e.g., show max 5 pages).
	maxPagesToShow := 5
	//nolint:gosec // G115: CurrentPage is pagination number, overflow is not a realistic concern
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
			//nolint:gosec // G115: CurrentPage is pagination number, overflow is not a realistic concern
			IsCurrent: i == int(p.CurrentPage),
		})
	}

	return p
}
