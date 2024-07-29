package frapi

import (
	"fmt"
)


// Example function that creates client object and calls its methods
func Example_test() {
	clientAPI, err := NewClient() // Initialize client
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	// Get and print the list of ISO4217 codes
	list, err := clientAPI.GetList()
	if err != nil {
		fmt.Printf("Error getting list: %v\n", err)
		return
	}
	fmt.Println(list)

	// Display the list of ISO4217 codes
	err = clientAPI.DisplayTheListOfISO4217()
	if err != nil {
		fmt.Printf("Error displaying ISO4217 list: %v\n", err)
		return
	}

	// Get the exchange rate
	err = clientAPI.GetRate("USD", "KZT")
	if err != nil {
		fmt.Printf("Error getting rate: %v\n", err)
		return
	}

	// Display the exchange rate
	if clientAPI.Resp != nil {
		err = clientAPI.DisplayTheRate()
		if err != nil {
			fmt.Printf("Error displaying rate: %v\n", err)
		}
	} else {
		fmt.Println("Exchange rate response object is nil")
	}
}
