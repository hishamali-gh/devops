package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"  // The 'http' module in the 'net' (network) package
	"os"
	"time"
)

type HealthResponse struct {  // Create a custom 'type' HealthResponse from the in-built type 'struct'
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

func main() {
	// 1. Read configuration from Environment Variables:
	// DevOps principle: never hardcode credentials, even port, in code!

	var port string = os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	var appVersion string = os.Getenv("APP_VERSION")

	if appVersion == "" {
		appVersion = "1.0.0"
	}

	// 2. The Health Check endpoint (/health):
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {  // *I, *II
		/* *I. If someone visits the provided endpoint pattern, handle the provided function.

		Here, in this case, if someone visits '/health' of the specified port (8000, locally) of our IP address (localhost, in this case)
		handle the provided function.

		NOTICE WE HAVEN'T CLOSE THE PARENTHESES OF THE HandleFunc FUNCTION UNTIL LAST?

		It's because we provided an anonymous function in line.

		We could also do that by creating a regular named function 'healthHandler' and just put there directly,
		like 'http.HandleFunc("/health", healthHandler)'. */

		/* *II. A ResponseWriter is not a passive 'struct'; it is an active connection pipe 
		between the Go application and the OS network socket, wired underneath.

		In Go, it is an 'interface' with exactly three methods:
			1. Header:
				Returns a 'Header' map (under the hood: map[string][]string).
				It represents the HTTP headers (metadata) of the response.

				We use its .Set() method to configure headers before sending data.

				Example: w.Header().Set("Content-Type", "application/json")

				Analogy: Writing labels on the outside of an envelope (e.g., "Contains: JSON").

			2. WriteHeader(statusCode int):
				Sends the HTTP status code (e.g., 200, 404, 500) to the client.

				We typically pass Go's built-in constants like 'http.StatusOK'.

				*CRUCIAL RULE*: Once WriteHeader is called, headers are locked! 
				The envelope is officially sealed and sent down the wire.

				Analogy: Stamping the official status seal on the envelope ("Approved / 200 OK").

			3. Write([]byte) (int, error):
				Takes raw bytes and pumps them directly down the network socket to the client.

				Because it has this method, it satisfies Go's universal 'io.Writer' contract.
				This is why helper tools like 'json.NewEncoder(w)' and 'fmt.Fprintf(w, ...)' 
				can stream data directly into 'w'.

				Analogy: Pumping the actual letter/content through the pipe to the client. */

		var res HealthResponse = HealthResponse{
			Status: "healthy",
			Timestamp: time.Now().UTC().Format(time.RFC3339),  // Turn 'time.Now()' to a string using the formatter
			Version: appVersion,  // A comma, even at the last line, is strictly required.
		}

		w.Header().Set("Content-Type", "application/json")  // Set the header
		w.WriteHeader(http.StatusOK)  // Stamp it with a status code

		json.NewEncoder(w).Encode(res)  // *III

		/* *III. As mentioned earlier, a ResponseWriter interface has a method 'Write' to insert data in.
		But it requires 'bytes', not strings or any other kind.

		We could manually convert text to bytes and call w.Write() yourself, but writing raw byte arrays by hand is tedious.

		json.NewEncoder creates an Encoder object that holds a reference to whatever ResponseWriter it is bound to,
		and using the method 'Encode' it does both jobs: 'write' the data by converting them to bytes and
		ship the data as output through the ResponseWriter connection pipe. */

		// Print request log to terminal (stdout)
		log.Printf("[INFO] %s %s - 200 OK\n", r.Method, r.URL.Path)  // 'log' is similar to 'fmt', but it automatically adds the time, date, and timestamp.
	})

	// 3. The Root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)

			return
		}

		fmt.Fprintf(w, "DevOps Go Microservice is up and running!\n")
		log.Printf("[INFO] %s %s - 200 OK\n", r.Method, r.URL.Path)
	})

	// 4. Start the HTTP server
	log.Printf("[INFO] Server starting on port %s (Version: %s)...\n", port, appVersion)

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatalf("[FATAL] Server failed: %s\n", err)
	}
}
