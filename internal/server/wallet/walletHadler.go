package wallet

import (
	walletDomain "exampleApp/internal/domain/wallet"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletService interface {
	GetBalanceByID(string) (float64, error)
	ChageAmount(walletDomain.WalletRequest) (walletDomain.Wallet, error)
}

type WalletHandler struct {
	walletService WalletService
}

func NewWalletHandler(walletService WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

func (wh *WalletHandler) GetBalanceByID(ctx *gin.Context) {
	uuid := ctx.Param("id")
	balance, err := wh.walletService.GetBalanceByID(uuid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"amount": balance})
}

func (wh *WalletHandler) ChageAmount(ctx *gin.Context) {
	var req walletDomain.WalletRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	w, err := wh.walletService.ChageAmount(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, w)
}
