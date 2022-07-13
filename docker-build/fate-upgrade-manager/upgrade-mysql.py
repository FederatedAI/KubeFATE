import mysql.connector
import sys
import os


def get_script_list(start_ver, target_ver):
    sql_script_names = os.listdir("sql")
    sql_script_names.sort()
    res = []
    start_index = -1
    end_index = -1
    for i in range(len(sql_script_names)):
        if sql_script_names[i].startswith(start_ver):
            start_index = i
        if sql_script_names[i].replace(".sql", "").endswith(target_ver):
            end_index = i
    if start_index == -1 or end_index == -1 or start_index > end_index:
        return res
    res = sql_script_names[start_index:end_index+1]
    print("will run scripts:")
    print(res)
    return res


def preprocess_script(script):
    queries = []
    query = ''
    delimiter = ';'
    with open("sql/%s" % script, "r") as sql_file:
        for line in sql_file.readlines():
            line = line.strip()
            if line.startswith('DELIMITER'):
                delimiter = line[10:]
            else:
                query += line+'\n'
                if line.endswith(delimiter):
                    # Get rid of the delimiter, remove any blank lines and add this query to our list
                    queries.append(query.strip().strip(delimiter))
                    query = ''
    return queries


def run_script(script, cursor):
    queries = preprocess_script(script)
    for query in queries:
        if not query.strip():
            continue
        print("execute query %s" % query)
        cursor.execute(query)


if __name__ == '__main__':
    _, user, password, start_ver, end_ver = sys.argv
    scripts_to_run = get_script_list(start_ver, end_ver)
    mydb = mysql.connector.connect(
        host="mysql",
        user=user,
        password=password,
        database="eggroll_meta"
    )
    cursor = mydb.cursor()
    for script in scripts_to_run:
        run_script(script, cursor)
    cursor.close()
