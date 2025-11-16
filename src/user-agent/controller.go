package user_agent

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"
)

func readUserAgentCtx(ctx context.Context, r *bufio.Reader, propmt string) (string, error) {
	// This function is required for user to cmd + c from terminal in the middle of prompt session
	fmt.Print(propmt)
	lineCh := make(chan struct {
		str string
		err error
	}, 1)

	go func() {
		str, err := r.ReadString('\n')
		lineCh <- struct {
			str string
			err error
		}{str: str, err: err}
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-lineCh:
		if res.err != nil {
			return "", res.err
		}
		return strings.TrimSpace(res.str), nil
	}
}

var UserAgentCache *UserAgent

func SetUserAgentCache(ctx context.Context, db *sql.DB) error {
	userAgent, userAgentErr := GetUserAgent(ctx, db)
	if userAgentErr != nil {
		return userAgentErr
	}

	if userAgent == nil {
		setUserAgentErr := SetUserAgent(ctx, db)
		if setUserAgentErr != nil {
			return fmt.Errorf("failed to set user agent: %w", setUserAgentErr)
		}
	}

	if UserAgentCache == nil {
		UserAgentCache = &UserAgent{}
	}

	UserAgentCache.name = userAgent.name
	UserAgentCache.email = userAgent.email

	return nil
}

func GetUserAgentCache(ctx context.Context, db *sql.DB) (*UserAgent, error) {
	if UserAgentCache == nil {
		userAgentErr := SetUserAgentCache(ctx, db)
		if userAgentErr != nil {
			return nil, userAgentErr
		}
	}

	return UserAgentCache, nil
}

func InitDB(ctx context.Context, db *sql.DB) {
	pingErr := db.Ping()
	if pingErr != nil {
		if strings.HasPrefix(pingErr.Error(), "unable to open database file") {
			mkdirErr := os.Mkdir("./data", 0755)
			if mkdirErr != nil {
				log.Fatal("Failed to create data directory:", mkdirErr)
			}

			newDbPingErr := db.Ping()
			if newDbPingErr != nil {
				log.Println(newDbPingErr)
				os.Exit(1)
			}

			CreateUserAgentTable(db)
			SetUserAgent(ctx, db)
			log.Println("Successfuly created database")

		} else {
			log.Println(pingErr.Error())
			os.Exit(1)
		}
	} else {
		userAgentErr := SetUserAgentCache(ctx, db)
		if userAgentErr != nil {
			log.Fatal(userAgentErr)
		}

		if UserAgentCache == nil {
			SetUserAgent(ctx, db)
			userAgentErr := SetUserAgentCache(ctx, db)

			if UserAgentCache == nil || userAgentErr != nil {
				log.Fatalf("Failed to get user agent from db: %v", userAgentErr)
			}
		}

		log.Printf("Retrieved user agent data: %v", UserAgentCache)
	}
}

func SetUserAgent(ctx context.Context, db *sql.DB) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("User information is required by the SEC to use their API.")

	name, nameErr := readUserAgentCtx(ctx, reader, "Provide your name: ")
	if nameErr != nil {
		log.Printf("Abored with error %v: ", nameErr)
	}

	var email string
	var emailErr error
	for {
		select {
		case <-ctx.Done():
			log.Printf("Aborted: %v", ctx.Err())
			return nil
		default:
		}

		email, emailErr = readUserAgentCtx(ctx, reader, "Porivde email address: ")
		if emailErr != nil {
			log.Println(emailErr)
			continue
		}

		_, emailParseErr := mail.ParseAddress(email)
		if emailParseErr != nil {
			log.Println(emailParseErr)
		} else {
			break
		}

	}

	userAgentErr := InsertUserAgent(ctx, db, name, email)
	if userAgentErr != nil {
		log.Fatalf("failed to set user agent: %v", userAgentErr)
		return nil
	}

	userAgent := fmt.Sprintf(
		"local-market-data (https://github.com/ArturAronov/local-market-data) \n user: %s email: %s",
		name,
		email,
	)
	fmt.Println(userAgent)

	return nil
}

/*

func InitRecords() {
	url := "https://www.sec.gov/files/company_tickers.json"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("User-Agent", "MyFinanceApp (hello@monkey.com)")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Raw JSON length: %d bytes\n", len(body))

	// Option 1: print the raw JSON
	fmt.Println(string(body))

	// Option 2: parse the JSON into a Go map
	var data map[string]map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}

	// Example: print some sample tickers
	for i := 0; i < 5; i++ {
		item := data[fmt.Sprintf("%d", i)]
		fmt.Println(item["ticker"], item["cik_str"], item["title"])
	}
}

*/
