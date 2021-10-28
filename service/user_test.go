package service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/mock"
	"github.com/alpacahq/ribbit-backend/mock/mockdb"
	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/repository/user"
	"github.com/alpacahq/ribbit-backend/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestListUsers(t *testing.T) {
	type listResponse struct {
		Users []model.User `json:"users"`
		Page  int          `json:"page"`
	}
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *listResponse
		userRepo   *mockdb.User
		rbac       *mock.RBAC
		auth       *mock.Auth
	}{
		{
			name:       "Invalid request",
			req:        `?limit=2222&page=-1`,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "Fail on query list",

			req: `?limit=100&page=1`,

			auth: &mock.Auth{
				UserFn: func(c *gin.Context) *model.AuthUser {

					return &model.AuthUser{

						ID:    1,
						Role:  model.UserRole,
						Email: "john@mail.com",
					}
				}},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",

			req: `?limit=100&page=1`,

			auth: &mock.Auth{
				UserFn: func(c *gin.Context) *model.AuthUser {

					return &model.AuthUser{

						ID:    1,
						Role:  model.SuperAdminRole,
						Email: "john@mail.com",
					}
				}},
			userRepo: &mockdb.User{
				ListFn: func(q *model.ListQuery, p *model.Pagination) ([]model.User, error) {

					if p.Limit == 100 && p.Offset == 100 {

						return []model.User{

							{
								Base: model.Base{
									CreatedAt: mock.TestTime(2001),
									UpdatedAt: mock.TestTime(2002),
								},
								ID:        10,
								FirstName: "John",

								LastName: "Doe",

								Email:               "john@mail.com",
								AccountID:           "",
								AccountNumber:       "",
								AccountCurrency:     "",
								AccountStatus:       "",
								DOB:                 "",
								City:                "",
								State:               "",
								Country:             "",
								TaxIDType:           "",
								TaxID:               "",
								FundingSource:       "",
								EmploymentStatus:    "",
								InvestingExperience: "",
								PublicShareholder:   "",
								AnotherBrokerage:    "",
								DeviceID:            "",
								ProfileCompletion:   "",
								ReferralCode:        "",

								Role: &model.Role{
									ID: 1,

									AccessLevel: 1,

									Name: "SUPER_ADMIN",
								},
							},
							{
								Base: model.Base{
									CreatedAt: mock.TestTime(2004),
									UpdatedAt: mock.TestTime(2005),
								},
								ID:        11,
								FirstName: "Joanna",

								LastName: "Dye",

								Email: "joanna@mail.com",

								AccountID:           "",
								AccountNumber:       "",
								AccountCurrency:     "",
								AccountStatus:       "",
								DOB:                 "",
								City:                "",
								State:               "",
								Country:             "",
								TaxIDType:           "",
								TaxID:               "",
								FundingSource:       "",
								EmploymentStatus:    "",
								InvestingExperience: "",
								PublicShareholder:   "",
								AnotherBrokerage:    "",
								DeviceID:            "",
								ProfileCompletion:   "",
								ReferralCode:        "",

								Role: &model.Role{
									ID: 2,

									AccessLevel: 2,

									Name: "ADMIN",
								},
							},
						}, nil

					}
					return nil, apperr.DB

				},
			},
			wantStatus: http.StatusOK,
			wantResp: &listResponse{
				Users: []model.User{
					{
						Base: model.Base{
							CreatedAt: mock.TestTime(2001),
							UpdatedAt: mock.TestTime(2002),
						},
						ID:        10,
						FirstName: "John",

						LastName: "Doe",

						Email: "john@mail.com",

						AccountID:           "",
						AccountNumber:       "",
						AccountCurrency:     "",
						AccountStatus:       "",
						DOB:                 "",
						City:                "",
						State:               "",
						Country:             "",
						TaxIDType:           "",
						TaxID:               "",
						FundingSource:       "",
						EmploymentStatus:    "",
						InvestingExperience: "",
						PublicShareholder:   "",
						AnotherBrokerage:    "",
						DeviceID:            "",
						ProfileCompletion:   "",
						ReferralCode:        "",

						Role: &model.Role{
							ID: 1,

							AccessLevel: 1,

							Name: "SUPER_ADMIN",
						},
					},
					{
						Base: model.Base{
							CreatedAt: mock.TestTime(2004),
							UpdatedAt: mock.TestTime(2005),
						},
						ID:        11,
						FirstName: "Joanna",

						LastName: "Dye",

						Email: "joanna@mail.com",

						AccountID:           "",
						AccountNumber:       "",
						AccountCurrency:     "",
						AccountStatus:       "",
						DOB:                 "",
						City:                "",
						State:               "",
						Country:             "",
						TaxIDType:           "",
						TaxID:               "",
						FundingSource:       "",
						EmploymentStatus:    "",
						InvestingExperience: "",
						PublicShareholder:   "",
						AnotherBrokerage:    "",
						DeviceID:            "",
						ProfileCompletion:   "",
						ReferralCode:        "",

						Role: &model.Role{
							ID: 2,

							AccessLevel: 2,

							Name: "ADMIN",
						},
					},
				}, Page: 1},
		},
	}
	gin.SetMode(gin.TestMode)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			rg := r.Group("/v1")
			userService := user.NewUserService(tt.userRepo, tt.auth, tt.rbac)
			service.UserRouter(userService, rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(listResponse)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestViewUser(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *model.User
		udb        *mockdb.User
		rbac       *mock.RBAC
		auth       *mock.Auth
	}{
		{
			name:       "Invalid request",
			req:        `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			req:  `1`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(*gin.Context, int) bool {
					return false
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `1`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(*gin.Context, int) bool {
					return true
				},
			},
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					return &model.User{
						Base: model.Base{
							CreatedAt: mock.TestTime(2000),
							UpdatedAt: mock.TestTime(2000),
						},
						ID:        1,
						FirstName: "John",
						LastName:  "Doe",
						Username:  "JohnDoe",
					}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &model.User{
				Base: model.Base{
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2000),
				},
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Username:  "JohnDoe",
			},
		},
	}
	gin.SetMode(gin.TestMode)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			rg := r.Group("/v1")
			userService := user.NewUserService(tt.udb, tt.auth, tt.rbac)
			service.UserRouter(userService, rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users/" + tt.req
			res, err := http.Get(path)
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

func TestUpdateUser(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		id         string
		wantStatus int
		wantResp   *model.User
		udb        *mockdb.User
		rbac       *mock.RBAC
		auth       *mock.Auth
	}{
		{
			name:       "Invalid request",
			id:         `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			id:   `1`,
			req:  `{"first_name":"jj","last_name":"okocha","mobile":"123456","phone":"321321","address":"home"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(*gin.Context, int) bool {
					return false
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			id:   `1`,
			req:  `{"first_name":"jj","last_name":"okocha","phone":"321321","address":"home"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(*gin.Context, int) bool {
					return true
				},
			},
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					return &model.User{
						Base: model.Base{
							CreatedAt: mock.TestTime(2000),
							UpdatedAt: mock.TestTime(2000),
						},
						ID:        1,
						FirstName: "John",
						LastName:  "Doe",
						Username:  "JohnDoe",
						Address:   "Work",
						Mobile:    "332223",
					}, nil
				},
				UpdateFn: func(usr *model.User) (*model.User, error) {
					usr.UpdatedAt = mock.TestTime(2010)
					usr.Mobile = "991991"
					return usr, nil
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &model.User{
				Base: model.Base{
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2010),
				},
				ID:        1,
				FirstName: "jj",
				LastName:  "okocha",
				Username:  "JohnDoe",
				Address:   "home",
				Mobile:    "991991",
			},
		},
	}
	gin.SetMode(gin.TestMode)
	client := http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			rg := r.Group("/v1")
			userService := user.NewUserService(tt.udb, tt.auth, tt.rbac)
			service.UserRouter(userService, rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users/" + tt.id
			req, _ := http.NewRequest("PATCH", path, bytes.NewBufferString(tt.req))
			res, err := client.Do(req)
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

func TestDeleteUser(t *testing.T) {
	cases := []struct {
		name       string
		id         string
		wantStatus int
		udb        *mockdb.User
		rbac       *mock.RBAC
		auth       *mock.Auth
	}{
		{
			name:       "Invalid request",
			id:         `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			id:   `1`,
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					return &model.User{
						Role: &model.Role{
							AccessLevel: model.SuperAdminRole,
						},
					}, nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(*gin.Context, model.AccessRole) bool {
					return false
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			id:   `1`,
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					return &model.User{
						Role: &model.Role{
							AccessLevel: model.SuperAdminRole,
						},
					}, nil
				},
				DeleteFn: func(*model.User) error {
					return nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(*gin.Context, model.AccessRole) bool {
					return true
				},
			},
			wantStatus: http.StatusOK,
		},
	}
	gin.SetMode(gin.TestMode)
	client := http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			rg := r.Group("/v1")
			userService := user.NewUserService(tt.udb, tt.auth, tt.rbac)
			service.UserRouter(userService, rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users/" + tt.id
			req, _ := http.NewRequest("DELETE", path, nil)
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
