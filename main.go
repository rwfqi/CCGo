package main

import (
	"html/template"
	"net/http"
	"strconv"
)

type Data struct {
	Amount     float64
	From       string
	To         string
	Result     float64
	Currencies []string
}

func currencySymbol(currency string) string {
	switch currency {
	case "USD":
		return "$"
	case "IDR":
		return "Rp."
	case "EUR":
		return "€"
	case "GBP":
		return "£"
	case "JPY":
		return "¥"
	default:
		return currency
	}
}

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.New("index").Funcs(template.FuncMap{"currencySymbol": currencySymbol}).Parse(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Currency Converter</title>
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-T3c6CoIi6uLrA9TneNEoa7RxnatzjcDSCmG1MXxSR1GAsXEV/Dwwykc2MPK8M2HN" crossorigin="anonymous">
		<style>
			body {
				background-color: #f8f9fa;
			}

			.container {
				max-width: 600px;
				background-color: #fff;
				padding: 20px;
				border-radius: 8px;
				box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
			}

			h1 {
				color: #007bff;
			}

			form {
				margin-top: 20px;
			}

			.btn-primary {
				background-color: #007bff;
				border-color: #007bff;
			}

			.btn-primary:hover {
				background-color: #0056b3;
				border-color: #0056b3;
			}

			.mt-3 {
				margin-top: 1rem;
			}
		</style>
		</head>
	<body class="container mt-5">
		<h1 class="mb-4">Currency Converter</h1>
		<form method="post">
			<div class="mb-3">
				<label for="amount" class="form-label">Amount:</label>
				<input type="text" class="form-control" id="amount" name="amount" value="{{.Amount}}" required>
			</div>
			<div class="mb-3">
				<label for="from" class="form-label">From Currency:</label>
				<select class="form-select" id="from" name="from">
					{{range .Currencies}}
						<option value="{{.}}" {{if eq . $.From}}selected{{end}}>{{.}} ({{currencySymbol .}})</option>
					{{end}}
				</select>
			</div>
			<div class="mb-3">
				<label for="to" class="form-label">To Currency:</label>
				<select class="form-select" id="to" name="to">
					{{range .Currencies}}
						<option value="{{.}}" {{if eq . $.To}}selected{{end}}>{{.}} ({{currencySymbol .}})</option>
					{{end}}
				</select>
			</div>
			<button type="submit" class="btn btn-primary">Convert</button>
		</form>
		<p class="mt-3">Result: {{currencySymbol .To}}{{.Result}}</p>
	</body>
	</html>
	`)
	if err != nil {
		panic(err)
	}
}

func convertFromUSD(amount float64, toCurrency string) float64 {
	exchangeRates := map[string]float64{
		"USD": 1.0,
		"EUR": 0.85,
		"GBP": 0.75,
		"JPY": 110.0,
		"IDR": 14000.0,
	}

	if rate, ok := exchangeRates[toCurrency]; ok {
		result := amount * rate
		return result
	}

	return 0.0
}

func convertFromGBP(amount float64, toCurrency string) float64 {
	exchangeRates := map[string]float64{
		"USD": 1.33,
		"EUR": 1.18,
		"GBP": 1.0,
		"JPY": 145.72,
		"IDR": 18400.0,
	}

	if rate, ok := exchangeRates[toCurrency]; ok {
		result := amount * rate
		return result
	}

	return 0.0
}

func convertFromIDR(amount float64, toCurrency string) float64 {
	exchangeRates := map[string]float64{
		"USD": 0.0000714,
		"EUR": 0.0000679,
		"GBP": 0.0000543,
		"JPY": 0.0079,
		"IDR": 1.0,
	}

	if rate, ok := exchangeRates[toCurrency]; ok {
		result := amount * rate
		return result
	}

	return 0.0
}

func convertFromEUR(amount float64, toCurrency string) float64 {
	exchangeRates := map[string]float64{
		"USD": 1.18,
		"EUR": 1.0,
		"GBP": 0.85,
		"JPY": 125.5,
		"IDR": 14700.0,
	}

	if rate, ok := exchangeRates[toCurrency]; ok {
		result := amount * rate
		return result
	}

	return 0.0
}

func convertFromJPY(amount float64, toCurrency string) float64 {
	exchangeRates := map[string]float64{
		"USD": 0.0091,
		"EUR": 0.00797,
		"GBP": 0.00688,
		"JPY": 1.0,
		"IDR": 127.0,
	}

	if rate, ok := exchangeRates[toCurrency]; ok {
		result := amount * rate
		return result
	}

	return 0.0
}

func convert(amount float64, from string, to string) float64 {
	switch from {
	case "USD":
		return convertFromUSD(amount, to)
	case "GBP":
		return convertFromGBP(amount, to)
	case "IDR":
		return convertFromIDR(amount, to)
	case "EUR":
		return convertFromEUR(amount, to)
	case "JPY":
		return convertFromJPY(amount, to)
	default:
		return 0.0
	}
}

func handlerCurrencyConverter(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method == http.MethodPost {
		amountStr := r.FormValue("amount")
		fromCurrency := r.FormValue("from")
		toCurrency := r.FormValue("to")

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}

		result := convert(amount, fromCurrency, toCurrency)

		data := Data{
			Amount:     amount,
			From:       fromCurrency,
			To:         toCurrency,
			Result:     result,
			Currencies: []string{"USD", "EUR", "GBP", "JPY", "IDR"},
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		data := Data{
			Currencies: []string{"USD", "EUR", "GBP", "JPY", "IDR"},
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	http.HandleFunc("/", handlerCurrencyConverter)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
