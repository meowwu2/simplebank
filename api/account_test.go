package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	mockdb "simplebank/db/mock"
	db "simplebank/db/sqlc"
	"simplebank/util"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAccountAPI(t *testing.T) {
	account := randomAccount()
	testCase := []struct{
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{{
		name:"ok",
		accountID: account.ID,
		buildStubs: func(store *mockdb.MockStore) {
			store.EXPECT().
			GetAccount(gomock.Any(),gomock.Eq(account.ID)).
			Times(1).
			Return(account,nil)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t,http.StatusOK,recorder.Code)
			requireBodyMatchAccount(t,recorder.Body,account)
		},
	}}
	for i:=range testCase{
		tc :=testCase[i]
		t.Run(tc.name,func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d",account.ID)
			request,err := http.NewRequest(http.MethodGet,url,nil)
			require.NoError(t,err)
			server.router.ServeHTTP(recorder,request)
			tc.checkResponse(t,recorder)
		})
		
	}
	
	
}

func randomAccount() db.Account {
	return db.Account{
		ID: util.RandomInt(1,1000),
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T,body *bytes.Buffer,account db.Account){
	data,err := io.ReadAll(body)
	require.NoError(t,err)
	var gotaccount db.Account
	err =json.Unmarshal(data,&gotaccount)
	require.NoError(t,err)
	require.Equal(t,account,gotaccount)
}