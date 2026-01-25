TODO:

- Get all CIK, name, ticker info & store in DB
- Ask user if they want to fetch all market data at once, cache company data, or no cache
- Create update logic

## Requirements:

- Go 1.24
- SQLite
- GCC (used by go-sqlite3)

## Terminal commands

- init -e / --email email@email.com
- ci -t / --ticker amazon
- ci -c / --cik 1018724

SEC endpoints:

- https://data.sec.gov/submissions/CIK0000320193.json (reports data)
- https://data.sec.gov/api/xbrl/companyfacts/CIK0000320193.json (all us-gaap info)
- https://data.sec.gov/api/xbrl/companyconcept/CIK0000320193/us-gaap/NetIncomeLoss.json
