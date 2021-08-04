# Modifier

Provides two different functions.

## Update from Query
```
go run go-modifier -o update -q=false
```
Will look at the queries specified in the .env file and attempt to modify the data that these queries return.
The -q flag (query only) indicates that you want to just run the query phase, not the update. 
So if you want to run this to just see how many records would be impacted, set that flag to true (or leave it off, it defaults to true)

If your query in the .env file is 
```
select id, Name, CreatedDate from Account
``` 
this command would attempt to generate random lipsum for the Acccount Name only (other two fields are not updateable in the metadata)

You can also specify multiple SOQL queries in this environment variable.
```
QUERIES=select Id, Name, Industry, Type from Account, select Id, FirstName from Contact, select Id, Subject, Status from case where Closed=false
``` 

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
One of the very helpful parts of this tool will update the parent reference and owner Ids to randomly selected IDs from the org you point it at. 

So for the above, it will also fetch random Account Ids to add to the CSV before insert. 
It will also fetch a list of users (that are standard and active) to set the ownerId field.

There is an optional switch on this command -fetch (fetchOnly). 
This will call to Mockaroo and fetch the data, update the relationship fields but not update the data. 