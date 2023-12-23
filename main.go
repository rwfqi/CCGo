package main

import (
	"html/template"
	"net/http"
	"strconv"
)

// Data struct digunakan untuk menyimpan data yang akan digunakan dalam template HTML
type Data struct {
	Amount     float64
	From       string
	To         string
	Result     float64
	Currencies []string
}

func convert(amount float64, from string, to string) float64 {
	// Definisikan nilai tukar mata uang
	exchangeRates := map[string]float64{
		"USD": 1.0,
		"EUR": 0.85,
		"GBP": 0.75,
		"JPY": 110.0,
		"IDR": 14000.0, // Angka ini bersifat contoh, sesuaikan dengan nilai tukar terkini
	}

	// Periksa apakah mata uang yang diinginkan terdaftar dalam nilai tukar
	if rate, ok := exchangeRates[to]; ok {
		// Hitung hasil konversi
		result := amount * rate
		return result
	}

	// Default: jika mata uang tidak ditemukan, kembalikan 0
	return 0.0
}

// handlerCurrencyConverter menghandle permintaan HTTP untuk konverter mata uang
func handlerCurrencyConverter(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Ambil data dari formulir HTML
		amountStr := r.FormValue("amount")
		fromCurrency := r.FormValue("from")
		toCurrency := r.FormValue("to")

		// Konversi string ke float64
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}

		// Hitung hasil konversi
		result := convert(amount, fromCurrency, toCurrency)

		// Persiapkan data untuk ditampilkan di template HTML
		data := Data{
			Amount:     amount,
			From:       fromCurrency,
			To:         toCurrency,
			Result:     result,
			Currencies: []string{"USD", "EUR", "GBP", "JPY", "IDR"}, // Daftar mata uang yang didukung
		}

		// Render template HTML dengan data yang disiapkan
		tmpl, err := template.New("index").Parse(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
		  <meta charset="utf-8">
		  <meta name="viewport" content="width=device-width, initial-scale=1">
		  <title>Currency Converter</title>
		  <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-T3c6CoIi6uLrA9TneNEoa7RxnatzjcDSCmG1MXxSR1GAsXEV/Dwwykc2MPK8M2HN" crossorigin="anonymous">
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
				  <option value="{{.}}" {{if eq . $.From}}selected{{end}}>{{.}}</option>
				{{end}}
			  </select>
			</div>
			<div class="mb-3">
			  <label for="to" class="form-label">To Currency:</label>
			  <select class="form-select" id="to" name="to">
				{{range .Currencies}}
				  <option value="{{.}}" {{if eq . $.To}}selected{{end}}>{{.}}</option>
				{{end}}
			  </select>
			</div>
			<button type="submit" class="btn btn-primary">Convert</button>
		  </form>
		  <p class="mt-3">Result: {{.Result}}</p>
		  
		  <!-- Bootstrap JS (optional, if you need it) -->
		  <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/js/bootstrap.bundle.min.js" integrity="sha384-C6RzsynM9kWDrMNeT87bh95OGNyZPhcTNXj1NW7RuBCsyN/o0jlpcV8Qyq46cDfL" crossorigin="anonymous"></script>
		</body>
		</html>
		
		`)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Render template HTML ke ResponseWriter
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		// Jika bukan metode POST, tampilkan formulir kosong
		data := Data{
			Currencies: []string{"USD", "EUR", "GBP", "JPY", "IDR"},
		}

		tmpl, err := template.New("index").Parse(`
			<!DOCTYPE html>
			<html>
				<head>
					<title>Currency Converter</title>
				</head>
				<body>
					<h1>Currency Converter</h1>
					<form method="post">
						<label>Amount: </label>
						<input type="text" name="amount" required><br>
						<label>From Currency: </label>
						<select name="from">
							{{range .Currencies}}
								<option value="{{.}}">{{.}}</option>
							{{end}}
						</select><br>
						<label>To Currency: </label>
						<select name="to">
							{{range .Currencies}}
								<option value="{{.}}">{{.}}</option>
							{{end}}
						</select><br>
						<input type="submit" value="Convert">
					</form>
				</body>
			</html>
		`)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Render template HTML ke ResponseWriter
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	// Setelah menjalankan program, buka browser dan kunjungi http://localhost:8080/
	http.HandleFunc("/", handlerCurrencyConverter)
	http.ListenAndServe(":8080", nil)
}
