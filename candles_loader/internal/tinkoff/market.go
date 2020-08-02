package tinkoff

import (
	"context"
	"log"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
)

// GetAllMarkets creates list of all Tinkoff markets
func GetAllMarkets(token string) (instruments []sdk.Instrument, err error) {
	// Create REST Client
	client := sdk.NewRestClient(token)

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Get currency instruments")
	// Example: USD000UTSTOM - USD, EUR_RUB__TOM - EUR
	currencies, err := client.Currencies(ctx)
	if err != nil {
		return nil, err
	}
	// Add currencies
	instruments = append(instruments, currencies...)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Fet etf instruments")
	// Example: FXMM - Казначейские облигации США, FXGD - золото
	etfs, err := client.ETFs(ctx)
	if err != nil {
		return nil, err
	}
	// Add etfs
	instruments = append(instruments, etfs...)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Get bond instruments")
	// Example: SU24019RMFS0 - ОФЗ 24019
	bonds, err := client.Bonds(ctx)
	if err != nil {
		return nil, err
	}
	// Add bonds
	instruments = append(instruments, bonds...)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Add stock instruments")
	// Example: SBUX - Starbucks Corporation
	stocks, err := client.Stocks(ctx)
	if err != nil {
		return nil, err
	}
	// Add stocks
	instruments = append(instruments, stocks...)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return
}

// GetFigiByTicket receive FIGI of instruments by ticket
func GetFigiByTicket(token string, ticket string) {
	// Create REST Client
	client := sdk.NewRestClient(token)

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение инструмента по тикеру, возвращает массив инструментов потому что тикер уникален только в рамках одной биржи
	// но может совпадать на разных биржах у разных кампаний
	// Например: https://www.moex.com/ru/issue.aspx?code=FIVE и https://www.nasdaq.com/market-activity/stocks/FIVE
	// В этом примере получить нужную компанию можно проверив поле Currency
	instruments, _ := client.InstrumentByTicker(ctx, ticket)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	log.Printf("Instrument for tiket '%s': %+v\n", ticket, instruments)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}
