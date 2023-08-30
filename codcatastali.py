#!/usr/bin/env python3
import sys
import requests
from datetime import datetime
from bs4 import BeautifulSoup

def gen_header_comment():
    return "/* Schema tabella generato automaticamente \n" \
            f" * {datetime.date(datetime.now())}:{datetime.time(datetime.now())}\n" \
            " * Non modificare */\n"

def gen_create_query(tab_name, city_field, code_field):
    return f"CREATE TABLE {tab_name} (\n" \
           f"    {city_field} VARCHAR(64) NOT NULL,\n" \
           f"    {code_field} VARCHAR(4) NOT NULL PRIMARY KEY\n);\n\n"

def gen_insert_query(row, tab_name, city_field, code_field):
    return f"INSERT INTO {tab_name}({city_field},{code_field}) VALUES (\"{row[1].text}\",\"{row[5].text}\");\n"

def main():
    if len(sys.argv) != 2:
        print(f"Usage: ./{sys.argv[0]} <OUTPUT_FILE>")
        sys.exit(1)

    head = {
        "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"
    }
    url = "https://dait.interno.gov.it/territorio-e-autonomie-locali/sut/elenco_codici_comuni.php"
    page = requests.get(url, headers=head)
    content = BeautifulSoup(page.content, "html.parser")
    table = content.find("table", {"class": "table-striped"})

    with open(sys.argv[1], 'w') as fsql:
        fsql.write(gen_header_comment())
        fsql.write(gen_create_query("codCatastali", "City", "Code"))
        for row in table.tbody.find_all("tr"):
            td = row.find_all("td")
            fsql.write(gen_insert_query(td, "codCatastali", "City", "Code"))


if __name__ == "__main__":
    main()