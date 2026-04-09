package filter_test

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/filter"
)

func intSlice(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}
	return s
}

func TestPaginate_NoLimit(t *testing.T) {
	items := intSlice(10)
	p := filter.Paginate(items, 0, 0)
	if len(p.Items) != 10 {
		t.Fatalf("expected 10 items, got %d", len(p.Items))
	}
	if p.HasMore {
		t.Fatal("expected HasMore=false when no limit")
	}
}

func TestPaginate_FirstPage(t *testing.T) {
	items := intSlice(20)
	p := filter.Paginate(items, 0, 5)
	if len(p.Items) != 5 {
		t.Fatalf("expected 5 items, got %d", len(p.Items))
	}
	if !p.HasMore {
		t.Fatal("expected HasMore=true")
	}
	if p.Total != 20 {
		t.Fatalf("expected Total=20, got %d", p.Total)
	}
}

func TestPaginate_LastPage(t *testing.T) {
	items := intSlice(7)
	p := filter.Paginate(items, 5, 5)
	if len(p.Items) != 2 {
		t.Fatalf("expected 2 items on last page, got %d", len(p.Items))
	}
	if p.HasMore {
		t.Fatal("expected HasMore=false on last page")
	}
}

func TestPaginate_OffsetBeyondTotal(t *testing.T) {
	items := intSlice(3)
	p := filter.Paginate(items, 10, 5)
	if len(p.Items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(p.Items))
	}
	if p.HasMore {
		t.Fatal("expected HasMore=false")
	}
}

func TestPaginate_NegativeOffset(t *testing.T) {
	items := intSlice(5)
	p := filter.Paginate(items, -3, 2)
	if p.Offset != 0 {
		t.Fatalf("expected offset normalised to 0, got %d", p.Offset)
	}
	if len(p.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(p.Items))
	}
}

func TestPaginate_ExactPage(t *testing.T) {
	items := intSlice(6)
	p := filter.Paginate(items, 0, 6)
	if len(p.Items) != 6 {
		t.Fatalf("expected 6 items, got %d", len(p.Items))
	}
	if p.HasMore {
		t.Fatal("expected HasMore=false when items fit exactly")
	}
}
