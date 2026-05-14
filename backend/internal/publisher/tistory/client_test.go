package tistory

import "testing"

func TestParseCategoriesFromWindowConfig(t *testing.T) {
	body := []byte(`
	<script>
	window.Config = {
		blog: {
			categories: [{"id":1550884,"name":"알고리즘 문제 풀이"},{"id":1550885,"name":"취준"}]
		}
	};
	</script>`)

	items, err := ParseCategories(body)
	if err != nil {
		t.Fatalf("ParseCategories() error = %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("len(items) = %d, want 3", len(items))
	}
	if items[1].CategoryID != "1550884" || items[1].Label != "알고리즘 문제 풀이" {
		t.Fatalf("items[1] = %+v", items[1])
	}
}
