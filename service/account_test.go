package service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alpacahq/ribbit-backend/mock"
	"github.com/alpacahq/ribbit-backend/mock/mockdb"
	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/repository/account"
	"github.com/alpacahq/ribbit-backend/secret"
	"github.com/alpacahq/ribbit-backend/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name        string
		req         string
		wantStatus  int
		wantResp    *model.User
		accountRepo *mockdb.Account
		rbac        *mock.RBAC
	}{
		{
			name:       "Invalid request",
			req:        `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","email":"johndoe@gmail.com","role_id":3}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on userSvc",
			req:  `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","email":"johndoe@gmail.com","role_id":3}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c *gin.Context, roleID int) bool {
					return false
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","email":"johndoe@gmail.com","role_id":3}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c *gin.Context, roleID int) bool {
					return true
				},
			},
			accountRepo: &mockdb.Account{
				CreateFn: func(usr *model.User) (*model.User, error) {
					usr.ID = 1
					usr.CreatedAt = mock.TestTime(2018)
					usr.UpdatedAt = mock.TestTime(2018)
					return usr, nil
				},
			},
			wantResp: &model.User{
				Base: model.Base{
					CreatedAt: mock.TestTime(2018),
					UpdatedAt: mock.TestTime(2018),
				},
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Username:  "juzernejm",
				Email:     "johndoe@gmail.com",
			},
			wantStatus: http.StatusOK,
		},
	}

	gin.SetMode(gin.TestMode)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			rg := r.Group("/v1")
			accountService := account.NewAccountService(nil, tt.accountRepo, tt.rbac, secret.New())
			service.AccountRouter(accountService, rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users"
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(model.User)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestChangePassword(t *testing.T) {
	cases := []struct {
		name        string
		req         string
		wantStatus  int
		id          string
		userRepo    *mockdb.User
		accountRepo *mockdb.Account
		rbac        *mock.RBAC
	}{
		{
			name:       "Invalid request",
			req:        `{"new_password":"new_password","old_password":"my_old_password"}`,
			wantStatus: http.StatusBadRequest,
			id:         "1",
		},
		{
			name: "Fail on RBAC",
			req:  `{"new_password":"newpassw","old_password":"oldpassw"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(c *gin.Context, id int) bool {
					return false
				},
			},
			id:         "1",
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `{"new_password":"newpassw","old_password":"oldpassw"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(c *gin.Context, id int) bool {
					return true
				},
			},
			id: "1",
			userRepo: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					return &model.User{
						Password: secret.New().HashPassword("oldpassw"),
					}, nil
				},
			},
			accountRepo: &mockdb.Account{
				ChangePasswordFn: func(usr *model.User) error {
					return nil
				},
			},
			wantStatus: http.StatusOK,
		},
	}
	gin.SetMode(gin.TestMode)
	client := &http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			rg := r.Group("/v1")
			accountService := account.NewAccountService(tt.userRepo, tt.accountRepo, tt.rbac, secret.New())
			service.AccountRouter(accountService, rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users/" + tt.id + "/password"
			req, err := http.NewRequest("PATCH", path, bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
