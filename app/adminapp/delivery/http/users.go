package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ShowUsers renders the users list page
func (h Handler) ShowUsers(c echo.Context) error {
	// TODO: Fetch users from database
	// users, err := h.userService.GetUsers(ctx, filters)

	data := map[string]interface{}{
		"Title": "User Management",
		"Stats": map[string]interface{}{
			"Total":    1234,
			"Active":   1089,
			"Inactive": 145,
			"NewToday": 23,
		},
		// Add actual users data here when you have the service
	}

	return c.Render(http.StatusOK, "users_page", data)
}

// SearchUsers handles user search via HTMX
func (h Handler) SearchUsers(c echo.Context) error {
	query := c.QueryParam("q")

	// TODO: Search users in database
	// users, err := h.userService.SearchUsers(ctx, query)

	h.logger.Info("Searching users", "query", query)

	// Return only the table body rows
	return c.HTML(http.StatusOK, `
		<tr class="user-row">
			<td><input type="checkbox" class="checkbox-input"></td>
			<td>
				<div class="user-cell">
					<div class="user-avatar" style="background: linear-gradient(135deg, #06b6d4, #8b5cf6);">
						<span>SR</span>
					</div>
					<div class="user-info">
						<span class="user-name">Search Result</span>
						<span class="user-id">#USR-999</span>
					</div>
				</div>
			</td>
			<td><span class="user-email">search@example.com</span></td>
			<td><span class="role-badge user">User</span></td>
			<td><span class="status-badge active">Active</span></td>
			<td><span class="user-date">Jan 25, 2024</span></td>
			<td>
				<div class="action-buttons">
					<button class="action-btn view" title="View Details">
						<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z" />
							<path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
						</svg>
					</button>
					<button class="action-btn edit" title="Edit User">
						<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
						</svg>
					</button>
					<button class="action-btn delete" title="Delete User">
						<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
						</svg>
					</button>
				</div>
			</td>
		</tr>
	`)
}

// ExportUsers handles user export
func (h Handler) ExportUsers(c echo.Context) error {
	// TODO: Generate CSV/Excel export
	h.logger.Info("Exporting users")

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Export functionality will be implemented",
	})
}

// ShowCreateUserForm shows the create user modal
func (h Handler) ShowCreateUserForm(c echo.Context) error {
	// TODO: Return modal HTML for creating user
	return c.HTML(http.StatusOK, `
		<div class="modal-backdrop" id="user-modal">
			<div class="modal-content">
				<h2>Add New User</h2>
				<p>Create user form will be here</p>
				<button onclick="document.getElementById('user-modal').remove()">Close</button>
			</div>
		</div>
	`)
}
