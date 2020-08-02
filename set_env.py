import sys 
import re

"""This script get PROJECT_ID, STORAGE_LOGIN and STORAGE_PASSWORD 
environment variables used on app engine from Github Secrets and
replace on app.yaml."""

app_engine_file = "app.yaml"

if __name__ == "__main__":
    if len(sys.argv) != 6:
        sys.exit("invalid number of arguments: {}".format(len(sys.argv)))
    project_id = sys.argv[1]
    storage_login = sys.argv[2]
    storage_password = sys.argv[3]
    email = sys.argv[4]
    password = sys.argv[5]
    file_content = ""
    with open (app_engine_file, "r") as file:
        app_engine_file_content = file.read()
        line = re.sub(r"##PROJECT_ID", project_id, app_engine_file_content) 
        line = re.sub(r"##STORAGE_LOGIN", storage_login, line)
        line = re.sub(r"##STORAGE_PASSWORD", storage_password, line)
        line = re.sub(r"##EMAIL", email, line)
        line = re.sub(r"##PASSWORD", password, line)
        file_content = line
    with open (app_engine_file, "w") as file:
        file.write(file_content)
