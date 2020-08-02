package tinkoff

import (
	"context"
	"log"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
)

// GetAllMarketsMap creates map of all Tinkoff markets where name is Key
func GetAllMarketsMap(token string) (instrumentNameMap map[string]sdk.Instrument) {
	// Create REST Client
	client := sdk.NewRestClient(token)

	// Init map
	instrumentNameMap = make(map[string]sdk.Instrument)

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Get currency instruments")
	// Example: USD000UTSTOM - USD, EUR_RUB__TOM - EUR
	instruments, err := client.Currencies(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	// Add markets into map
	for _, instrument := range instruments {
		instrumentNameMap[instrument.Name] = instrument
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Fet etf instruments")
	// Example: FXMM - Казначейские облигации США, FXGD - золото
	instruments, err = client.ETFs(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	// Add markets into map
	for _, instrument := range instruments {
		instrumentNameMap[instrument.Name] = instrument
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Get obligation instruments")
	// Example: SU24019RMFS0 - ОФЗ 24019
	instruments, err = client.Bonds(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	// Add markets into map
	for _, instrument := range instruments {
		instrumentNameMap[instrument.Name] = instrument
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Add stock instruments")
	// Example: SBUX - Starbucks Corporation
	instruments, err = client.Stocks(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	// Add markets into map
	for _, instrument := range instruments {
		instrumentNameMap[instrument.Name] = instrument
	}

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
