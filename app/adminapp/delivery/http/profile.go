package http

import (
	"net/http"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userpb"
)

type User struct {
	ID           string
	Username     string
	Fullname     string
	Email        string
	PhoneNumber  string
	Avatar       string
	Roles        []string
	LastOnlineAt string
}

// ShowProfilePage renders the user's profile page
func (h Handler) ShowProfilePage(c echo.Context) error {
	ctx := grpcContext(c)

	userPb, aErr := h.userAd.UserProfile(ctx, &empty.Empty{})
	if aErr != nil {
		return h.renderGRPCError(c, "ShowProfilePage", aErr)
	}

	data := map[string]interface{}{
		"User":         convertUserPbToUser(userPb),
		"TemplateName": "profile_page",
	}

	if isHTMX(c) {
		return c.Render(http.StatusOK, "profile_page", data)
	}

	return c.Render(http.StatusOK, "main_layout", data)
}

func (h Handler) ShowProfileView(c echo.Context) error {
	ctx := grpcContext(c)

	userPb, aErr := h.userAd.UserProfile(ctx, &empty.Empty{})
	if aErr != nil {
		return h.renderGRPCError(c, "ShowProfileView", aErr)
	}

	data := map[string]interface{}{
		"User": convertUserPbToUser(userPb),
	}

	return c.Render(http.StatusOK, "profile_view", data)
}

// ShowProfileEditForm renders the EDITABLE profile form partial
// (Called by 'Edit' button)
func (h Handler) ShowProfileEditForm(c echo.Context) error {
	ctx := grpcContext(c)

	userPb, aErr := h.userAd.UserProfile(ctx, &empty.Empty{})
	if aErr != nil {
		return h.renderGRPCError(c, "ShowProfileEditForm", aErr)
	}

	data := map[string]interface{}{
		"User": convertUserPbToUser(userPb),
	}

	return c.Render(http.StatusOK, "profile_edit_form", data)
}

// UpdateProfile handles the submission of the edit form
// (Called by 'Save' button)
func (h Handler) UpdateProfile(c echo.Context) error {
	ctx := grpcContext(c)

	username := c.FormValue("username")
	fullname := c.FormValue("fullname")
	email := c.FormValue("email")
	phoneNumber := c.FormValue("phone_number")

	userPb, aErr := h.userAd.UserUpdateFromOwn(ctx, &userpb.UserUpdateFromOwnRequest{
		Username:    username,
		Fullname:    fullname,
		Email:       email,
		PhoneNumber: phoneNumber,
	})
	if aErr != nil {
		return h.renderGRPCError(c, "UpdateProfile", aErr)
	}

	data := map[string]interface{}{
		"User": convertUserPbToUser(userPb),
	}

	return c.Render(http.StatusOK, "profile_view", data)
}
