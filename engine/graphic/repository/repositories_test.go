package repository

import (
	gresource "engine/graphic/resource"
	"testing"
)

func TestUniformBlockRepository(t *testing.T) {
	t.Run("should put and find by id", func(t *testing.T) {
		repo := NewUniformBlockRepository()
		ub := gresource.UniformBlock{ID: 1, Name: "ProjView", BindingPoint: 0}

		if !repo.Put(ub) {
			t.Fatal("expected Put to succeed")
		}

		got, ok := repo.ByID(1)
		if !ok {
			t.Fatal("expected ByID to find block")
		}
		if got != ub {
			t.Errorf("ByID = %+v, want %+v", got, ub)
		}
	})

	t.Run("should return false for missing id", func(t *testing.T) {
		if _, ok := NewUniformBlockRepository().ByID(999); ok {
			t.Error("expected ByID missing to return false")
		}
	})

	t.Run("should find by name", func(t *testing.T) {
		repo := NewUniformBlockRepository()
		repo.Put(gresource.UniformBlock{ID: 7, Name: "ProjView"})

		got, ok := repo.ByName("ProjView")
		if !ok {
			t.Fatal("expected ByName to find block")
		}
		if got.ID != 7 {
			t.Errorf("ByName ID = %d, want 7", got.ID)
		}
		if _, ok := repo.ByName("missing"); ok {
			t.Error("expected ByName missing to return false")
		}
	})

	t.Run("should reject duplicate by id", func(t *testing.T) {
		repo := NewUniformBlockRepository()
		repo.Put(gresource.UniformBlock{ID: 1, Name: "a"})

		if repo.Put(gresource.UniformBlock{ID: 1, Name: "b"}) {
			t.Error("expected duplicate Put by ID to fail")
		}
	})

	t.Run("should delete from both indexes", func(t *testing.T) {
		repo := NewUniformBlockRepository()
		repo.Put(gresource.UniformBlock{ID: 1, Name: "ProjView"})

		if !repo.Delete(1) {
			t.Fatal("expected Delete to succeed")
		}
		if _, ok := repo.ByID(1); ok {
			t.Error("expected ByID to be empty after delete")
		}
		if _, ok := repo.ByName("ProjView"); ok {
			t.Error("expected ByName to be empty after delete")
		}
		if repo.Delete(1) {
			t.Error("expected second Delete to fail")
		}
	})

	t.Run("should list all stored blocks", func(t *testing.T) {
		repo := NewUniformBlockRepository()
		repo.Put(gresource.UniformBlock{ID: 1, Name: "a"})
		repo.Put(gresource.UniformBlock{ID: 2, Name: "b"})

		if got := len(repo.All()); got != 2 {
			t.Errorf("All len = %d, want 2", got)
		}
	})
}

func TestShaderRepository(t *testing.T) {
	t.Run("should put and find by id", func(t *testing.T) {
		repo := NewShaderRepository()
		sh := gresource.Shader{ID: 1, Name: "Tilemap"}

		if !repo.Put(sh) {
			t.Fatal("expected Put to succeed")
		}

		got, ok := repo.ByID(1)
		if !ok {
			t.Fatal("expected ByID to find shader")
		}
		if got != sh {
			t.Errorf("ByID = %+v, want %+v", got, sh)
		}
	})

	t.Run("should return false for missing id", func(t *testing.T) {
		if _, ok := NewShaderRepository().ByID(999); ok {
			t.Error("expected ByID missing to return false")
		}
	})

	t.Run("should find by name", func(t *testing.T) {
		repo := NewShaderRepository()
		repo.Put(gresource.Shader{ID: 3, Name: "Tilemap"})

		got, ok := repo.ByName("Tilemap")
		if !ok {
			t.Fatal("expected ByName to find shader")
		}
		if got.ID != 3 {
			t.Errorf("ByName ID = %d, want 3", got.ID)
		}
		if _, ok := repo.ByName("missing"); ok {
			t.Error("expected ByName missing to return false")
		}
	})

	t.Run("should reject duplicate by id", func(t *testing.T) {
		repo := NewShaderRepository()
		repo.Put(gresource.Shader{ID: 1, Name: "a"})

		if repo.Put(gresource.Shader{ID: 1, Name: "b"}) {
			t.Error("expected duplicate Put by ID to fail")
		}
	})

	t.Run("should delete from both indexes", func(t *testing.T) {
		repo := NewShaderRepository()
		repo.Put(gresource.Shader{ID: 1, Name: "Tilemap"})

		if !repo.Delete(1) {
			t.Fatal("expected Delete to succeed")
		}
		if _, ok := repo.ByID(1); ok {
			t.Error("expected ByID to be empty after delete")
		}
		if _, ok := repo.ByName("Tilemap"); ok {
			t.Error("expected ByName to be empty after delete")
		}
		if repo.Delete(1) {
			t.Error("expected second Delete to fail")
		}
	})

	t.Run("should list all stored shaders", func(t *testing.T) {
		repo := NewShaderRepository()
		repo.Put(gresource.Shader{ID: 1, Name: "a"})
		repo.Put(gresource.Shader{ID: 2, Name: "b"})

		if got := len(repo.All()); got != 2 {
			t.Errorf("All len = %d, want 2", got)
		}
	})
}
