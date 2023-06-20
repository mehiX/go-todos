package todos

import "strings"

type Todo struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

// CleanTags returns a comma separated list of deduplicated tags (small caps)
func (t Todo) CleanTags() string {
	tags := make(map[string]bool)
	for _, tag := range t.Tags {
		cleaned := strings.ToLower(strings.TrimSpace(tag))
		tags[cleaned] = true
	}

	var unique []string
	for k := range tags {
		unique = append(unique, k)
	}

	return strings.Join(unique, ",")
}
