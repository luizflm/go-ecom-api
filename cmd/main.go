package main

func main() {
	cfg := config{
		addr: ":8080",
		db:   dbConfig{},
	}

	app := application{
		config: cfg,
	}
}
