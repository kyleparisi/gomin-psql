import psycopg2
import os

db_database = os.getenv("DB_DATABASE")
dbconnect = psycopg2.connect(database=db_database)

cursor = dbconnect.cursor()
cursor.execute(f"""
SELECT
	table_name,
	ordinal_position,
	column_name,
	udt_name AS data_type,
	numeric_precision,
	datetime_precision,
	numeric_scale,
	character_maximum_length AS data_length,
	is_nullable,
	column_name AS CHECK,
	column_name AS check_constraint,
	column_default,
	column_name AS foreign_key
FROM
	information_schema.columns
WHERE
	table_schema = 'public'
	order by table_name, ordinal_position;
""")


def convert_to_go(name):
    if not name:
      return
    names = name.split("_")
    k = 0
    for v in names:
        names[k] = v.capitalize()
        k=k+1
    return "".join(names)


objects = {}

data = cursor.fetchall()
include_time = False
for row in data:
    table_name = convert_to_go(row[0])
    column_name = convert_to_go(row[2])
    nullable = row[9]
    data_type = row[3]
    if objects.get(table_name) is None:
        objects[table_name] = {}
    objects[table_name][column_name] = {}
    objects[table_name][column_name]["nullable"] = nullable == "NO"
    if data_type in ["longtext", "varchar"]:
        objects[table_name][column_name]["data_type"] = "string"

    if data_type in ["int"]:
        objects[table_name][column_name]["data_type"] = "int"

    if data_type in ["bigint", "int8"]:
        objects[table_name][column_name]["data_type"] = "int64"

    if data_type in ["float"]:
        objects[table_name][column_name]["data_type"] = "float64"

    if data_type in ["date", "datetime"]:
        include_time = True
        objects[table_name][column_name]["data_type"] = "time.Time"

print("""package app""")
if include_time:
    print('''
    import (
        "time"
    )''')

for key in objects:
    print(f"""
type {key} struct {{""")
    struct = objects[key]
    for attribute in struct:
        if struct[attribute].get("data_type") is None:
            continue
        data_type = struct[attribute]["data_type"]
        print(f"""  {attribute} {data_type}""")

    print(f"""}}""")

dbconnect.close()