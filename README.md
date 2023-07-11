# Transaction Tool API

With this tool, you can report bank transactions to customers. It works by passing it a csv file with all year's transactions, and then it sends an email with a summary.

## Software requirements

It is necessary to have installed:

- Docker
- Docker-compose

## How it is work?

For now, the tool only offers one service (to report a summary by mail). This can be consumed throughout the curl:

`curl -X POST -H 'Content-Type:text/csv' --data-binary @example.csv http://localhost:8080/transaction-tool/resume/{user_id}
`

where "example.csv" is the csv transaction file and "user_id" is the id associated with the user to whom it wants to send the report.

## How does it launch the application?

You only need to go to the root of the project and do:

`docker-compose up
`

This will create two Docker containers (and Docker images), one for the Golang application that will be available throughout port 8080 and another for the MySQL database that will be available throughout port 3306. If some of these ports are not available, you need to change them in the docker-compose.yml file. 

### Before launching the application

It is necessary to set up the mail sender configuration by defining four environment variables that are established in docker-compose.yml file:

- NOTIFIER_SENDER: It is the email address that will send the report.
- NOTIFIER_PASSWORD: It is the password associated with the email address that will send the report. For Gmail, it has to be an app password ([how do I create one?](https://support.google.com/mail/answer/185833?hl=en)), but for others, you must find out.
- NOTIFIER_HOST: Host of the email address. By default, the Gmail host is established.
- NOTIFIER_PORT: Port of the email address. By default, the Gmail port is established.

### How does it know the receiver's email address?

The service needs a user ID; however, how does it know the receiver's email address to which I want to send the report? The API assumes that users are registered in a database because it would never send a bank transaction summary to an unknown user. For that reason, there is a table called user in which you need to register the user to whom you want to send the summary. You can register the user in two ways:

- Before launching the application: There is a file called init.sql (at the root of the project), which contains all previous statements that need the database after its Docker building. At the bottom, you will see a user insert statement for a ramdon user. You can replicate this for other users by specifying the user name, email address, and ID that you need for the report service.
- After launching the application: The tool does not expose services for user creation, but you can connect to it by using the database name, host, port, and user fields that are transaction, localhost, 3306, and root, respectively. then insert statements like previously mentioned.

### CSV file format

The file must contain only two columns: transaction amount, which is represented by a non-zero float number, and transaction date, which is represented by an RFC3339-formatted datetime. All transactions must belong to the same year. At the root of the project is an example called example.csv.

### What does the summary output contain?

- The general balance
- The transaction credit average 
- The transaction debit average 
- The total transactions by month