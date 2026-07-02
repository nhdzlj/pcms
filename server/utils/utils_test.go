package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// ============ JWT Tests ============

func TestGenerateAndParseToken(t *testing.T) {
	InitJWT("my-test-secret")

	token, err := GenerateToken(1, "admin")
	if err != nil {
		t.Fatalf("生成 token 失败: %v", err)
	}
	if token == "" {
		t.Fatal("token 不应为空")
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("解析 token 失败: %v", err)
	}
	if claims.UserID != 1 {
		t.Fatalf("期望 UserID=1, 得到 %d", claims.UserID)
	}
	if claims.Username != "admin" {
		t.Fatalf("期望 admin, 得到 %s", claims.Username)
	}
	if claims.Issuer != "pcms" {
		t.Fatalf("期望 Issuer=pcms, 得到 %s", claims.Issuer)
	}
}

func TestParseToken_Invalid(t *testing.T) {
	InitJWT("test-secret")

	t.Run("无效token", func(t *testing.T) {
		_, err := ParseToken("invalid.token.here")
		if err == nil {
			t.Fatal("期望报错")
		}
	})

	t.Run("不同密钥", func(t *testing.T) {
		InitJWT("key-a")
		token, _ := GenerateToken(1, "user")
		InitJWT("key-b")
		_, err := ParseToken(token)
		if err == nil {
			t.Fatal("不同密钥应报错")
		}
	})
}

// ============ Pagination Tests ============

func TestGetPagination(t *testing.T) {
	tests := []struct {
		query        string
		expectPage   int
		expectSize   int
	}{
		{"", 1, 20},
		{"page=3&page_size=10", 3, 10},
		{"page=0&page_size=0", 1, 20},
		{"page=-1&page_size=-5", 1, 20},
		{"page=1&page_size=200", 1, 20},
		{"page=abc&page_size=xyz", 1, 20},
	}

	for _, tc := range tests {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test?"+tc.query, nil)

		page, size := GetPagination(c)
		if page != tc.expectPage || size != tc.expectSize {
			t.Errorf("query=%q: 期望 (%d,%d), 得到 (%d,%d)",
				tc.query, tc.expectPage, tc.expectSize, page, size)
		}
	}
}

// ============ Response Tests ============

func TestResponseSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	Success(c, map[string]string{"key": "val"})
	if w.Code != http.StatusOK {
		t.Fatalf("期望 200, 得到 %d", w.Code)
	}
}

func TestResponseCreated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	Created(c, map[string]string{"id": "1"})
	if w.Code != http.StatusCreated {
		t.Fatalf("期望 201, 得到 %d", w.Code)
	}
}

func TestResponseErrors(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(c *gin.Context, msg string)
		expected int
	}{
		{"BadRequest", BadRequest, 400},
		{"Unauthorized", Unauthorized, 401},
		{"NotFound", NotFound, 404},
		{"InternalError", InternalError, 500},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			tc.fn(c, "test message")
			if w.Code != tc.expected {
				t.Fatalf("期望 %d, 得到 %d", tc.expected, w.Code)
			}
		})
	}
}

func TestResponseError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	Error(c, http.StatusForbidden, "forbidden")
	if w.Code != 403 {
		t.Fatalf("期望 403, 得到 %d", w.Code)
	}
}
