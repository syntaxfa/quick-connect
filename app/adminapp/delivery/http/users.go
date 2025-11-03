package http

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userpb"
)

const defaultPageSize = 10

// PaginationPage struct for template page loops.
type PaginationPage struct {
	PageNumber int
	URL        string
	IsCurrent  bool
}

// PaginationData struct for template.
type PaginationData struct {
	TotalNumber   uint64
	CurrentPage   uint64
	TotalPage     uint64
	HasPrev       bool
	PrevPage      uint64
	HasNext       bool
	NextPage      uint64
	Pages         []PaginationPage
	BaseURL       string
	Query         string
	SortDirection int
}

// ShowUsersPage renders the main user page shell
// It fetches the total user count for the stat card.
func (h Handler) ShowUsersPage(c echo.Context) error {
	ctx := grpcContext(c)
	user, _ := getUserFromContext(c)

	// Make a light request just to get the total count for the stat card
	listReq := &userpb.UserListRequest{
		CurrentPage: 1,
		PageSize:    1, // We only need the total count
	}

	listResp, err := h.userAd.UserList(ctx, listReq)
	totalUsers := uint64(0)
	if err != nil {
		// Don't fail the whole page, just log it and render with 0
		h.logger.Error("Failed to get user count for stats", "error", err)
	} else {
		totalUsers = listResp.GetTotalNumber()
	}

	data := map[string]interface{}{
		"TotalUsers":   totalUsers,
		"TemplateName": "users_page",
		"User":         user,
	}

	if isHTMX(c) {
		return c.Render(http.StatusOK, "users_page", data)
	}

	return c.Render(http.StatusOK, "main_layout", data)
}

// ListUsersPartial renders the table and pagination (called by HTMX).
func (h Handler) ListUsersPartial(c echo.Context) error {
	ctx := grpcContext(c)

	// 1. Parse Query Parameters
	page, _ := strconv.ParseUint(c.QueryParam("page"), 10, 64)
	if page == 0 {
		page = 1
	}

	sortDirInt, _ := strconv.ParseInt(c.QueryParam("sort_direction"), 10, 32)
	if sortDirInt == 0 {
		sortDirInt = 2 // Default to DESC (Newest First)
	}

	usernameQuery := c.QueryParam("username")

	listReq := &userpb.UserListRequest{
		CurrentPage:   page,
		PageSize:      defaultPageSize, // Define page size (e.g., 10)
		SortDirection: userpb.SortDirection(sortDirInt),
		Username:      usernameQuery,
	}

	listResp, err := h.userAd.UserList(ctx, listReq)
	if err != nil {
		return h.renderGRPCError(c, "ListUsersPartial", err)
	}

	users := make([]User, len(listResp.GetUsers()))
	for i, pbUser := range listResp.GetUsers() {
		users[i] = convertUserPbToUser(pbUser)
	}

	pagination := h.buildPaginationData(listResp, usernameQuery, int(sortDirInt))

	data := map[string]interface{}{
		"Users":      users,
		"Pagination": pagination,
	}

	return c.Render(http.StatusOK, "users_list_partial", data)
}

// ShowDeleteUserConfirm renders the delete confirmation modal.
func (h Handler) ShowDeleteUserConfirm(c echo.Context) error {
	userID := c.QueryParam("id")
	username := c.QueryParam("username")

	data := map[string]interface{}{
		"ID":       userID,
		"Username": username,
	}

	return c.Render(http.StatusOK, "delete_user_confirm_modal", data)
}

// DeleteUser handles deleting a user (called by HTMX).
func (h Handler) DeleteUser(c echo.Context) error {
	ctx := grpcContext(c)
	userID := c.Param("id")

	if userID == "" {
		return c.String(http.StatusBadRequest, "User ID is required")
	}

	_, err := h.userAd.UserDelete(ctx, &userpb.UserDeleteRequest{UserId: userID})
	if err != nil {
		// TODO: Return an error partial or message
		return h.renderGRPCError(c, "DeleteUser", err)
	}

	return c.NoContent(http.StatusOK)
}

// buildPaginationData is a helper to create the pagination struct.
func (h Handler) buildPaginationData(resp *userpb.UserListResponse, query string, sortDir int) PaginationData {
	p := PaginationData{
		TotalNumber:   resp.GetTotalNumber(),
		CurrentPage:   resp.GetCurrentPage(),
		TotalPage:     resp.GetTotalPage(),
		HasPrev:       resp.GetCurrentPage() > 1,
		HasNext:       resp.GetCurrentPage() < resp.GetTotalPage(),
		BaseURL:       "/users/list", // Base URL for pagination links
		Query:         query,
		SortDirection: sortDir,
	}

	if p.HasPrev {
		p.PrevPage = p.CurrentPage - 1
	}
	if p.HasNext {
		p.NextPage = p.CurrentPage + 1
	}

	// Build page number links (e.g., show max 5 pages)
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
			URL:        fmt.Sprintf("%s?page=%d&username=%s&sort_direction=%d", p.BaseURL, i, p.Query, p.SortDirection),
			//nolint:gosec // G115: CurrentPage is pagination number, overflow is not a realistic concern
			IsCurrent: i == int(p.CurrentPage),
		})
	}

	return p
}

func (h Handler) DetailUser(c echo.Context) error {
	ctx := grpcContext(c)

	userPb, aErr := h.userAd.UserDetail(ctx, &userpb.UserDetailRequest{UserId: c.Param("id")})
	if aErr != nil {
		return h.renderGRPCError(c, "DetailUser", aErr)
	}

	data := map[string]interface{}{
		"User": convertUserPbToUser(userPb),
	}

	return c.Render(http.StatusOK, "user_detail_modal", data)
}

// ShowEditUserModal fetches user details and renders the edit form modal.
func (h Handler) ShowEditUserModal(c echo.Context) error {
	ctx := grpcContext(c)

	userPb, aErr := h.userAd.UserDetail(ctx, &userpb.UserDetailRequest{UserId: c.Param("id")})
	if aErr != nil {
		return h.renderGRPCError(c, "ShowEditUserModal", aErr)
	}

	data := map[string]interface{}{
		"User":     convertUserPbToUser(userPb),
		"AllRoles": GetAllRoles(),
	}

	return c.Render(http.StatusOK, "user_edit_modal", data)
}

// UpdateUser handles the submission of the edit user form.
func (h Handler) UpdateUser(c echo.Context) error {
	ctx := grpcContext(c)
	userID := c.Param("id")

	username := c.FormValue("username")
	fullname := c.FormValue("fullname")
	email := c.FormValue("email")
	phoneNumber := c.FormValue("phone_number")
	roleStrings := c.Request().Form["roles"]

	roles := ParseRolesFromForm(roleStrings)

	req := &userpb.UserUpdateFromSuperUserRequest{
		UserId:      userID,
		Username:    username,
		Fullname:    fullname,
		Email:       email,
		PhoneNumber: phoneNumber,
		Roles:       roles,
	}

	userPb, aErr := h.userAd.UserUpdateFromSuperuser(ctx, req)
	if aErr != nil {
		return h.renderGRPCError(c, "UpdateUser", aErr)
	}

	c.Response().Header().Set("HX-Trigger", "userListChanged")

	data := map[string]interface{}{
		"User": convertUserPbToUser(userPb),
	}

	return c.Render(http.StatusOK, "user_detail_modal", data)
}

// ShowCreateUserModal renders the create user form modal.
func (h Handler) ShowCreateUserModal(c echo.Context) error {
	data := map[string]interface{}{
		"AllRoles": GetAllRoles(),
	}

	return c.Render(http.StatusOK, "user_create_modal", data)
}

// CreateUser handles the submission of the new user form.
func (h Handler) CreateUser(c echo.Context) error {
	ctx := grpcContext(c)

	req := &userpb.CreateUserRequest{
		Username:    c.FormValue("username"),
		Password:    c.FormValue("password"),
		Fullname:    c.FormValue("fullname"),
		Email:       c.FormValue("email"),
		PhoneNumber: c.FormValue("phone_number"),
		Roles:       ParseRolesFromForm(c.Request().Form["roles"]),
	}

	_, aErr := h.userAd.CreateUser(ctx, req)
	if aErr != nil {
		return h.renderGRPCError(c, "CreateUser", aErr)
	}

	c.Response().Header().Set("HX-Trigger", "userListChanged")

	return c.NoContent(http.StatusCreated)
}
