package frapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Hello prints a greeting message.
func Hello() {
    fmt.Println("Hello from frapi!")
}

// Client is the main UI struct that contains currency data and response information.
type Client struct {
	MapOfcurrency map[string]currency
	Resp *responceFromClient
}

// currency represents the structure for JSON decoding of ISO 4217 codes.
type currency struct {
    Code        string
    Name        string
    NumericCode string
}

// exchangeRate represents the structure for JSON decoding of exchange rates.
type exchangeRate struct {
	Code        string  `json:"code"`
	AlphaCode   string  `json:"alphaCode"`
	NumericCode string  `json:"numericCode"`
	Name        string  `json:"name"`
	Rate        float64 `json:"rate"`
	Date        string  `json:"date"`
	InverseRate float64 `json:"inverseRate"`
}

// responseFromClient holds all the information used for gaining exchange rates.
type responceFromClient struct {
	From string
	To string
	Rate float64
	InverseRate float64
	Result string
	Date string
}

// getURL generates the URL for fetching exchange rates based on the alpha code.
func getUrl (alphaCode string) string {
	return fmt.Sprintf("http://www.floatrates.com/daily/%s.json", strings.ToLower(alphaCode))
}


// fetchExchangeRate fetches and decodes the JSON file for exchange rates based on the alpha code.
func FetchOfExchangeRate (alphaCode string) (map[string]exchangeRate, error) {
	var url string = getUrl(alphaCode)
	buffer, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error of GET command: %v", err)
	}
	defer buffer.Body.Close()

	if buffer.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %v", buffer.StatusCode)
	}

	var rates map[string]exchangeRate
	err = json.NewDecoder(buffer.Body).Decode(&rates)
	if err != nil {
		return nil, fmt.Errorf("error decoding of json: %v", err)
	}

	return rates, nil
}


// GetList returns a list of currency codes and their names.
func (obj *Client) GetList () ([]string, error) {
	if obj.MapOfcurrency == nil {
        return nil, fmt.Errorf("Client has no data")
    }

    var codes []string
    for _, currency := range obj.MapOfcurrency {
        codes = append(codes, currency.Code +   " : " + currency.Name)
    }

	return codes, nil
}

// getMap fetches and maps usable ISO 4217 currencies
func (obj *Client) getMap() error {
	buffer, err := FetchOfExchangeRate("usd")
	if err != nil {
		fmt.Printf("Error fetching exchange rates: %v", err)
		return err
	}

	obj.MapOfcurrency = make(map[string]currency)
	
	for _, valueCurr := range buffer {
		obj.MapOfcurrency[valueCurr.AlphaCode] = currency{
			Code:        valueCurr.Code,
			Name:        valueCurr.Name,
			NumericCode: valueCurr.NumericCode,
		}
	}

	obj.MapOfcurrency["USD"] = currency{
		Code: "USD",
		Name: "United States dollar",
		NumericCode: "840",
	}

	return nil
}

// NewClient creates a new Client and initializes its data.
func NewClient () (*Client, error) {
	var objOfClient *Client = &Client{}
	err := objOfClient.getMap()
	if err != nil {
		return nil, err
	}


	return objOfClient, nil
}

// CheckExistenceOfISO checks if an ISO 4217 code is available in FloatRates.
func (obj *Client) CheckExistencefISO (code string) bool {
	_, ok := obj.MapOfcurrency[strings.ToUpper(code)]

	return ok
}

// GetRate gets the exchange rate and inverse rate from one currency to another.
func (objOfClient *Client) GetRate(Fromcurrency, Tocurrency string) error {
	fromCurrency := strings.ToUpper(Fromcurrency)
	toCurrency := strings.ToUpper(Tocurrency)

	if !objOfClient.CheckExistencefISO(fromCurrency) {
		return fmt.Errorf("NON EXISTING FROM ISO CODE : %s", fromCurrency)
	} else if !objOfClient.CheckExistencefISO(toCurrency) {
		return fmt.Errorf("NON EXISTING TO ISO CODE: %s", toCurrency)
	}

	responce, err := FetchOfExchangeRate(fromCurrency)
	if err != nil {
		fmt.Printf("Error of fetching exchange rate: %v", err)
		return err
	}

	if responce == nil {
		return fmt.Errorf("FetchOfExchangeRate returned nil")
	}

	exchangeRateResp, exists := responce[strings.ToLower(Tocurrency)]
	if !exists {
		return fmt.Errorf("NON EXISTING ISO CODE FETCH")
	}

	objOfClient.Resp = &responceFromClient{
		From:        fromCurrency,
		To:          toCurrency,
		Rate:        exchangeRateResp.Rate,
		InverseRate: exchangeRateResp.InverseRate,
		Date:        exchangeRateResp.Date,
		Result: fmt.Sprintf("1 %s = %f %s %s\n1 %s = %f %s %s",
			fromCurrency, exchangeRateResp.Rate, toCurrency, exchangeRateResp.Date,
			toCurrency, exchangeRateResp.InverseRate, fromCurrency, exchangeRateResp.Date),
	}

	return nil
}


// DisplayTheRate prints the exchange rate result.
func (obj *Client) DisplayTheRate() error {
	if obj == nil {
        return fmt.Errorf("client is nil")
    }

	if obj.Resp == nil {
		return fmt.Errorf("responceFromClient is nil")
	}

	fmt.Println(obj.Resp.Result)
	return nil
}

// DisplayTheListOfISO4217 prints the list of ISO 4217 codes and their names.
func (obj *Client) DisplayTheListOfISO4217() error {
	if obj == nil {
        return fmt.Errorf("client is nil")
    }

	for _, v := range obj.MapOfcurrency {
		fmt.Printf("%s : %s\n", v.Code, v.Name)
	}

	return nil
}