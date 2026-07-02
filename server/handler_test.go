package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"pcms/config"
	"pcms/handlers"
	"pcms/services"
	"pcms/store"
	"pcms/utils"

	"github.com/gin-gonic/gin"
)

// ============ Test Helpers ============

func handlerTestEnv() (*gin.Engine, *store.MemStore) {
	gin.SetMode(gin.TestMode)
	utils.InitJWT("test-secret-key")
	config.Load() // 初始化 config.AppConfig

	mem := store.NewMemStore()
	authSvc := services.NewAuthService(mem)
	authSvc.Register(services.RegisterInput{
		Username: "testuser",
		Password: "test123456",
		Email:    "test@test.com",
	})

	r := setupHandlerRouter(mem, authSvc)
	return r, mem
}

func setupHandlerRouter(mem *store.MemStore, authSvc *services.AuthService) *gin.Engine {
	r := gin.New()
	// 静态文件
	// 注册认证路由（不需要JWT）
	auth := r.Group("/api/v1/auth")
	{
		authHandler := handlers.NewAuthHandler(authSvc)
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// 需要JWT认证
	protected := r.Group("/api/v1")
	protected.Use(authMiddleware())
	{
		protected.GET("/auth/me", handlers.NewAuthHandler(authSvc).Me)
		protected.GET("/categories", handlers.NewCategoryHandler(services.NewCategoryService(mem)).GetTree)
		protected.POST("/categories", handlers.NewCategoryHandler(services.NewCategoryService(mem)).Create)
		protected.PUT("/categories/:id", handlers.NewCategoryHandler(services.NewCategoryService(mem)).Update)
		protected.DELETE("/categories/:id", handlers.NewCategoryHandler(services.NewCategoryService(mem)).Delete)
		protected.PUT("/categories/:id/move", handlers.NewCategoryHandler(services.NewCategoryService(mem)).Move)

		docH := handlers.NewDocumentHandler(services.NewDocumentService(mem))
		protected.GET("/documents", docH.List)
		protected.GET("/documents/search", docH.Search)
		protected.POST("/documents", docH.Create)
		protected.GET("/documents/:id", docH.Get)
		protected.PUT("/documents/:id", docH.Update)
		protected.DELETE("/documents/:id", docH.Delete)
		protected.GET("/documents/:id/versions", docH.GetVersions)
		protected.GET("/documents/:id/versions/:vid", docH.GetVersion)

		tagH := handlers.NewTagHandler(services.NewTagService(mem))
		protected.GET("/tags", tagH.List)
		protected.POST("/tags", tagH.Create)
		protected.DELETE("/tags/:id", tagH.Delete)

		attH := handlers.NewAttachmentHandler(services.NewAttachmentService(mem))
		protected.GET("/attachments", attH.List)
		protected.POST("/attachments", attH.Create)
		protected.DELETE("/attachments/:id", attH.Delete)
		protected.PUT("/attachments/:id/bind", attH.BindDocument)

		protected.POST("/files/upload", handlers.Upload)
	}

	// 静态文件（无认证）
	r.GET("/api/v1/files/*filepath", handlers.ServeFile)
	return r
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ah := c.GetHeader("Authorization")
		if ah == "" || len(ah) < 8 || ah[:7] != "Bearer " {
			c.JSON(401, gin.H{"code": -1, "message": "未授权"})
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(ah[7:])
		if err != nil {
			c.JSON(401, gin.H{"code": -1, "message": "令牌无效"})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func doReq(r *gin.Engine, method, path, token string, body interface{}) *httptest.ResponseRecorder {
	var bb []byte
	if body != nil {
		bb, _ = json.Marshal(body)
	}
	var req *http.Request
	if bb != nil {
		req = httptest.NewRequest(method, path, bytes.NewReader(bb))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func respJSON(w *httptest.ResponseRecorder) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &m)
	return m
}

func respData(w *httptest.ResponseRecorder) map[string]interface{} {
	m := respJSON(w)
	if d, ok := m["data"].(map[string]interface{}); ok {
		return d
	}
	return nil
}

func testToken() string {
	utils.InitJWT("test-secret-key")
	t, _ := utils.GenerateToken(1, "testuser")
	return t
}

// ============ Auth Handler Tests ============

func TestAuth_Register(t *testing.T) {
	r, _ := handlerTestEnv()

	t.Run("success", func(t *testing.T) {
		w := doReq(r, "POST", "/api/v1/auth/register", "", map[string]string{
			"username": "new", "password": "pass1234", "email": "x@x.com",
		})
		if w.Code != 201 {
			t.Fatalf("want 201, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("duplicate", func(t *testing.T) {
		w := doReq(r, "POST", "/api/v1/auth/register", "", map[string]string{
			"username": "testuser", "password": "pass",
		})
		if w.Code == 201 {
			t.Fatal("should fail for duplicate")
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		w := doReq(r, "POST", "/api/v1/auth/register", "", map[string]string{})
		if w.Code == 201 {
			t.Fatal("should fail")
		}
	})
}

func TestAuth_Login(t *testing.T) {
	r, _ := handlerTestEnv()

	t.Run("success", func(t *testing.T) {
		w := doReq(r, "POST", "/api/v1/auth/login", "", map[string]string{
			"username": "testuser", "password": "test123456",
		})
		if w.Code != 200 {
			t.Fatalf("want 200, got %d", w.Code)
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		w := doReq(r, "POST", "/api/v1/auth/login", "", map[string]string{
			"username": "testuser", "password": "wrong",
		})
		if w.Code == 200 {
			t.Fatal("should fail")
		}
	})

	t.Run("no user", func(t *testing.T) {
		w := doReq(r, "POST", "/api/v1/auth/login", "", map[string]string{
			"username": "noone", "password": "pass",
		})
		if w.Code == 200 {
			t.Fatal("should fail")
		}
	})
}

func TestAuth_Me(t *testing.T) {
	r, _ := handlerTestEnv()
	token := testToken()

	t.Run("success", func(t *testing.T) {
		w := doReq(r, "GET", "/api/v1/auth/me", token, nil)
		if w.Code != 200 {
			t.Fatalf("want 200, got %d", w.Code)
		}
	})

	t.Run("no auth", func(t *testing.T) {
		w := doReq(r, "GET", "/api/v1/auth/me", "", nil)
		if w.Code != 401 {
			t.Fatalf("want 401, got %d", w.Code)
		}
	})
}

// ============ Category Handler Tests ============

func TestCategory_CRUD(t *testing.T) {
	r, _ := handlerTestEnv()
	token := testToken()

	t.Run("create", func(t *testing.T) {
		w := doReq(r, "POST", "/api/v1/categories", token, map[string]string{"name": "Tech"})
		if w.Code != 201 {
			t.Fatalf("want 201, got %d", w.Code)
		}
	})

	t.Run("tree", func(t *testing.T) {
		w := doReq(r, "GET", "/api/v1/categories", token, nil)
		if w.Code != 200 {
			t.Fatalf("want 200, got %d", w.Code)
		}
	})

	t.Run("update", func(t *testing.T) {
		w := doReq(r, "PUT", "/api/v1/categories/1", token, map[string]string{"name": "New"})
		if w.Code != 200 {
			t.Fatalf("want 200, got %d", w.Code)
		}
	})

	t.Run("delete", func(t *testing.T) {
		w := doReq(r, "DELETE", "/api/v1/categories/1", token, nil)
		if w.Code != 200 {
			t.Fatalf("want 200, got %d", w.Code)
		}
	})
}

func TestCategory_Move(t *testing.T) {
	r, _ := handlerTestEnv()
	token := testToken()

	doReq(r, "POST", "/api/v1/categories", token, map[string]string{"name": "A"})
	doReq(r, "POST", "/api/v1/categories", token, map[string]string{"name": "B"})

	w := doReq(r, "PUT", "/api/v1/categories/2/move", token, map[string]interface{}{
		"parent_id": 1, "sort_order": 0,
	})
	if w.Code != 200 {
		t.Fatalf("want 200, got %d", w.Code)
	}
}

func TestCategory_Unauth(t *testing.T) {
	r, _ := handlerTestEnv()

	w := doReq(r, "GET", "/api/v1/categories", "", nil)
	if w.Code != 401 {
		t.Fatalf("want 401, got %d", w.Code)
	}
}

// ============ Document Handler Tests ============

func TestDocument_Create(t *testing.T) {
	r, _ := handlerTestEnv()
	token := testToken()

	w := doReq(r, "POST", "/api/v1/documents", token, map[string]string{
		"title": "Doc", "content": "Hello",
	})
	if w.Code != 201 {
		t.Fatalf("want 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDocument_List(t *testing.T) {
	r, _ := handlerTestEnv()
	token := testToken()

	doReq(r, "POST", "/api/v1/documents", token, map[string]string{"title": "A", "content": "a"})
	doReq(r, "POST", "/api/v1/documents", token, map[string]string{"title": "B", "content": "b"})

	w := doReq(r, "GET", "/api/v1/documents", token, nil)
	if w.Code != 200 {
		t.Fatalf("want 200, got %d", w.Code)
	}
}

func TestDocument_Get(t *testing.T) {
	r, _ := handlerTestEnv()
	token := testToken()

	doReq(r, "POST", "/api/v1/documents", token, map[string]string{"title": "Mine", "content": "c"})

	w := doReq(r, "GET", "/api/v1/documents/1", token, nil)
	if w.Code != 200 {
		t.Fatalf("want 200, got %d", w.Code)
	}

	w = doReq(r, "GET", "/api/v1/documents/999", token, nil)
	if w.Code != 404 {
		t.Fatalf("want 404, got %d", w.Code)
	}
}

func TestDocument_UpdateDelete(t *testing.T) {
	r, _ := handlerTestEnv()
	token := testToken()

	doReq(r, "POST", "/api/v1/documents", token, map[string]string{"title": "Old", "content": "old"})

	w := doReq(r, "PUT", "/api/v1/documents/1", token, map[string]string{"title": "New", "content": "new"})
	if w.Code != 200 {
		t.Fatalf("want 200, got %d", w.Code)
	}

	w = doReq(r, "DELETE", "/api/v1/documents/1", token, nil)
	if w.Code != 200 {
		t.Fatalf("want 200, got %d", w.Code)
	}
}

func TestDocument_Search(t *testing.T) {
	r, _ := handlerTestEnv()
	token := testToken()

	doReq(r, "POST", "/api/v1/documents", token, map[string]string{"title": "Go", "content": "programming"})

	w := doReq(r, "GET", "/api/v1/documents/search?keyword=Go", token, nil)
	if w.Code != 200 {
		t.Fatalf("want 200, got %d", w.Code)
	}

	w = doReq(r, "GET", "/api/v1/documents/search?keyword=", token, nil)
	if w.Code == 200 {
		t.Fatal("should fail without keyword")
	}
}

// ============ Tag Handler Tests ============

func TestTag_CRUD(t *testing.T) {
	r, _ := handlerTestEnv()
	token := testToken()

	t.Run("create", func(t *testing.T) {
		w := doReq(r, "POST", "/api/v1/tags", token, map[string]string{"name": "K8s"})
		if w.Code != 201 {
			t.Fatalf("want 201, got %d", w.Code)
		}
	})

	t.Run("list", func(t *testing.T) {
		w := doReq(r, "GET", "/api/v1/tags", token, nil)
		if w.Code != 200 {
			t.Fatalf("want 200, got %d", w.Code)
		}
	})

	t.Run("delete", func(t *testing.T) {
		w := doReq(r, "DELETE", "/api/v1/tags/1", token, nil)
		if w.Code != 200 {
			t.Fatalf("want 200, got %d", w.Code)
		}
	})
}

// ============ Attachment Handler Tests ============

func TestAttachment_CRUD(t *testing.T) {
	r, _ := handlerTestEnv()
	token := testToken()

	body := map[string]interface{}{
		"file_name": "test.pdf", "file_path": "/uploads/test.pdf",
		"file_size": 1024, "mime_type": "application/pdf",
	}

	t.Run("create", func(t *testing.T) {
		w := doReq(r, "POST", "/api/v1/attachments", token, body)
		if w.Code != 201 {
			t.Fatalf("want 201, got %d", w.Code)
		}
	})

	t.Run("list", func(t *testing.T) {
		w := doReq(r, "GET", "/api/v1/attachments", token, nil)
		if w.Code != 200 {
			t.Fatalf("want 200, got %d", w.Code)
		}
	})

	t.Run("delete", func(t *testing.T) {
		w := doReq(r, "DELETE", "/api/v1/attachments/1", token, nil)
		if w.Code != 200 {
			t.Fatalf("want 200, got %d", w.Code)
		}
	})
}

// ============ File Handler Tests ============

func TestFile_Serve(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("UPLOAD_DIR", dir)
	p := filepath.Join(dir, "2024", "01", "01")
	os.MkdirAll(p, 0755)
	os.WriteFile(filepath.Join(p, "test.txt"), []byte("hello"), 0644)

	r, _ := handlerTestEnv()

	t.Run("found", func(t *testing.T) {
		w := doReq(r, "GET", "/api/v1/files/2024/01/01/test.txt", "", nil)
		if w.Code != 200 || w.Body.String() != "hello" {
			t.Fatalf("want 200 'hello', got %d %q", w.Code, w.Body.String())
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doReq(r, "GET", "/api/v1/files/nonexist.txt", "", nil)
		if w.Code != 404 {
			t.Fatalf("want 404, got %d", w.Code)
		}
	})
}
