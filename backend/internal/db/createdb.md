1. Login to PostgreSQL as superuser
You might have to entire your password, on success you should
see a new prompt like: postgres=#
```bash
sudo -u postgres psql
```

2. Create a new PostgreSQL user
Replace user and password fields with your desired user/pass.
Make sure your password is in single quotes (ie 'realPassword123'). 
```sql
CREATE USER mason WITH PASSWORD 'realPassword123'; 
``` 

3. Create the mana database 
You can name it something other than mana_db if you wish. 
For this tutorial, replace mana_db with what you called it.
```sql
CREATE DATABASE mana_db;
```

4. Grant privileges to the user for this database 
```sql
GRANT ALL PRIVILEGES ON DATABASE mana_db TO mason;
```
5. Connect to the database and grant user Schema privileges
```sql
\c mana_db 
GRANT ALL PRIVILEGES ON SCHEMA public TO mason;
```

6. Exit PostgreSQL shell
```sql
\q 
```

7. Update PostgreSQL host and port 
```bash
cd /etc/postgresql
ls
```
You should see a directory in there, a number
for example if it was 17:
```bash
cd 17
```
Inside of this directory you will find a postgresql.conf file.
Scroll to the comment saying connection settings and you will you 
should see listen_addresses and port. Comment out their given
listen_addresses and port (add # in front) and on newlines for each value,
enter your desired listen_addresses and port. 

8. Update your .env file to reflect these changes
DB_NAME, DB_USER, DB_PASSWORD will changed based on what you entered 