package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/Viczdera/bank-gorm/db/sqlc"
	"github.com/Viczdera/bank-gorm/token"
	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccount int64 `json:"from_account" binding:"required"`
	ToAccount   int64 `json:"to_account" binding:"required"`
	// Currency    string `json:"currency" binding:"required,oneof=USD EUR"`
	Currency string `json:"currency" binding:"required,currency"`
	Amount   int64  `json:"amount" binding:"required"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccount, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(AUTH_PAYLOAD_KEY).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccount, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccount: fromAccount.ID,
		ToAccount:   req.ToAccount,
		Amount:      req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return account, false
	}
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return account, false
	}
	return account, true
}
