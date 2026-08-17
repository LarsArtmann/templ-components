package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/layout"
)

type demoUser struct {
	name, email, role, status string
}

func demoUsers() []demoUser {
	return []demoUser{
		{"Alice Smith", "alice@example.com", "Admin", "active"},
		{"Bob Jones", "bob@example.com", "User", "invited"},
		{"Carla Nguyen", "carla@example.com", "User", "active"},
		{"David Kim", "david@example.com", "Editor", "active"},
		{"Elena Petrova", "elena@example.com", "User", "suspended"},
		{"Frank Miller", "frank@example.com", "Admin", "active"},
		{"Grace Lee", "grace@example.com", "User", "invited"},
		{"Hugo Oliveira", "hugo@example.com", "Editor", "active"},
	}
}

const usersPageSize = 5

// sortDemoUsers returns the user list sorted by the DataTable sort key
// ("Name" or "email") in the requested direction.
func sortDemoUsers(users []demoUser, sortKey string, dir display.SortDirection) []demoUser {
	out := make([]demoUser, len(users))
	copy(out, users)
	byEmail := sortKey == "email"
	sort.Slice(out, func(i, j int) bool {
		less := out[i].name < out[j].name
		if byEmail {
			less = out[i].email < out[j].email
		}
		if dir == display.SortDesc {
			return !less
		}
		return less
	})
	return out
}

// pageDemoUsers slices the sorted list for the requested 1-based page.
func pageDemoUsers(users []demoUser, page int) []demoUser {
	if page < 1 {
		page = 1
	}
	start := (page - 1) * usersPageSize
	if start >= len(users) {
		return nil
	}
	end := min(start+usersPageSize, len(users))
	return users[start:end]
}

func demoUserRows(users []demoUser) []display.TableRow {
	rows := make([]display.TableRow, len(users))
	for i, u := range users {
		rows[i] = display.SimpleTableRow(u.name, u.email, u.role, u.status)
	}
	return rows
}

func usersTotalPages(total int) uint {
	pages := (total + usersPageSize - 1) / usersPageSize
	if pages < 1 {
		pages = 1
	}
	return uint(pages)
}

// usersSortColumnLabel maps a DataTable sort key back to its column Label,
// which is what DataTableProps.ActiveSortColumn expects.
func usersSortColumnLabel(sortKey string) string {
	if sortKey == "email" {
		return "Email"
	}
	return "Name"
}

// filterDemoUsers applies the FilterDropdown demo filters: status (exact
// match, "all"/"" matches everything) and sort ("name" keeps insertion
// order, anything else reverses it as a stand-in for "newest first").
func filterDemoUsers(status, sortKey string) []demoUser {
	users := demoUsers()
	filtered := make([]demoUser, 0, len(users))
	for _, u := range users {
		if status == "" || status == "all" || u.status == status {
			filtered = append(filtered, u)
		}
	}
	if sortKey == "date" || sortKey == "created" {
		for i, j := 0, len(filtered)-1; i < j; i, j = i+1, j-1 {
			filtered[i], filtered[j] = filtered[j], filtered[i]
		}
	}
	return filtered
}

// renderUsersPage serves the /users list page: a real server-driven round
// trip for DataTable sorting and Pagination — the pattern the library is
// built for.
func renderUsersPage(w http.ResponseWriter, r *http.Request) {
	sortKey := r.URL.Query().Get("sort")
	if sortKey != "email" {
		sortKey = "Name"
	}
	dir := display.SortDirection(r.URL.Query().Get("dir"))
	if dir != display.SortDesc {
		dir = display.SortAsc
	}
	page := 1
	if v, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && v > 0 {
		page = v
	}

	sorted := sortDemoUsers(demoUsers(), sortKey, dir)
	pageFunc := func(props layout.PageProps) templ.Component {
		return usersDemoPage(props, sortKey, dir, uint(page), usersTotalPages(len(sorted)))
	}
	renderPage(w, r, "Users - templ-components", "Server-driven data table with sorting and pagination", pageFunc)
}
