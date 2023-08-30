#!/usr/bin/env python3
import csv
import sys
from datetime import datetime

def gen_header_comment():
    return "/* Schema tabella generato automaticamente \n" \
            f" * {datetime.date(datetime.now())}:{datetime.time(datetime.now())}\n" \
            " * Non modificare */\n"

def gen_create_query(tab_name, city_field, code_field):
    return f"CREATE TABLE {tab_name} (\n" \
           f"    {city_field} VARCHAR(64) NOT NULL,\n" \
           f"    {code_field} VARCHAR(4) NOT NULL PRIMARY KEY\n);\n\n"

def gen_insert_query(row, tab_name, city_field, code_field):
    return f"INSERT INTO {tab_name}({city_field},{code_field}) VALUES (\"{row[6]}\",\"{row[9]}\");\n"

def main():
    if len(sys.argv) != 3:
        print(f"Usage: ./{sys.argv[0]} <INPUT_FILE> <OUTPUT_FILE>")
        sys.exit(1)

    with open(sys.argv[1], 'r') as ftab, open(sys.argv[2], 'w') as fsql:
        fsql.write(gen_header_comment())
        fsql.write(gen_create_query("codNazioni", "State", "Code"))
        for row in csv.reader(ftab, delimiter=','):
            if row[9] == "n.d." or row[6] == "Denominazione IT":
                continue
            fsql.write(gen_insert_query(row, "codNazioni", "State", "Code"))

if __name__ == "__main__":
    main()