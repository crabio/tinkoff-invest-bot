# Global Rank Parser

This is Python tool for downloading, parsing and saving into Parquet file Forbes global companies rank.

## Presequencies

* Python 3.7+

## Run

For running execute command: `python3 main.py`

If tool succesfully parsed all data, 200 code will be printed.

Output data by default will be saved into '../data/companies_rank.parquet'.

## TODOs

* Add argument to config limit of parsed companies
* Add argument to specify output file path
* Add loading to PostgreSQL DB
* Add argument flag to loading into DB or file
