// Copyright 2022 Board of Trustees of the University of Illinois.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"fmt"
	"log"
	"net/http"
	"wellness/core"
	"wellness/core/model"

	"github.com/rokwire/core-auth-library-go/authorization"
	"github.com/rokwire/core-auth-library-go/authservice"
	"github.com/rokwire/core-auth-library-go/tokenauth"
	"github.com/rokwire/logging-library-go/errors"
	"github.com/rokwire/logging-library-go/logutils"
)

//Authorization is an interface for auth types
type Authorization interface {
	check(req *http.Request) (int, *tokenauth.Claims, error)
	start()
}

// Auth handler
type Auth struct {
	coreAuth     *CoreAuth
	internalAuth *InternalAuth
}

// NewAuth creates new auth handler
func NewAuth(app *core.Application, config model.Config, authService *authservice.AuthService) *Auth {
	coreAuth := NewCoreAuth(app, config, authService)
	internalAuth := newInternalAuth(config.InternalAPIKey)

	auth := Auth{coreAuth: coreAuth, internalAuth: internalAuth}
	return &auth
}

// CoreAuth implementation
type CoreAuth struct {
	app       *core.Application
	tokenAuth *tokenauth.TokenAuth

	permissionsAuth *PermissionsAuth
	userAuth        *UserAuth
	standardAuth    *StandardAuth
}

// NewCoreAuth creates new CoreAuth
func NewCoreAuth(app *core.Application, config model.Config, authService *authservice.AuthService) *CoreAuth {

	adminPermissionAuth := authorization.NewCasbinAuthorization("driver/web/authorization_model.conf",
		"driver/web/authorization_policy.csv")
	tokenAuth, err := tokenauth.NewTokenAuth(true, authService, adminPermissionAuth, nil)
	if err != nil {
		log.Fatalf("Error intitializing token auth: %v", err)
	}
	permissionsAuth := newPermissionsAuth(tokenAuth)
	usersAuth := newUserAuth(tokenAuth)
	standardAuth := newStandardAuth(tokenAuth)

	auth := CoreAuth{app: app, tokenAuth: tokenAuth, permissionsAuth: permissionsAuth,
		userAuth: usersAuth, standardAuth: standardAuth}
	return &auth
}

//PermissionsAuth entity
//This enforces that the user has permissions matching the policy
type PermissionsAuth struct {
	tokenAuth *tokenauth.TokenAuth
}

func (a *PermissionsAuth) start() {}

func (a *PermissionsAuth) check(req *http.Request) (int, *tokenauth.Claims, error) {
	claims, err := a.tokenAuth.CheckRequestTokens(req)
	if err != nil {
		return http.StatusUnauthorized, nil, errors.WrapErrorAction("typeCheckServicesAuthRequestToken", logutils.TypeToken, nil, err)
	}

	if err == nil && claims != nil {
		err = a.tokenAuth.AuthorizeRequestPermissions(claims, req)
		if err != nil {
			return http.StatusForbidden, nil, errors.WrapErrorAction("check permissions", logutils.TypeRequest, nil, err)
		}
	}

	return http.StatusOK, claims, err
}

func newPermissionsAuth(tokenAuth *tokenauth.TokenAuth) *PermissionsAuth {
	permissionsAuth := PermissionsAuth{tokenAuth: tokenAuth}
	return &permissionsAuth
}

//UserAuth entity
// This enforces that the user is not anonymous
type UserAuth struct {
	tokenAuth *tokenauth.TokenAuth
}

func (a *UserAuth) start() {}

func (a *UserAuth) check(req *http.Request) (int, *tokenauth.Claims, error) {
	claims, err := a.tokenAuth.CheckRequestTokens(req)
	if err != nil {
		return http.StatusUnauthorized, nil, errors.WrapErrorAction("typeCheckServicesAuthRequestToken", logutils.TypeToken, nil, err)
	}

	if err == nil && claims != nil {
		if claims.Anonymous {
			return http.StatusForbidden, nil, errors.New("token must not be anonymous")
		}
	}

	return http.StatusOK, claims, err
}

func newUserAuth(tokenAuth *tokenauth.TokenAuth) *UserAuth {
	userAuth := UserAuth{tokenAuth: tokenAuth}
	return &userAuth
}

//StandardAuth entity
// This enforces standard auth check
type StandardAuth struct {
	tokenAuth *tokenauth.TokenAuth
}

func (a *StandardAuth) start() {}

func (a *StandardAuth) check(req *http.Request) (int, *tokenauth.Claims, error) {
	claims, err := a.tokenAuth.CheckRequestTokens(req)
	if err != nil {
		return http.StatusUnauthorized, nil, errors.WrapErrorAction("typeCheckServicesAuthRequestToken", logutils.TypeToken, nil, err)
	}

	return http.StatusOK, claims, err
}

func newStandardAuth(tokenAuth *tokenauth.TokenAuth) *StandardAuth {
	standartAuth := StandardAuth{tokenAuth: tokenAuth}
	return &standartAuth
}

//InternalAuth entity
// This enforces Internal API Key auth check
type InternalAuth struct {
	internalAPIKey string
}

func (a *InternalAuth) check(req *http.Request) (int, error) {
	internalAPIKey := req.Header.Get("INTERNAL-API-KEY")
	if internalAPIKey != a.internalAPIKey {
		return http.StatusUnauthorized, errors.WrapErrorAction("typeCheckServicesInternalAPIKey", logutils.TypeRequest, nil, fmt.Errorf("wrong or missing INTERNAL-API-KEY request header"))
	}

	return http.StatusOK, nil
}

func newInternalAuth(internalAPIKey string) *InternalAuth {
	auth := InternalAuth{internalAPIKey: internalAPIKey}
	return &auth
}
