package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	db "github.com/Viczdera/bank-gorm/db/sqlc"
	"github.com/Viczdera/bank-gorm/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	//get reqBody
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(AUTH_PAYLOAD_KEY).(*token.Payload)
	// 	ctx.MustGet() is a method from the Gin framework that retrieves a value from the context. Unlike the regular Get() method:

	// It panics if the key doesn't exist
	// It's used when you're absolutely certain the value should be there
	// The key AUTH_PAYLOAD_KEY is presumably set by an authentication middleware
	// .(*token.Payload) is performing a type assertion:

	// MustGet returns an interface{}
	// The type assertion converts it to a *token.Payload
	// If the assertion fails, this will panic
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			codeName := pqErr.Code.Name()
			log.Println(codeName)
			switch codeName {
			case " foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	authPayload := ctx.MustGet(AUTH_PAYLOAD_KEY).(*token.Payload)
	// 	ctx.MustGet() is a method from the Gin framework that retrieves a value from the context. Unlike the regular Get() method:

	// It panics if the key doesn't exist
	// It's used when you're absolutely certain the value should be there
	// The key AUTH_PAYLOAD_KEY is presumably set by an authentication middleware
	// .(*token.Payload) is performing a type assertion:

	// MustGet returns an interface{}
	// The type assertion converts it to a *token.Payload
	// If the assertion fails, this will panic
	account, err := server.store.GetAccount(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if account.Owner != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, errResponse(fmt.Errorf("account does not belong to the authenticated user")))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountsListRequest struct {
	CurrentPage int32 `form:"current_page" binding:"required,min=1"`
	PageSize    int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) getAccountsList(ctx *gin.Context) {
	var req getAccountsListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(AUTH_PAYLOAD_KEY).(*token.Payload)

	args := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.CurrentPage - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, args)

	// log.Printf("Accounts: %+v", accounts)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
