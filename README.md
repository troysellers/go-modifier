# Modifier

Provides two different functions.

## Update from Query
```
go run go-modifier -o update -q=false
```
Will look at the queries specified in the .env file and attempt to modify the data that these queries return.
If the field in the query is modifiable. 

e.g. if your query in the .env file is 
```
select id, Name, CreatedDate from Account
``` 
this command would attempt to generate random lipsum for the Acccount Name only (other two fields are not updateable in the metadata)

## Creat from Mockaroo
There is a [mockaroo project](https://www.mockaroo.com/projects/25058) that has some default data sets defined, standard objects and fields 
* account
* contact
* case
* casecomment
* opportunity
* lead
* task
* event

You can specify the amount of records (up to 5000) and the type of object to create. 

```
go run go-modifier -o create -c 500 -s contact
```
will get 500 fake contacts from Mockaroo and attempt to insert them into Salesforce. 

