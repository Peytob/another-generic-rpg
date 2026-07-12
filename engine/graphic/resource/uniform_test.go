package resource

import "testing"

func TestUniformVariables_SetAndLookup(t *testing.T) {
	t.Run("should return variable after set", func(t *testing.T) {
		uv := NewUniformVariables().Set("u_proj", 64, 0)

		v, ok := uv.Lookup("u_proj")
		if !ok {
			t.Fatal("expected variable to be found")
		}
		if v.Name != "u_proj" || v.Size != 64 || v.Offset != 0 {
			t.Errorf("unexpected variable: %+v", v)
		}
	})

	t.Run("should return false for missing variable", func(t *testing.T) {
		if _, ok := NewUniformVariables().Lookup("missing"); ok {
			t.Error("expected missing variable to not be found")
		}
	})

	t.Run("should allow chained set", func(t *testing.T) {
		uv := NewUniformVariables().Set("a", 4, 0).Set("b", 16, 4)

		if _, ok := uv.Lookup("a"); !ok {
			t.Error("expected first variable after chaining")
		}
		if _, ok := uv.Lookup("b"); !ok {
			t.Error("expected second variable after chaining")
		}
	})
}

func TestUniformVariables_TotalSize(t *testing.T) {
	t.Run("should return max of size+offset across variables", func(t *testing.T) {
		uv := NewUniformVariables().
			Set("a", 4, 0).
			Set("b", 16, 32).
			Set("c", 8, 8)

		if got := uv.TotalSize(); got != 48 {
			t.Errorf("TotalSize = %d, want 48", got)
		}
	})

	t.Run("should return 0 for empty variables", func(t *testing.T) {
		if got := NewUniformVariables().TotalSize(); got != 0 {
			t.Errorf("TotalSize = %d, want 0", got)
		}
	})

	t.Run("should be nil-safe", func(t *testing.T) {
		var uv *UniformVariables

		if got := uv.TotalSize(); got != 0 {
			t.Errorf("TotalSize on nil = %d, want 0", got)
		}
		if _, ok := uv.Lookup("x"); ok {
			t.Error("expected Lookup on nil to return false")
		}
	})
}
